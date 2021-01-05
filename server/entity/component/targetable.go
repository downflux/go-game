package targetable

import (
	"github.com/downflux/game/engine/curve/common/delta"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/server/entity/component/positionable"
)

type Component interface {
	positionable.Component

	Health(t id.Tick) float64
	HealthCurve() *delta.Curve
}

type Base struct {
	health *delta.Curve
}

func New(hp *delta.Curve) *Base {
	return &Base{
		health: hp,
	}
}

func (c *Base) Health(t id.Tick) float64  { return c.health.Get(t).(float64) }
func (c *Base) HealthCurve() *delta.Curve { return c.health }
