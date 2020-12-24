package move

import (
	"sync"

	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/fsm/fsm"
	"github.com/downflux/game/fsm/instance"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/status"
	"google.golang.org/protobuf/proto"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/fsm/api/constants_go_proto"
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

type Instance struct {
	*instance.Base

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

// New constructs a new Instance FSM instance.
//
// TODO(minkezhang): Add executionTick arg to allow for scheduling in the
// future.
func New(
	e entity.Entity,
	dfStatus *status.Status,
	destination *gdpb.Position) *Instance {
	t := dfStatus.Tick()
	return &Instance{
		Base:          instance.New(FSM, pending),
		e:             e,
		dfStatus:      dfStatus,
		tick:          t,
		executionTick: t,
		destination:   destination,
	}
}

func (n *Instance) Accept(v visitor.Visitor) error { return v.Visit(n) }
func (n *Instance) Entity() entity.Entity          { return n.e }
func (n *Instance) ID() id.InstanceID              { return id.InstanceID(n.e.ID()) }

// SchedulePartialMove allows us to mutate the FSM instance to deal with
// partial moves. This allows us to know when the visitor should make the next
// meaningful calculation.
func (n *Instance) SchedulePartialMove(t id.Tick) error {
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

func (n *Instance) Precedence(i instance.Instance) bool {
	if i.Type() != fsmType {
		return false
	}

	return n.tick > i.(*Instance).tick && !proto.Equal(n.Destination(), i.(*Instance).Destination())
}

// TODO(minkezhang): Return a cloned instance instead.
func (n *Instance) Destination() *gdpb.Position { return n.destination }

func (n *Instance) Cancel() error {
	n.mux.Lock()
	defer n.mux.Unlock()

	s, err := n.stateUnsafe()
	if err != nil {
		return err
	}

	return n.To(s, canceled, false)
}

func (n *Instance) State() (fsm.State, error) {
	n.mux.Lock()
	defer n.mux.Unlock()

	return n.stateUnsafe()
}

func (n *Instance) stateUnsafe() (fsm.State, error) {
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
