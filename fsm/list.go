package list

import (
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/fsm/fsm"
	"github.com/downflux/game/fsm/instance"
	"github.com/downflux/game/server/id"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	vcpb "github.com/downflux/game/engine/visitor/api/constants_go_proto"
	fcpb "github.com/downflux/game/fsm/api/constants_go_proto"
)

const (
	agentType = vcpb.AgentType_AGENT_TYPE_FSM_LIST
)

var (
	notImplemented = status.Errorf(
		codes.Unimplemented, "given function has not been implemented yet")
)

type List struct {
	fsmType   fcpb.FSMType
	instances map[id.InstanceID]instance.Instance
}

func New(fsmType fcpb.FSMType) *List {
	return &List{
		instances: map[id.InstanceID]instance.Instance{},
		fsmType:   fsmType,
	}
}

func (l *List) AgentType() vcpb.AgentType               { return agentType }
func (l *List) Type() fcpb.FSMType                      { return l.fsmType }
func (l *List) Get(iid id.InstanceID) instance.Instance { return l.instances[iid] }

func (l *List) Clear() error {
	for iid, i := range l.instances {
		s, err := i.State()
		if err != nil {
			// TODO(minkezhang): Log and move on here.
			return err
		}
		if s == fsm.State(fcpb.CommonState_COMMON_STATE_CANCELED.String()) || s == fsm.State(fcpb.CommonState_COMMON_STATE_FINISHED.String()) {
			delete(l.instances, iid)
		}
	}
	return nil
}

func (l *List) Accept(v visitor.Visitor) error {
	if err := v.Visit(l); err != nil {
		return err
	}

	var eg errgroup.Group
	for _, i := range l.instances {
		i := i
		eg.Go(func() error { return i.Accept(v) })
	}
	return eg.Wait()
}

func (l *List) Merge(j *List) error {
	// TODO(minkezhang): Consider making this concurrent.
	for _, i := range j.instances {
		if err := l.Add(i); err != nil {
			return err
		}
	}
	return nil
}

func (l *List) Add(i instance.Instance) error {
	if l.Type() != i.Type() {
		return status.Errorf(codes.FailedPrecondition, "cannot add instance of type %v to a list of type %v", i.Type(), l.Type())
	}

	if l.instances == nil {
		l.instances = map[id.InstanceID]instance.Instance{}
	}

	j, found := l.instances[i.ID()]
	if !found || (found && i.Precedence(j)) {
		// Cancel any conflicting move commands.
		if found {
			if err := j.Cancel(); err != nil {
				return err
			}
		}
		l.instances[i.ID()] = i
	}
	return nil
}

func (l *List) Remove(iid id.InstanceID) error {
	delete(l.instances, iid)
	return nil
}
