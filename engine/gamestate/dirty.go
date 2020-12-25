// Package dirty encapsulates logic necessary for marking specific Curve and
// Entity instances as having been modified during the current game tick.
package dirty

import (
	"sync"

	"github.com/downflux/game/engine/id/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

// Curve represents a curve.Curve instance which was altered in the current
// tick and will need to be broadcast to all clients.
//
// The Entity UUID and EntityProperty uniquely identifies a curve.
type Curve struct {
	EntityID id.EntityID
	Property gcpb.EntityProperty
}

// Entity represents a visitor.Entity instance which was added in the current
// tick. Curves which were altered do not need to create a dirty Entity entry.
type Entity struct {
	ID id.EntityID
}

// List is an abstract cache of game state mutations over a period of time.
type List struct {
	mux      sync.Mutex
	curves   map[id.EntityID]map[gcpb.EntityProperty]bool
	entities map[id.EntityID]bool
}

// New creates a new List instance.
func New() *List {
	return &List{}
}

func (l *List) Curves() []Curve {
	l.mux.Lock()
	defer l.mux.Unlock()

	var curves []Curve
	for eid, properties := range l.curves {
		for property := range properties {
			curves = append(curves, Curve{
				EntityID: eid,
				Property: property,
			})
		}
	}

	return curves
}

func (l *List) Entities() []Entity {
	l.mux.Lock()
	defer l.mux.Unlock()

	var entities []Entity
	for eid := range l.entities {
		entities = append(entities, Entity{ID: eid})
	}

	return entities
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

// AddCurve marks the specified curve as dirty.
func (l *List) AddCurve(c Curve) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.curves == nil {
		l.curves = map[id.EntityID]map[gcpb.EntityProperty]bool{}
	}
	if l.curves[c.EntityID] == nil {
		l.curves[c.EntityID] = map[gcpb.EntityProperty]bool{}
	}

	l.curves[c.EntityID][c.Property] = true

	return nil
}

// Pop returns a clone of the current dirty list and resets the internal cache.
func (l *List) Pop() *List {
	l.mux.Lock()
	defer l.mux.Unlock()

	nl := &List{
		curves:   l.curves,
		entities: l.entities,
	}

	l.curves = nil
	l.entities = nil

	return nl
}
