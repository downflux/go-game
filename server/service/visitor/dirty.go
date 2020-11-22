package dirty

import (
	"sync"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

// Curve represents a curve.Curve instance which was altered in the current
// tick and will need to be broadcast to all clients.
//
// The Entity UUID and CurveCategory uniquely identifies a curve.
type Curve struct {
	EntityID string
	Category gcpb.CurveCategory
}

type List struct {
	curvesMux sync.Mutex
	curves    map[string]map[gcpb.CurveCategory]bool
}

func New() *List {
	return &List{}
}

func (l *List) Add(c Curve) error {
	l.curvesMux.Lock()
	defer l.curvesMux.Unlock()

	if l.curves == nil {
		l.curves = map[string]map[gcpb.CurveCategory]bool{}
	}
	if l.curves[c.EntityID] == nil {
		l.curves[c.EntityID] = map[gcpb.CurveCategory]bool{}
	}

	l.curves[c.EntityID][c.Category] = true

	return nil
}

func (l *List) Pop() []Curve {
	l.curvesMux.Lock()
	defer l.curvesMux.Unlock()

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
