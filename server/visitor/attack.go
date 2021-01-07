package attack

import (
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/fsm/attack"
	"github.com/downflux/game/server/fsm/commonstate"

	vcpb "github.com/downflux/game/engine/visitor/api/constants_go_proto"
)

const (
	visitorType = vcpb.VisitorType_VISITOR_TYPE_ATTACK
)

type Visitor struct {
	visitor.BaseVisitor

	status status.ReadOnlyStatus
}

func New(dfStatus status.ReadOnlyStatus) *Visitor {
	return &Visitor{
		BaseVisitor: *visitor.NewBaseVisitor(visitorType),
		status:      dfStatus,
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
