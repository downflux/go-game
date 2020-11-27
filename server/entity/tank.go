package tank

import (
	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/curve/linearmove"
	"github.com/downflux/game/server/service/visitor/entity/entity"
	"github.com/downflux/game/server/service/visitor/visitor"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

type Tank struct {
	entity.BaseEntity

	id          string
	curveLookup map[gcpb.CurveCategory]curve.Curve
}

func New(eid string, t float64, p *gdpb.Position) *Tank {
	mc := linearmove.New(eid, t)
	mc.Add(t, p)

	return &Tank{
		id: eid,
		curveLookup: map[gcpb.CurveCategory]curve.Curve{
			gcpb.CurveCategory_CURVE_CATEGORY_MOVE: mc,
		},
	}
}

func (e *Tank) ID() string { return e.id }
func (e *Tank) CurveCategories() []gcpb.CurveCategory {
	return []gcpb.CurveCategory{gcpb.CurveCategory_CURVE_CATEGORY_MOVE}
}

// TODO(minkezhang): Decide if we should return default value.
func (e *Tank) Curve(t gcpb.CurveCategory) curve.Curve { return e.curveLookup[t] }

func (e *Tank) Type() gcpb.EntityType          { return gcpb.EntityType_ENTITY_TYPE_TANK }
func (e *Tank) Accept(v visitor.Visitor) error { return v.Visit(e) }
