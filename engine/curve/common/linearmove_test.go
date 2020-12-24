package linearmove

import (
	"testing"

	"github.com/downflux/game/engine/id/id"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

func TestDatumBefore(t *testing.T) {
	testConfigs := []struct {
		name   string
		d1, d2 datum
		want   bool
	}{
		{name: "TrivialCompareBefore", d1: datum{tick: 0}, d2: datum{tick: 1}, want: true},
		{name: "TrivialCompareAfter", d1: datum{tick: 1}, d2: datum{tick: 0}, want: false},
		{name: "TrivialCompareNotBefore", d1: datum{tick: 0}, d2: datum{tick: 0}, want: false},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := datumBefore(c.d1, c.d2); got != c.want {
				t.Errorf("datumBefore() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	testConfigs := []struct {
		name string
		data []datum
		d    datum
		want []datum
	}{
		{name: "TrivialInsert", data: nil, d: datum{tick: 1}, want: []datum{{tick: 1}}},
		{name: "InsertBefore", data: []datum{{tick: 1}}, d: datum{tick: 0}, want: []datum{{tick: 0}, {tick: 1}}},
		{name: "InsertAfter", data: []datum{{tick: 0}}, d: datum{tick: 1}, want: []datum{{tick: 0}, {tick: 1}}},
		{
			name: "InsertOverride",
			data: []datum{{tick: 0, value: &gdpb.Position{X: 0, Y: 0}}},
			d:    datum{tick: 0, value: &gdpb.Position{X: 1, Y: 1}},
			want: []datum{{tick: 0, value: &gdpb.Position{X: 1, Y: 1}}},
		},
		{name: "InsertBetween", data: []datum{{tick: 0}, {tick: 2}}, d: datum{tick: 1}, want: []datum{{tick: 0}, {tick: 1}, {tick: 2}}},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			got := insert(c.data, c.d)
			if diff := cmp.Diff(got, c.want, cmp.AllowUnexported(datum{}), protocmp.Transform()); diff != "" {
				t.Errorf("insert() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestReplaceTail(t *testing.T) {
	replaceC1 := New("eid", 0)
	replaceC1.Add(0, &gdpb.Position{X: 0, Y: 0})
	replaceC1.Add(1, &gdpb.Position{X: 10, Y: 10})
	replaceC1.Add(2, &gdpb.Position{X: 20, Y: 20})
	replaceC2 := New("eid", 1)
	replaceC2.Add(1, &gdpb.Position{X: 1, Y: 1})

	replaceSameTickC1 := New("eid", 0)
	replaceSameTickC1.Add(0, &gdpb.Position{X: 0, Y: 0})
	replaceSameTickC1.Add(1, &gdpb.Position{X: 10, Y: 10})
	replaceSameTickC1.Add(2, &gdpb.Position{X: 20, Y: 20})
	replaceSameTickC2 := New("eid", 0)
	replaceSameTickC2.Add(1, &gdpb.Position{X: 1, Y: 1})

	tooStaleC1 := New("eid", 1)
	tooStaleC1.Add(0, &gdpb.Position{X: 0, Y: 0})
	tooStaleC1.Add(1, &gdpb.Position{X: 10, Y: 10})
	tooStaleC1.Add(2, &gdpb.Position{X: 20, Y: 20})
	tooStaleC2 := New("eid", 0)
	tooStaleC2.Add(1, &gdpb.Position{X: 1, Y: 1})

	updateTickC1 := New("eid", 0)
	updateTickC2 := New("eid", 2)
	updateTickC2.Add(0, &gdpb.Position{X: 0, Y: 0})
	updateTickC2.Add(1, &gdpb.Position{X: 10, Y: 10})
	updateTickC1.ReplaceTail(updateTickC2)
	updateTickC3 := New("eid", 1)
	updateTickC3.Add(2, &gdpb.Position{X: 20, Y: 20})

	testConfigs := []struct {
		name string
		c1   *Curve
		c2   *Curve
		tick id.Tick
		want *gdpb.Position
	}{
		{
			name: "ReplaceTailNormal",
			c1:   replaceC1,
			c2:   replaceC2,
			tick: 0.7,
			want: &gdpb.Position{X: 0.7, Y: 0.7},
		},
		{
			name: "ReplaceTailSameTick",
			c1:   replaceSameTickC1,
			c2:   replaceSameTickC2,
			tick: 0.7,
			want: &gdpb.Position{X: 0.7, Y: 0.7},
		},
		{
			name: "ReplaceTooStale",
			c1:   tooStaleC1,
			c2:   tooStaleC2,
			tick: 0.7,
			want: &gdpb.Position{X: 7, Y: 7},
		},
		{
			name: "ReplaceUpdateTick",
			c1:   updateTickC1,
			c2:   updateTickC3,
			tick: 2,
			want: &gdpb.Position{X: 10, Y: 10},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			c.c1.ReplaceTail(c.c2)
			got := c.c1.Get(c.tick)
			if diff := cmp.Diff(got, c.want, protocmp.Transform()); diff != "" {
				t.Errorf("Get() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestExportTail(t *testing.T) {
	const eid = "eid"
	cSimple := New(eid, 0)
	cSimple.Add(0, &gdpb.Position{X: 0, Y: 0})
	cSimple.Add(1, &gdpb.Position{X: 1, Y: 1})
	cSimple.Add(2, &gdpb.Position{X: 2, Y: 2})

	testConfigs := []struct {
		name string
		c    *Curve
		t    id.Tick
		want *gdpb.Curve
	}{
		{
			name: "ExportTailSimple",
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
						Datum: &gdpb.CurveDatum_PositionDatum{cSimple.Get(0).(*gdpb.Position)},
					},
					{
						Tick:  1,
						Datum: &gdpb.CurveDatum_PositionDatum{cSimple.Get(1).(*gdpb.Position)},
					},
					{
						Tick:  2,
						Datum: &gdpb.CurveDatum_PositionDatum{cSimple.Get(2).(*gdpb.Position)},
					},
				},
			},
		},
		{
			name: "ExportTailPartial",
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
						Datum: &gdpb.CurveDatum_PositionDatum{cSimple.Get(1).(*gdpb.Position)},
					},
					{
						Tick:  2,
						Datum: &gdpb.CurveDatum_PositionDatum{cSimple.Get(2).(*gdpb.Position)},
					},
				},
			},
		},
		{
			name: "ExportTailOffsetIndex",
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
						Datum: &gdpb.CurveDatum_PositionDatum{cSimple.Get(1).(*gdpb.Position)},
					},
					{
						Tick:  2,
						Datum: &gdpb.CurveDatum_PositionDatum{cSimple.Get(2).(*gdpb.Position)},
					},
				},
			},
		},
		{
			name: "ExportTailPastLastDataPoint",
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
						Datum: &gdpb.CurveDatum_PositionDatum{cSimple.Get(2).(*gdpb.Position)},
					},
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			got := c.c.ExportTail(c.t)
			if diff := cmp.Diff(got, c.want, protocmp.Transform()); diff != "" {
				t.Errorf("ExportTail() mismatch (-want, +got):\n%v", diff)
			}
		})
	}
}

func TestGet(t *testing.T) {
	testConfigs := []struct {
		name string
		c    *Curve
		t    id.Tick
		want *gdpb.Position
	}{
		{
			name: "GetNull",
			c:    &Curve{},
			t:    1,
			want: &gdpb.Position{},
		},
		{
			name: "GetBeforeCreation",
			c:    &Curve{data: []datum{{tick: 1, value: &gdpb.Position{X: 1, Y: 1}}}},
			t:    0,
			want: &gdpb.Position{X: 1, Y: 1},
		},
		{
			name: "GetAlreadyKnown",
			c:    &Curve{data: []datum{{tick: 1, value: &gdpb.Position{X: 1, Y: 1}}}},
			t:    1,
			want: &gdpb.Position{X: 1, Y: 1},
		},
		{
			name: "GetAfterLastKnown",
			c:    &Curve{data: []datum{{tick: 0, value: &gdpb.Position{X: 1, Y: 1}}}},
			t:    1,
			want: &gdpb.Position{X: 1, Y: 1},
		},
		{
			name: "GetInterpolatedValue",
			c: &Curve{data: []datum{
				{tick: 0, value: &gdpb.Position{X: 0, Y: 0}},
				{tick: 1, value: &gdpb.Position{X: 1, Y: 1}},
			}},
			t:    0.7,
			want: &gdpb.Position{X: 0.7, Y: 0.7},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := c.c.Get(c.t); !proto.Equal(got.(*gdpb.Position), c.want) {
				t.Fatalf("Get() = %v, want = %v", got, c.want)
			}
		})
	}
}
