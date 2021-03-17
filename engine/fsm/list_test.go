package list

import (
	"testing"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/fsm/mock/simple"
	"github.com/downflux/game/engine/id/id"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_MOVE
)

func TestDelete(t *testing.T) {
}

func TestNew(t *testing.T) {
	l := New(fsmType)

	if got := l.Type(); got != fsmType {
		t.Fatalf("Type() = %v, want = %v", got, fsmType)
	}
}

func TestAddError(t *testing.T) {
	l := New(fcpb.FSMType_FSM_TYPE_UNKNOWN)
	i := simple.New(id.ActionID("action-id"), 0)

	if err := l.Add(i); err == nil {
		t.Error("Add() = nil, want a non-nil error")
	}
}

func TestAdd(t *testing.T) {
	aid := id.ActionID("action-id")

	l := New(fsmType)
	i := simple.New(aid, 0)

	if err := l.Add(i); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got := l.Get(aid).ID(); got != aid {
		t.Errorf("ID() = %v, want = %v", got, aid)
	}
}

func TestAddCancel(t *testing.T) {
	aid := id.ActionID("action-id")

	l := New(fsmType)

	a1 := simple.New(aid, 0)
	a2 := simple.New(aid, 1)

	if err := l.Add(a1); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}
	if err := l.Add(a2); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	want := fsm.State(simple.Canceled)
	if got, err := a1.State(); err != nil || got != want {
		t.Fatalf("State() = %v, %v, want = %v, nil", got, err, want)
	}

	if diff := cmp.Diff(
		a2,
		l.Get(a2.ID()),
		cmp.AllowUnexported(simple.Action{}, action.Base{}),
		cmpopts.IgnoreFields(action.Base{}, "fsm"),
	); diff != "" {
		t.Errorf("Get() mismatch (-want +got):\n%v", diff)
	}
}
