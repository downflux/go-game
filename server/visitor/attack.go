package attack

import (
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/fsm/attack"

	vcpb "github.com/downflux/game/engine/visitor/api/constants_go_proto"
)

const (
	visitorType = vcpb.VisitorType_VISITOR_TYPE_CHASE
)

type Visitor struct {
	visitor.BaseVisitor
}

func New() *Visitor {
	return &Visitor{
		BaseVisitor: *visitor.NewBaseVisitor(visitorType),
	}
}

func (v *Visitor) visitFSM(node *attack.Action) error { return nil }

func (v *Visitor) Visit(a visitor.Agent) error {
	if node, ok := a.(*attack.Action); ok {
		return v.visitFSM(node)
	}
	return nil
}
