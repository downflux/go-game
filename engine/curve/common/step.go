package step

import (
	"reflect"
	"sync"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/curve/data"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

const (
	curveType = gcpb.CurveType_CURVE_TYPE_STEP_FLOAT
)

type Curve struct {
	curve.Base

	// mux guards the tick and data properties.
	mux  sync.Mutex
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
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.tick
}

func (c *Curve) Get(tick id.Tick) interface{} {
	c.mux.Lock()
	defer c.mux.Unlock()

	return nil
}

func (c *Curve) Add(tick id.Tick, value interface{}) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.data.Set(tick, value.(float64))
	return nil
}
