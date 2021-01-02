package tank

import (
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/server/entity/component/attackable"
	"github.com/downflux/game/server/entity/component/moveable"
)

var (
	_ entity.Entity        = &Entity{}
	_ moveable.Component   = &Entity{}
	_ attackable.Component = &Entity{}
)
