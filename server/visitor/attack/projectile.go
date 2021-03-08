package projectile

import (
	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/fsm/attack/projectile"
	"github.com/downflux/game/server/fsm/commonstate"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_PROJECTILE_SHOOT
)

type Visitor struct {
	visitor.Base                       // Read-only.
	status       status.ReadOnlyStatus // Read-only.
	dirty        *dirty.List
}

func New(s status.ReadOnlyStatus, d *dirty.List) *Visitor {
	return &Visitor{
		Base:   *visitor.New(fsmType),
		status: s,
		dirty:  d,
	}
}

func (v Visitor) Visit(i visitor.Agent) error {
	if node, ok := i.(*projectile.Action); ok {
		return v.visitFSM(node)
	}
	return nil
}

func (v Visitor) visitFSM(i *projectile.Action) error {
	s, err := i.State()
	if err != nil {
		return err
	}

	tick := v.status.Tick()

	switch s {
	case commonstate.Executing:
		c := i.Target().TargetHealthCurve()
		if err := c.Add(
			tick,
			-1*i.Source().AttackStrength(),
		); err != nil {
			return err
		}

		if err := v.dirty.AddCurve(dirty.Curve{
			EntityID: i.Target().ID(),
			Property: c.Property(),
		}); err != nil {
			return err
		}

		if err := i.To(s, commonstate.Finished, false); err != nil {
			return err
		}
	}
	return nil
}
