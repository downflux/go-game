package visitor

import (
	"github.com/downflux/game/curve/curve"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type Entity interface {
	Accept(v Visitor) error
	Type() gcpb.EntityType
}

type Visitor interface {
	// Schedule adds a Visitor-specific command to the Visitor. This
	// function will be called concurrently by the game engine.
	Schedule(args interface{}) error

	// Visit will run appropriate commands for the current tick. If
	// a timeout occurs, the function will return early. This function
	// may be called concurrently by the game engine.
	Visit(e Entity) ([]curve.Curve, error)
}
