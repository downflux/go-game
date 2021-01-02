// Package tank encapsulates logic for a basic tank unit.
package tank

import (
	"github.com/downflux/game/engine/curve/common/linearmove"
	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/curve/list"
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/server/entity/component/attackable"
	"github.com/downflux/game/server/entity/component/moveable"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

const (
	// velocity is measured in tiles per second.
	velocity = 2

	strength = 2
)

type moveComponent = moveable.Base
type attackComponent = attackable.Base

// Entity implements the entity.Entity interface and represents a simple armored
// unit.
type Entity struct {
	entity.LifeCycle
	*moveComponent
	*attackComponent

	// eid is a UUID of the Entity.
	eid id.EntityID

	// curves is a list of Curves tracking the Entity properties.
	curves *list.List
}

// New constructs a new instance of the Tank.
func New(eid id.EntityID, t id.Tick, p *gdpb.Position) (*Entity, error) {
	mc := linearmove.New(eid, t)
	mc.Add(t, p)

	curves, err := list.New([]curve.Curve{mc})
	if err != nil {
		return nil, err
	}

	return &Entity{
		moveComponent:   moveable.New(mc, velocity),
		attackComponent: attackable.New(strength),
		eid:             eid,
		curves:          curves,
	}, nil
}

// ID returns the UUID of the Tank.
func (e *Entity) ID() id.EntityID { return e.eid }

func (e *Entity) Curves() *list.List { return e.curves }

// Type returns the registered EntityType.
func (e *Entity) Type() gcpb.EntityType { return gcpb.EntityType_ENTITY_TYPE_TANK }
