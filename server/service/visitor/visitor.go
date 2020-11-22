package visitor

import (
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
	//
	// Visitors should never return an unimplemented error -- return
	// a no-op instead. This ensures Entity objects do not have to do
	// conditional branches in the Accept function.
	Visit(e Entity) error
}
