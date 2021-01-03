package timer

import (
	"testing"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

var (
	_ curve.Curve = &Curve{}
)

func TestAdd(t *testing.T) {
	const eid = "entity-id"
	const cooldown = 10
	const t0 = 0
	const t1 = t0 + cooldown + cooldown
	const property = gcpb.EntityProperty_ENTITY_PROPERTY_UNKNOWN

	testConfigs := []struct {
		name    string
		tick    id.Tick
		success bool
	}{
		{name: "AddBeforeSuccess", tick: t0 - cooldown, success: true},
		{name: "AddBeforeFailure", tick: t0 - (cooldown - 1), success: false},
		{name: "AddAfterSuccess", tick: t1 + cooldown, success: true},
		{name: "AddAfterFailure", tick: t1 + (cooldown - 1), success: false},
		{name: "AddBetweenSuccess", tick: t0 + cooldown, success: true},
		{name: "AddBetweenFailure", tick: t0 + cooldown + 1, success: false},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			timerCurve := New(eid, t0, cooldown, property)
			timerCurve.Data().Set(t0, true)
			timerCurve.Data().Set(t1, true)

			err := timerCurve.Add(c.tick, true)
			if c.success && err != nil {
				t.Errorf("Add() = %v, want = nil", err)
			} else if !c.success && err == nil {
				t.Error("Add() = nil, want a non-nil error")
			}
		})
	}
}

func TestGet(t *testing.T) {
	const eid = "entity-id"
	const cooldown = 10
	const t0 = 0
	const t1 = t0 + cooldown
	const t2 = t1 + 1
	const property = gcpb.EntityProperty_ENTITY_PROPERTY_UNKNOWN

	timerCurve := New(eid, t0, cooldown, property)
	timerCurve.Data().Set(t0, true)
	timerCurve.Data().Set(t2, true)

	testConfigs := []struct {
		name string
		tick id.Tick
		want bool
	}{
		{name: "GetBefore", tick: t0 - 1, want: false},
		{name: "GetImmediatelyAfter", tick: t0 + 1, want: true},
		{name: "GetAfterCooldown", tick: t1, want: false},
		{name: "GetAfter", tick: t1, want: false},
		{name: "GetEnd", tick: t2 + cooldown, want: false},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := timerCurve.Get(c.tick); got != c.want {
				t.Errorf("Get() = %v, want = %v", got, c.want)
			}
		})
	}
}
