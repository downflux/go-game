// Package entity is a shared utility library for common implementations of the
// visitor.Entity interface. Concrete implementations of the Entity interface
// may inherit from this package as necessary.
//
// Example
//
//  type ConcreteEntity struct {
//    entity.BaseEntity
//    ...
//  }
//
//  func (e *ConcreteEntity) Curve(...) { ... }
package entity

import (
	"sync"

	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/server/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

// NoCurveEntity implements a subset of the Entity interface and is used by
// Entity implementations which do not have any properties that may be tracked
// by curves.
//
// TODO(minkezhang): Rename NoCurve.
type NoCurveEntity struct{}

// Curve returns a Curve instance for the given CurveCategory. In the
// NoCurveEntity implementation, this returns a trivially true nil value for
// all categories.
func (e *NoCurveEntity) Curve(c gcpb.CurveCategory) curve.Curve { return nil }

// CurveCategory returns a list of registered CurveCategory instances tracked
// by the Entity implementation. NoCurveEntity will return an empty list.
func (e *NoCurveEntity) CurveCategories() []gcpb.CurveCategory { return nil }

// BaseEntity implements a subset of the Entity interface concerned with
// tracking the lifecycle of the Entity. Entities such as tanks are created
// inside a factory, and are destroyed at the end of the game or when attacked
// by another Entity.
//
// TODO(minkezhang): Rename Lifecycle.
type BaseEntity struct {
	lifetimeMux sync.RWMutex
	start       id.Tick
	end         id.Tick
}

// Start returns the tick at which the Entity is spawned. The tick is set in
// the constructor (delegated to each concrete impementation).
func (e *BaseEntity) Start() id.Tick {
	e.lifetimeMux.RLock()
	defer e.lifetimeMux.RUnlock()

	return e.start
}

// End returns the tick at which the Entity was destroyed. Since the game state
// is append-only, the instance itself is not removed from the internal list,
// hence the need for this marker.
func (e *BaseEntity) End() id.Tick {
	e.lifetimeMux.RLock()
	defer e.lifetimeMux.RUnlock()

	return e.end
}

// Delete marks the target Entity as having been destroyed.
func (e *BaseEntity) Delete(tick id.Tick) {
	e.lifetimeMux.Lock()
	defer e.lifetimeMux.Unlock()

	e.end = tick
}
