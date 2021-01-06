// Package tank encapsulates logic for a basic tank unit.
package tank

import (
	"reflect"
	"log"

	"github.com/downflux/game/engine/curve/common/delta"
	"github.com/downflux/game/engine/curve/common/linearmove"
	"github.com/downflux/game/engine/curve/common/step"
	"github.com/downflux/game/engine/curve/common/timer"
	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/curve/list"
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/server/entity/component/attackable"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/entity/component/positionable"
	"github.com/downflux/game/server/entity/component/targetable"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

const (
	// velocity is measured in tiles per second.
	velocity = 2

	// cooloff is measured in ticks.
	// TODO(minkezhang): Refactor to be in terms of seconds instead.
	cooloff = id.Tick(10)

	strength    = 2
	attackRange = 2

	health = 100
)

type moveComponent = moveable.Base
type attackComponent = attackable.Base
type targetComponent = targetable.Base
type positionComponent = positionable.Base

// Entity implements the entity.Entity interface and represents a simple armored
// unit.
type Entity struct {
	entity.LifeCycle
	moveComponent
	attackComponent
	targetComponent
	positionComponent

	// eid is a UUID of the Entity.
	eid id.EntityID

	// curves is a list of Curves tracking the Entity properties.
	curves *list.List
}

// New constructs a new instance of the Tank.
func New(eid id.EntityID, t id.Tick, p *gdpb.Position) (*Entity, error) {
	mc := linearmove.New(eid, t)
	mc.Add(t, p)
	ac := timer.New(eid, t, cooloff, gcpb.EntityProperty_ENTITY_PROPERTY_ATTACK_TIMER)
	hp := delta.New(step.New(eid, t, gcpb.EntityProperty_ENTITY_PROPERTY_HEALTH, reflect.TypeOf(float64(0))))
	log.Println(hp.Add(t, health))
	log.Printf("Debug: new tank health == %v", hp.Get(t))

	curves, err := list.New([]curve.Curve{mc, ac, hp})
	if err != nil {
		return nil, err
	}

	return &Entity{
		moveComponent:     *moveable.New(velocity),
		attackComponent:   *attackable.New(strength, attackRange, ac),
		targetComponent:   *targetable.New(hp),
		positionComponent: *positionable.New(mc),
		eid:               eid,
		curves:            curves,
	}, nil
}

// ID returns the UUID of the Tank.
func (e *Entity) ID() id.EntityID { return e.eid }

func (e *Entity) Curves() *list.List { return e.curves }

// Type returns the registered EntityType.
func (e *Entity) Type() gcpb.EntityType { return gcpb.EntityType_ENTITY_TYPE_TANK }
