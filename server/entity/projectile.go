package projectile

import (
	"github.com/downflux/game/engine/curve/list"
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/entity/component/positionable"
)

type (
	moveComponent     = moveable.Base
	positionComponent = positionable.Base
)

type Entity struct {
	entity.Base
	entity.LifeCycle
	moveComponent
	positionComponent

	// curves is a list of Curves tracking the Entity properties.
	curves *list.List
}

func (e *Entity) Curves() *list.List { return e.curves }
