package produce

import (
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/entity/list"
	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/produce"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
	serverstatus "github.com/downflux/game/engine/status/status"
)

const (
	// entityIDLen is the length of the randomly generated UUID of new
	// Entity objects.
	entityIDLen = 8

	// fsmType is the registered FSMType of the produce visitor.
	fsmType = fcpb.FSMType_FSM_TYPE_PRODUCE
)

// unsupportedEntityType creates an appropriate error to return when a given
// function cannot handle the EntityType.
func unsupportedEntityType(t gcpb.EntityType) error {
	return status.Errorf(codes.Unimplemented, "creating a new %v entity is not supported", t)
}

// Visitor adds a new Entity to the global state. This struct implements the
// visitor.Visitor interface.
type Visitor struct {
	visitor.Base

	entities *list.List

	// dirty is a reference to the global cache of mutated Curve and
	// Entity instances.
	dirty *dirty.List

	// status is reference to the global Executor status struct.
	status serverstatus.ReadOnlyStatus
}

// New creates a new instance of the Visitor struct.
func New(dfStatus serverstatus.ReadOnlyStatus, entities *list.List, dcs *dirty.List) *Visitor {
	return &Visitor{
		Base:     *visitor.New(fsmType),
		entities: entities,
		dirty:    dcs,
		status:   dfStatus,
	}
}

// TODO(minkezhang): Delete this function.
func (v *Visitor) Schedule(args interface{}) error { return nil }

func (v *Visitor) visitFSM(node *produce.Action) error {
	s, err := node.State()
	if err != nil {
		return err
	}

	tick := v.status.Tick()

	switch s {
	case commonstate.Executing:
		defer node.Finish()

		var eid id.EntityID = id.EntityID(id.RandomString(entityIDLen))
		for v.entities.Get(eid) != nil {
			eid = id.EntityID(id.RandomString(entityIDLen))
		}

		var ne entity.Entity
		switch entityType := node.EntityType(); entityType {
		case gcpb.EntityType_ENTITY_TYPE_TANK:
			ne, err := tank.New(eid, tick, node.SpawnPosition(), node.SpawnClientID())
			if err != nil {
				return err
			}
			if err := v.dirty.AddEntity(dirty.Entity{ID: eid}); err != nil {
				return err
			}
			if err := v.entities.Append(ne); err != nil {
				return err
			}
		default:
			return unsupportedEntityType(entityType)
		}

		if ne != nil {
			for _, property := range ne.Curves().Properties() {
				if err := v.dirty.AddCurve(dirty.Curve{EntityID: eid, Property: property}); err != nil {
					return err
				}
			}
		}
	default:
		return nil
	}

	return nil
}

// Visit mutates an entity.List with a new Entity.
func (v *Visitor) Visit(a visitor.Agent) error {
	if node, ok := a.(*produce.Action); ok {
		return v.visitFSM(node)
	}
	return nil
}
