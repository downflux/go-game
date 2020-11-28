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

type Entity struct {
	ID id.EntityID
}

type List struct {
	mux      sync.Mutex
	curves   map[id.EntityID]map[gcpb.CurveCategory]bool
	entities map[id.EntityID]bool
}

func New() *List {
	return &List{}
}

func (l *List) AddEntity(e Entity) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.entities == nil {
		l.entities = map[id.EntityID]bool{}
	}

	l.entities[e.ID] = true

	return nil
}

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
