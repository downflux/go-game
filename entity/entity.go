package entity

import (
	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/curve/linearmove"
	"github.com/downflux/game/server/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

const idLen = 8

type Entity interface {
	ID() string
	Type() gcpb.EntityType
	Curve(t gcpb.CurveCategory) curve.Curve
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
	mc := linearmove.New(id.RandomString(idLen), eid)
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

// TODO(minkezhang): Decide if we should return default value.
func (e *SimpleEntity) Curve(t gcpb.CurveCategory) curve.Curve { return e.curveLookup[t] }
