package schedule

import (
	"sync"

	"github.com/downflux/game/fsm/instance"
	"github.com/downflux/game/fsm/list"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fcpb "github.com/downflux/game/fsm/api/constants_go_proto"
)

var (
	ListTypes = []fcpb.FSMType{
		fcpb.FSMType_FSM_TYPE_MOVE,
	}
)

type Schedule struct {
	listTypes []fcpb.FSMType

	mux       sync.Mutex
	instances map[fcpb.FSMType]*list.List
}

func New(listTypes []fcpb.FSMType) *Schedule {
	s := &Schedule{
		listTypes: listTypes,
	}
	s.resetInstancesUnsafe()
	return s
}

func (s *Schedule) resetInstancesUnsafe() {
	if s.instances == nil {
		s.instances = map[fcpb.FSMType]*list.List{}
	}

	for _, t := range s.listTypes {
		s.instances[t] = list.New(t)
	}
}

func (s *Schedule) Pop() *Schedule {
	s.mux.Lock()
	defer s.mux.Unlock()

	ns := &Schedule{
		instances: s.instances,
	}

	s.resetInstancesUnsafe()

	return ns
}

func (s *Schedule) Add(i instance.Instance) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, found := s.instances[i.Type()]; !found {
		return status.Errorf(
			codes.FailedPrecondition,
			"schedule does not accept %v FSM instances",
			i.Type())
	}

	return s.instances[i.Type()].Add(i)
}

func (s *Schedule) Merge(t *Schedule) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	// TODO(minkezhang): Consider if we should make this parallel.
	for _, t := range s.listTypes {
		if err := s.instances[t].Merge(t.Get(t)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schedule) Get(fsmType fcpb.FSMType) *list.List {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.instances[fsmType]
}
