// Package attackable impart a primary attack function to the entity.
//
// An attack may have an optional travel time. Of attacks with a travel time,
// the attack may hit a stationary position or chase a target within some attack
// range.
package attackable

import (
	"github.com/downflux/game/engine/curve/common/step"
	"github.com/downflux/game/engine/curve/common/timer"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/server/entity/component/positionable"
)

type Component interface {
	positionable.Component

	ID() id.EntityID
	AttackStrength() float64

	// AttackRange specifies the distance at which an attacker may start
	// attacking the target.
	AttackRange() float64

	// AttackVelocity expresses the speed of the projectile. A velocity may
	// be of speed positive infinity, representing a hitscan weapon.
	AttackVelocity() float64

	AttackTargetCurve() *step.Curve
	AttackTimerCurve() *timer.Curve
}

type Base struct {
	strength    float64
	attackRange float64
	velocity    float64
	targetCurve *step.Curve
	attackTimer *timer.Curve
}

func New(s float64, r float64, v float64, t *step.Curve, c *timer.Curve) *Base {
	return &Base{
		strength:    s,
		attackRange: r,
		velocity:    v,
		targetCurve: t,
		attackTimer: c,
	}
}

func (c Base) AttackStrength() float64        { return c.strength }
func (c Base) AttackRange() float64           { return c.attackRange }
func (c Base) AttackVelocity() float64        { return c.velocity }
func (c Base) AttackTimerCurve() *timer.Curve { return c.attackTimer }
func (c Base) AttackTargetCurve() *step.Curve { return c.targetCurve }
