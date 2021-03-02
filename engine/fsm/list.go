// Package list implements a collection of FSM actions.
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
)

var (
	notImplemented = status.Errorf(
		codes.Unimplemented, "given function has not been implemented yet")
)

type List struct {
	fsmType fcpb.FSMType
	actions map[id.ActionID]action.Action
}

func New(fsmType fcpb.FSMType) *List {
	return &List{
		actions: map[id.ActionID]action.Action{},
		fsmType: fsmType,
	}
}

// TODO(minkezhang): Rename Action.
func (l *List) Get(iid id.ActionID) action.Action { return l.actions[iid] }
func (l *List) Type() fcpb.FSMType                { return l.fsmType }

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

// TODO(minkezhang): Rename to make clear List is not an FSM agent.
func (l *List) Accept(v visitor.Visitor) error {
	var eg errgroup.Group
	for _, i := range l.actions {
		i := i
		eg.Go(func() error { return i.Accept(v) })
	}
	return eg.Wait()
}

func (l *List) Merge(j *List) error {
	// TODO(minkezhang): Consider making this concurrent.
	for _, i := range j.actions {
		if err := l.Add(i); err != nil {
			return err
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
