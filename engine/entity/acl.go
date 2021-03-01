package acl

import (
	"github.com/downflux/game/engine/curve/common/step"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

type Permission uint32

const (
	ClientWritable = 1 << iota
	PublicWritable
)

type ACL struct {
	permission    Permission
	clientIDCurve *step.Curve
}

func New(cidc *step.Curve, p Permission) *ACL {
	return &ACL{
		clientIDCurve: cidc,
		permission:    p,
	}
}

func (a ACL) PublicWritable() bool { return a.permission&PublicWritable == PublicWritable }
func (a ACL) ClientWritable() bool { return a.permission&ClientWritable == ClientWritable }

func (a ACL) Export() *gdpb.ACL {
	return &gdpb.ACL{
		PublicWritable: a.PublicWritable(),
		ClientWritable: a.ClientWritable(),
	}
}
