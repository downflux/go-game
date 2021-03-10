package attack

import (
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/schedule"
	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/fsm/attack/attack"
	"github.com/downflux/game/server/fsm/attack/projectile"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move/move"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_ATTACK
)

type Visitor struct {
	visitor.Base                       // Read-only.
	status       status.ReadOnlyStatus // Read-only.

	schedule *schedule.Schedule
	dirty    *dirty.List
}

func New(
	dfStatus status.ReadOnlyStatus,
	dirtystate *dirty.List,
	schedule *schedule.Schedule) *Visitor {
	return &Visitor{
		Base:     *visitor.New(fsmType),
		status:   dfStatus,
		dirty:    dirtystate,
		schedule: schedule,
	}
}

func (v *Visitor) visitFSM(node *attack.Action) error {
	s, err := node.State()
	if err != nil {
		return err
	}

	tick := v.status.Tick()
	switch s {
	case commonstate.Executing:
		// TODO(minkezhang): Implement a string step curve for
		// recording targets, ENTITY_PROPERTY_ATTACK_TARGET.
		dcs := []dirty.Curve{
			{node.Source().ID(), node.Source().AttackTimerCurve().Property()},
		}
		for _, c := range dcs {
			if err := v.dirty.AddCurve(c); err != nil {
				return err
			}
		}

		if err := node.Source().AttackTimerCurve().Add(tick, true); err != nil {
			return err
		}

		// TODO(minkezhang): Move generator functions to attack FSM
		// instead. See fsm/chase for example.
		moveAction := move.New(
			node.Source().AttackProjectile(),
			v.status,
			node.Target().Position(tick),
			move.Direct)
		projectileAction := projectile.New(node.Source(), node.Target(), moveAction)

		node.SetProjectileMove(projectileAction)
		return v.schedule.Extend([]action.Action{moveAction, projectileAction})
	}
	return nil
}

func (v *Visitor) Visit(a visitor.Agent) error {
	if node, ok := a.(*attack.Action); ok {
		return v.visitFSM(node)
	}
	return nil
}
