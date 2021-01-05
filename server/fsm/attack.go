// Package attack defines the Action used for carrying out the Attack command.
//
// A Pending state indicates the attack target is out of range or the attack
// ability is not off cooldown yet.
//
// An Executing state indicates the attack target is within range and is off
// cooldown.
//
// A Finished state indicates the target is dead.
package attack

import (
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/map/utils"
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
		{From: commonstate.Executing, To: commonstate.Canceled},
	}

	FSM = fsm.New(transitions, fsmType)
)

type Action struct {
	*action.Base
	chase *chase.Action // Read-only.
	tick  id.Tick       // Read-only.

	status     status.ReadOnlyStatus // Read-only.
	attackable attackable.Component  // Read-only.
	target     targetable.Component  // Read-only.
}

func New(
	chase *chase.Action,
	t id.Tick,
	dfStatus status.ReadOnlyStatus,
	attackable attackable.Component,
	target targetable.Component) *Action {
	return &Action{
		Base:       action.New(FSM, commonstate.Pending),
		attackable: attackable,
		target:     target,
		tick:       t,
		status:     dfStatus,
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
	if s, err := a.chase.State(); err != nil {
		return commonstate.Unknown, err
	} else if s == commonstate.Canceled {
		return s, nil
	}

	s, err := a.Base.State()
	if err != nil {
		return commonstate.Unknown, err
	}

	switch s {
	case commonstate.Pending:
		tick := a.status.Tick()
		if a.target.Health(tick) <= 0 {
			return commonstate.Finished, a.To(s, commonstate.Finished, true)
		}
		if a.attackable.AttackTimerCurve().Ok(tick) && utils.Euclidean(
			a.attackable.Position(tick),
			a.target.Position(tick),
		) <= a.attackable.Range() {
			return commonstate.Executing, a.To(s, commonstate.Executing, true)
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

	if err := a.chase.Cancel(); err != nil {
		return err
	}
	return a.To(s, commonstate.Canceled, false)
}
