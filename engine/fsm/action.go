package action

import (
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/visitor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

type Action interface {
	visitor.Agent

	Cancel() error
	ID() id.ActionID
	Precedence(i Action) bool
	State() (fsm.State, error)
	To(f fsm.State, t fsm.State, virtual bool) error
	Type() fcpb.FSMType
}

type Base struct {
	fsm   *fsm.FSM
	state fsm.State
}

func New(fsm *fsm.FSM, state fsm.State) *Base {
	return &Base{
		fsm:   fsm,
		state: state,
	}
}

func (n *Base) State() (fsm.State, error) { return n.state, nil }
func (n *Base) Type() fcpb.FSMType        { return n.fsm.Type() }

func (n *Base) To(f fsm.State, t fsm.State, virtual bool) error {
	exists, virtualOnly := n.fsm.Exists(f, t)
	if !exists {
		return status.Errorf(codes.FailedPrecondition, "no transition exists between the %v and %v states", f, t)
	}

	if !virtual && virtualOnly {
		return status.Errorf(
			codes.FailedPrecondition,
			"real transition between %v -> %v cannot occur for a virtual-only edge",
			f,
			t,
		)
	}

	if !virtual {
		n.state = t
	}
	return nil
}
