package attack

import (
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/entity/component/attackable"
	"github.com/downflux/game/server/entity/component/targetable"
	"github.com/downflux/game/server/fsm/chase"
	"github.com/downflux/game/server/fsm/commonstate"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_ATTACK
)

var (
	transitions = []fsm.Transition{
		{From: commonstate.Pending, To: commonstate.Executing, VirtualOnly: true},
		{From: commonstate.Pending, To: commonstate.Finished, VirtualOnly: true},
		{From: commonstate.Pending, To: commonstate.Canceled},
	}

	FSM = fsm.New(transitions, fsmType)
)

type Action struct {
	*action.Base
	chase *chase.Action
	tick  id.Tick

	attackable attackable.Component
	target     targetable.Component
}

func New(chase *chase.Action, t id.Tick, attackable attackable.Component, target targetable.Component) *Action {
	return &Action{
		Base:       action.New(FSM, commonstate.Pending),
		attackable: attackable,
		target:     target,
		tick:       t,
	}
}

func (a *Action) Accept(v visitor.Visitor) error { return v.Visit(a) }
func (a *Action) ID() id.ActionID                { return id.ActionID(a.attackable.ID()) }

func (a *Action) Precedence(o action.Action) bool {
	if a.Type() != fsmType {
		return false
	}

	b := o.(*Action)

	return a.tick >= b.tick && a.target != b.target
}

func (a *Action) State() (fsm.State, error) {
	s, err := a.Base.State()
	if err != nil {
		return commonstate.Unknown, err
	}

	switch s {
	case commonstate.Pending:
		// TODO(minkezhang): Implement
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
