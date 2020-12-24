package schedule

import (
	"sync"

	"github.com/downflux/game/engine/fsm/instance"
	"github.com/downflux/game/engine/fsm/list"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

var (
	FSMTypes = []fcpb.FSMType{
		fcpb.FSMType_FSM_TYPE_MOVE,
		fcpb.FSMType_FSM_TYPE_PRODUCE,
	}
)

type Schedule struct {
	fsmTypes map[fcpb.FSMType]bool

	mux       sync.Mutex
	instances map[fcpb.FSMType]*list.List
}

func New(fsmTypes []fcpb.FSMType) *Schedule {
	s := &Schedule{
		fsmTypes: map[fcpb.FSMType]bool{},
	}
	for _, fsmType := range fsmTypes {
		s.fsmTypes[fsmType] = true
	}
	return s
}

func (s *Schedule) Pop() *Schedule {
	s.mux.Lock()
	defer s.mux.Unlock()

	ns := &Schedule{
		instances: s.instances,
	}
	s.instances = nil
	return ns
}

func (s *Schedule) Add(i instance.Instance) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.instances == nil {
		s.instances = map[fcpb.FSMType]*list.List{}
	}

	fsmType := i.Type()

	if !s.fsmTypes[fsmType] {
		return status.Errorf(codes.FailedPrecondition, "schedule does not accept %v FSM instances", fsmType)
	}

	if _, found := s.instances[fsmType]; !found {
		s.instances[fsmType] = list.New(fsmType)
	}

	return s.instances[fsmType].Add(i)
}

func (s *Schedule) Merge(t *Schedule) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.instances == nil {
		s.instances = map[fcpb.FSMType]*list.List{}
	}

	// TODO(minkezhang): Consider if we should make this parallel.
	for fsmType := range s.fsmTypes {
		if l := t.Get(fsmType); l != nil {
			if _, found := s.instances[fsmType]; !found {
				s.instances[fsmType] = list.New(fsmType)
			}

			if err := s.instances[fsmType].Merge(l); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Schedule) Get(fsmType fcpb.FSMType) *list.List {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.instances == nil {
		s.instances = map[fcpb.FSMType]*list.List{}
	}

	if _, found := s.instances[fsmType]; !found {
		s.instances[fsmType] = list.New(fsmType)
	}

	return s.instances[fsmType]
}

func (s *Schedule) Clear() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, l := range s.instances {
		if err := l.Clear(); err != nil {
			return err
		}
	}
	return nil
}
