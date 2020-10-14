package entity

import (
	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type Entity interface {
	ID() string
	Type() gcpb.EntityType
	CurveID(t gcpb.CurveType) string
}
