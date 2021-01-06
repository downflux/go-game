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

type datum struct {
	tick  id.Tick
	value float64
}

func TestAdd(t *testing.T) {
	testConfigs := []struct {
		name string
		add  []datum
		want []datum
	}{
		{name: "AddNegative", add: []datum{{0, 101}, {0, -1}}, want: []datum{{0, 100}}},
		{name: "AddFromNil", add: []datum{{0, 100}}, want: []datum{{0, 100}, {-1, 0}, {1, 100}}},
		{
			name: "AddAfter",
			add:  []datum{{0, 101}, {100, 101}},
			want: []datum{{0, 101}, {1, 101}, {99, 101}, {100, 202}, {101, 202}},
		},
		{
			name: "AddBefore",
			add:  []datum{{100, 101}, {0, 101}},
			want: []datum{
				{0, 101},
				{1, 101},
				{99, 101},
				{100, 202},
				{101, 202},
			},
		},
		{
			name: "AddBetween",
			add:  []datum{{0, 101}, {100, 101}, {50, 101}},
			want: []datum{
				{0, 101},
				{1, 101},
				{49, 101},
				{50, 202},
				{51, 202},
				{99, 202},
				{100, 303},
				{101, 303},
			},
		},
	}
	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			stepCurve := step.New(
				"entity-id",
				0,
				gcpb.EntityProperty_ENTITY_PROPERTY_UNKNOWN,
				reflect.TypeOf(float64(0)),
			)
			testCurve := New(stepCurve)
			for _, d := range c.add {
				if err := testCurve.Add(d.tick, d.value); err != nil {
					t.Fatalf("Add() = %v, want = %v", err)
				}
			}
			for _, d := range c.want {
				if got := testCurve.Get(d.tick); got != d.value {
					t.Errorf("Get() = %v, want = %v", got, d.value)
				}
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
