package timer

import (
	"reflect"

	"github.com/downflux/game/engine/curve/common/step"
	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/curve/data"
	"github.com/downflux/game/engine/id/id"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

const (
	curveType = gcpb.CurveType_CURVE_TYPE_TIMER
)

var (
	datumType       = reflect.TypeOf(false)
	activateTooSoon = status.Error(codes.FailedPrecondition, "cannot activate timer so soon")
)

type Curve struct {
	*step.Curve

	cooloff id.Tick

	tick id.Tick
	data *data.Data
}

func New(
	eid id.EntityID,
	tick id.Tick,
	cooloff id.Tick,
	property gcpb.EntityProperty) *Curve {
	return &Curve{
		Curve:   step.New(eid, tick, property, datumType),
		cooloff: cooloff,
	}
}

func (c *Curve) Add(tick id.Tick, value interface{}) error {
	if c.Get(tick).(bool) {
		return activateTooSoon
	}

	if i := c.Curve.Data().Search(tick); i < c.Curve.Data().Len() {
		if c.Curve.Data().Tick(i)-tick < c.cooloff {
			return activateTooSoon
		}
	}

	c.Data().Set(tick, true)
	return nil
}

func (c *Curve) Ok(tick id.Tick) bool { return !c.Get(tick).(bool) }

// Get returns the value of the underlying curve. A true value implies the
// timer was recently reset.
func (c *Curve) Get(tick id.Tick) interface{} {
	v := c.Curve.Get(tick).(bool)

	if i := c.Curve.Data().Search(tick); i > 0 {
		if tick-c.Curve.Data().Tick(i-1) >= c.cooloff {
			v = false
		}

	}
	return v
}

func (c *Curve) Export(tick id.Tick) *gdpb.Curve {
	pb := c.Curve.Export(tick)
	pb.Type = curveType
	return pb
}

func (c *Curve) Merge(o curve.Curve) error {
	return status.Error(codes.Unimplemented, "Merge is not implemented for a Timer curve")
}
