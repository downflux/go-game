package list

import (
	"log"
	"testing"

	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/fsm/instance"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/fsm/move"
	"github.com/downflux/game/server/entity/tank"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"

	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_MOVE
)

func TestNew(t *testing.T) {
	l := New(fsmType)

	if got := l.Type(); got != fsmType {
		t.Fatalf("Type() = %v, want = %v", got, fsmType)
	}

	if got := l.AgentType(); got != agentType {
		t.Errorf("AgentType() = %v, want = %v", got, agentType)
	}
}

func TestAddError(t *testing.T) {
	l := New(fcpb.FSMType_FSM_TYPE_UNKNOWN)
	i := move.New(tank.New("entity-id", 0, nil), status.New(0), nil)

	log.Println(l.Type(), i.Type())

	if err := l.Add(i); err == nil {
		t.Error("Add() = nil, want a non-nil error")
	}
}

func TestAdd(t *testing.T) {
	const iid = "entity-id"

	l := New(fsmType)
	i := move.New(tank.New(iid, 0, nil), status.New(0), nil)

	if err := l.Add(i); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got := l.Get(iid).ID(); got != iid {
		t.Errorf("ID() = %v, want = %v", got, iid)
	}
}

func TestAddCancel(t *testing.T) {
	const iid = "entity-id"

	l := New(fsmType)

	e := tank.New(iid, 0, nil)
	dfStatus := status.New(0)
	i1 := move.New(e, dfStatus, &gdpb.Position{X: 0, Y: 0})
	dfStatus.IncrementTick()
	i2 := move.New(e, dfStatus, &gdpb.Position{X: 1, Y: 1})

	if err := l.Add(i1); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}
	if err := l.Add(i2); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	want := fsm.State(fcpb.CommonState_COMMON_STATE_CANCELED.String())
	if got, err := i1.State(); err != nil || got != want {
		t.Fatalf("State() = %v, %v, want = %v, nil", got, err, want)
	}

	if diff := cmp.Diff(
		i2,
		l.Get(i2.ID()),
		cmp.AllowUnexported(move.Instance{}, instance.Base{}),
		cmpopts.IgnoreFields(instance.Base{}, "mux", "fsm"),
		cmpopts.IgnoreFields(move.Instance{}, "dfStatus", "mux", "e"),
		protocmp.Transform(),
	); diff != "" {
		t.Errorf("Get() mismatch (-want +got):\n%v", diff)
	}
}
