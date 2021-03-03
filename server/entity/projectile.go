package projectile

import (
	"github.com/downflux/game/engine/entity/component/curve"
	"github.com/downflux/game/engine/entity/component/lifecycle"
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/entity/component/positionable"
)

type (
	moveComponent      = moveable.Base
	positionComponent  = positionable.Base
	curveComponent     = curve.Component
	lifecycleComponent = lifecycle.Component
)

type Entity struct {
	entity.Base
	curveComponent
	lifecycleComponent
	moveComponent
	positionComponent
}
