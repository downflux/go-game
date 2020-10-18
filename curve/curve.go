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
	ID() string
	DatumType() reflect.Type
	EntityID() string

	// Add takes a value and copies it into the curve.
	Add(t float64, v interface{}) error

	// Get returns a copy of the interal value at a given tick.
	Get(t float64) (interface{}, error)

	// Mutate, not copy constructor.
	// TODO(minkezhang): Make this replace, not merge.
	Merge(c Curve) error

	ExportDelta() (*gdpb.Curve, error)
	// TODO(minkezhang): Implement the following.
	/**
	 * Export() (*gdpb.Curve, error)
	 * Extract(t1, t2 float32) Curve // Same ID
	 * Contains(c Curve) bool
	 * Hash() string // placeholder type
	 */
}
