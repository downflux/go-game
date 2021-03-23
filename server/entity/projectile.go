package projectile

import (
	"reflect"

	"github.com/downflux/game/engine/curve/common/linearmove"
	"github.com/downflux/game/engine/curve/common/step"
	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/curve/list"
	"github.com/downflux/game/engine/entity/component/lifecycle"
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/entity/component/positionable"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	curvecomponent "github.com/downflux/game/engine/entity/component/curve"
)

const (
	// moveVelocity is measured in tiles per second.
	moveVelocity = 20
)

type (
	moveComponent      = moveable.Base
	positionComponent  = positionable.Base
	curveComponent     = curvecomponent.Component
	lifecycleComponent = lifecycle.Component
)

type Entity struct {
	entity.Base
	curveComponent
	lifecycleComponent
	moveComponent
	positionComponent
}

func New(eid id.EntityID, t id.Tick, pos *gdpb.Position, cid id.ClientID) (*Entity, error) {
	mc := linearmove.New(eid, t)
	mc.Add(t, pos)

	cidc := step.New(
		eid,
		t,
		gcpb.EntityProperty_ENTITY_PROPERTY_CLIENT_ID,
		reflect.TypeOf(id.ClientID("")),
	)
	cidc.Add(t, cid)

	curves, err := list.New([]curve.Curve{mc, cidc})
	if err != nil {
		return nil, err
	}

	return &Entity{
		Base: *entity.New(
			gcpb.EntityType_ENTITY_TYPE_TANK_PROJECTILE, eid, cidc),
		moveComponent:     *moveable.New(moveVelocity),
		positionComponent: *positionable.New(mc),
		curveComponent:    *curvecomponent.New(curves),
	}, nil
}
