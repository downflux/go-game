package linearmove

import (
	"testing"

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

func TestGetError(t *testing.T) {
	c := &Curve{}
	if got, err := c.Get(0); err == nil {
		t.Errorf("Get() = %v, nil, want a non-nil error", got)
	}
}

func TestReplaceTail(t *testing.T) {
	c1 := New("c1", "eid")
	c1.Add(0, &gdpb.Position{X: 0, Y: 0})
	c1.Add(1, &gdpb.Position{X: 10, Y: 10})
	c1.Add(2, &gdpb.Position{X: 20, Y: 20})

	c2 := New("", "eid")
	c2.Add(1, &gdpb.Position{X: 1, Y: 1})

	c1.ReplaceTail(c2)

	want := &gdpb.Position{X: 0.5, Y: 0.5}
	got, err := c1.Get(0.5)
	if err != nil {
		t.Fatalf("Get() = _, %v, want = nil", err)
	}

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Get() mismatch (-want +got):\n%v", diff)
	}
}

func TestGet(t *testing.T) {
	testConfigs := []struct {
		name string
		c    *Curve
		t    float64
		want *gdpb.Position
	}{
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
			t:    0.5,
			want: &gdpb.Position{X: 0.5, Y: 0.5},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			got, err := c.c.Get(c.t)
			if err != nil {
				t.Fatalf("Get() = _, %v, want = _, nil", err)
			}
			if !proto.Equal(got.(*gdpb.Position), c.want) {
				t.Fatalf("Get() = %v, want = %v", got, c.want)
			}
		})
	}
}
