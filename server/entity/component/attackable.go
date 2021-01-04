package attackable

import (
	"github.com/downflux/game/engine/curve/common/timer"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/server/entity/component/positionable"
)

type Component interface {
	positionable.Component

	ID() id.EntityID
	Strength() float64
	Range() float64
	AttackTimerCurve() *timer.Curve
}

type Base struct {
	strength    float64
	attackRange float64
	attackTimer *timer.Curve
}

func New(s float64, r float64, c *timer.Curve) *Base {
	return &Base{
		strength:    s,
		attackRange: r,
		attackTimer: c,
	}
}

func (c Base) Strength() float64              { return c.strength }
func (c Base) Range() float64                 { return c.attackRange }
func (c Base) AttackTimerCurve() *timer.Curve { return c.attackTimer }
