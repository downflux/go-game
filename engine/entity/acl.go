package acl

import (
	"github.com/downflux/game/engine/id/id"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

type ACLType uint32

const (
	ClientWritable = 1 << iota
	PublicWritable
)

type ACL struct {
	clientID id.ClientID
	acl      ACLType
}

func New(cid id.ClientID, acl ACLType) *ACL {
	return &ACL{
		clientID: cid,
		acl:      acl,
	}
}

func (a ACL) PublicWritable() bool { return a.acl&PublicWritable == PublicWritable }
func (a ACL) ClientWritable() bool { return a.acl&ClientWritable == ClientWritable }

func (a ACL) Validate(c *gdpb.ClientID) bool {
	return a.PublicWritable() || (a.validateClientWritable(
		id.ClientID(c.GetClientId())))
}

func (a ACL) validateClientWritable(c id.ClientID) bool {
	return a.ClientWritable() && c == a.clientID
}

func (a ACL) Export() *gdpb.ACL {
	return &gdpb.ACL{
		Owner: &gdpb.ClientID{
			ClientId: a.clientID.Value(),
		},
		PublicWritable: a.PublicWritable(),
		ClientWritable: a.ClientWritable(),
	}
}
