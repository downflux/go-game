// Package tank encapsulates logic for a basic tank unit.
package tank

import (
	"github.com/downflux/game/engine/curve/common/linearmove"
	"github.com/downflux/game/engine/curve/curve"
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

// Tank implements the entity.Entity interface and represents a simple armored
// unit.
type Tank struct {
	entity.LifeCycle

	// eid is a UUID of the Entity.
	eid id.EntityID

	// curves is a list of Curves tracking the Entity properties.
	curves map[gcpb.EntityProperty]curve.Curve
}

// New constructs a new instance of the Tank.
func New(eid id.EntityID, t id.Tick, p *gdpb.Position) *Tank {
	mc := linearmove.New(eid, t)
	mc.Add(t, p)

	return &Tank{
		eid: eid,
		curves: map[gcpb.EntityProperty]curve.Curve{
			gcpb.EntityProperty_ENTITY_PROPERTY_POSITION: mc,
		},
	}
}

// ID returns the UUID of the Tank.
func (e *Tank) ID() id.EntityID { return e.eid }

// Properties returns the list of registered properties tracked by the
// Tank instance.
func (e *Tank) Properties() []gcpb.EntityProperty {
	return []gcpb.EntityProperty{gcpb.EntityProperty_ENTITY_PROPERTY_POSITION}
}

// Curve returns the Curve instance for a specific EntityProperty.
func (e *Tank) Curve(t gcpb.EntityProperty) curve.Curve { return e.curves[t] }

// Type returns the registered EntityType.
func (e *Tank) Type() gcpb.EntityType { return gcpb.EntityType_ENTITY_TYPE_TANK }
