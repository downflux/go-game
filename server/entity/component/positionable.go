package positionable

import (
	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/id/id"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

type Component interface {
	Position(t id.Tick) *gdpb.Position
	PositionCurve() curve.Curve
}

type Base struct {
	curve curve.Curve
}

func New(c curve.Curve) *Base {
	return &Base{
		curve: c,
	}
}

func (c Base) Position(t id.Tick) *gdpb.Position { return c.curve.Get(t).(*gdpb.Position) }
func (c Base) PositionCurve() curve.Curve        { return c.curve }
