package attack

import (
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/fsm/attack"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/engine/gamestate/dirty"

	vcpb "github.com/downflux/game/engine/visitor/api/constants_go_proto"
)

const (
	visitorType = vcpb.VisitorType_VISITOR_TYPE_ATTACK
)

type Visitor struct {
	visitor.BaseVisitor

	status status.ReadOnlyStatus
	dirty  *dirty.List
}

func New(dfStatus status.ReadOnlyStatus, dirties *dirty.List) *Visitor {
	return &Visitor{
		BaseVisitor: *visitor.NewBaseVisitor(visitorType),
		status:      dfStatus,
		dirty: dirties,
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
		dirtyCurves := []dirty.Curve{
			{node.Source().ID(), node.Source().AttackTimerCurve().Property()},
			{node.Target().ID(), node.Target().HealthCurve().Property()},
		}
		for _, c := range dirtyCurves {
	                if err := v.dirty.AddCurve(c); err != nil {
	                        return err
	                }
		}

		if err := node.Source().AttackTimerCurve().Add(tick, true); err != nil {
			return err
		}
		return node.Target().HealthCurve().Add(tick, -1*node.Source().Strength())
	}
	return nil
}

func (v *Visitor) Visit(a visitor.Agent) error {
	if node, ok := a.(*attack.Action); ok {
		return v.visitFSM(node)
	}
	return nil
}
