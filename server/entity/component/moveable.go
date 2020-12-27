package moveable

import (
	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/id/id"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

type Component interface {
	Position(t id.Tick) *gdpb.Position
	PositionCurve() curve.Curve
	Velocity() float64
}

type ComponentImpl struct {
	curve    curve.Curve
	velocity float64
}

func New(c curve.Curve, v float64) *ComponentImpl {
	return &ComponentImpl{
		curve:    c,
		velocity: v,
	}
}

func (c ComponentImpl) Position(t id.Tick) *gdpb.Position { return c.curve.Get(t).(*gdpb.Position) }
func (c ComponentImpl) PositionCurve() curve.Curve        { return c.curve }
func (c ComponentImpl) Velocity() float64                 { return c.velocity }
