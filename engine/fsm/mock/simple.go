package simple

import (
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/visitor"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType  = fcpb.FSMType_FSM_TYPE_MOVE
	Pending  = "PENDING"
	Canceled = "CANCELED"
)

var (
	transitions = []fsm.Transition{
		{From: Pending, To: Canceled},
	}

	FSM = fsm.New(transitions, fsmType)
)

type Action struct {
	*action.Base

	id       id.ActionID
	priority int
}

func New(aid id.ActionID, priority int) *Action {
	return &Action{
		Base:     action.New(FSM, Pending),
		id:       aid,
		priority: priority,
	}
}

func (n *Action) Accept(v visitor.Visitor) error { return v.Visit(n) }
func (n *Action) ID() id.ActionID                { return n.id }

func (n *Action) Precedence(i action.Action) bool {
	if i.Type() != fsmType || n.ID() != i.ID() {
		return false
	}

	m := i.(*Action)
	return n.priority >= m.priority
}

func (n *Action) Cancel() error {
	s, err := n.State()
	if err != nil {
		return err
	}

	return n.To(s, Canceled, false)
}
