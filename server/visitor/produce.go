// Package produce impements logic for adding new Entity instances.
package produce

import (
	"sync"

	"github.com/downflux/game/server/entity/entity"
	"github.com/downflux/game/server/entity/entitylist"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/visitor/dirty"
	"github.com/downflux/game/server/visitor/visitor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
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

// Args is an external-facing struct to be passed into Schedule. This
// represents the scheduled time and Entity that will be created.
type Args struct {
	ScheduledTick id.Tick
	EntityType    gcpb.EntityType
	SpawnPosition *gdpb.Position
}

// cacheRow represents a scheduled add command. The tick at which this command
// executes is stored in the wrapping map object.
type cacheRow struct {
	entityType    gcpb.EntityType
	spawnPosition *gdpb.Position
}

// Visitor adds a new Entity to the global state. This struct implements the
// visitor.Visitor interface.
type Visitor struct {
	visitor.Base
	visitor.Leaf

	// dirties is a reference to the global cache of mutated Curve and
	// Entity instances.
	dirties *dirty.List

	// dfStatus is reference to the global Executor status struct.
	dfStatus *serverstatus.Status

	// cacheMux guards the cache property.
	cacheMux sync.Mutex

	// cache is an internal collection of scheduled add Entity commands.
	cache map[id.Tick][]cacheRow
}

// New creates a new instance of the Visitor struct.
func New(dfStatus *serverstatus.Status, dirties *dirty.List) *Visitor {
	return &Visitor{
		dirties:  dirties,
		dfStatus: dfStatus,
	}
}

// Type returns the registered VisitorType.
func (v *Visitor) Type() vcpb.VisitorType { return visitorType }

// Schedule adds an add Entity command to the internal cache.
func (v *Visitor) Schedule(args interface{}) error {
	argsImpl := args.(Args)

	v.cacheMux.Lock()
	defer v.cacheMux.Unlock()

	if v.cache == nil {
		v.cache = map[id.Tick][]cacheRow{}
	}

	v.cache[argsImpl.ScheduledTick] = append(
		v.cache[argsImpl.ScheduledTick],
		cacheRow{
			entityType:    argsImpl.EntityType,
			spawnPosition: argsImpl.SpawnPosition,
		})

	return nil
}

// Visit mutates an EntityList with a new Entity.
func (v *Visitor) Visit(a visitor.Agent) error {
	if a.AgentType() != vcpb.AgentType_AGENT_TYPE_ENTITY {
		return nil
	}

	e := a.(entity.Entity)
	if e.Type() != gcpb.EntityType_ENTITY_TYPE_ENTITY_LIST {
		return nil
	}
	tick := v.dfStatus.Tick()

	v.cacheMux.Lock()
	defer v.cacheMux.Unlock()

	var err error
	for t := id.Tick(0); t <= tick; t++ {
		if te := func() error {
			defer delete(v.cache, t)

			var err error
			for i, cRow := range v.cache[t] {
				var eid id.EntityID = id.EntityID(id.RandomString(entityIDLen))
				for e.(*entitylist.List).Get(eid) != nil {
					eid = id.EntityID(id.RandomString(entityIDLen))
				}

				if te := func() error {
					defer func() { v.cache[t][i] = cacheRow{} }()

					var ne entity.Entity
					switch entityType := cRow.entityType; entityType {
					case gcpb.EntityType_ENTITY_TYPE_TANK:
						ne = tank.New(eid, tick, cRow.spawnPosition)
						if err := v.dirties.AddEntity(dirty.Entity{ID: eid}); err != nil {
							return err
						}
						if err := e.(*entitylist.List).Add(ne); err != nil {
							return err
						}
					default:
						return unsupportedEntityType(entityType)
					}

					for _, curveCategory := range ne.CurveCategories() {
						if err := v.dirties.Add(dirty.Curve{EntityID: eid, Category: curveCategory}); err != nil {
							return err
						}
					}

					return nil
				}(); te != nil && err == nil {
					err = te
				}
			}
			return err
		}(); te != nil && err == nil {
			err = te
		}
	}
	return err
}
