package move

import (
	"sync"

	"github.com/downflux/game/fsm/fsm"
	"github.com/downflux/game/fsm/instance"
	"github.com/downflux/game/server/entity/entity"
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
	unknown   = fsm.State(fcpb.MoveState_MOVE_STATE_UNKNOWN)
	pending   = fsm.State(fcpb.MoveState_MOVE_STATE_PENDING)
	executing = fsm.State(fcpb.MoveState_MOVE_STATE_EXECUTING)
	canceled  = fsm.State(fcpb.MoveState_MOVE_STATE_CANCELED)
	finished  = fsm.State(fcpb.MoveState_MOVE_STATE_FINISHED)
)

var (
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

	scheduledTick id.Tick        // Read-only.
	dfStatus      *status.Status // Read-only.
	destination   *gdpb.Position // Read-only.

	// TODO(minkezhang): Use moveable.Moveable instead.
	e entity.Entity // Read-only.

	// mux guards the Base and nextTick properties.
	mux sync.Mutex

	// TODO(minkezhang): Move nextTick and destination into
	// separate external cache.
	nextTick id.Tick
}

func New(
	e entity.Entity,
	dfStatus *status.Status,
	destination *gdpb.Position) *Instance {
	t := dfStatus.Tick()
	return &Instance{
		Base:          instance.New(FSM, pending),
		e:             e,
		dfStatus:      dfStatus,
		scheduledTick: t,
		nextTick:      t,
		destination:   destination,
	}
}

func (n *Instance) Entity() entity.Entity { return n.e }

func (n *Instance) ID() id.InstanceID { return id.InstanceID(n.e.ID()) }

func (n *Instance) Schedule(t id.Tick) error {
	n.mux.Lock()
	defer n.mux.Unlock()

	s, err := n.stateUnsafe()
	if err != nil {
		return err
	}

	if err := n.To(s, pending, false); err != nil {
		return err
	}

	n.nextTick = t
	return nil
}

func (n *Instance) Precedence(i instance.Instance) bool {
	if i.Type() != fcpb.FSMType_FSM_TYPE_MOVE {
		return false
	}

	return !proto.Equal(n.destination, i.(*Instance).destination)
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
		var t fsm.State

		if proto.Equal(n.destination, c.Get(tick).(*gdpb.Position)) {
			t = finished
		} else if n.nextTick <= tick {
			t = executing
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
