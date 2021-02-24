package simple

import (
	"github.com/downflux/game/engine/visitor/visitor"

	vcpb "github.com/downflux/game/engine/visitor/api/constants_go_proto"
)

type Visitor struct {
	*visitor.BaseVisitor
}

func New() *Visitor {
	return &Visitor{
		BaseVisitor: visitor.NewBaseVisitor(vcpb.VisitorType_VISITOR_TYPE_MOVE),
	}
}

func (v Visitor) Visit(a visitor.Agent) error { return nil }
