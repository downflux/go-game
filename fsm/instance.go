package instance

import (
	"sync"

	"github.com/downflux/game/fsm/fsm"
	"github.com/downflux/game/server/visitor/visitor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fcpb "github.com/downflux/game/fsm/api/constants_go_proto"
	vcpb "github.com/downflux/game/server/visitor/api/constants_go_proto"
)

type Instance interface {
	visitor.Agent

	Type() fcpb.FSMType
	State() (fsm.State, error)
	To(t fsm.State, virtual bool) error
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

func (n *Base) AgentType() vcpb.AgentType      { return vcpb.AgentType_AGENT_TYPE_FSM }
func (n *Base) Type() fcpb.FSMType             { return n.fsm.Type() }
func (n *Base) Accept(v visitor.Visitor) error { return v.Visit(n) }

func (n *Base) To(t fsm.State, virtual bool) error {
	n.mux.Lock()
	defer n.mux.Unlock()

	f := n.state
	exists, virtualOnly := n.fsm.Exists(f, t)
	if !exists {
		return status.Errorf(codes.FailedPrecondition, "no transition exists between the %v and %v states", f, t)
	}

	if !virtual && virtualOnly {
		return status.Errorf(
			codes.FailedPrecondition,
			"real transition between %v -> %v cannot occur for a virtual-only edge",
		)
	}

	if !virtual {
		n.state = t
	}
	return nil
}

func (n *Base) State() (fsm.State, error) {
	n.mux.Lock()
	defer n.mux.Unlock()

	return n.state, nil
}
