package schedule

import (
	"testing"

	"github.com/downflux/game/fsm/instance"
	"github.com/downflux/game/fsm/move"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/service/status"

	fcpb "github.com/downflux/game/fsm/api/constants_go_proto"
)

func TestAddError(t *testing.T) {
	s := New(nil)
	stat := status.New(0)

	i := move.New(tank.New("entity-id", 0, nil), status.New(0), stat.Tick(), nil)
	if err := s.Add(i); err == nil {
		t.Errorf("Add() = nil, want a non-nil error")
	}
}

func TestAdd(t *testing.T) {
	const iid = "entity-id"

	s := New([]fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE})
	stat := status.New(0)

	i := move.New(tank.New(iid, 0, nil), status.New(0), stat.Tick(), nil)
	if err := s.Add(i); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got := s.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(iid); got.ID() != iid {
		t.Errorf("ID() = %v, want = %v", got, iid)
	}
}

func TestMerge(t *testing.T) {
	const iid = "entity-id"

	simpleMergeStat := status.New(0)
	mergeFilterStat := status.New(0)

	testConfigs := []struct {
		name      string
		s1Types   []fcpb.FSMType
		s2Types   []fcpb.FSMType
		instances []instance.Instance
		want      []instance.Instance
	}{
		{
			name:    "TestSimpleMerge",
			s1Types: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
			s2Types: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
			instances: []instance.Instance{
				move.New(tank.New(iid, 0, nil), simpleMergeStat, simpleMergeStat.Tick(), nil),
			},
			want: []instance.Instance{
				move.New(tank.New(iid, 0, nil), simpleMergeStat, simpleMergeStat.Tick(), nil),
			},
		},
		{
			name:    "TestMergeFilter",
			s1Types: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
			s2Types: nil,
			instances: []instance.Instance{
				move.New(tank.New(iid, 0, nil), mergeFilterStat, mergeFilterStat.Tick(), nil),
			},
			want: nil,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			s1 := New(c.s1Types)
			s2 := New(c.s2Types)

			for _, i := range c.instances {
				if err := s1.Add(i); err != nil {
					t.Fatalf("Add() = %v, want = nil", err)
				}
			}

			if err := s2.Merge(s1); err != nil {
				t.Fatalf("Merge() = %v, want = nil", err)
			}

			for _, i := range c.want {
				if got := s2.Get(i.Type()).Get(i.ID()); got.ID() != iid {
					t.Fatalf("ID() = %v, want = %v", got, iid)
				}
			}
		})
	}
}

func TestPop(t *testing.T) {
	const iid = "entity-id"

	s1 := New([]fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE})
	stat := status.New(0)

	i := move.New(tank.New(iid, 0, nil), stat, stat.Tick(), nil)
	if err := s1.Add(i); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	s2 := s1.Pop()
	if got := s2.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(iid); got.ID() != iid {
		t.Fatalf("ID() = %v, want = %v", got, iid)
	}
	if got := s1.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(iid); got != nil {
		t.Errorf("Get() = %v, want = nil", got)
	}
}
