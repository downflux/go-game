// Package entity declares the game Entity interface and common shared
// implementation details.
//
// Example
//
//  type (
//    lifecycleComponent = lifecycle.Component
//    curveComponent = curve.Component
//  )
//  type ConcreteEntity struct {
//    entity.Base
//    lifecycleComponent
//    curveComponent
//    ...
//  }
package entity

import (
	"github.com/downflux/game/engine/curve/common/step"
	"github.com/downflux/game/engine/curve/list"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

type Entity interface {
	// Type returns the registered EntityType of the Entity.
	Type() gcpb.EntityType

	// ID returns the UUID of the Entity.
	ID() id.EntityID

	Curves() *list.List
	Export() *gdpb.Entity

	// Start returns the game tick at which the Entity was created.
	Start() id.Tick

	// End returns the game tick at which the Entity was destroyed.
	// Entities are never deleted by the game. Entities which are marked
	// as destroyed must not be mutated again -- for units like revived
	// units, we should instead either create new entities, or make sure
	// revivable units do not actually call End.
	End() id.Tick

	// Delete marks the Entity as permanently non-relevant for the current
	// game. This may occur when the HP is set to zero, etc.
	Delete(tick id.Tick)
}

type Base struct {
	entityType gcpb.EntityType // Read-only.
	id         id.EntityID     // Read-only.
	cidc       *step.Curve
}

func New(t gcpb.EntityType, eid id.EntityID, cidc *step.Curve) *Base {
	return &Base{
		entityType: t,
		id:         eid,
		cidc:       cidc,
	}
}

func (e Base) Type() gcpb.EntityType { return e.entityType }
func (e Base) ID() id.EntityID       { return e.id }

// Export converts the static properties of the entity into a gdpb.Entity
// object. Note that dynamic properties (e.g. position) are not considered here.
// These properties must be manually converted via Curve.Export instead.
func (e Base) Export() *gdpb.Entity {
	return &gdpb.Entity{
		EntityId: e.ID().Value(),
		Type:     e.Type(),
	}
}
