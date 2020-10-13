package curve

import (
	"reflect"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type Curve interface {
	Type() gcpb.CurveType
	DatumType() reflect.Type
	ClientID() string
	CurveID() string
	Get(t float32) interface{} // interpolation
}

func IsSubsetOf(c1, c2 Curve) bool
func Hash(c Curve) bool // placeholder type
func Merge(c1, c2 Curve) error
func Extract(t1, t2 float32) Curve // CurveID is the same
