package projectile

import (
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/entity/component/positionable"
)

var (
	_ entity.Entity          = &Entity{}
	_ moveable.Component     = &Entity{}
	_ positionable.Component = &Entity{}
)
