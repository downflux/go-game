package gamestate

import (
	"github.com/downflux/game/engine/entity/list"
	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

type GameState struct {
	// status represents the current Executor state metadata.
	status *status.Status

	// entities is a list of all Entity instances for the current game.
	// An Entity is an arbitrary stateful object -- it may not be a
	// physical game object like a tank; the entitylist.List object
	// itself is implements the Entity interface.
	//
	// Entity object states are mutated by Visitor instances.
	entities *list.List
}

func New(status *status.Status, entities *list.List) *GameState {
	return &GameState{
		status:   status,
		entities: entities,
	}
}

func (s *GameState) Entities() *list.List   { return s.entities }
func (s *GameState) Status() *status.Status { return s.status }

func (s *GameState) NoFilter() *dirty.List {
	filter := dirty.New()

	for _, e := range s.entities.Iter() {
		filter.AddEntity(dirty.Entity{ID: e.ID()})
		for _, property := range e.Properties() {
			filter.AddCurve(dirty.Curve{EntityID: e.ID(), Property: property})
		}
	}

	return filter
}

func (s *GameState) Export(tick id.Tick, filter *dirty.List) *gdpb.GameState {
	state := &gdpb.GameState{}

	for _, e := range filter.Entities() {
		state.Entities = append(
			state.GetEntities(),
			&gdpb.Entity{
				EntityId: e.ID.Value(),
				Type:     s.entities.Get(e.ID).Type(),
			},
		)
	}

	for _, c := range filter.Curves() {
		state.Curves = append(
			state.GetCurves(),
			s.entities.Get(c.EntityID).Curve(c.Property).ExportTail(tick),
		)
	}

	return state
}
