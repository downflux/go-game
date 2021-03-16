// Package chase defines the Action used for carrying out the Chase command.
//
// A Pending state indicates the underlying move command is in Pending state
// and a move is scheduled for future execution.
//
// An OutOfRange state indicates the target is too far away.
package chase

import (
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/entity/component/targetable"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move/move"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_CHASE

	// chaseRadius is in units of tiles.
	chaseRadius = 3
)

var (
	OutOfRange = fsm.State(fcpb.ChaseState_CHASE_STATE_OUT_OF_RANGE.String())

	transitions = []fsm.Transition{
		{From: commonstate.Pending, To: commonstate.Executing, VirtualOnly: true},
		{From: commonstate.Pending, To: OutOfRange, VirtualOnly: true},
		{From: commonstate.Pending, To: commonstate.Canceled},
	}

	FSM = fsm.New(transitions, fsmType)
)

type chaseRange struct {
	// start indicates the minimum distance between source and destination
	// at which the source should start chasing.
	start float64
}

type Action struct {
	*action.Base

	source      moveable.Component    // Read-only.
	destination targetable.Component  // Read-only.
	chaseRadius float64               // Read-only.
	status      status.ReadOnlyStatus // Read-only.

	move *move.Action
}

func New(dfStatus status.ReadOnlyStatus, source moveable.Component, destination targetable.Component) *Action {
	return &Action{
		Base:        action.New(FSM, commonstate.Pending),
		source:      source,
		destination: destination,
		chaseRadius: chaseRadius,
		status:      dfStatus,
	}
}

func GenerateMove(a *Action) *move.Action {
	return move.New(a.Source(), a.Status(), a.Destination().Position(a.Status().Tick()), move.Default)
}
func (a *Action) Accept(v visitor.Visitor) error    { return v.Visit(a) }
func (a *Action) Source() moveable.Component        { return a.source }
func (a *Action) Destination() targetable.Component { return a.destination }
func (a *Action) ID() id.ActionID                   { return id.ActionID(a.source.ID()) }
func (a *Action) Status() status.ReadOnlyStatus     { return a.status }

func (a *Action) SetMove(m *move.Action) error {
	a.move = m
	return nil
}

func (a *Action) Precedence(other action.Action) bool {
	if other.Type() != fsmType {
		return false
	}

	// Move is only set during execution -- it's possible that the visitor
	// has not yet scheduled a move.
	if a.move == nil || other.(*Action).move == nil {
		return true
	}

	return a.move.Precedence(other.(*Action).move)
}

func (a *Action) State() (fsm.State, error) {
	var err error
	moveState := commonstate.Finished
	if a.move != nil {
		moveState, err = a.move.State()
		if err != nil {
			return commonstate.Unknown, err
		}
	}

	s, err := a.Base.State()
	if err != nil {
		return commonstate.Unknown, err
	}
	if moveState == commonstate.Canceled {
		return moveState, a.To(s, moveState, true)
	}

	tick := a.status.Tick()

	switch s {
	case commonstate.Pending:
		if d := utils.Euclidean(
			a.source.Position(tick),
			a.destination.Position(tick)); moveState == commonstate.Finished && d > a.chaseRadius {
			return OutOfRange, a.To(s, OutOfRange, true)
		}

		return s, nil
	default:
		return s, nil
	}
}

func (a *Action) Cancel() error {
	s, err := a.State()
	if err != nil {
		return err
	}

	if a.move != nil {
		if err := a.move.Cancel(); err != nil {
			return err
		}
	}
	return a.To(s, commonstate.Canceled, false)
}
