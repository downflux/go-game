package simple

import (
	"github.com/downflux/game/engine/curve/list"
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type Entity struct {
	entity.Base
	entity.LifeCycle
}

func New(eid id.EntityID) *Entity {
	return &Entity{
		Base: *entity.New(gcpb.EntityType_ENTITY_TYPE_TANK, eid),
	}
}

func (e *Entity) Curves() *list.List { return nil }
