package visitor

import (
	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/server/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	vcpb "github.com/downflux/game/server/visitor/api/constants_go_proto"
)

type Entity interface {
	Accept(v Visitor) error
	Type() gcpb.EntityType

	ID() id.EntityID

	// TODO(minkezhang): Decide if we should return default value.
	Curve(t gcpb.CurveCategory) curve.Curve

	// CurveCategories returns list of curve categories defined in a specific
	// entity. This list is created at init time and is immutable.
	CurveCategories() []gcpb.CurveCategory

	Start() id.Tick
	End() id.Tick

	Delete(tick id.Tick)
}

type Visitor interface {
	Type() vcpb.VisitorType

	// Schedule adds a Visitor-specific command to the Visitor. This
	// function will be called concurrently by the game engine.
	Schedule(args interface{}) error

	// Visit will run appropriate commands for the current tick. If
	// a timeout occurs, the function will return early. This function
	// may be called concurrently by the game engine.
	//
	// Visitors should never return an unimplemented error -- return
	// a no-op instead. This ensures Entity objects do not have to do
	// conditional branches in the Accept function.
	Visit(e Entity) error
}
