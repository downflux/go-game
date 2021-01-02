package chase

import (
	"sync"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/entity/component/positionable"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_CHASE
)

var (
	Waiting = fsm.State(fcpb.ChaseState_CHASE_STATE_WAITING.String())

	transitions = []fsm.Transition{
		{From: commonstate.Pending, To: commonstate.Executing, VirtualOnly: true},
		{From: commonstate.Pending, To: commonstate.Canceled},
		{From: commonstate.Pending, To: Waiting, VirtualOnly: true},
		{From: commonstate.Executing, To: commonstate.Pending, VirtualOnly: true},
		{From: commonstate.Executing, To: commonstate.Canceled},
		{From: Waiting, To: commonstate.Canceled},
	}

	FSM = fsm.New(transitions, fsmType)
)

type Action struct {
	*action.Base

	source      moveable.Component     // Read-only.
	destination positionable.Component // Read-only.

	// mux guards the Base and move properties.
	mux  sync.Mutex
	move *move.Action
}

func New(source moveable.Component, destination positionable.Component, moveaction *move.Action) *Action {
	return &Action{
		Base:        action.New(FSM, commonstate.Pending),
		source:      source,
		destination: destination,
		move:        moveaction,
	}
}

func (a *Action) Accept(v visitor.Visitor) error      { return v.Visit(a) }
func (a *Action) Source() moveable.Component          { return a.source }
func (a *Action) Destination() positionable.Component { return a.destination }
func (a *Action) ID() id.ActionID                     { return id.ActionID(a.source.ID()) }

func (a *Action) SetMove(m *move.Action) error {
	a.mux.Lock()
	defer a.mux.Unlock()

	a.move = m
	return nil
}

func (a *Action) Precedence(other action.Action) bool {
	if other.Type() != fsmType {
		return false
	}

	return a.move.Precedence(other.(*Action).move)
}

func (a *Action) State() (fsm.State, error) {
	a.mux.Lock()
	defer a.mux.Unlock()

	return a.stateUnsafe()
}

func (a *Action) stateUnsafe() (fsm.State, error) {
	s, err := a.Base.State()
	if err != nil {
		return commonstate.Unknown, err
	}

	switch s {
	case commonstate.Pending:
		moveState, err := a.move.State()
		if err != nil {
			return commonstate.Unknown, err
		}
		switch moveState {
		case commonstate.Finished:
			return Waiting, a.To(s, Waiting, true)
		default:
			return moveState, a.To(s, moveState, true)
		}
	default:
		return s, nil
	}
}

func (a *Action) Cancel() error {
	a.mux.Lock()
	defer a.mux.Unlock()

	s, err := a.stateUnsafe()
	if err != nil {
		return err
	}

	return a.To(s, commonstate.Canceled, false)
}
