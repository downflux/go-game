// Package entity declares the game Entity interface and common shared
// implementation details.
//
// Example
//
//  type ConcreteEntity struct {
//    entity.LifeCycle
//    ...
//  }
//
//  func (e *ConcreteEntity) Curve(...) { ... }
package entity

import (
	"sync"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/server/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type Entity interface {
	// Type returns the registered EntityType of the Entity.
	Type() gcpb.EntityType

	// ID returns the UUID of the Entity.
	ID() id.EntityID

	// Curve returns a Curve instance of a specific mutable property,
	// e.g. HP or position.
	//
	// TODO(minkezhang): Decide if we should return default value.
	Curve(t gcpb.EntityProperty) curve.Curve

	// Properties() returns list of entity properties defined in a specific
	// entity. This list is created at init time and is immutable.
	Properties() []gcpb.EntityProperty

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

// NoCurve implements a subset of the Entity interface and is used by
// Entity implementations which do not have any properties that may be tracked
// by curves.
type NoCurve struct{}

// Curve returns a Curve instance for the given EntityProperty. In the
// NoCurve implementation, this returns a trivially true nil value for
// all categories.
func (e *NoCurve) Curve(c gcpb.EntityProperty) curve.Curve { return nil }

// Properties returns a list of registered CurveCategory instances tracked
// by the Entity implementation. NoCurve will return an empty list.
func (e *NoCurve) Properties() []gcpb.EntityProperty { return nil }

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
