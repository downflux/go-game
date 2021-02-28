package acl

import (
	gdpb "github.com/downflux/game/api/data_go_proto"
)

type ACLType uint32
const (
	PublicWritable = 1 << iota
	TeamWritable
	ClientWritable
)

type ACL struct {
	acl ACLType
}

func New(acl ACLType) *ACL { return &ACL{ acl: acl } }
func (a ACL) PublicWritable() bool { return a.acl & PublicWritable == PublicWritable }
func (a ACL) TeamWritable() bool { return a.acl & TeamWritable == TeamWritable }
func (a ACL) ClientWritable() bool { return a.acl & ClientWritable == ClientWritable }

func (a ACL) Export() *gdpb.ACL {
	return &gdpb.ACL{
		PublicWritable: a.PublicWritable(),
		TeamWritable: a.TeamWritable(),
		ClientWritable: a.ClientWritable(),
	}
}
