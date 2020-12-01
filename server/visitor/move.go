// Package move implements logic for entity position mutations.
package move

import (
	"sync"

	"github.com/downflux/game/curve/linearmove"
	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/pathing/hpf/astar"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/status"
	"github.com/downflux/game/server/visitor/dirty"
	"github.com/downflux/game/server/visitor/visitor"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	vcpb "github.com/downflux/game/server/visitor/api/constants_go_proto"
)

const (
	// ticksPerTile represents the number of ticks necessary to move an
	// Entity instance one full tile width.
	//
	// TODO(minkezhang): Make this a property of the entity.
	ticksPerTile = id.Tick(10)

	// visitorType is the registered VisitorType of the move visitor.
	visitorType = vcpb.VisitorType_VISITOR_TYPE_MOVE
)

// coordinate transforms a gdpb.Position instance into a gdpb.Coordinate
// instance. We're assuming the position values are sane and don't overflow
// int32.
func coordinate(p *gdpb.Position) *gdpb.Coordinate {
	return &gdpb.Coordinate{
		X: int32(p.GetX()),
		Y: int32(p.GetY()),
	}
}

// position transforms a gdpb.Coordinate instance into a gdpb.Position
// instance.
func position(c *gdpb.Coordinate) *gdpb.Position {
	return &gdpb.Position{
		X: float64(c.GetX()),
		Y: float64(c.GetY()),
	}
}

// cacheRow represents a scheduled move command. The EntityID is stored in the
// map key.
type cacheRow struct {
	// scheduledTick represents the tick at which the path should be
	// calculated.
	scheduledTick id.Tick

	// destination represents the move target for the Entity.
	destination *gdpb.Position
}

// Args is an external-facing struct used to specify the Entity being moved.
type Args struct {
	// Tick is the tick at which the entity should start moving.
	Tick id.Tick

	// EntityID is the UUID of the Entity instance.
	EntityID id.EntityID

	// Destination is the Entity move target.
	Destination *gdpb.Position
}

// Visitor mutates the Entity position Curve. This struct implements the
// visitor.Visitor interface.
type Visitor struct {
	// tileMap is the underlying Map object used for the game.
	tileMap *tile.Map

	// abstractGraph is the underlying abstracted pathing logic data layer
	// for the associated Map.
	abstractGraph *graph.Graph

	// dfStatus is a shared object with the game engine and indicates
	// current tick, etc.
	dfStatus *status.Status

	// dirties is a shared object between the game engine and the
	// Visitor.
	dirties *dirty.List

	// minPathLength represents the minimum lookahead path length to
	// calculate, where the path is a list of tile.Map coordinates.
	//
	// Longer calculated paths is discouraged, as these paths become
	// outdated once a new move command is issued for the Entity, which
	// may happen frequently in an RTS game.
	minPathLength int

	// cacheMux guards the cache property from concurrent access.
	cacheMux sync.Mutex

	// cache is the list of scheduled move commands to execute.
	cache map[id.EntityID]cacheRow
}

// New constructs a new move Visitor instance.
func New(
	tileMap *tile.Map,
	abstractGraph *graph.Graph,
	dfStatus *status.Status,
	dirties *dirty.List,
	minPathLength int) *Visitor {
	return &Visitor{
		tileMap:       tileMap,
		abstractGraph: abstractGraph,
		dfStatus:      dfStatus,
		dirties:       dirties,
		cache:         map[id.EntityID]cacheRow{},
		minPathLength: minPathLength,
	}
}

// Type returns the registered VisitorType.
func (v *Visitor) Type() vcpb.VisitorType { return visitorType }

// scheduleUnsafe adds a move command to the cache.
func (v *Visitor) scheduleUnsafe(tick id.Tick, eid id.EntityID, dest *gdpb.Position) error {
	if v.cache == nil {
		v.cache = map[id.EntityID]cacheRow{}
	}

	if v.cache[eid].scheduledTick <= tick {
		v.cache[eid] = cacheRow{
			scheduledTick: tick,
			destination:   dest,
		}
	}

	return nil
}

// Schedule adds a move command to the internal schedule.
func (v *Visitor) Schedule(args interface{}) error {
	argsImpl := args.(Args)

	v.cacheMux.Lock()
	defer v.cacheMux.Unlock()

	return v.scheduleUnsafe(argsImpl.Tick, argsImpl.EntityID, argsImpl.Destination)
}

// Visit mutates the specified entity's position curve.
func (v *Visitor) Visit(e visitor.Entity) error {
	if e.Type() != gcpb.EntityType_ENTITY_TYPE_TANK {
		return nil
	}

	tick := v.dfStatus.Tick()

	// TODO(minkezhang): Make this concurrent.
	v.cacheMux.Lock()
	defer v.cacheMux.Unlock()

	cRow, found := v.cache[e.ID()]
	if !found {
		return nil
	}

	if cRow.scheduledTick > tick {
		return nil
	}

	c := e.Curve(gcpb.CurveCategory_CURVE_CATEGORY_MOVE)
	if c == nil {
		return nil
	}

	// TODO(minkezhang): proto.Clone the return values in map.astar.Path.
	p, _, err := astar.Path(
		v.tileMap,
		v.abstractGraph,
		utils.MC(coordinate(c.Get(tick).(*gdpb.Position))),
		utils.MC(coordinate(cRow.destination)),
		v.minPathLength)
	if err != nil {
		// TODO(minkezhang): Handle error by logging and continuing.
		return err
	}

	cv := linearmove.New(e.ID(), tick)
	for i, tile := range p {
		cv.Add(tick+id.Tick(i)*ticksPerTile, position(tile.Val.GetCoordinate()))
	}
	if err := v.dirties.Add(dirty.Curve{
		EntityID: e.ID(),
		Category: c.Category(),
	}); err != nil {
		return err
	}
	if err := c.ReplaceTail(cv); err != nil {
		return err
	}

	// Check for partial moves and delay next lookup iteration until a
	// suitable time in the future.
	lastPosition := position(p[len(p)-1].Val.GetCoordinate())
	if lastPosition == cRow.destination {
		delete(v.cache, e.ID())
	} else {
		if err := v.scheduleUnsafe(tick+ticksPerTile*id.Tick(len(p)), e.ID(), cRow.destination); err != nil {
			// TODO(minkezhang): Handle error by logging and continuing.
			return err
		}
	}

	return nil
}
