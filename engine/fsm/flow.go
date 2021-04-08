package flow

import (
	"fmt"
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/fsm/group/group"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/visitor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_FLOW
)

var (
	Unknown   = fsm.State(fcpb.CommonState_COMMON_STATE_UNKNOWN.String())
	Executing = fsm.State(fcpb.CommonState_COMMON_STATE_EXECUTING.String())
	Canceled  = fsm.State(fcpb.CommonState_COMMON_STATE_CANCELED.String())
	Finished  = fsm.State(fcpb.CommonState_COMMON_STATE_FINISHED.String())

	transitions = []fsm.Transition{
		{From: Executing, To: Finished},
		{From: Executing, To: Canceled},
		{From: Canceled, To: Canceled},
	}
	FSM = fsm.New(transitions, fsmType)
)

type Action struct {
	*action.Base

	actions map[id.ActionID]action.Action

	id    id.ActionID // Read-only.
	tick  id.Tick     // Read-only.
	group group.Group // Read-only.
}

func New(tick id.Tick, group group.Group, suffix id.ActionID) *Action {
	return &Action{
		Base:    action.New(FSM, Executing),
		group:   group,
		id:      id.ActionID(fmt.Sprintf("%s/%s", group.ID(), suffix)),
		tick:    tick,
		actions: map[id.ActionID]action.Action{},
	}
}

func (a *Action) Append(i action.Action) error {
	if a.group.Contains(i.Type()) {
		a.actions[i.ID()] = i
		return nil
	}
	return status.Errorf(
		codes.FailedPrecondition,
		"cannot append %v action to flow %v",
		i.Type(),
		a.ID(),
	)
}

func (a *Action) Accept(v visitor.Visitor) error {
	return v.Visit(a)
}

func (a *Action) Precedence(i action.Action) bool {
	return a.ID() == i.ID() && a.group.ID() == i.(*Action).group.ID() && a.tick >= i.(*Action).tick
}

func (a *Action) State() (fsm.State, error) {
	if _, err := a.Base.State(); err != nil {
		return Unknown, err
	}

	done := false
	for _, i := range a.actions {
		s, err := i.State()
		if err != nil {
			return Unknown, err
		}
		done = done && (s == Finished)
	}

	if done {
		return Finished, nil
	}
	return Executing, nil
}

func (a *Action) Cancel() error {
	s, err := a.Base.State()
	if err != nil {
		return err
	}
	if err := a.To(s, Canceled, false); err != nil {
		return err
	}

	for _, i := range a.actions {
		if err := i.Cancel(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Action) ID() id.ActionID {
	return a.id
}
