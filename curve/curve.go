package curve

import (
	"reflect"

	gcpb "github.com/downflux/game/api/constants_go_proto"
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
	Type() gcpb.CurveType
	ID() string
	DatumType() reflect.Type
	EntityID() string

	// Add takes a value and copies it into the curve.
	Add(t float64, v interface{}) error

	// Get returns a copy of the interal value at a given tick.
	Get(t float64) (interface{}, error)

	// TODO(minkezhang): Implement the following.
	/**
	 * Merge(c Curve) error
	 * Extract(t1, t2 float32) Curve // Same ID
	 * Contains(c Curve) bool
	 * Hash() bool // placeholder type
	 */
}
