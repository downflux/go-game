package timer

import (
	"reflect"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/curve/data"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

const (
	curveType = gcpb.CurveType_CURVE_TYPE_TIMER
)

var (
	datumType = reflect.TypeOf(false)
)

type Curve struct {
	curve.Base

	interval id.Tick

	tick id.Tick
	data *data.Data
}

func New(
	eid id.EntityID,
	tick id.Tick,
	interval id.Tick,
	property gcpb.EntityProperty,
	datumType reflect.Type) *Curve {
	return &Curve{
		Base:     *curve.New(eid, curveType, datumType, property),
		interval: interval,
		tick:     tick,
		data:     data.New(nil),
	}
}

func (c *Curve) DatumType() reflect.Type { return datumType }
func (c *Curve) Data() *data.Data        { return c.data }
func (c *Curve) Tick() id.Tick           { return c.tick }

func (c *Curve) Add(tick id.Tick, value interface{}) error {
	i := c.Data().Search(tick)

	if i < c.Data().Len() {
		if tick-c.Data().Tick(i) > c.interval {
			return nil
		}
	}
	if i < c.Data().Len()-1 {
		if c.Data().Tick(i+1)-tick > c.interval {
			return nil
		}
	}

	c.Data().Set(tick, true)
	return nil
}
