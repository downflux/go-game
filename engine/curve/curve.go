// Package curve defines a server implementation of the parametric property
// curve. This is used by the server to describe the property of an entity as
// an evolving time series instead of a discrete set of points.
package curve

import (
	"reflect"

	"github.com/downflux/game/engine/curve/data"
	"github.com/downflux/game/engine/id/id"

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
//   Get(t Tick) T
// }
//
// We can remove the DatumType function once this is implemented.
type Curve interface {

	// EntityID links back to the specific entity that uses this curve.
	EntityID() id.EntityID

	// Type indicates the type of the curve itself, e.g. if the curve is
	// a linear interpolation, a delta graph, or else.
	Type() gcpb.CurveType

	// Property indicates the property type of the curve, e.g. if
	// this should be interpreted as an HP curve, position, etc.
	Property() gcpb.EntityProperty

	// DatumType indicates the data type of the time-series values,
	// e.g. Coordinates, bool, etc.
	DatumType() reflect.Type

	Data() *data.Data

	// Tick indicates the last time at which the curve was updated by the
	// server.
	Tick() id.Tick

	// Add takes a value and copies it into the curve.
	Add(t id.Tick, v interface{}) error

	// Get returns a copy of the interal value at a given tick.
	Get(t id.Tick) interface{}

	// Merge conditionally mutates the last N values of the curve
	// with the values specified in the input, as long as the input curve
	// was updated after the source curve.
	Merge(c Curve) error

	// Export returns the last N values of the curve as a protobuf,
	// ready to be sent on wire. Setting tick = 0 will export the entire
	// curve.
	Export(t id.Tick) *gdpb.Curve
}

type Base struct {
	eid       id.EntityID         // Read-only.
	curveType gcpb.CurveType      // Read-only.
	property  gcpb.EntityProperty // Read-only.
	datumType reflect.Type        // Read-only.
}

func New(
	eid id.EntityID,
	curveType gcpb.CurveType,
	datumType reflect.Type,
	property gcpb.EntityProperty) *Base {
	return &Base{
		eid:       eid,
		curveType: curveType,
		property:  property,
		datumType: datumType,
	}
}

func (c Base) Type() gcpb.CurveType          { return c.curveType }
func (c Base) Property() gcpb.EntityProperty { return c.property }
func (c Base) DatumType() reflect.Type       { return c.datumType }
func (c Base) EntityID() id.EntityID         { return c.eid }

func (c Base) Export(t id.Tick) *gdpb.Curve {
	return &gdpb.Curve{
		Type:     c.Type(),
		Property: c.Property(),
		EntityId: c.EntityID().Value(),
		Tick:     t.Value(),
	}
}
