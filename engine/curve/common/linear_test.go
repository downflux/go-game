package linearmove

import (
	"testing"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/curve/data"
	"github.com/downflux/game/engine/id/id"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

var (
	_ curve.Curve = &Curve{}
)

func newCurve(eid id.EntityID, t id.Tick) *Curve {
	return New(eid, t, gcpb.EntityProperty_ENTITY_PROPERTY_POSITION)
}

func TestMerge(t *testing.T) {
	replaceC1 := newCurve("eid", 0)
	replaceC1.Add(0, float64(0))
	replaceC1.Add(1, float64(10))
	replaceC1.Add(2, float64(20))
	replaceC2 := newCurve("eid", 1)
	replaceC2.Add(1, float64(1))

	replaceSameTickC1 := newCurve("eid", 0)
	replaceSameTickC1.Add(0, float64(0))
	replaceSameTickC1.Add(1, float64(10))
	replaceSameTickC1.Add(2, float64(20))
	replaceSameTickC2 := newCurve("eid", 0)
	replaceSameTickC2.Add(1, float64(1))

	tooStaleC1 := newCurve("eid", 1)
	tooStaleC1.Add(0, float64(0))
	tooStaleC1.Add(1, float64(10))
	tooStaleC1.Add(2, float64(20))
	tooStaleC2 := newCurve("eid", 0)
	tooStaleC2.Add(1, float64(1))

	updateTickC1 := newCurve("eid", 0)
	updateTickC2 := newCurve("eid", 2)
	updateTickC2.Add(0, float64(0))
	updateTickC2.Add(1, float64(10))
	updateTickC1.Merge(updateTickC2)
	updateTickC3 := newCurve("eid", 1)
	updateTickC3.Add(2, float64(20))

	testConfigs := []struct {
		name string
		c1   *Curve
		c2   *Curve
		tick id.Tick
		want float64
	}{
		{
			name: "MergeNormal",
			c1:   replaceC1,
			c2:   replaceC2,
			tick: 0.7,
			want: 0.7,
		},
		{
			name: "MergeSameTick",
			c1:   replaceSameTickC1,
			c2:   replaceSameTickC2,
			tick: 0.7,
			want: 0.7,
		},
		{
			name: "ReplaceTooStale",
			c1:   tooStaleC1,
			c2:   tooStaleC2,
			tick: 0.7,
			want: 7,
		},
		{
			name: "ReplaceUpdateTick",
			c1:   updateTickC1,
			c2:   updateTickC3,
			tick: 2,
			want: 10,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			c.c1.Merge(c.c2)
			if got := c.c1.Get(c.tick).(float64); got != c.want {
				t.Errorf("Get() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestExport(t *testing.T) {
	const eid = "eid"
	cSimple := newCurve(eid, 0)
	cSimple.Add(0, float64(0))
	cSimple.Add(1, float64(1))
	cSimple.Add(2, float64(2))

	testConfigs := []struct {
		name string
		c    *Curve
		t    id.Tick
		want *gdpb.Curve
	}{
		{
			name: "ExportSimple",
			c:    cSimple,
			t:    0,
			want: &gdpb.Curve{
				EntityId: eid,
				Tick:     cSimple.Tick().Value(),
				Property: cSimple.Property(),
				Type:     cSimple.Type(),
				Data: []*gdpb.CurveDatum{
					{
						Tick:  0,
						Datum: &gdpb.CurveDatum_DoubleDatum{cSimple.Get(0).(float64)},
					},
					{
						Tick:  1,
						Datum: &gdpb.CurveDatum_DoubleDatum{cSimple.Get(1).(float64)},
					},
					{
						Tick:  2,
						Datum: &gdpb.CurveDatum_DoubleDatum{cSimple.Get(2).(float64)},
					},
				},
			},
		},
		{
			name: "ExportPartial",
			c:    cSimple,
			t:    1,
			want: &gdpb.Curve{
				EntityId: eid,
				Tick:     cSimple.Tick().Value(),
				Property: cSimple.Property(),
				Type:     cSimple.Type(),
				Data: []*gdpb.CurveDatum{
					{
						Tick:  1,
						Datum: &gdpb.CurveDatum_DoubleDatum{cSimple.Get(1).(float64)},
					},
					{
						Tick:  2,
						Datum: &gdpb.CurveDatum_DoubleDatum{cSimple.Get(2).(float64)},
					},
				},
			},
		},
		{
			name: "ExportOffsetIndex",
			c:    cSimple,
			t:    1.1,
			want: &gdpb.Curve{
				EntityId: eid,
				Tick:     cSimple.Tick().Value(),
				Property: cSimple.Property(),
				Type:     cSimple.Type(),
				Data: []*gdpb.CurveDatum{
					{
						Tick:  1,
						Datum: &gdpb.CurveDatum_DoubleDatum{cSimple.Get(1).(float64)},
					},
					{
						Tick:  2,
						Datum: &gdpb.CurveDatum_DoubleDatum{cSimple.Get(2).(float64)},
					},
				},
			},
		},
		{
			name: "ExportPastLastDataPoint",
			c:    cSimple,
			t:    2.1,
			want: &gdpb.Curve{
				EntityId: eid,
				Tick:     cSimple.Tick().Value(),
				Property: cSimple.Property(),
				Type:     cSimple.Type(),
				Data: []*gdpb.CurveDatum{
					{
						Tick:  2,
						Datum: &gdpb.CurveDatum_DoubleDatum{cSimple.Get(2).(float64)},
					},
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			got := c.c.Export(c.t)
			if diff := cmp.Diff(got, c.want, protocmp.Transform()); diff != "" {
				t.Errorf("Export() mismatch (-want, +got):\n%v", diff)
			}
		})
	}
}

func TestGet(t *testing.T) {
	testConfigs := []struct {
		name string
		c    *Curve
		t    id.Tick
		want float64
	}{
		{
			name: "GetNull",
			c:    &Curve{},
			t:    1,
			want: 0,
		},
		{
			name: "GetBeforeCreation",
			c: &Curve{
				data: data.New(map[id.Tick]interface{}{
					1: float64(1),
				})},
			t:    0,
			want: 1,
		},
		{
			name: "GetAlreadyKnown",
			c: &Curve{
				data: data.New(map[id.Tick]interface{}{
					1: float64(1),
				})},
			t:    1,
			want: 1,
		},
		{
			name: "GetAfterLastKnown",
			c: &Curve{
				data: data.New(map[id.Tick]interface{}{
					0: float64(1),
				})},
			t:    1,
			want: 1,
		},
		{
			name: "GetInterpolatedValue",
			c: &Curve{
				data: data.New(map[id.Tick]interface{}{
					0: float64(0),
					1: float64(1),
				})},
			t:    0.7,
			want: 0.7,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := c.c.Get(c.t).(float64); got != c.want {
				t.Fatalf("Get() = %v, want = %v", got, c.want)
			}
		})
	}
}
