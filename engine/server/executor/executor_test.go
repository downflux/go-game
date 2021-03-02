// Package executor tests the game executor.
//
// TODO(minkezhang): Add testing for broadcasting to clients here.
package executor

import (
	"testing"
	"time"

	"github.com/downflux/game/engine/fsm/schedule"
	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/gamestate/gamestate"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/mock/simple"
	"github.com/downflux/game/engine/visitor/visitor"

	entitylist "github.com/downflux/game/engine/entity/list"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
	simpleaction "github.com/downflux/game/engine/fsm/mock/simple"
	serverstatus "github.com/downflux/game/engine/status/status"
	visitorlist "github.com/downflux/game/engine/visitor/list"
)

var (
	tickDuration = 100 * time.Millisecond
)

func newExecutor(t *testing.T) *Executor {
	visitors, err := visitorlist.New([]visitor.Visitor{simple.New()})
	if err != nil {
		t.Fatalf("New() = %v, want = nil", err)
	}

	return New(
		visitors,
		gamestate.New(
			serverstatus.New(tickDuration),
			entitylist.New(),
		),
		dirty.New(),
		schedule.New([]fcpb.FSMType{fcpb.FSMType_FSM_TYPE_MOVE}),
	)
}

func TestSchedule(t *testing.T) {
	const priority = 0
	aid := id.ActionID("action-id")

	e := newExecutor(t)
	if err := e.Schedule(simpleaction.New(aid, priority)); err != nil {
		t.Errorf("Schedule() = %v, want = nil", err)
	}
}

func TestDoTick(t *testing.T) {
	const priority = 0
	aid := id.ActionID("action-id")

	e := newExecutor(t)

	tick := e.Status().GetTick()

	// Executor should have added action to the cache.
	if err := e.Schedule(simpleaction.New(aid, priority)); err != nil {
		t.Fatalf("Schedule() = %v, want = nil", err)
	}
	if a := e.scheduleCache.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(aid); a == nil {
		t.Fatal("Get() = nil, want a non-nil value")
	}

	mock := e.visitors.Visitor(fcpb.FSMType_FSM_TYPE_MOVE).(*simple.Visitor)
	count := mock.Count()

	e.doTick()
	if got := e.Status().GetTick(); got != tick+1 {
		t.Fatalf("GetTick() = %v, want = %v", got, tick+1)
	}

	if a := e.scheduleCache.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(aid); a != nil {
		t.Fatalf("Get() = %v, want = nil", a)
	}
	if a := e.schedule.Get(fcpb.FSMType_FSM_TYPE_MOVE).Get(aid); a == nil {
		t.Fatal("Get() = nil, want a non-nil value")
	}

	// TODO(minkezhang): Change to 1 when we remove Visit(List).
	if get := mock.Count(); get != count+1 {
		t.Errorf("Count() = %v, want = %v", get, count+1)
	}
}
