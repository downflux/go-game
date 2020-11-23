package produce

import (
	"sync"

	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/visitor/entity/entitylist"
	"github.com/downflux/game/server/service/visitor/entity/tank"
	"github.com/downflux/game/server/service/visitor/dirty"
	"github.com/downflux/game/server/service/visitor/visitor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	serverstatus "github.com/downflux/game/server/service/status"
	vcpb "github.com/downflux/game/server/service/visitor/api/constants_go_proto"
)

const (
	entityIDLen = 8

	visitorType = vcpb.VisitorType_VISITOR_TYPE_PRODUCE
)

func unsupportedEntityType(t gcpb.EntityType) error {
	return status.Errorf(codes.Unimplemented, "creating a new %v entity is not supported", t)
}

type Args struct {
	ScheduledTick float64
	EntityType gcpb.EntityType
	SpawnPosition *gdpb.Position
}

type cacheRow struct {
	entityType gcpb.EntityType
	spawnPosition *gdpb.Position
}

type Visitor struct {
	dirties *dirty.List
        dfStatus *serverstatus.Status


	cacheMux sync.Mutex
	cache map[float64][]cacheRow
}

func New(entities *entitylist.List) *Visitor { return &Visitor{} }

func (v *Visitor) Type() vcpb.VisitorType { return visitorType }

func (v *Visitor) Schedule(args interface{}) error {
	argsImpl := args.(Args)

	v.cacheMux.Lock()
	defer v.cacheMux.Unlock()

	if v.cache == nil {
		v.cache = map[float64][]cacheRow{}
	}

	v.cache[argsImpl.ScheduledTick] = append(
		v.cache[argsImpl.ScheduledTick],
		cacheRow{
			entityType: argsImpl.EntityType,
			spawnPosition: argsImpl.SpawnPosition,
		})

	return nil
}

func (v *Visitor) Visit(e visitor.Entity) error {
	if e.Type() != gcpb.EntityType_ENTITY_TYPE_ENTITY_LIST {
		return nil
	}

	v.cacheMux.Lock()
	defer v.cacheMux.Unlock()

	tick := v.dfStatus.Tick()

	var eid string
	for eid = id.RandomString(entityIDLen); e.(*entitylist.List).Get(eid) != nil; eid = id.RandomString(entityIDLen) {}

	var ne visitor.Entity
	for _, cRow := range v.cache[tick] {
		switch t := cRow.entityType; t {
		case gcpb.EntityType_ENTITY_TYPE_TANK:
			ne = tank.New(eid, tick, cRow.spawnPosition)
			if err := v.dirties.AddEntity(dirty.Entity{ID: eid}); err != nil {
				return err
			}
			if err := e.(*entitylist.List).Add(ne); err != nil {
				return err
			}
		default:
			return unsupportedEntityType(t)
		}
	}

	for _, curveCategory := range ne.CurveCategories() {
		if err := v.dirties.Add(dirty.Curve{EntityID: eid, Category: curveCategory}); err != nil {
			return err
		}
	}

	return nil
}
