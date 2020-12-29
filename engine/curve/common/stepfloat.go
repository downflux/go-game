package stepfloat

import (
	"reflect"
	"sync"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

const (
	curveType = gcpb.CurveType_CURVE_TYPE_STEP_FLOAT
)

var (
	datumType = reflect.TypeOf(float64(0))
)

type Curve struct {
	curve.Base

	// mux guards the tick property.
	mux  sync.Mutex
	tick id.Tick
}

func New(eid id.EntityID, tick id.Tick, property gcpb.EntityProperty) *Curve {
	return &Curve{
		Base: *curve.New(eid, curveType, datumType, property),
		tick: tick,
	}
}

func (c *Curve) Tick() id.Tick { return c.tick }
