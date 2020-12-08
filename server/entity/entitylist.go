// Package entitylist implements the Entity interface for a list tracking all
// game Entity instances.
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

// List implements the Entity interface for tracking all entities in a game.
type List struct {
	entity.Base
	entity.NoCurve

	// eid is the UUID of the entity.
	eid id.EntityID

	// entities is the list of all registered game entities, other than the
	// List instance itself.
	entities map[id.EntityID]entity.Entity
}

// New constructs a new instance of the List.
func New(eid id.EntityID) *List {
	return &List{
		entities: map[id.EntityID]entity.Entity{},
		eid:      eid,
	}
}

// ID returns the UUID of the List.
func (l *List) ID() id.EntityID { return l.eid }

// Get returns a specific Entity instance, given the UUID. Get returns nil if
// the UUID specified is the ID of the List.
func (l *List) Get(eid id.EntityID) entity.Entity {
	return l.entities[eid]
}

// Iter returns the list of Entity instances tracked by the List. This is used
// for loop ranges.
func (l *List) Iter() []entity.Entity {
	var entities []entity.Entity
	for _, e := range l.entities {
		entities = append(entities, e)
	}

	return entities
}

// Add tracks a new Entity instance.
//
// TODO(minkezhang): Rename to Append.
func (l *List) Add(e entity.Entity) error {
	if _, found := l.entities[e.ID()]; found {
		return status.Error(codes.AlreadyExists, "an entity already exists with the given ID")
	}

	l.entities[e.ID()] = e
	return nil
}

// Type returns the registered EntityType of the List.
func (l *List) Type() gcpb.EntityType { return gcpb.EntityType_ENTITY_TYPE_ENTITY_LIST }

// Accept registers a Visitor instance and defines the order in which managed
// entities will be mutated by the input.
//
// This is part of the visitor pattern.
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
