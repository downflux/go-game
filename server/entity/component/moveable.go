package moveable

import (
	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/server/entity/component/positionable"
)

type Component interface {
	positionable.Component

	ID() id.EntityID
	Velocity() float64
}

type Base struct {
	positionable.Base
	velocity float64
}

func New(c curve.Curve, v float64) *Base {
	return &Base{
		Base:     *positionable.New(c),
		velocity: v,
	}
}

func (c Base) Velocity() float64 { return c.velocity }
