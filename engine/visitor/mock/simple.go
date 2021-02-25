package simple

import (
	"github.com/downflux/game/engine/visitor/visitor"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

type Visitor struct {
	*visitor.Base

	counter int
}

func New() *Visitor {
	return &Visitor{
		Base: visitor.New(fcpb.FSMType_FSM_TYPE_MOVE),
	}
}

func (v *Visitor) Visit(a visitor.Agent) error {
	v.counter += 1
	return nil
}

func (v *Visitor) Count() int { return v.counter }
