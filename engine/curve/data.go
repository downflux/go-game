package data

import (
	"sort"

	"github.com/downflux/game/engine/id/id"
)

func Before(t id.Tick, u id.Tick) bool { return t < u }

type Data struct {
	ticks sort.Float64Slice
	data  map[id.Tick]interface{}
}

func New(data map[id.Tick]interface{}) *Data {
	d := &Data{}
	for tick, value := range data {
		d.Set(tick, value)
	}
	return d
}

func (d *Data) Len() int                     { return d.ticks.Len() }
func (d *Data) Search(tick id.Tick) int      { return d.ticks.Search(tick.Value()) }
func (d *Data) Tick(i int) id.Tick           { return id.Tick(d.ticks[i]) }
func (d *Data) Get(tick id.Tick) interface{} { return d.data[tick] }

func (d *Data) Set(tick id.Tick, value interface{}) {
	i := d.ticks.Search(tick.Value())

	if d.data == nil {
		d.data = map[id.Tick]interface{}{}
	}

	// Overwrite existing value if there is a value already stored at the
	// given tick.
	d.data[tick] = value

	// Add new tick to ordered tick lookup table.
	if i == d.ticks.Len() || tick != d.Tick(i) {
		d.ticks = append(d.ticks, 0)
		copy(d.ticks[i+1:], d.ticks[i:])
		d.ticks[i] = tick.Value()
	}
}

func (d *Data) Truncate(tick id.Tick) {
	i := d.ticks.Search(tick.Value())
	for j := i; j < d.Len(); j++ {
		delete(d.data, d.Tick(j))
	}
	d.ticks = d.ticks[:i]
}

func (d *Data) Clone(tick id.Tick) *Data {
	clone := New(nil)

	for i := d.ticks.Search(tick.Value()); i < d.Len(); i++ {
		tick := d.Tick(i)
		// TODO(minkezhang): Ensure value is either a primitive or
		// should be deep copied.
		clone.Set(tick, d.Get(tick))
	}

	return clone
}

func (d *Data) Merge(other *Data) error {
	d.Truncate(other.Tick(0))

	for i := 0; i < other.Len(); i++ {
		tick := other.Tick(i)
		// TODO(minkezhang): Check for memory leaks from curve o.
		d.Set(tick, other.Get(tick))
	}
	return nil
}
