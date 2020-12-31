package data

import (
	"sort"
	"testing"

	"github.com/downflux/game/engine/id/id"
)

var (
	referenceData = New(map[id.Tick]interface{}{
		100: 101,
		200: 201,
		300: 301,
	})
)

func TestTick(t *testing.T) {
	d := &Data{
		ticks: []float64{100, 101, 102},
	}

	for i, tick := range d.ticks {
		want := id.Tick(tick)
		if got := id.Tick(d.Tick(i)); got != want {
			t.Errorf("Tick() = %v, want = %v", got, want)
		}
	}
}

func TestGet(t *testing.T) {
	for tick, want := range referenceData.data {
		if got := referenceData.Get(tick); want != got {
			t.Errorf("Get() = %v, want = %v", got, want)
		}
	}
}

func TestTruncate(t *testing.T) {
	d := referenceData.Clone(referenceData.Tick(0))
	d.Truncate(d.Tick(1))

	if got := d.Len(); got != 1 {
		t.Fatalf("Len() = %v, want = %v", got, 1)
	}

	want := referenceData.Get(referenceData.Tick(0))
	if got := d.Get(d.Tick(0)); got != want {
		t.Fatalf("Get() = %v, want = %v", got, want)
	}

	if got := d.Get(referenceData.Tick(1)); got != nil {
		t.Errorf("Get() = %v, want = %v", got, nil)
	}
}

func TestSearch(t *testing.T) {
	testConfigs := []struct {
		name string
		tick id.Tick
		want int
	}{
		{name: "SearchExists", tick: 100, want: 0},
		{name: "SearchBig", tick: 301, want: referenceData.Len()},
		{name: "SearchBetween", tick: 101, want: 1},
		{name: "SearchSmall", tick: 99, want: 0},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := referenceData.Search(c.tick); c.want != got {
				t.Fatalf("Search() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestSearchAfter(t *testing.T) {
	testConfigs := []struct {
		name string
		tick id.Tick
		want int
	}{
		{name: "SearchExists", tick: 100, want: 0},
		{name: "SearchBig", tick: 301, want: referenceData.Len()},
		{name: "SearchBetween", tick: 101, want: 1},
		{name: "SearchSmall", tick: 99, want: 0},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := sort.Search(
				referenceData.Len(),
				func(i int) bool { return !Before(referenceData.Tick(i), c.tick) },
			); c.want != got {
				t.Errorf("Search() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestClone(t *testing.T) {
	testConfigs := []struct {
		name     string
		original *Data
		tick     id.Tick
	}{
		{name: "CompleteClone", original: referenceData, tick: 0},
		{name: "PartialClone", original: referenceData, tick: 101},
		{name: "NullClone", original: New(nil), tick: 500},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			copied := c.original.Clone(c.tick)
			for i := 0; i < c.original.Len(); i++ {
				if !Before(c.original.Tick(i), c.tick) {
					tick := c.original.Tick(i)
					want := c.original.Get(tick)
					if got := copied.Get(tick); got != want {
						t.Errorf("Get() = %v, want = %v", got, want)
					}
				}
			}
		})
	}
}

func TestSet(t *testing.T) {
	testConfigs := []struct {
		name  string
		d     *Data
		tick  id.Tick
		value int
	}{
		{name: "SetEmpty", d: New(nil), tick: 100, value: 101},
		{
			name: "SetOverwrite",
			d: New(map[id.Tick]interface{}{
				100: 100,
			}),
			tick:  100,
			value: 101,
		},
		{
			name: "SetBefore",
			d: New(map[id.Tick]interface{}{
				100: 101,
			}),
			tick:  99,
			value: 199,
		},
		{
			name: "SetAfter",
			d: New(map[id.Tick]interface{}{
				100: 101,
			}),
			tick:  200,
			value: 201,
		},
		{
			name: "SetBetween",
			d: New(map[id.Tick]interface{}{
				100: 101,
				300: 301,
			}),
			tick:  200,
			value: 201,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			c.d.Set(c.tick, c.value)
			if got := c.d.Get(c.tick); c.value != got {
				t.Fatalf("Get() = %v, want = %v", got, c.value)
			}

			if isSorted := sort.IsSorted(c.d.ticks); !isSorted {
				t.Errorf("IsSorted() = %v, want = %v", isSorted, true)
			}
		})
	}
}
