// Package tank encapsulates logic for a basic tank unit.
package tank

import (
	"reflect"

	"github.com/downflux/game/engine/curve/common/delta"
	"github.com/downflux/game/engine/curve/common/linearmove"
	"github.com/downflux/game/engine/curve/common/step"
	"github.com/downflux/game/engine/curve/common/timer"
	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/curve/list"
	"github.com/downflux/game/engine/entity/component/lifecycle"
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/server/entity/component/attackable"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/entity/component/positionable"
	"github.com/downflux/game/server/entity/component/targetable"
	"github.com/downflux/game/server/entity/projectile"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	curvecomponent "github.com/downflux/game/engine/entity/component/curve"
)

const (
	// moveVelocity is measured in tiles per second.
	moveVelocity = 2

	// attackVelocity is measured in tiles per second.
	attackVelocity = 10

	// cooloff is measured in ticks.
	// TODO(minkezhang): Refactor to be in terms of seconds instead.
	cooloff = id.Tick(10)

	strength    = 2
	attackRange = 2

	health = float64(100)
)

type (
	moveComponent      = moveable.Base
	attackComponent    = attackable.Base
	targetComponent    = targetable.Base
	positionComponent  = positionable.Base
	lifecycleComponent = lifecycle.Component
	curveComponent     = curvecomponent.Component
)

// Entity implements the entity.Entity interface and represents a simple armored
// unit.
type Entity struct {
	entity.Base
	moveComponent
	attackComponent
	targetComponent
	positionComponent
	lifecycleComponent
	curveComponent
}

// New constructs a new instance of the Tank.
func New(
	eid id.EntityID,
	t id.Tick,
	pos *gdpb.Position,
	cid id.ClientID,
	proj *projectile.Entity) (*Entity, error) {
	mc := linearmove.New(eid, t)
	mc.Add(t, pos)
	ac := timer.New(eid, t, cooloff, gcpb.EntityProperty_ENTITY_PROPERTY_ATTACK_TIMER)
	tc := step.New(eid, t, gcpb.EntityProperty_ENTITY_PROPERTY_ATTACK_TARGET, reflect.TypeOf(id.ClientID("")))

	cidc := step.New(
		eid,
		t,
		gcpb.EntityProperty_ENTITY_PROPERTY_CLIENT_ID,
		reflect.TypeOf(id.ClientID("")),
	)
	cidc.Add(t, cid)

	hp := delta.New(step.New(eid, t, gcpb.EntityProperty_ENTITY_PROPERTY_HEALTH, reflect.TypeOf(float64(0))))
	if err := hp.Add(t, health); err != nil {
		return nil, err
	}

	curves, err := list.New([]curve.Curve{mc, ac, hp, cidc})
	if err != nil {
		return nil, err
	}

	return &Entity{
		Base: *entity.New(
			gcpb.EntityType_ENTITY_TYPE_TANK, eid, cidc),

		moveComponent: *moveable.New(moveVelocity),
		attackComponent: *attackable.New(
			strength, attackRange, attackVelocity, tc, ac, proj),
		targetComponent:   *targetable.New(hp),
		positionComponent: *positionable.New(mc),
		curveComponent:    *curvecomponent.New(curves),
	}, nil
}
