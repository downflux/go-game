package instance

import (
	"sync"

	"github.com/downflux/game/fsm/fsm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Instance interface {
	State() (fsm.State, error)
	To(t fsm.State) error
}

type Base struct {
	fsm *fsm.FSM

	mux   sync.Mutex
	state fsm.State
}

func New(fsm *fsm.FSM, state fsm.State) *Base {
	return &Base{
		fsm:   fsm,
		state: state,
	}
}

func (n *Base) To(t fsm.State) error {
	n.mux.Lock()
	defer n.mux.Unlock()

	f := n.state
	if !n.fsm.Exists(f, t) {
		return status.Errorf(codes.FailedPrecondition, "no transition exists between the %v and %v states", f, t)
	}

	return nil
}

func (n *Base) State() (fsm.State, error) {
	n.mux.Lock()
	defer n.mux.Unlock()

	return n.state, nil
}
