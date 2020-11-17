// Package curve defines a server implementation of the parametric property
// curve. This is used by the server to describe the property of an entity as
// an evolving time series instead of a discrete set of points.
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

	// Type indicates the type of the curve itself, e.g. if the curve is
	// a linear interpolation, a delta graph, or else.
	Type() gcpb.CurveType

	// Category indicates the property type of the curve, e.g. if this is
	// should be interpreted as an HP curve, position, or else.
	Category() gcpb.CurveCategory

	// DatumType indicates the data type of the time-series values,
	// e.g. Coordinates, bool, etc.
	DatumType() reflect.Type

	// Tick indicates the last time at which the curve was updated by the
	// server.
	Tick() float64

	// EntityID links back to the specific entity that uses this curve.
	EntityID() string

	// Add takes a value and copies it into the curve.
	Add(t float64, v interface{}) error

	// Get returns a copy of the interal value at a given tick.
	Get(t float64) interface{}

	// ReplaceTail conditionally mutates the last N values of the curve
	// with the values specified in the input, as long as the input curve
	// was updated after the source curve.
	ReplaceTail(c Curve) error

	// ExportTail returns the last N values of the curve as a protobuf,
	// ready to be sent on wire. Setting tick = 0 will export the entire
	// curve.
	ExportTail(tick float64) *gdpb.Curve
}
