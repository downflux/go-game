package dirty

import (
	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type Curve struct {
	eid       string
	curveType gcpb.CurveCategory
}
