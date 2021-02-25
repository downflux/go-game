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
	vcpb "github.com/downflux/game/engine/visitor/api/constants_go_proto"
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
		map[vcpb.VisitorType]fcpb.FSMType{
			vcpb.VisitorType_VISITOR_TYPE_MOVE: fcpb.FSMType_FSM_TYPE_MOVE,
		},
	)
}

func TestSchedule(t *testing.T) {
	const priority = 0
	aid := id.ActionID("action-id")
	e := newExecutor(t)
	if err := e.Schedule(simpleaction.New(aid, priority)); err != nil {
		t.Fatalf("Schedule() = %v, want = nil", err)
	}
}
