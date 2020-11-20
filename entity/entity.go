package entity

import (
	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/curve/linearmove"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

const idLen = 8

// TODO(minkezhang): Migrate to typed string.
type EntityID string

type Entity interface {
	ID() string
	Type() gcpb.EntityType
	Curve(t gcpb.CurveCategory) curve.Curve

	// CurveCategories returns list of curve categories defined in a specific
	// entity. This list is created at init time and is immutable.
	CurveCategories() []gcpb.CurveCategory

	// TODO(minkezhang): Implement these methods.
	/*
	 * Start() float64
	 * End() float64
	 * Delete()
	 */
}

type SimpleEntity struct {
	id          string
	curveLookup map[gcpb.CurveCategory]curve.Curve
}

// TODO(minkezhang): Make this client-friendly too.
func NewSimpleEntity(eid string, t float64, p *gdpb.Position) *SimpleEntity {
	mc := linearmove.New(eid, t)
	mc.Add(t, p)

	return &SimpleEntity{
		id: eid,
		curveLookup: map[gcpb.CurveCategory]curve.Curve{
			gcpb.CurveCategory_CURVE_CATEGORY_MOVE: mc,
		},
	}
}

func (e *SimpleEntity) ID() string            { return e.id }
func (e *SimpleEntity) Type() gcpb.EntityType { return gcpb.EntityType_ENTITY_TYPE_TANK }
func (e *SimpleEntity) CurveCategories() []gcpb.CurveCategory {
	return []gcpb.CurveCategory{gcpb.CurveCategory_CURVE_CATEGORY_MOVE}
}

// TODO(minkezhang): Decide if we should return default value.
func (e *SimpleEntity) Curve(t gcpb.CurveCategory) curve.Curve { return e.curveLookup[t] }
