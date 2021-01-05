package positionable

import (
	"github.com/downflux/game/engine/curve/common/linearmove"
	"github.com/downflux/game/engine/id/id"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

type Component interface {
	Position(t id.Tick) *gdpb.Position
	PositionCurve() *linearmove.Curve
}

type Base struct {
	curve *linearmove.Curve
}

func New(c *linearmove.Curve) *Base {
	return &Base{
		curve: c,
	}
}

func (c Base) Position(t id.Tick) *gdpb.Position { return c.curve.Get(t).(*gdpb.Position) }
func (c Base) PositionCurve() *linearmove.Curve  { return c.curve }
