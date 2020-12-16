package list

import (
	"github.com/downflux/game/fsm/instance"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/visitor/visitor"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fcpb "github.com/downflux/game/fsm/api/constants_go_proto"
	vcpb "github.com/downflux/game/server/visitor/api/constants_go_proto"
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

func (l *List) Clear(v visitor.Visitor) error { return notImplemented }

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

func (l *List) Add(i instance.Instance) error {
	if i.Type() != l.Type() {
		status.Errorf(
			codes.FailedPrecondition,
			"cannot add instance of type %v to a list of type %v",
			i.Type(),
			l.Type())
	}

	if l.instances == nil {
		l.instances = map[id.InstanceID]instance.Instance{}
	}

	j, found := l.instances[i.ID()]
	if found && i.Precedence(j) {
		if err := j.Cancel(); err != nil {
			return err
		}
		l.instances[i.ID()] = i
	}
	l.instances[i.ID()] = i
	return nil
}

func (l *List) Remove(iid id.InstanceID) error {
	delete(l.instances, iid)
	return nil
}
