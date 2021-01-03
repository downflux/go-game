package step

import (
	"reflect"
	"testing"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

var (
	_ curve.Curve = &Curve{}
)

func TestGet(t *testing.T) {
	const t0 = 100
	const t1 = 200
	const v0 = float64(101)
	const v1 = float64(201)

	referenceCurve := New(
		"entity-id",
		0,
		gcpb.EntityProperty_ENTITY_PROPERTY_UNKNOWN,
		reflect.TypeOf(v0),
	)
	referenceCurve.Add(t0, v0)
	referenceCurve.Add(t1, v1)

	testConfigs := []struct {
		name string
		c    *Curve
		tick id.Tick
		want float64
	}{
		{
			name: "NoDataGet",
			c: New(
				"entity-id",
				0,
				gcpb.EntityProperty_ENTITY_PROPERTY_UNKNOWN,
				reflect.TypeOf(float64(0)),
			),
			tick: 10,
			want: 0,
		},
		{
			name: "GetBeforeTick",
			c:    referenceCurve,
			tick: t0 - 1,
			want: 0,
		},
		{
			name: "GetAtTick",
			c:    referenceCurve,
			tick: t0,
			want: v0,
		},
		{
			name: "GetBetweenTick",
			c:    referenceCurve,
			tick: t0 + 1,
			want: v0,
		},
		{
			name: "GetAfterTick",
			c:    referenceCurve,
			tick: t1 + 1,
			want: v1,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := c.c.Get(c.tick).(float64); got != c.want {
				t.Errorf("Get() = %v, want = %v", got, c.want)
			}
		})
	}
}
