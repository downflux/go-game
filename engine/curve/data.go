package data

import (
	"sort"

	"github.com/downflux/game/engine/id/id"
)

type Data struct {
	ticks sort.Float64Slice
	data  map[id.Tick]interface{}
}

func (d *Data) Search(tick id.Tick) int      { return d.ticks.Search(tick.Value()) }
func (d *Data) Tick(i int) id.Tick           { return id.Tick(d.ticks[i]) }
func (d *Data) Get(tick id.Tick) interface{} { return d.data[tick] }

func (d *Data) Insert(tick id.Tick, value interface{}) {
	i := d.ticks.Search(tick.Value())

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
