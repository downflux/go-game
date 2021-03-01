// Package entity declares the game Entity interface and common shared
// implementation details.
//
// Example
//
//  type ConcreteEntity struct {
//    entity.Base
//    entity.LifeCycle
//    ...
//  }
package entity

import (
	"sync"

	"github.com/downflux/game/engine/curve/common/step"
	"github.com/downflux/game/engine/curve/list"
	"github.com/downflux/game/engine/entity/acl"
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
	ACL() acl.ACL

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
	acl        acl.ACL         // Read-only.
	cidc       *step.Curve
}

func New(t gcpb.EntityType, eid id.EntityID, cidc *step.Curve, permission acl.Permission) *Base {
	return &Base{
		entityType: t,
		id:         eid,
		cidc:       cidc,
		acl:        *acl.New(cidc, permission),
	}
}

func (e Base) Type() gcpb.EntityType { return e.entityType }
func (e Base) ID() id.EntityID       { return e.id }
func (e Base) ACL() acl.ACL          { return e.acl }

// Export converts the static properties of the entity into a gdpb.Entity
// object. Note that dynamic properties (e.g. position) are not considered here.
// These properties must be manually converted via Curve.Export instead.
func (e Base) Export() *gdpb.Entity {
	return &gdpb.Entity{
		EntityId: e.ID().Value(),
		Type:     e.Type(),
		Acl:      e.ACL().Export(),
	}
}

// LifeCycle implements a subset of the Entity interface concerned with
// tracking the lifecycle of the Entity. Entities such as tanks are created
// inside a factory, and are destroyed at the end of the game or when attacked
// by another Entity.
type LifeCycle struct {
	lifetimeMux sync.RWMutex
	start       id.Tick
	end         id.Tick
}

// Start returns the tick at which the Entity is spawned. The tick is set in
// the constructor (delegated to each concrete impementation).
func (e *LifeCycle) Start() id.Tick {
	e.lifetimeMux.RLock()
	defer e.lifetimeMux.RUnlock()

	return e.start
}

// End returns the tick at which the Entity was destroyed. Since the game state
// is append-only, the instance itself is not removed from the internal list,
// hence the need for this marker.
func (e *LifeCycle) End() id.Tick {
	e.lifetimeMux.RLock()
	defer e.lifetimeMux.RUnlock()

	return e.end
}

// Delete marks the target Entity as having been destroyed.
func (e *LifeCycle) Delete(tick id.Tick) {
	e.lifetimeMux.Lock()
	defer e.lifetimeMux.Unlock()

	e.end = tick
}
