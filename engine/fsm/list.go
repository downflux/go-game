package list

import (
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/visitor"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
	vcpb "github.com/downflux/game/engine/visitor/api/constants_go_proto"
)

const (
	agentType = vcpb.AgentType_AGENT_TYPE_FSM_LIST
)

var (
	notImplemented = status.Errorf(
		codes.Unimplemented, "given function has not been implemented yet")
)

type List struct {
	fsmType fcpb.FSMType  // Read-only.
	dependencies map[fcpb.FSMType]bool  // Read-only.
	actions map[id.ActionID]action.Action
}

func New(fsmType fcpb.FSMType, deps []fcpb.FSMType) *List {
	dependencies := map[fcpb.FSMType]bool{}
	for _, d := range deps {
		dependencies[d] = true
	}

	return &List{
		actions: map[id.ActionID]action.Action{},
		fsmType: fsmType,
		dependencies: dependencies,
	}
}

func (l *List) AgentType() vcpb.AgentType         { return agentType }
func (l *List) Type() fcpb.FSMType                { return l.fsmType }
func (l *List) Get(iid id.ActionID) action.Action { return l.actions[iid] }

func (l *List) Clear() error {
	for iid, i := range l.actions {
		s, err := i.State()
		if err != nil {
			// TODO(minkezhang): Log and move on here.
			return err
		}
		if s == fsm.State(fcpb.CommonState_COMMON_STATE_CANCELED.String()) || s == fsm.State(fcpb.CommonState_COMMON_STATE_FINISHED.String()) {
			delete(l.actions, iid)
		}
	}
	return nil
}

func (l *List) Accept(v visitor.Visitor) error {
	if err := v.Visit(l); err != nil {
		return err
	}

	var eg errgroup.Group
	for _, i := range l.actions {
		i := i
		eg.Go(func() error { return i.Accept(v) })
	}
	return eg.Wait()
}

func (l *List) Merge(j *List) error {
	if j.Type() != l.Type() {
		return nil
	}

	// TODO(minkezhang): Consider making this concurrent.
	for _, i := range j.actions {
		if err := l.Add(i); err != nil {
			return err
		}
	}
	return nil
}

// Cancel unconditionally cancels the target action with the given ID.
//
// We don't want to complicate the FSM Precedence() function by adding in
// external deps, so the next-best place to put cross-FSM canceling
// functionality to be in the List struct.
//
// Here, the input list has higher precedence (i.e., we cancel the corresponding
// Action object in l).
func (l *List) Cancel(j *List) error{
	if !l.dependencies[j.Type()] {
		return nil
	}

	for _, i := range j.actions {
		if target, found := l.actions[i.ID()]; found {
			if err := target.Cancel(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *List) Add(i action.Action) error {
	if l.Type() != i.Type() {
		return status.Errorf(codes.FailedPrecondition, "cannot add instance of type %v to a list of type %v", i.Type(), l.Type())
	}

	if l.actions == nil {
		l.actions = map[id.ActionID]action.Action{}
	}

	j, found := l.actions[i.ID()]
	if !found || (found && i.Precedence(j)) {
		// Cancel any conflicting move commands.
		if found {
			if err := j.Cancel(); err != nil {
				return err
			}
		}
		l.actions[i.ID()] = i
	}
	return nil
}

func (l *List) Remove(iid id.ActionID) error {
	delete(l.actions, iid)
	return nil
}
