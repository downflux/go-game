package chase

import (
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/schedule"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/fsm/chase"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
	vcpb "github.com/downflux/game/engine/visitor/api/constants_go_proto"
)

const (
	visitorType = vcpb.VisitorType_VISITOR_TYPE_CHASE
)

type Visitor struct {
	schedule *schedule.Schedule
	status   *status.Status
}

func New(status *status.Status, schedule *schedule.Schedule) *Visitor {
	return &Visitor{
		schedule: schedule,
		status:   status,
	}
}

func (v *Visitor) Type() vcpb.VisitorType {
	return visitorType
}

func (v *Visitor) visitFSM(a action.Action) error {
	if a.Type() != fcpb.FSMType_FSM_TYPE_CHASE {
		return nil
	}

	s, err := a.State()
	if err != nil {
		return err
	}

	c := a.(*chase.Action)
	switch s {
	case chase.OutOfRange:
		m := chase.GenerateMove(c)
		if err := v.schedule.Add(m); err != nil {
			return err
		}
		if err := c.SetMove(m); err != nil {
			return err
		}
	}

	return nil
}

func (v *Visitor) Visit(a visitor.Agent) error {
	switch t := a.AgentType(); t {
	case vcpb.AgentType_AGENT_TYPE_FSM:
		return v.visitFSM(a.(action.Action))
	default:
		return nil
	}
}
