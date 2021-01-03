package tank

import (
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/server/entity/component/attackable"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/entity/component/positionable"
	"github.com/downflux/game/server/entity/component/targetable"
)

var (
	_ entity.Entity          = &Entity{}
	_ moveable.Component     = &Entity{}
	_ attackable.Component   = &Entity{}
	_ targetable.Component   = &Entity{}
	_ positionable.Component = &Entity{}
)
