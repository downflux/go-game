// Package tank encapsulates logic for a basic tank unit.
package tank

import (
	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/curve/linearmove"
	"github.com/downflux/game/server/entity/entity"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/visitor/visitor"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

// Tank implements the visitor.Entity interface and represents a simple armored
// unit.
type Tank struct {
	entity.BaseEntity

	// eid is a UUID of the Entity.
	eid id.EntityID

	// curveLookup is a list of Curves tracking the Entity properties.
	curveLookup map[gcpb.CurveCategory]curve.Curve
}

// New constructs a new instance of the Tank.
func New(eid id.EntityID, t id.Tick, p *gdpb.Position) *Tank {
	mc := linearmove.New(eid, t)
	mc.Add(t, p)

	return &Tank{
		eid: eid,
		curveLookup: map[gcpb.CurveCategory]curve.Curve{
			gcpb.CurveCategory_CURVE_CATEGORY_MOVE: mc,
		},
	}
}

// ID returns the UUID of the Tank.
func (e *Tank) ID() id.EntityID { return e.eid }

// CurveCategories returns the list of registered properties tracked by the
// Tank instance.
func (e *Tank) CurveCategories() []gcpb.CurveCategory {
	return []gcpb.CurveCategory{gcpb.CurveCategory_CURVE_CATEGORY_MOVE}
}

// Curve returns the Curve instance for a specific CurveCategory.
func (e *Tank) Curve(t gcpb.CurveCategory) curve.Curve { return e.curveLookup[t] }

// Type returns the registered EntityType.
func (e *Tank) Type() gcpb.EntityType { return gcpb.EntityType_ENTITY_TYPE_TANK }

// Accept allows the input Visitor instance to mutate the internal state of the
// Tank.
func (e *Tank) Accept(v visitor.Visitor) error { return v.Visit(e) }
