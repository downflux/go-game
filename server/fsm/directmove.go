package directmove

import (
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/fsm/commonstate"
	"google.golang.org/protobuf/proto"

	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_DIRECT_MOVE
)

var (
	transitions = []fsm.Transition{
		{
			From:        commonstate.Pending,
			To:          commonstate.Executing,
			VirtualOnly: true,
		},
		{
			From: commonstate.Pending,
			To:   commonstate.Canceled,
		},
		{
			From:        commonstate.Pending,
			To:          commonstate.Finished,
			VirtualOnly: true,
		},
		{
			From: commonstate.Executing,
			To:   commonstate.Canceled,
		},
	}

	FSM = fsm.New(transitions, fsmType)
)

type Action struct {
	*action.Base

	tick id.Tick // Read-only.

	e           moveable.Component    // Read-only.
	status      status.ReadOnlyStatus // Read-only.
	destination *gdpb.Position        // Read-only.
}

func (n *Action) Destination() *gdpb.Position    { return n.destination }
func (n *Action) Accept(v visitor.Visitor) error { return v.Visit(n) }
func (n *Action) ID() id.ActionID                { return id.ActionID(n.e.ID()) }

func (n *Action) Precedence(i action.Action) bool {
	if i.Type() != fsmType {
		return false
	}

	m := i.(*Action)
	return n.tick >= m.tick && !proto.Equal(n.Destination(), m.Destination())
}

func (n *Action) Cancel() error {
	s, err := n.State()
	if err != nil {
		return err
	}

	return n.To(s, commonstate.Canceled, false)
}

func New(
	e moveable.Component,
	dfStatus status.ReadOnlyStatus,
	destination *gdpb.Position) *Action {
	t := dfStatus.Tick()
	return &Action{
		Base:        action.New(FSM, commonstate.Pending),
		e:           e,
		tick:        t,
		status:      dfStatus,
		destination: destination,
	}
}
