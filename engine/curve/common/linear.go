// Package linear implements a linear curve type, i.e. that of a position
// curve traveling at constant velocity.
package linearmove

import (
	"reflect"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/curve/data"
	"github.com/downflux/game/engine/id/id"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

const (
	curveType = gcpb.CurveType_CURVE_TYPE_LINEAR
)

var (
	datumType = reflect.TypeOf(float64(0))
)

type Curve struct {
	curve.Base

	tick id.Tick
	data *data.Data
}

// New constructs an instance of a Curve.
func New(eid id.EntityID, tick id.Tick, property gcpb.EntityProperty) *Curve {
	return &Curve{
		Base: *curve.New(eid, curveType, datumType, property),
		tick: tick,
		data: data.New(nil),
	}
}

func (c *Curve) Data() *data.Data { return c.data }
func (c *Curve) Tick() id.Tick    { return c.tick }

// Add inserts a single datum point into the Curve.
//
// TODO(minkezhang): Add duplicate removal.
// TODO(minkezhang): Add point interpolation removal (if a < b < c have the
// same slopes, remove b).
func (c *Curve) Add(t id.Tick, v interface{}) error {
	c.data.Set(t, v)
	return nil
}

func (c *Curve) Merge(o curve.Curve) error {
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

func (c *Curve) Get(t id.Tick) interface{} {
	if c.data == nil {
		return float64(0)
	}

	if data.Before(t, c.data.Tick(0)) {
		return c.data.Get(c.data.Tick(0))
	}
	if data.Before(c.data.Tick(c.data.Len()-1), t) {
		return c.data.Get(c.data.Tick(c.data.Len() - 1))
	}

	i := c.data.Search(t)
	if i == 0 {
		return c.data.Get(c.data.Tick(0))
	}

	t0 := c.data.Tick(i - 1)
	t1 := c.data.Tick(i)
	p0 := c.data.Get(t0).(float64)
	p1 := c.data.Get(t1).(float64)

	tickDelta := t.Value() - t0.Value()

	dx := p1 - p0
	dt := t1.Value() - t0.Value()

	return p0 + dx*(tickDelta/dt)
}

func (c *Curve) Export(tick id.Tick) *gdpb.Curve {
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
			Datum: &gdpb.CurveDatum_DoubleDatum{c.data.Get(c.data.Tick(j)).(float64)},
		})
	}

	return pb
}
