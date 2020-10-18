package linearmove

import (
	"reflect"
	"sort"
	"sync"

	"github.com/downflux/game/curve/curve"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

const (
	curveType = gcpb.CurveType_CURVE_TYPE_LINEAR_MOVE
)

var (
	datumType      = reflect.TypeOf(&gdpb.Position{})
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

type datum struct {
	tick float64
	// value should be a clone of the input and considered immutable
	value *gdpb.Position
}

func datumBefore(d1, d2 datum) bool {
	return d1.tick < d2.tick
}

func insert(data []datum, d datum) []datum {
	i := sort.Search(len(data), func(i int) bool { return !datumBefore(data[i], d) })
	// Override existing value if the given input will result in an invalid
	// function.
	if i < len(data) && data[i].tick == d.tick {
		data[i] = d
	} else {
		data = append(data, datum{})
		copy(data[i+1:], data[i:])
		data[i] = d
	}
	return data
}

// TODO(minkezhang): Rename to Curve.
type Curve struct {
	id       string
	entityID string
	data     []datum

	deltaMux sync.Mutex
	delta    []datum
}

func New(id, eid string) *Curve {
	return &Curve{
		id:       id,
		entityID: eid,
	}
}

func (c *Curve) Type() gcpb.CurveType    { return curveType }
func (c *Curve) ID() string              { return c.id }
func (c *Curve) DatumType() reflect.Type { return datumType }
func (c *Curve) EntityID() string        { return c.entityID }

// TODO(minkezhang): Add duplicate removal here / somewhere.
func (c *Curve) Add(t float64, v interface{}) error {
	// TODO(minkezhang): Decide if copying v is necessary here.
	d := datum{tick: t, value: v.(*gdpb.Position)}

	c.data = insert(c.data, d)

	// Add to delta cache for broadcasting.
	c.deltaMux.Lock()
	c.delta = append(c.delta, d)
	c.deltaMux.Unlock()

	// TODO(minkezhang): Add data validation.
	return nil
}

func (c *Curve) Merge(o curve.Curve) error {
	switch o.Type() {
	case gcpb.CurveType_CURVE_TYPE_LINEAR_MOVE:
		for _, d := range o.(*Curve).data {
			c.Add(d.tick, d.value)
		}
	default:
		return notImplemented
	}
	return nil
}

func (c *Curve) Get(t float64) (interface{}, error) {
	if c.data == nil || datumBefore(datum{tick: t}, c.data[0]) {
		return nil, status.Error(codes.OutOfRange, "given tick occurs before the curve existed")
	}

	i := sort.Search(len(c.data), func(i int) bool { return !datumBefore(c.data[i], datum{tick: t}) })

	if i == len(c.data) {
		return proto.Clone(c.data[len(c.data)-1].value).(*gdpb.Position), nil
	}
	if i == 0 {
		return proto.Clone(c.data[0].value).(*gdpb.Position), nil
	}

	return &gdpb.Position{
		X: (c.data[i].value.GetX() + c.data[i-1].value.GetX()) / 2,
		Y: (c.data[i].value.GetY() + c.data[i-1].value.GetY()) / 2,
	}, nil
}

func (c *Curve) ExportDelta() (*gdpb.Curve, error) {
	c.deltaMux.Lock()
	delta := c.delta
	c.delta = nil
	c.deltaMux.Unlock()

	crv := &gdpb.Curve{
		CurveId: c.ID(),
		Type: c.Type(),
		EntityId: c.EntityID(),
	}

	for _, d := range delta {
		crv.Data = append(crv.GetData(), &gdpb.CurveDatum{
			Tick: d.tick,
			Datum: &gdpb.CurveDatum_PositionDatum{d.value},
		})
	}

	return crv, nil
}
