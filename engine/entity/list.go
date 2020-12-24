// Package entitylist implements the Entity interface for a list tracking all
// game Entity instances.
package list

import (
	"github.com/downflux/game/engine/entity/entity"
	"github.com/downflux/game/engine/id/id"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// List implements the Entity interface for tracking all entities in a game.
type List struct {
	// entities is the list of all registered game entities, other than the
	// List instance itself.
	entities map[id.EntityID]entity.Entity
}

// New constructs a new instance of the List.
func New(eid id.EntityID) *List {
	return &List{
		entities: map[id.EntityID]entity.Entity{},
	}
}

// Get returns a specific Entity instance, given the UUID. Get returns nil if
// the UUID specified is the ID of the List.
func (l *List) Get(eid id.EntityID) entity.Entity {
	return l.entities[eid]
}

// Iter returns the list of Entity instances tracked by the List. This is used
// for loop ranges.
//
// TODO(minkezhang): Determine if this is deprecated or not.
func (l *List) Iter() []entity.Entity {
	var entities []entity.Entity
	for _, e := range l.entities {
		entities = append(entities, e)
	}

	return entities
}

// Append tracks a new Entity instance.
func (l *List) Append(e entity.Entity) error {
	if _, found := l.entities[e.ID()]; found {
		return status.Error(codes.AlreadyExists, "an entity already exists with the given ID")
	}

	l.entities[e.ID()] = e
	return nil
}
