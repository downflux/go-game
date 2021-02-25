// Package visitor defines interfaces necessary for the visitor design pattern.
//
// See https://en.wikipedia.org/wiki/Visitor_pattern for more information.
package visitor

import (
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

type Agent interface {
	// Accept conditionally allows the Visitor to mutate the Agent.
	//
	// Example:
	//  func (a *ConcreteAgent) Accept(v Vistor) { return v.Visit(a) }
	Accept(v Visitor) error
}

// Visitor defines the list of functions necessary for a process regularly
// mutating arbitrary Entity instances.
type Visitor interface {
	// Type returns a registered FSMType.
	Type() fcpb.FSMType

	// Visit will run appropriate commands for the current tick. If
	// a timeout occurs, the function will return early. This function
	// may be called concurrently by the game engine.
	//
	// TODO(minkezhang): implement timeout behavior.
	//
	// Visitors should never return an unimplemented error -- return
	// a no-op instead. This ensures Entity objects do not have to do
	// conditional branches in the Accept function.
	Visit(a Agent) error
}

type Base struct {
	fsmType fcpb.FSMType
}

func New(t fcpb.FSMType) *Base    { return &Base{fsmType: t} }
func (v Base) Type() fcpb.FSMType { return v.fsmType }
