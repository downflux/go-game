package linearmove

import (
	"reflect"
	"sort"
	"sync"

	"github.com/downflux/game/curve/curve"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

const (
	curveType = gcpb.CurveType_CURVE_TYPE_LINEAR_MOVE
)

var (
	datumType      = reflect.TypeOf(&gdpb.Position{})
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

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
	tick float64

	// value must be a clone of the input and is considered immutable
	value *gdpb.Position
}

// datumBefore compares two data points and checks if d1 precedes d2.
func datumBefore(d1, d2 datum) bool {
	return d1.tick < d2.tick
}

// Curve implements a curve.Curve which represents the physical location
// of a specific entity.
type Curve struct {
	// id is read-only and not alterable after construction
	id string

	// entityID is read-only and not alterable after construction
	entityID string

	dataMux sync.RWMutex
	data    []datum

	deltaMux sync.Mutex
	// delta keeps track of data that have not yet been communicated to the
	// client yet.
	//
	// TODO(minkezhang): add enableDelta bool field -- clients should not
	// be keeping track of this field.
	delta []datum
}

// New constructs an instance of a Curve.
func New(id, eid string) *Curve {
	return &Curve{
		id:       id,
		entityID: eid,
	}
}

// Type returns the type of the Curve, which govners e.g. the interpolation,
// data interpretation, etc.
func (c *Curve) Type() gcpb.CurveType { return curveType }

// ID returns the Curve ID.
func (c *Curve) ID() string { return c.id }

// EntityID returns the ID of the parent Entity.
func (c *Curve) EntityID() string { return c.entityID }

// DatumType returns the type of the datum value.
func (c *Curve) DatumType() reflect.Type { return datumType }

// addDatumUnsafe adds a single datum point into the Curve, but does not hold
// the required c.dataMux; the caller is responsible for acquiring this lock.
//
// TODO(minkezhang): Add duplicate removal.
//
// TODO(minkezhang): Add point interpolation removal (if a < b < c have the
// same slopes, remove b).
func (c *Curve) addDatumUnsafe(t float64, v interface{}) error {
	d := datum{tick: t, value: proto.Clone(v.(*gdpb.Position)).(*gdpb.Position)}

	c.data = insert(c.data, d)

	// Add to delta cache for broadcasting.
	c.deltaMux.Lock()
	c.delta = append(c.delta, d)
	c.deltaMux.Unlock()

	// TODO(minkezhang): Add data validation.
	return nil
}

// cloneData exposes a concurrency-safe copy of the internal Curve data.
func (c *Curve) cloneData() []datum {
	c.dataMux.RLock()
	defer c.dataMux.RUnlock()

	res := make([]datum, len(c.data))
	copy(res, c.data)

	return res
}

// Add inserts a single datum point into the Curve.
func (c *Curve) Add(t float64, v interface{}) error {
	c.dataMux.Lock()
	defer c.dataMux.Unlock()

	return c.addDatumUnsafe(t, v)
}

// ReplaceTail takes as input another Curve of the same type and replaces any
// data in the original Curve which occurs after the earliest element of the
// replacement Curve. In the game, this will occur when the original Curve
// predicts too far in the future.
func (c *Curve) ReplaceTail(o curve.Curve) error {
	c.dataMux.Lock()
	defer c.dataMux.Unlock()

	switch o.Type() {
	case gcpb.CurveType_CURVE_TYPE_LINEAR_MOVE:
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
	default:
		return notImplemented
	}
	return nil
}

// Get queries the Curve at a specific point for an interpolated value.
func (c *Curve) Get(t float64) interface{} {
	c.dataMux.RLock()
	defer c.dataMux.RUnlock()

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

	tickDelta := t - c.data[i-1].tick
	return &gdpb.Position{
		X: c.data[i-1].value.GetX() + (c.data[i].value.GetX()-c.data[i-1].value.GetX())/(c.data[i].tick-c.data[i-1].tick)*tickDelta,
		Y: c.data[i-1].value.GetY() + (c.data[i].value.GetY()-c.data[i-1].value.GetY())/(c.data[i].tick-c.data[i-1].tick)*tickDelta,
	}
}

// ExportDelta builds a gdpb.Curve instance for data yet to be communicated
// to the client.
func (c *Curve) ExportDelta() (*gdpb.Curve, error) {
	c.deltaMux.Lock()
	delta := c.delta
	c.delta = nil
	c.deltaMux.Unlock()

	pb := &gdpb.Curve{
		CurveId:  c.ID(),
		Type:     c.Type(),
		EntityId: c.EntityID(),
	}

	for _, d := range delta {
		pb.Data = append(pb.GetData(), &gdpb.CurveDatum{
			Tick:  d.tick,
			Datum: &gdpb.CurveDatum_PositionDatum{d.value},
		})
	}

	return pb, nil
}
