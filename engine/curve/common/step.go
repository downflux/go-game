package step

import (
	"reflect"
	"sync"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/curve/data"
	"github.com/downflux/game/engine/id/id"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

const (
	curveType = gcpb.CurveType_CURVE_TYPE_STEP_FLOAT
)

type Curve struct {
	curve.Base

	// mux guards the tick and data properties.
	mux  sync.RWMutex
	tick id.Tick
	data *data.Data
}

func New(eid id.EntityID, tick id.Tick, property gcpb.EntityProperty, datumType reflect.Type) *Curve {
	return &Curve{
		Base: *curve.New(eid, curveType, datumType, property),
		tick: tick,
		data: data.New(nil),
	}
}

func (c *Curve) Tick() id.Tick {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.tick
}

func (c *Curve) Get(tick id.Tick) interface{} {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if c.data == nil || c.data.Len() == 0 {
		return reflect.Zero(c.Base.DatumType()).Interface()
	}

	i := c.data.Search(tick)
	if i == c.data.Len() {
		i = i - 1
	} else if c.data.Tick(i) != tick {
		if i == 0 {
			return reflect.Zero(c.Base.DatumType()).Interface()
		} else {
			return c.data.Get(c.data.Tick(i - 1))
		}
	}
	return c.data.Get(c.data.Tick(i))
}

func (c *Curve) Add(tick id.Tick, value interface{}) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.data.Set(tick, value.(float64))
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

	return c.data.Merge(o.(*Curve).data)
}

// Export builds a gdpb.Curve instance for data yet to be communicated
// to the client.
//
// Export tail will include in the Curve returned a single point before the
// tick -- this allows clients to extrapolate the current position of an
// entity if input tick does not fall on an exact data point.
func (c *Curve) Export(tick id.Tick) *gdpb.Curve {
	c.mux.RLock()
	defer c.mux.RUnlock()

	pb := c.Base.Export()
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
