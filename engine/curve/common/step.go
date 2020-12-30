package step

import (
	"reflect"
	"sort"
	"sync"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

const (
	curveType = gcpb.CurveType_CURVE_TYPE_STEP_FLOAT
)

type ComparableList interface {
	sort.Interface
	Get(i int) interface{}
}

type Curve struct {
	curve.Base

	// mux guards the tick and data properties.
	mux  sync.Mutex
	tick id.Tick
	data ComparableList
}

func New(eid id.EntityID, tick id.Tick, property gcpb.EntityProperty, datumType reflect.Type) *Curve {
	return &Curve{
		Base: *curve.New(eid, curveType, datumType, property),
		tick: tick,
	}
}

func (c *Curve) Tick() id.Tick { return c.tick }

func (c *Curve) Get(tick id.Tick) interface{} {
	c.mux.Lock()
	defer c.mux.Unlock()
}
