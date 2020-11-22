package entitylist

import (
	"github.com/downflux/game/server/service/visitor/entity/entity"
	"github.com/downflux/game/server/service/visitor/visitor"
	"golang.org/x/sync/errgroup"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type EntityList struct {
	entity.BaseEntity

	entities map[string]visitor.Entity
}

func New() *EntityList {
	return &EntityList{
		entities: map[string]visitor.Entity{},
	}
}

func (l *EntityList) Type() gcpb.EntityType { return gcpb.EntityType_ENTITY_TYPE_ENTITY_LIST }
func (l *EntityList) Accept(v visitor.Visitor) error {
	if err := v.Visit(l); err != nil {
		return err
	}

	var eg errgroup.Group
	for _, e := range l.entities {
		e := e
		eg.Go(func() error { return e.Accept(v) })
	}
	return eg.Wait()
}
