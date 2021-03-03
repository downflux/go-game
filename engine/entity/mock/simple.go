package simple

import (
	"github.com/downflux/game/engine/entity/component/curve"
	"github.com/downflux/game/engine/entity/component/lifecycle"
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type (
	lifecycleComponent = lifecycle.Component
	curveComponent     = curve.Component
)

type Entity struct {
	entity.Base
	lifecycleComponent
	curveComponent
}

func New(eid id.EntityID) *Entity {
	return &Entity{
		Base: *entity.New(gcpb.EntityType_ENTITY_TYPE_TANK, eid, nil),
	}
}
