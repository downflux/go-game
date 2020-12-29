// Package linearmove implements a specific curve type, i.e. that of a
// position curve traveling at constant velocity.
package linearmove

import (
	"reflect"
	"sort"
	"sync"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/id/id"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

const (
	// TODO(minkezhang): Untether the static property with this curve type.
	property = gcpb.EntityProperty_ENTITY_PROPERTY_POSITION

	curveType = gcpb.CurveType_CURVE_TYPE_LINEAR_MOVE
)

var (
	datumType      = reflect.TypeOf(&gdpb.Position{})
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

// Curve implements a curve.Curve which represents the physical location
// of a specific entity.
type Curve struct {
	curve.Base

	// mux guards the tick and data properties.
	mux  sync.RWMutex
	tick id.Tick
	data []datum
}

// New constructs an instance of a Curve.
func New(eid id.EntityID, tick id.Tick) *Curve {
	return &Curve{
		Base: *curve.New(eid, curveType, datumType, property),
		tick: tick,
	}
}

// Tick returns the last server tick at which the curve was updated and
// current. Values along the parametric curve past this tick should be
// considered non-authoritative.
func (c *Curve) Tick() id.Tick {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.tick
}

// Add inserts a single datum point into the Curve.
func (c *Curve) Add(t id.Tick, v interface{}) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.addDatumUnsafe(t, v)
}

// ReplaceTail takes as input another Curve of the same type and replaces any
// data in the original Curve which occurs after the earliest element of the
// replacement Curve. In the game, this will occur when the original Curve
// predicts too far in the future.
func (c *Curve) ReplaceTail(o curve.Curve) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.tick > o.Tick() {
		return nil
	}
	c.tick = o.Tick()

	if o.Type() != gcpb.CurveType_CURVE_TYPE_LINEAR_MOVE {
		return notImplemented
	}

	data := o.(*Curve).cloneData()
	// We need to delete the struct because of memory leaks from
	// the pointer stored at datum.value to gdpb.Position.
	//
	// See https://github.com/golang/go/wiki/SliceTricks.
	if len(data) > 0 {
		i := sort.Search(len(c.data), func(i int) bool { return !datumBefore(c.data[i], datum{tick: data[0].tick}) })
		if i < len(c.data) {
			for j := i; j < len(c.data); j++ {
				c.data[j] = datum{}
			}
		}
		c.data = c.data[:i]
	}
	for _, d := range data {
		c.addDatumUnsafe(d.tick, d.value)
	}
	return nil
}

// Get queries the Curve at a specific point for an interpolated value.
func (c *Curve) Get(t id.Tick) interface{} {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if c.data == nil {
		return &gdpb.Position{}
	}

	if datumBefore(datum{tick: t}, c.data[0]) {
		return proto.Clone(c.data[0].value).(*gdpb.Position)
	}

	if datumBefore(c.data[len(c.data)-1], datum{tick: t}) {
		return proto.Clone(c.data[len(c.data)-1].value).(*gdpb.Position)
	}

	i := sort.Search(len(c.data), func(i int) bool { return !datumBefore(c.data[i], datum{tick: t}) })

	if i == 0 {
		return proto.Clone(c.data[0].value).(*gdpb.Position)
	}

	tickDelta := t.Value() - c.data[i-1].tick.Value()
	return &gdpb.Position{
		X: c.data[i-1].value.GetX() + (c.data[i].value.GetX()-c.data[i-1].value.GetX())/(c.data[i].tick.Value()-c.data[i-1].tick.Value())*tickDelta,
		Y: c.data[i-1].value.GetY() + (c.data[i].value.GetY()-c.data[i-1].value.GetY())/(c.data[i].tick.Value()-c.data[i-1].tick.Value())*tickDelta,
	}
}

// ExportTail builds a gdpb.Curve instance for data yet to be communicated
// to the client.
//
// Export tail will include in the Curve returned a single point before the
// tick -- this allows clients to extrapolate the current position of an
// entity if input tick does not fall on an exact data point.
func (c *Curve) ExportTail(tick id.Tick) *gdpb.Curve {
	c.mux.RLock()
	defer c.mux.RUnlock()

	pb := &gdpb.Curve{
		Type:     c.Type(),
		Property: c.Property(),
		EntityId: c.EntityID().Value(),
		Tick:     c.Tick().Value(),
	}

	i := sort.Search(len(c.data), func(i int) bool { return !datumBefore(c.data[i], datum{tick: tick}) })
	// If tick is a very large number, still include at minimum the last
	// known position of an entity.
	if i > len(c.data)-1 {
		i = len(c.data) - 1
	}
	// If the tick falls in between two indices, return the smaller index
	// as we still need to interpolate the position until time passes to
	// the larger tick.
	if (c.data[i].tick > tick) && (i > 0) {
		i -= 1
	}

	for i := i; i < len(c.data); i++ {
		pb.Data = append(pb.GetData(), &gdpb.CurveDatum{
			Tick:  c.data[i].tick.Value(),
			Datum: &gdpb.CurveDatum_PositionDatum{c.data[i].value},
		})
	}

	return pb
}

// insert adds a datum object into a sorted list of data.
func insert(l []datum, d datum) []datum {
	i := sort.Search(len(l), func(i int) bool { return !datumBefore(l[i], d) })

	// Override existing value if the given input will result in an invalid
	// function.
	if i < len(l) && l[i].tick == d.tick {
		l[i] = d
	} else {
		l = append(l, datum{})
		copy(l[i+1:], l[i:])
		l[i] = d
	}
	return l
}

// datum represents a specific metric at a specific tick.
type datum struct {
	tick id.Tick

	// value must be a clone of the input and is considered immutable
	value *gdpb.Position
}

// datumBefore compares two data points and checks if d1 precedes d2.
func datumBefore(d1, d2 datum) bool {
	return d1.tick < d2.tick
}

// addDatumUnsafe adds a single datum point into the Curve, but does not hold
// the required c.mux; the caller is responsible for acquiring this lock.
//
// TODO(minkezhang): Add duplicate removal.
//
// TODO(minkezhang): Add point interpolation removal (if a < b < c have the
// same slopes, remove b).
func (c *Curve) addDatumUnsafe(t id.Tick, v interface{}) error {
	d := datum{tick: t, value: proto.Clone(v.(*gdpb.Position)).(*gdpb.Position)}

	c.data = insert(c.data, d)

	// TODO(minkezhang): Add data validation.
	return nil
}

// cloneData exposes a concurrency-safe copy of the internal Curve data.
func (c *Curve) cloneData() []datum {
	c.mux.RLock()
	defer c.mux.RUnlock()

	res := make([]datum, len(c.data))
	copy(res, c.data)

	return res
}
