package moveable

import (
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/server/entity/component/positionable"
)

type Component interface {
	positionable.Component

	ID() id.EntityID
	MoveVelocity() float64
}

type Base struct {
	velocity float64
}

func New(v float64) *Base {
	return &Base{
		velocity: v,
	}
}

func (c Base) MoveVelocity() float64 { return c.velocity }
