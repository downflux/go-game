package targetable

import (
	"github.com/downflux/game/engine/curve/common/delta"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/server/entity/component/positionable"
)

type Component interface {
	positionable.Component

	ID() id.EntityID
	TargetHealth(t id.Tick) float64
	TargetHealthCurve() *delta.Curve
}

type Base struct {
	health *delta.Curve
}

func New(hp *delta.Curve) *Base {
	return &Base{
		health: hp,
	}
}

func (c *Base) TargetHealth(t id.Tick) float64  { return c.health.Get(t).(float64) }
func (c *Base) TargetHealthCurve() *delta.Curve { return c.health }
