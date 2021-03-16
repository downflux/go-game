package move

import (
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/fsm/commonstate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
	serverstatus "github.com/downflux/game/engine/status/status"
)

type MoveType int

const (
	Default MoveType = iota
	Direct

	fsmType = fcpb.FSMType_FSM_TYPE_MOVE
)

var (
	transitions = []fsm.Transition{
		{From: commonstate.Pending, To: commonstate.Executing, VirtualOnly: true},
		{From: commonstate.Pending, To: commonstate.Canceled},
		{From: commonstate.Pending, To: commonstate.Finished, VirtualOnly: true},
		{From: commonstate.Executing, To: commonstate.Canceled},
	}

	FSM = fsm.New(transitions, fsmType)
)

type Action struct {
	*action.Base

	// tick is the tick at which the command was originally
	// scheduled.
	tick id.Tick // Read-only.

	status      serverstatus.ReadOnlyStatus // Read-only.
	destination *gdpb.Position              // Read-only.
	moveType    MoveType                    // Read-only.

	e moveable.Component // Read-only.

	// TODO(minkezhang): Move executionTick and destination into
	// separate external cache.
	executionTick id.Tick
}

// New constructs a new Action FSM action.
//
// TODO(minkezhang): Add executionTick arg to allow for scheduling in the
// future.
func New(
	e moveable.Component,
	dfStatus serverstatus.ReadOnlyStatus,
	destination *gdpb.Position,
	moveType MoveType) *Action {
	t := dfStatus.Tick()
	return &Action{
		Base:          action.New(FSM, commonstate.Pending),
		e:             e,
		status:        dfStatus,
		tick:          t,
		executionTick: t,
		// TODO(minkezhang): Change to original position after mesh nav
		// migration.
		destination:   &gdpb.Position{
			X: float64(int(destination.GetX())),
			Y: float64(int(destination.GetY()))},
		moveType:      moveType,
	}
}

func (n *Action) Accept(v visitor.Visitor) error { return v.Visit(n) }
func (n *Action) Component() moveable.Component  { return n.e }
func (n *Action) ID() id.ActionID                { return id.ActionID(n.e.ID()) }
func (n *Action) MoveType() MoveType             { return n.moveType }

// SchedulePartialMove allows us to mutate the FSM action to deal with
// partial moves. This allows us to know when the visitor should make the next
// meaningful calculation.
func (n *Action) SchedulePartialMove(t id.Tick) error {
	if n.MoveType() == Direct {
		return status.Error(
			codes.FailedPrecondition,
			"cannot schedule partial move for a direct move")
	}
	n.executionTick = t
	return nil
}

func (n *Action) Precedence(i action.Action) bool {
	if i.Type() != n.Type() {
		return false
	}

	m := i.(*Action)
	return n.tick >= m.tick && !proto.Equal(n.Destination(), m.Destination())
}

// TODO(minkezhang): Return a cloned instance instead.
func (n *Action) Destination() *gdpb.Position { return n.destination }

func (n *Action) Cancel() error {
	s, err := n.State()
	if err != nil {
		return err
	}

	return n.To(s, commonstate.Canceled, false)
}

func (n *Action) State() (fsm.State, error) {
	tick := n.status.Tick()

	s, err := n.Base.State()
	if err != nil {
		return commonstate.Unknown, err
	}

	switch s {
	case commonstate.Pending:
		var t fsm.State = commonstate.Unknown

		if n.executionTick <= tick {
			t = commonstate.Executing
			if proto.Equal(n.destination, n.e.Position(tick)) {
				t = commonstate.Finished
			}
		}

		if t != commonstate.Unknown {
			if err := n.To(s, t, true); err != nil {
				return commonstate.Unknown, err
			}
			return t, nil
		}

		return commonstate.Pending, nil
	default:
		return s, nil
	}
}
