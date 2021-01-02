package delta

import (
	"reflect"
	"testing"

	"github.com/downflux/game/engine/curve/common/step"
	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

const (
	t0 = id.Tick(100)
	v0 = float64(101)
)

var (
	_ curve.Curve = &Curve{}
)

func newTestCurve() *Curve {
	stepCurve := step.New(
		"entity-id",
		0,
		gcpb.EntityProperty_ENTITY_PROPERTY_UNKNOWN,
		reflect.TypeOf(float64(0)),
	)
	stepCurve.Add(t0, v0)
	return New(stepCurve)
}

func TestAdd(t *testing.T) {
	testConfigs := []struct {
		name  string
		tick  id.Tick
		delta float64
		want  float64
	}{
		{name: "TestBefore", tick: t0 - 1, delta: 1, want: v0 + 1},
		{name: "TestAfter", tick: t0 + 1, delta: 1, want: v0},
		{name: "TestAt", tick: t0, delta: 1, want: v0 + 1},
		{name: "TestNegative", tick: t0, delta: -1, want: v0 + -1},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			testCurve := newTestCurve()
			if err := testCurve.Add(c.tick, c.delta); err != nil {
				t.Fatalf("Add() = %v, want = nil", err)
			}

			if got := testCurve.Get(t0); got != c.want {
				t.Fatalf("Get() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestExport(t *testing.T) {
	c := newTestCurve()
	pb := c.Export(0)
	if got := pb.GetType(); got != c.Curve.Type() {
		t.Errorf("GetType() = %v, want = %v", got, c.Curve.Type())
	}
}
