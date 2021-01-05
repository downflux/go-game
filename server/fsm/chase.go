package chase

import (
	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/entity/component/targetable"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_CHASE
)

var (
	InRange = fsm.State(fcpb.ChaseState_CHASE_STATE_IN_RANGE.String())
	OutOfRange = fsm.State(fcpb.ChaseState_CHASE_STATE_OUT_OF_RANGE.String())

	transitions = []fsm.Transition{
		{From: commonstate.Pending, To: commonstate.Executing, VirtualOnly: true},
		{From: commonstate.Pending, To: InRange, VirtualOnly: true},
		{From: commonstate.Pending, To: OutOfRange, VirtualOnly: true},
		{From: commonstate.Pending, To: commonstate.Canceled},
	}

	FSM = fsm.New(transitions, fsmType)

	presetChaseRange = chaseRange{ start: 3, stop: 1 }
)

type chaseRange struct {
	// start indicates the distance between source and destination at
	// which the source starts chasing. The start property should be
	// strictly greater than the stop distance.
	start float64

	// stop indicates the distance between source and destination at which
	// the source stops chasing.
	stop float64
}

type Action struct {
	*action.Base

	source      moveable.Component   // Read-only.
	destination targetable.Component // Read-only.
	chaseRange  chaseRange           // Read-only.
	status      *status.Status       // Read-only.

	move *move.Action
}

func New(dfStatus *status.Status, source moveable.Component, destination targetable.Component, moveaction *move.Action) *Action {
	return &Action{
		Base:        action.New(FSM, commonstate.Pending),
		source:      source,
		destination: destination,
		move:        moveaction,
		chaseRange:  presetChaseRange,
		status:      dfStatus,
	}
}

func (a *Action) Accept(v visitor.Visitor) error    { return v.Visit(a) }
func (a *Action) Source() moveable.Component        { return a.source }
func (a *Action) Destination() targetable.Component { return a.destination }
func (a *Action) ID() id.ActionID                   { return id.ActionID(a.source.ID()) }

func (a *Action) SetMove(m *move.Action) error {
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
	s, err := a.Base.State()
	if err != nil {
		return commonstate.Unknown, err
	}

	tick := a.status.Tick()

	switch s {
	case commonstate.Pending:
		moveState, err := a.move.State()
		if err != nil {
			return commonstate.Unknown, err
		}

		if d := utils.Euclidean(
			a.source.Position(tick),
			a.destination.Position(tick)); d < a.chaseRange.stop {
			return InRange, a.To(s, InRange, true)
		} else if moveState == commonstate.Finished && d > a.chaseRange.start {
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

	return a.To(s, commonstate.Canceled, false)
}
