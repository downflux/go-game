// Package visitor defines interfaces necessary for the visitor design pattern.
//
// See https://en.wikipedia.org/wiki/Visitor_pattern for more information.
package visitor

import (
	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/server/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	vcpb "github.com/downflux/game/server/visitor/api/constants_go_proto"
)

// Entity defines the list of functions necessary for working with the Visitor
// interface.
type Entity interface {
	// Accept conditionally allows the Visitor to mutate the Entity.
	//
	// Example:
	//  func (e *ConcreteEntity) Accept(v Vistor) { return v.Visit(e) }
	Accept(v Visitor) error

	// Type returns the registered EntityType of the Entity. This is useful
	// for the Visitor when determining Entity-specific mutations,
	// e.g. attacking infantry vs. attacking tank.
	Type() gcpb.EntityType

	// ID returns the UUID of the Entity.
	ID() id.EntityID

	// Curve returns a Curve instance of a specific mutable property,
	// e.g. HP or position. Visitors mutate these curves to reflect an
	// Entity change. The Curve changes are broadcasted to all clients once
	// per game tick.
	//
	// TODO(minkezhang): Decide if we should return default value.
	Curve(t gcpb.CurveCategory) curve.Curve

	// CurveCategories returns list of curve categories defined in a specific
	// entity. This list is created at init time and is immutable.
	CurveCategories() []gcpb.CurveCategory

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

// Visitor defines the list of functions necessary for a process regularly
// mutating arbitrary Entity instances.
type Visitor interface {
	// Type returns a registered VisitorType.
	Type() vcpb.VisitorType

	// Schedule adds a Visitor-specific command to the Visitor. This
	// function will be called concurrently by the game engine.
	Schedule(args interface{}) error

	// Visit will run appropriate commands for the current tick. If
	// a timeout occurs, the function will return early. This function
	// may be called concurrently by the game engine.
	//
	// TODO(minkezhang): implement timeout behavior.
	//
	// Visitors should never return an unimplemented error -- return
	// a no-op instead. This ensures Entity objects do not have to do
	// conditional branches in the Accept function.
	Visit(e Entity) error
}
