// Package dirty encapsulates logic necessary for marking specific Curve and
// Entity instances as having been modified during the current game tick.
package dirty

import (
	"sync"

	"github.com/downflux/game/server/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

// Curve represents a curve.Curve instance which was altered in the current
// tick and will need to be broadcast to all clients.
//
// The Entity UUID and CurveCategory uniquely identifies a curve.
type Curve struct {
	EntityID id.EntityID
	Category gcpb.CurveCategory
}

// Entity represents a visitor.Entity instance which was added in the current
// tick. Curves which were altered do not need to create a dirty Entity entry.
type Entity struct {
	ID id.EntityID
}

// List is an abstract cache of game state mutations over a period of time.
type List struct {
	mux      sync.Mutex
	curves   map[id.EntityID]map[gcpb.CurveCategory]bool
	entities map[id.EntityID]bool
}

// New creates a new List instance.
func New() *List {
	return &List{}
}

// AddEntity marks the specified entity as dirty.
func (l *List) AddEntity(e Entity) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.entities == nil {
		l.entities = map[id.EntityID]bool{}
	}

	l.entities[e.ID] = true

	return nil
}

// Add marks the specified curve as dirty.
func (l *List) Add(c Curve) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.curves == nil {
		l.curves = map[id.EntityID]map[gcpb.CurveCategory]bool{}
	}
	if l.curves[c.EntityID] == nil {
		l.curves[c.EntityID] = map[gcpb.CurveCategory]bool{}
	}

	l.curves[c.EntityID][c.Category] = true

	return nil
}

// PopEntities returns a list of all mutated entities for the given interval.
// The internal cache is then cleared.
//
// TODO(minkezhang): Consolidate PopEntities and Pop.
func (l *List) PopEntities() []Entity {
	l.mux.Lock()
	defer l.mux.Unlock()

	var entities []Entity
	for eid := range l.entities {
		entities = append(entities, Entity{ID: eid})
	}

	l.entities = nil
	return entities
}

// Pop returns a list of all mutated curves for the given interval.
// The internal cache is then cleared.
//
// TODO(minkezhang): Consolidate PopEntities and Pop.
func (l *List) Pop() []Curve {
	l.mux.Lock()
	defer l.mux.Unlock()

	var curves []Curve
	for eid, categories := range l.curves {
		for category := range categories {
			curves = append(curves, Curve{
				EntityID: eid,
				Category: category,
			})
		}
	}
	l.curves = nil

	return curves
}
