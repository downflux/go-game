package dependent

import (
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/fsm/mock/simple"
	"github.com/downflux/game/engine/id/id"
)

type (
	simpleAction = *simple.Action
)

type Action struct {
	simpleAction
	child *Action
}

func New(aid id.ActionID, priority int, child *Action) *Action {
	return &Action{
		simpleAction: simple.New(aid, priority),
		child:        child,
	}
}

func (n *Action) Child() *Action { return n.child }

func (n *Action) State() (fsm.State, error) {
	if n.child != nil {
		s, err := n.child.State()
		if err != nil || s == simple.Canceled {
			return s, err
		}
	}

	return n.simpleAction.State()
}

func (n *Action) Cancel() error {
	s, err := n.State()
	if err != nil {
		return err
	}
	if err := n.simpleAction.To(s, simple.Canceled, false); err != nil {
		return err
	}
	if n.child != nil {
		return n.child.Cancel()
	}
	return nil
}
