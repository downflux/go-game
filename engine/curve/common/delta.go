package delta

import (
	"reflect"

	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/id/id"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Curve struct {
	curve.Curve
}

func New(c curve.Curve) *Curve { return &Curve{Curve: c} }

func (c *Curve) Add(t id.Tick, v interface{}) error {
	if c.DatumType() != reflect.TypeOf(float64(0)) {
		return status.Errorf(codes.FailedPrecondition, "cannot add a curve delta for %v type curve", c.DatumType())
	}

	for i := c.Data().Search(t); i < c.Data().Len(); i++ {
		tick := c.Data().Tick(i)
		value := c.Data().Get(tick).(float64)
		c.Data().Set(tick, value+v.(float64))
	}
	return nil
}
