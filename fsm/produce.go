package produce

import (
	"github.com/downflux/game/fsm/fsm"
	"github.com/downflux/game/fsm/instance"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/status"
	"github.com/downflux/game/server/visitor/visitor"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/fsm/api/constants_go_proto"
)

const (
	fsmType  = fcpb.FSMType_FSM_TYPE_PRODUCE
	idLength = 16
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
		{From: executing, To: finished},
	}
	FSM = fsm.New(transitions, fsmType)
)

type Instance struct {
	*instance.Base

	id            id.InstanceID   // Read-only.
	tick          id.Tick         // Read-only.
	executionTick id.Tick         // Read-only.
	dfStatus      *status.Status  // Read-only.
	entityType    gcpb.EntityType // Read-only.
	spawnPosition *gdpb.Position  // read-only.
}

func New(
	dfStatus *status.Status,
	executionTick id.Tick,
	entityType gcpb.EntityType,
	spawnPosition *gdpb.Position) *Instance {
	return &Instance{
		Base:          instance.New(FSM, pending),
		id:            id.InstanceID(id.RandomString(idLength)),
		executionTick: executionTick,
		dfStatus:      dfStatus,
		entityType:    entityType,
		spawnPosition: spawnPosition,
	}
}

func (n *Instance) EntityType() gcpb.EntityType    { return n.entityType }
func (n *Instance) Accept(v visitor.Visitor) error { return v.Visit(n) }
func (n *Instance) ID() id.InstanceID              { return n.id }
func (n *Instance) SpawnPosition() *gdpb.Position  { return n.spawnPosition }

func (n *Instance) Precedence(i instance.Instance) bool {
	if i.Type() != fsmType {
		return false
	}

	return n.tick > i.(*Instance).tick
}

func (n *Instance) Finish() error {
	s, err := n.State()
	if err != nil {
		return err
	}

	return n.To(s, finished, false)
}

func (n *Instance) Cancel() error {
	s, err := n.State()
	if err != nil {
		return err
	}

	return n.To(s, canceled, false)
}

func (n *Instance) State() (fsm.State, error) {
	tick := n.dfStatus.Tick()

	s, err := n.Base.State()
	if err != nil {
		return unknown, err
	}

	switch s {
	case pending:
		if tick >= n.executionTick {
			if err := n.To(s, executing, true); err != nil {
				return unknown, err
			}
			return executing, nil
		}
		return pending, nil
	default:
		return s, nil
	}
}
