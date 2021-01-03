package attackable

import (
	"github.com/downflux/game/engine/id/id"
)

type Component interface {
	ID() id.EntityID
	Strength() float64
	Range() float64
}

type Base struct {
	strength    float64
	attackRange float64
}

func New(s float64, r float64) *Base {
	return &Base{
		strength:    s,
		attackRange: r,
	}
}

func (c Base) Strength() float64 { return c.strength }
func (c Base) Range() float64    { return c.attackRange }
