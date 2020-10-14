package entity

import (
	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type Entity interface {
	ID() string
	Type() gcpb.EntityType
	CurveID(t gcpb.CurveType) string
}

type SimpleEntity struct {
	id string
	curveLookup map[gcpb.CurveType]string
}

func NewSimpleEntity(id string) *SimpleEntity {
	return &SimpleEntity{id: id, curveLookup: map[gcpb.CurveType]string{}}
}
func (e *SimpleEntity) ID() string { return e.id }
func (e *SimpleEntity) Type() gcpb.EntityType { return gcpb.EntityType_ENTITY_TYPE_TANK }
func (e *SimpleEntity) CurveID(t gcpb.CurveType) string { return e.curveLookup[t] }
