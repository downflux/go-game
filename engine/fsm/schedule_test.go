package schedule

import (
	"testing"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/mock/simple"
	"github.com/downflux/game/engine/id/id"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

// TestAddError validates we cannot add unknown FSMs to the schedule.
func TestAddError(t *testing.T) {
	s := New(nil)

	a := simple.New(id.ActionID("action-id"), 0)
	if err := s.Add(a); err == nil {
		t.Errorf("Add() = nil, want a non-nil error %s", err)
	}
}

func TestAdd(t *testing.T) {
	aid := id.ActionID("action-id")

	s := New([]fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE})
	i := simple.New(aid, 0)
	if err := s.Add(i); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got := s.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(aid); got.ID() != aid {
		t.Errorf("ID() = %v, want = %v", got, aid)
	}
}

func TestMerge(t *testing.T) {
	const eid = "entity-id"
	aid := id.ActionID("action-id")

	testConfigs := []struct {
		name    string
		s1Types []fcpb.FSMType
		s2Types []fcpb.FSMType
		actions []action.Action
		want    []action.Action
	}{
		{
			name:    "TestSimpleMerge",
			s1Types: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
			s2Types: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
			actions: []action.Action{
				simple.New(aid, 0),
			},
			want: []action.Action{
				simple.New(aid, 0),
			},
		},
		{
			name:    "TestMergeFilter",
			s1Types: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
			s2Types: nil,
			actions: []action.Action{
				simple.New(aid, 0),
			},
			want: nil,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			s1 := New(c.s1Types)
			s2 := New(c.s2Types)

			for _, i := range c.actions {
				if err := s1.Add(i); err != nil {
					t.Fatalf("Add() = %v, want = nil", err)
				}
			}

			if err := s2.Merge(s1); err != nil {
				t.Fatalf("Merge() = %v, want = nil", err)
			}

			for _, i := range c.want {
				if got := s2.Get(i.Type()).Get(i.ID()); got.ID() != aid {
					t.Fatalf("ID() = %v, want = %v", got, eid)
				}
			}
		})
	}
}

func TestPop(t *testing.T) {
	aid := id.ActionID("entity-id")

	s1 := New([]fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE})
	a := simple.New(aid, 0)
	if err := s1.Add(a); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	s2 := s1.Pop()
	if got := s2.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(aid); got.ID() != aid {
		t.Fatalf("ID() = %v, want = %v", got, aid)
	}
	if got := s1.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(aid); got != nil {
		t.Errorf("Get() = %v, want = nil", got)
	}
}
