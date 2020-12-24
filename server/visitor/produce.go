package produce

import (
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/entity/list"
	"github.com/downflux/game/fsm/fsm"
	"github.com/downflux/game/fsm/instance"
	"github.com/downflux/game/fsm/produce"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/visitor/dirty"
	"github.com/downflux/game/server/visitor/visitor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	fcpb "github.com/downflux/game/fsm/api/constants_go_proto"
	serverstatus "github.com/downflux/game/server/service/status"
	vcpb "github.com/downflux/game/server/visitor/api/constants_go_proto"
)

const (
	// entityIDLen is the length of the randomly generated UUID of new
	// Entity objects.
	entityIDLen = 8

	// visitorType is the registered VisitorType of the produce visitor.
	visitorType = vcpb.VisitorType_VISITOR_TYPE_PRODUCE
)

// unsupportedEntityType creates an appropriate error to return when a given
// function cannot handle the EntityType.
func unsupportedEntityType(t gcpb.EntityType) error {
	return status.Errorf(codes.Unimplemented, "creating a new %v entity is not supported", t)
}

// Visitor adds a new Entity to the global state. This struct implements the
// visitor.Visitor interface.
type Visitor struct {
	entities *list.List

	// dirties is a reference to the global cache of mutated Curve and
	// Entity instances.
	dirties *dirty.List

	// dfStatus is reference to the global Executor status struct.
	dfStatus *serverstatus.Status
}

// New creates a new instance of the Visitor struct.
func New(dfStatus *serverstatus.Status, entities *list.List, dirties *dirty.List) *Visitor {
	return &Visitor{
		entities: entities,
		dirties:  dirties,
		dfStatus: dfStatus,
	}
}

// Type returns the registered VisitorType.
func (v *Visitor) Type() vcpb.VisitorType { return visitorType }

// TODO(minkezhang): Delete this function.
func (v *Visitor) Schedule(args interface{}) error { return nil }

func (v *Visitor) visitFSM(i instance.Instance) error {
	if i.Type() != fcpb.FSMType_FSM_TYPE_PRODUCE {
		return nil
	}

	s, err := i.State()
	if err != nil {
		return err
	}

	p := i.(*produce.Instance)

	tick := v.dfStatus.Tick()

	switch s {
	case fsm.State(fcpb.CommonState_COMMON_STATE_EXECUTING.String()):
		defer p.Finish()

		var eid id.EntityID = id.EntityID(id.RandomString(entityIDLen))
		for v.entities.Get(eid) != nil {
			eid = id.EntityID(id.RandomString(entityIDLen))
		}

		var ne entity.Entity
		switch entityType := p.EntityType(); entityType {
		case gcpb.EntityType_ENTITY_TYPE_TANK:
			ne = tank.New(eid, tick, p.SpawnPosition())
			if err := v.dirties.AddEntity(dirty.Entity{ID: eid}); err != nil {
				return err
			}
			if err := v.entities.Append(ne); err != nil {
				return err
			}
		default:
			return unsupportedEntityType(entityType)
		}
		for _, property := range ne.Properties() {
			if err := v.dirties.Add(dirty.Curve{EntityID: eid, Property: property}); err != nil {
				return err
			}
		}
	default:
		return nil
	}

	return nil
}

// Visit mutates an entity.List with a new Entity.
func (v *Visitor) Visit(a visitor.Agent) error {
	switch t := a.AgentType(); t {
	case vcpb.AgentType_AGENT_TYPE_FSM:
		return v.visitFSM(a.(instance.Instance))
	default:
		return nil
	}
}
