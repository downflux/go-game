package simple

import (
	"github.com/downflux/game/engine/curve/list"
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type Entity struct {
	entity.LifeCycle

	id id.EntityID
}

func New(eid id.EntityID) *Entity {
	return &Entity{
		id: eid,
	}
}

func (e *Entity) Type() gcpb.EntityType { return gcpb.EntityType_ENTITY_TYPE_TANK }
func (e *Entity) ID() id.EntityID       { return e.id }
func (e *Entity) Curves() *list.List    { return nil }
