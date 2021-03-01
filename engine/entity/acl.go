package acl

import (
	"github.com/downflux/game/engine/curve/common/step"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

type ACLType uint32

const (
	ClientWritable = 1 << iota
	PublicWritable
)

type ACL struct {
	acl           ACLType
	clientIDCurve *step.Curve
}

func New(acl ACLType, cidc *step.Curve) *ACL {
	return &ACL{
		clientIDCurve: cidc,
		acl:           acl,
	}
}

func (a ACL) PublicWritable() bool { return a.acl&PublicWritable == PublicWritable }
func (a ACL) ClientWritable() bool { return a.acl&ClientWritable == ClientWritable }

func (a ACL) Export() *gdpb.ACL {
	return &gdpb.ACL{
		PublicWritable: a.PublicWritable(),
		ClientWritable: a.ClientWritable(),
	}
}
