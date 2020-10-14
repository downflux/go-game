package entity

import (
	"math/rand"

	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/curve/linearmove"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

type Entity interface {
	ID() string
	Type() gcpb.EntityType
	Curve(t gcpb.CurveCategory) curve.Curve
}

type SimpleEntity struct {
	id          string
	curveLookup map[gcpb.CurveCategory]curve.Curve
}

// TODO(minkezhang): Export to shared for server.
func randID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// TODO(minkezhang): Make this client-friendly too.
func NewSimpleEntity(eid string, t float64, p *gdpb.Position) *SimpleEntity {
	mc := linearmove.New(randID(), eid)
	mc.Add(t, p)

	return &SimpleEntity{id: eid, curveLookup: map[gcpb.CurveCategory]curve.Curve{
		gcpb.CurveCategory_CURVE_CATEGORY_MOVE: mc,
	}}
}

func (e *SimpleEntity) ID() string                             { return e.id }
func (e *SimpleEntity) Type() gcpb.EntityType                  { return gcpb.EntityType_ENTITY_TYPE_TANK }
func (e *SimpleEntity) Curve(t gcpb.CurveCategory) curve.Curve { return e.curveLookup[t] }
