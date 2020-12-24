package move

import (
	"sync"

	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"google.golang.org/protobuf/proto"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_MOVE
)

var (
	unknown   = fsm.State(fcpb.CommonState_COMMON_STATE_UNKNOWN.String())
	pending   = fsm.State(fcpb.CommonState_COMMON_STATE_PENDING.String())
	executing = fsm.State(fcpb.CommonState_COMMON_STATE_EXECUTING.String())
	canceled  = fsm.State(fcpb.CommonState_COMMON_STATE_CANCELED.String())
	finished  = fsm.State(fcpb.CommonState_COMMON_STATE_FINISHED.String())

	transitions = []fsm.Transition{
		{From: pending, To: executing, VirtualOnly: true},
		{From: pending, To: canceled},
		{From: pending, To: finished, VirtualOnly: true},
		{From: executing, To: pending},
		{From: executing, To: canceled},
	}

	FSM = fsm.New(transitions, fsmType)
)

type Action struct {
	*action.Base

	// tick is the tick at which the command was originally
	// scheduled.
	tick id.Tick // Read-only.

	dfStatus    *status.Status // Read-only.
	destination *gdpb.Position // Read-only.

	// TODO(minkezhang): Use moveable.Moveable instead.
	e entity.Entity // Read-only.

	// mux guards the Base and executionTick properties.
	mux sync.Mutex

	// TODO(minkezhang): Move executionTick and destination into
	// separate external cache.
	executionTick id.Tick
}

// New constructs a new Action FSM action.
//
// TODO(minkezhang): Add executionTick arg to allow for scheduling in the
// future.
func New(
	e entity.Entity,
	dfStatus *status.Status,
	destination *gdpb.Position) *Action {
	t := dfStatus.Tick()
	return &Action{
		Base:          action.New(FSM, pending),
		e:             e,
		dfStatus:      dfStatus,
		tick:          t,
		executionTick: t,
		destination:   destination,
	}
}

func (n *Action) Accept(v visitor.Visitor) error { return v.Visit(n) }
func (n *Action) Entity() entity.Entity          { return n.e }
func (n *Action) ID() id.ActionID                { return id.ActionID(n.e.ID()) }

// SchedulePartialMove allows us to mutate the FSM action to deal with
// partial moves. This allows us to know when the visitor should make the next
// meaningful calculation.
func (n *Action) SchedulePartialMove(t id.Tick) error {
	n.mux.Lock()
	defer n.mux.Unlock()

	s, err := n.stateUnsafe()
	if err != nil {
		return err
	}

	if err := n.To(s, pending, false); err != nil {
		return err
	}

	n.executionTick = t
	return nil
}

func (n *Action) Precedence(i action.Action) bool {
	if i.Type() != fsmType {
		return false
	}

	return n.tick > i.(*Action).tick && !proto.Equal(n.Destination(), i.(*Action).Destination())
}

// TODO(minkezhang): Return a cloned instance instead.
func (n *Action) Destination() *gdpb.Position { return n.destination }

func (n *Action) Cancel() error {
	n.mux.Lock()
	defer n.mux.Unlock()

	s, err := n.stateUnsafe()
	if err != nil {
		return err
	}

	return n.To(s, canceled, false)
}

func (n *Action) State() (fsm.State, error) {
	n.mux.Lock()
	defer n.mux.Unlock()

	return n.stateUnsafe()
}

func (n *Action) stateUnsafe() (fsm.State, error) {
	tick := n.dfStatus.Tick()

	s, err := n.Base.State()
	if err != nil {
		return unknown, err
	}

	switch s {
	case pending:
		c := n.e.Curve(gcpb.EntityProperty_ENTITY_PROPERTY_POSITION)
		var t fsm.State = unknown

		if n.executionTick <= tick {
			t = executing
			if proto.Equal(n.destination, c.Get(tick).(*gdpb.Position)) {
				t = finished
			}
		}

		if t != unknown {
			if err := n.To(s, t, true); err != nil {
				return unknown, err
			}
			return t, nil
		}

		return pending, nil
	default:
		return s, nil
	}
}
