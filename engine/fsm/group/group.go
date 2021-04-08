// Package group establishes the concept of "sibling" FSMs. These are FSM types
// which may not directly influence one another, but nevertheless conflict and
// must establish precedence, e.g. the interdependence of Chase, Move, and
// Attack -- even of the Chase FSM does not currently have a Move scheduled, it
// may still conflict with a newly-scheduled Move action.
package group

import (
	"strings"

	"github.com/downflux/game/engine/id/id"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

func fsmTypesToString(ts []fcpb.FSMType) string {
	var types []string
	for _, t := range ts {
		types = append(types, t.String())
	}
	return strings.Join(types, "/")
}

type Group struct {
	id       id.GroupID            // Read-only.
	fsmTypes map[fcpb.FSMType]bool // Read-only.
}

func New(ts []fcpb.FSMType) *Group {
	tg := &Group{
		id: id.GroupID(fsmTypesToString(ts)),
	}
	for _, t := range ts {
		tg.fsmTypes[t] = true
	}
	return tg
}

func (g Group) ID() id.GroupID {
	return g.id
}

func (g Group) Contains(t fcpb.FSMType) bool {
	return g.fsmTypes[t]
}

func (g Group) Collide(o Group) bool {
	for t := range o.fsmTypes {
		if g.Contains(t) {
			return true
		}
	}
	return false
}
