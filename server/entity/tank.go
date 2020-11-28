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

type Tank struct {
	entity.BaseEntity

	eid         id.EntityID
	curveLookup map[gcpb.CurveCategory]curve.Curve
}

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

func (e *Tank) ID() id.EntityID { return e.eid }
func (e *Tank) CurveCategories() []gcpb.CurveCategory {
	return []gcpb.CurveCategory{gcpb.CurveCategory_CURVE_CATEGORY_MOVE}
}

// TODO(minkezhang): Decide if we should return default value.
func (e *Tank) Curve(t gcpb.CurveCategory) curve.Curve { return e.curveLookup[t] }

func (e *Tank) Type() gcpb.EntityType          { return gcpb.EntityType_ENTITY_TYPE_TANK }
func (e *Tank) Accept(v visitor.Visitor) error { return v.Visit(e) }
