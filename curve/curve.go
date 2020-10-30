package curve

import (
	"reflect"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

// Curve represents the evolution of a specific data metric over time, e.g. HP,
// unit orientation, position, etc.
//
// TODO(minkezhang): Translate to generics instead. See
// https://blog.golang.org/generics-next-step.
//
// type Curve[T any] interface {
//   Get(t float64) T
// }
//
// We can remove the DatumType function once this is implemented.
type Curve interface {
	// TODO(minkezhang): Consider if we need Start(), End() tick values.

	Type() gcpb.CurveType
	Category() gcpb.CurveCategory
	DatumType() reflect.Type
	Tick() float64

	EntityID() string

	// Add takes a value and copies it into the curve.
	Add(t float64, v interface{}) error

	// Get returns a copy of the interal value at a given tick.
	Get(t float64) interface{}

	// Mutate, not copy constructor.
	ReplaceTail(c Curve) error

	ExportTail(tick float64) *gdpb.Curve
}
