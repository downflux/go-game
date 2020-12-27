package schedule

import (
	"testing"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/fsm/move"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

func newTank(t *testing.T, eid id.EntityID, tick id.Tick) *tank.Entity {
	tankEntity, err := tank.New(eid, tick, nil)
	if err != nil {
		t.Fatalf("New() = %v, want = nil", err)
	}
	return tankEntity
}

func TestAddError(t *testing.T) {
	s := New(nil)

	i := move.New(newTank(t, id.EntityID("entity-id"), id.Tick(0)), status.New(0), nil)
	if err := s.Add(i); err == nil {
		t.Errorf("Add() = nil, want a non-nil error")
	}
}

func TestAdd(t *testing.T) {
	const eid = "entity-id"

	s := New([]fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE})
	i := move.New(newTank(t, eid, id.Tick(0)), status.New(0), nil)
	if err := s.Add(i); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got := s.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(eid); got.ID() != eid {
		t.Errorf("ID() = %v, want = %v", got, eid)
	}
}

func TestMerge(t *testing.T) {
	const eid = "entity-id"

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
				move.New(newTank(t, eid, id.Tick(0)), status.New(0), nil),
			},
			want: []action.Action{
				move.New(newTank(t, eid, id.Tick(0)), status.New(0), nil),
			},
		},
		{
			name:    "TestMergeFilter",
			s1Types: []fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE},
			s2Types: nil,
			actions: []action.Action{
				move.New(newTank(t, eid, id.Tick(0)), status.New(0), nil),
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
				if got := s2.Get(i.Type()).Get(i.ID()); got.ID() != eid {
					t.Fatalf("ID() = %v, want = %v", got, eid)
				}
			}
		})
	}
}

func TestPop(t *testing.T) {
	const eid = "entity-id"

	s1 := New([]fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE})
	i := move.New(newTank(t, eid, id.Tick(0)), status.New(0), nil)
	if err := s1.Add(i); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	s2 := s1.Pop()
	if got := s2.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(eid); got.ID() != eid {
		t.Fatalf("ID() = %v, want = %v", got, eid)
	}
	if got := s1.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(eid); got != nil {
		t.Errorf("Get() = %v, want = nil", got)
	}
}
