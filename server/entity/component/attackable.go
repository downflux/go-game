package attackable

import (
	// "github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/id/id"
	//  gdpb "github.com/downflux/game/api/data_go_proto"
)

type Component interface {
	ID() id.EntityID
	Strength() float64
}

type Base struct {
	strength float64
}

func New(s float64) *Base {
	return &Base{
		strength: s,
	}
}

func (c Base) Strength() float64 { return c.strength }
