package entitylist

import (
	"github.com/downflux/game/server/entity/entity"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/visitor/visitor"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type List struct {
	entity.BaseEntity
	entity.NoCurveEntity

	eid id.EntityID

	entities map[id.EntityID]visitor.Entity
}

func New(eid id.EntityID) *List {
	return &List{
		entities: map[id.EntityID]visitor.Entity{},
		eid:      eid,
	}
}

func (l *List) ID() id.EntityID { return l.eid }

func (l *List) Get(eid id.EntityID) visitor.Entity {
	return l.entities[eid]
}

func (l *List) Iter() []visitor.Entity {
	var entities []visitor.Entity
	for _, e := range l.entities {
		entities = append(entities, e)
	}

	return entities
}

func (l *List) Add(e visitor.Entity) error {
	if _, found := l.entities[e.ID()]; found {
		return status.Error(codes.AlreadyExists, "an entity already exists with the given ID")
	}

	l.entities[e.ID()] = e
	return nil
}

func (l *List) Type() gcpb.EntityType { return gcpb.EntityType_ENTITY_TYPE_ENTITY_LIST }
func (l *List) Accept(v visitor.Visitor) error {
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
