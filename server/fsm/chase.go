package chase

import (
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/entity/component/moveable"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_CHASE
)

var (
	unknown   = fsm.State(fcpb.CommonState_COMMON_STATE_UNKNOWN.String())
	pending   = fsm.State(fcpb.CommonState_COMMON_STATE_PENDING.String())
	executing = fsm.State(fcpb.CommonState_COMMON_STATE_EXECUTING.String())
	canceled  = fsm.State(fcpb.CommonState_COMMON_STATE_CANCELED.String())
	finished  = fsm.State(fcpb.CommonState_COMMON_STATE_FINISHED.String())

	transitions = []fsm.Transition{
		{From: pending, To: executing, VirtualOnly: true},
		{From: pending, To: canceled},
		{From: pending, To: finished, VirtualOnly: true},
		{From: executing, To: pending},
		{From: executing, To: canceled},
	}

	FSM = fsm.New(transitions, fsmType)
)

type Action struct {
	action.Base

	source      moveable.Component
	destination entity.Entity
}

func (a *Action) Cancel() error { return nil }

func (a *Action) Accept(v visitor.Visitor) error { return v.Visit(a) }
