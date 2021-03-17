package schedule

import (
	"testing"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/mock/simple"
	"github.com/downflux/game/engine/id/id"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

// TestExtendError validates we cannot add unknown FSMs to the schedule.
func TestExtendError(t *testing.T) {
	s := New(nil)

	a := simple.New(id.ActionID("action-id"), 0)
	if err := s.Extend([]action.Action{a}); err == nil {
		t.Errorf("Extend() = nil, want a non-nil error %s", err)
	}
}

func TestExtend(t *testing.T) {
	aid := id.ActionID("action-id")

	s := New([]fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE})
	i := simple.New(aid, 0)
	if err := s.Extend([]action.Action{i}); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got := s.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(aid); got.ID() != aid {
		t.Errorf("ID() = %v, want = %v", got, aid)
	}
}

func TestMerge(t *testing.T) {
	const eid = "entity-id"
	aid := id.ActionID("action-id")

	type sc struct {
		ts      []fcpb.FSMType
		actions []action.Action
	}

	testPrecedenceReferencedActionLow := simple.New(aid, 0)
	testPrecedenceReferencedActionHigh := simple.New(aid, 1)

	testConfigs := []struct {
		name string
		s1   sc
		s2   sc
		want []action.Action
	}{
		{
			name: "TestSimpleMerge",

			s1: sc{
				ts: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
				actions: []action.Action{
					simple.New(aid, 0),
				},
			},
			s2: sc{
				ts: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
			},

			want: []action.Action{
				simple.New(aid, 0),
			},
		},
		{
			name: "TestMergeFilter",

			s1: sc{
				ts: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
				actions: []action.Action{
					simple.New(aid, 0),
				},
			},
			s2: sc{ts: nil},

			want: nil,
		},
		// TODO(minkezhang): Move to list_test instead.
		{
			name: "TestPrecedence",

			s1: sc{
				ts: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
				actions: []action.Action{
					testPrecedenceReferencedActionHigh,
				},
			},
			s2: sc{
				ts: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
				actions: []action.Action{
					testPrecedenceReferencedActionLow,
				},
			},

			want: []action.Action{
				testPrecedenceReferencedActionHigh,
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			s1 := New(c.s1.ts)
			s2 := New(c.s2.ts)

			if err := s1.Extend(c.s1.actions); err != nil {
				t.Fatalf("Extend() = %v, want = nil", err)
			}
			if err := s2.Extend(c.s2.actions); err != nil {
				t.Fatalf("Extend() = %v, want = nil", err)
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

	if s, err := testPrecedenceReferencedActionLow.State(); err != nil || s != simple.Canceled {
		t.Errorf("State() = %v, %v, want = %v, nil", s, err, simple.Canceled)
	}
}

func TestPop(t *testing.T) {
	aid := id.ActionID("entity-id")

	s1 := New([]fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE})
	a := simple.New(aid, 0)
	if err := s1.Extend([]action.Action{a}); err != nil {
		t.Fatalf("Extend() = %v, want = nil", err)
	}

	s2 := s1.Pop()
	if got := s2.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(aid); got.ID() != aid {
		t.Fatalf("ID() = %v, want = %v", got, aid)
	}
	if got := s1.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(aid); got != nil {
		t.Errorf("Get() = %v, want = nil", got)
	}
}
