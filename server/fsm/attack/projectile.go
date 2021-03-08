package projectile

import (
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/entity/component/attackable"
	"github.com/downflux/game/server/entity/component/targetable"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move/move"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_PROJECTILE_SHOOT
)

var (
	transitions = []fsm.Transition{
		{From: commonstate.Pending, To: commonstate.Executing, VirtualOnly: true},
		{From: commonstate.Pending, To: commonstate.Canceled},
		{From: commonstate.Executing, To: commonstate.Canceled},
		{From: commonstate.Executing, To: commonstate.Finished},
	}

	FSM = fsm.New(transitions, fsmType)
)

type Action struct {
	*action.Base

	source attackable.Component // Read-only.
	target targetable.Component // Read-only.

	move *move.Action // Read-only.
}

func New(source attackable.Component, target targetable.Component, move *move.Action) *Action {
	return &Action{
		Base:   action.New(FSM, commonstate.Pending),
		source: source,
		target: target,
		move:   move,
	}
}

func (n *Action) Cancel() error {
	return status.Error(codes.Unimplemented, "cannot cancel a projectile in flight")
}

func (n *Action) Precedence(i action.Action) bool { return false }
func (n *Action) ID() id.ActionID                 { return n.move.ID() }
func (n *Action) Source() attackable.Component    { return n.source }
func (n *Action) Target() targetable.Component    { return n.target }
func (n *Action) Accept(v visitor.Visitor) error  { return v.Visit(n) }

func (n *Action) State() (fsm.State, error) {
	s, err := n.Base.State()
	if err != nil {
		return s, err
	}

	moveState, err := n.move.State()
	if err != nil || s == commonstate.Canceled {
		return s, err
	}

	switch s {
	case commonstate.Pending:
		if moveState == commonstate.Finished {
			return commonstate.Executing, nil
		}
		return s, nil
	default:
		return s, nil
	}
}
