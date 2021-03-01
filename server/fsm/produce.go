package produce

import (
	"github.com/downflux/game/engine/entity/acl"
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/fsm/commonstate"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	fsmType  = fcpb.FSMType_FSM_TYPE_PRODUCE
	idLength = 16
)

var (
	transitions = []fsm.Transition{
		{From: commonstate.Pending, To: commonstate.Executing, VirtualOnly: true},
		{From: commonstate.Pending, To: commonstate.Canceled},
		{From: commonstate.Executing, To: commonstate.Finished},
	}
	FSM = fsm.New(transitions, fsmType)
)

type Action struct {
	*action.Base

	id            id.ActionID           // Read-only.
	tick          id.Tick               // Read-only.
	executionTick id.Tick               // Read-only.
	status        status.ReadOnlyStatus // Read-only.
	entityType    gcpb.EntityType       // Read-only.
	spawnPosition *gdpb.Position        // Read-only.
	spawnClientID id.ClientID           // Read-only.
	permission    acl.Permission        // Read-only.
}

func New(
	dfStatus status.ReadOnlyStatus,
	executionTick id.Tick,
	entityType gcpb.EntityType,
	spawnPosition *gdpb.Position,
	spawnClientID id.ClientID,
	permission acl.Permission) *Action {
	return &Action{
		Base:          action.New(FSM, commonstate.Pending),
		id:            id.ActionID(id.RandomString(idLength)),
		executionTick: executionTick,
		status:        dfStatus,
		entityType:    entityType,
		spawnPosition: spawnPosition,
		spawnClientID: spawnClientID,
		permission:    permission,
	}
}

func (n *Action) EntityType() gcpb.EntityType    { return n.entityType }
func (n *Action) Accept(v visitor.Visitor) error { return v.Visit(n) }
func (n *Action) ID() id.ActionID                { return n.id }
func (n *Action) SpawnPosition() *gdpb.Position  { return n.spawnPosition }
func (n *Action) SpawnClientID() id.ClientID     { return n.spawnClientID }
func (n *Action) Permission() acl.Permission     { return n.permission }

func (n *Action) Precedence(i action.Action) bool {
	if i.Type() != fsmType {
		return false
	}

	return n.tick > i.(*Action).tick
}

func (n *Action) Finish() error {
	s, err := n.State()
	if err != nil {
		return err
	}

	return n.To(s, commonstate.Finished, false)
}

func (n *Action) Cancel() error {
	s, err := n.State()
	if err != nil {
		return err
	}

	return n.To(s, commonstate.Canceled, false)
}

func (n *Action) State() (fsm.State, error) {
	tick := n.status.Tick()

	s, err := n.Base.State()
	if err != nil {
		return commonstate.Unknown, err
	}

	switch s {
	case commonstate.Pending:
		if tick >= n.executionTick {
			if err := n.To(s, commonstate.Executing, true); err != nil {
				return commonstate.Unknown, err
			}
			return commonstate.Executing, nil
		}
		return commonstate.Pending, nil
	default:
		return s, nil
	}
}
