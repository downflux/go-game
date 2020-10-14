package linearmove

import (
	"reflect"
	"sort"

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
	datumType = reflect.TypeOf(&gdpb.Position{})
)

type datum struct {
	tick  float64
	value *gdpb.Position // clone, immutable
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

type LinearMoveCurve struct {
	id       string
	entityID string
	data     []datum
}

func New(id, eid string) *LinearMoveCurve {
	return &LinearMoveCurve{
		id:       id,
		entityID: eid,
	}
}

func (c *LinearMoveCurve) Type() gcpb.CurveType    { return curveType }
func (c *LinearMoveCurve) ID() string              { return c.id }
func (c *LinearMoveCurve) DatumType() reflect.Type { return datumType }
func (c *LinearMoveCurve) EntityID() string        { return c.entityID }

func (c *LinearMoveCurve) Add(t float64, v interface{}) error {
	c.data = insert(c.data, datum{tick: t, value: v.(*gdpb.Position)}) // copy
	// TODO(minkezhang): Add data validation.
	return nil
}

func (c *LinearMoveCurve) Get(t float64) (interface{}, error) {
	if c.data == nil || datumBefore(datum{tick: t}, c.data[0]) {
		return nil, status.Error(codes.OutOfRange, "given tick occurs before the curve existed")
	}

	i := sort.Search(len(c.data), func(i int) bool { return !datumBefore(c.data[i], datum{tick: t}) })

	if i == len(c.data) {
		return proto.Clone(c.data[len(c.data)-1].value).(*gdpb.Position), nil
	}
	if i == 0 {
		return proto.Clone(c.data[0].value).(*gdpb.Position), nil // proto.Clone
	}

	return &gdpb.Position{
		X: (c.data[i].value.GetX() + c.data[i-1].value.GetX()) / 2,
		Y: (c.data[i].value.GetY() + c.data[i-1].value.GetY()) / 2,
	}, nil
}
