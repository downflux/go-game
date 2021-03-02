package chase

import (
	"github.com/downflux/game/engine/fsm/schedule"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/fsm/chase"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_CHASE
)

type Visitor struct {
	visitor.Base

	schedule *schedule.Schedule
	status   status.ReadOnlyStatus
}

func New(dfStatus status.ReadOnlyStatus, schedule *schedule.Schedule) *Visitor {
	return &Visitor{
		Base:     *visitor.New(fsmType),
		schedule: schedule,
		status:   dfStatus,
	}
}

func (v *Visitor) visitFSM(node *chase.Action) error {
	s, err := node.State()
	if err != nil {
		return err
	}

	switch s {
	case chase.OutOfRange:
		m := chase.GenerateMove(node)
		if err := v.schedule.Add(m); err != nil {
			return err
		}
		if err := node.SetMove(m); err != nil {
			return err
		}
	}

	return nil
}

func (v *Visitor) Visit(a visitor.Agent) error {
	if node, ok := a.(*chase.Action); ok {
		return v.visitFSM(node)
	}
	return nil
}
