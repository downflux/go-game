// Package linearmove implements a specific curve type, i.e. that of a
// position curve traveling at constant velocity.
package linearmove

import (
	"reflect"
	"sync"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/curve/data"
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
	datumType = reflect.TypeOf(&gdpb.Position{})
)

// Curve implements a curve.Curve which represents the physical location
// of a specific entity.
type Curve struct {
	curve.Base

	// mux guards the tick and data properties.
	mux  sync.RWMutex
	tick id.Tick
	data *data.Data
}

// New constructs an instance of a Curve.
func New(eid id.EntityID, tick id.Tick) *Curve {
	return &Curve{
		Base: *curve.New(eid, curveType, datumType, property),
		tick: tick,
		data: data.New(nil),
	}
}

func (c *Curve) Data() *data.Data {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.data
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
//
// TODO(minkezhang): Add duplicate removal.
// TODO(minkezhang): Add point interpolation removal (if a < b < c have the
// same slopes, remove b).
func (c *Curve) Add(t id.Tick, v interface{}) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.data.Set(t, v)
	return nil
}

// Merge takes as input another Curve of the same type and replaces any
// data in the original Curve which occurs after the earliest element of the
// replacement Curve. In the game, this will occur when the original Curve
// predicts too far in the future.
//
// This is not technically thread-safe -- the mutex for the other curve is not
// acquired. Special care should be taken that the other input curve is a
// temporary struct.
func (c *Curve) Merge(o curve.Curve) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	// Only replace the tail if the candidate curve has been updated after
	// the current curve.
	if c.tick > o.Tick() {
		return nil
	}
	c.tick = o.Tick()

	if o.Type() != c.Type() {
		return status.Errorf(codes.FailedPrecondition, "cannot merge curves of type %v and %v", c.Type(), o.Type())
	}

	return c.data.Merge(o.Data())
}

// Get queries the Curve at a specific point for an interpolated value.
func (c *Curve) Get(t id.Tick) interface{} {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if c.data == nil {
		return &gdpb.Position{}
	}

	if data.Before(t, c.data.Tick(0)) {
		return proto.Clone(
			c.data.Get(c.data.Tick(0)).(*gdpb.Position),
		).(*gdpb.Position)
	}
	if data.Before(c.data.Tick(c.data.Len()-1), t) {
		return proto.Clone(
			c.data.Get(c.data.Tick(c.data.Len() - 1)).(*gdpb.Position),
		).(*gdpb.Position)
	}

	i := c.data.Search(t)
	if i == 0 {
		return proto.Clone(
			c.data.Get(c.data.Tick(0)).(*gdpb.Position),
		).(*gdpb.Position)
	}

	t0 := c.data.Tick(i - 1)
	t1 := c.data.Tick(i)
	p0 := c.data.Get(t0).(*gdpb.Position)
	p1 := c.data.Get(t1).(*gdpb.Position)

	tickDelta := t.Value() - t0.Value()

	dx := p1.GetX() - p0.GetX()
	dy := p1.GetY() - p0.GetY()
	dt := t1.Value() - t0.Value()

	return &gdpb.Position{
		X: p0.GetX() + dx*(tickDelta/dt),
		Y: p0.GetY() + dy*(tickDelta/dt),
	}
}

// Export builds a gdpb.Curve instance for data yet to be communicated
// to the client.
//
// Export will include in the Curve returned a single point before the
// tick -- this allows clients to extrapolate the current position of an
// entity if input tick does not fall on an exact data point.
func (c *Curve) Export(tick id.Tick) *gdpb.Curve {
	c.mux.RLock()
	defer c.mux.RUnlock()

	pb := c.Base.Export(tick)
	pb.Tick = c.Tick().Value()

	i := c.data.Search(tick)
	// If tick is a very large number, still include at minimum the last
	// known position of an entity.
	if i == c.data.Len() {
		i = c.data.Len() - 1
	}
	// If the tick falls in between two indices, return the smaller index
	// as we still need to interpolate the position until time passes to
	// the larger tick.
	if (c.data.Tick(i) > tick) && (i > 0) {
		i -= 1
	}

	for j := i; j < c.data.Len(); j++ {
		pb.Data = append(pb.GetData(), &gdpb.CurveDatum{
			Tick:  c.data.Tick(j).Value(),
			Datum: &gdpb.CurveDatum_PositionDatum{c.data.Get(c.data.Tick(j)).(*gdpb.Position)},
		})
	}

	return pb
}
