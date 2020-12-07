// Package move implements logic for entity position mutations.
package move

import (
	"log"
	"math"
	"sync"

	"github.com/downflux/game/curve/linearmove"
	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/pathing/hpf/astar"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/status"
	"github.com/downflux/game/server/visitor/dirty"
	"github.com/downflux/game/server/visitor/visitor"
	"google.golang.org/protobuf/proto"

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

func d(a, b *gdpb.Position) float64 {
	return math.Sqrt(math.Pow(a.GetX() - b.GetX(), 2) + math.Pow(a.GetY() - b.GetY(), 2))
}

// position transforms a gdpb.Coordinate instance into a gdpb.Position
// instance.
func position(c *gdpb.Coordinate) *gdpb.Position {
	return &gdpb.Position{
		X: float64(c.GetX()),
		Y: float64(c.GetY()),
	}
}

// Args is an external-facing struct used to specify the Entity being moved.
type Args struct {
	// Tick is the tick at which the entity should start moving.
	Tick id.Tick

	// EntityID is the UUID of the Entity instance.
	EntityID id.EntityID

	// Destination is the Entity move target.
	Destination *gdpb.Position

	// IsExternal indicates if the move request was issued by an external
	// client.
	// TODO(minkezhang): Delete this. This is an alias for IsPartial, which
	// may only be scheduled by the move Visitor itself. If we want a IsBot
	// bool, we should add that separately.
	IsExternal bool
}

type partialCacheRow struct {
	scheduledTick id.Tick

	// TODO(minkezhang): Rename to isPartial.
	isExternal bool
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

	partialCache map[id.EntityID]partialCacheRow
	destinationCache map[id.EntityID]*gdpb.Position
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
		minPathLength: minPathLength,
		partialCache: map[id.EntityID]partialCacheRow{},
		destinationCache: map[id.EntityID]*gdpb.Position{},
	}
}

// Type returns the registered VisitorType.
func (v *Visitor) Type() vcpb.VisitorType { return visitorType }

// scheduleUnsafe adds a move command to the cache.
func (v *Visitor) scheduleUnsafe(tick id.Tick, eid id.EntityID, dest *gdpb.Position, isExternal bool) error {
	if v.partialCache == nil {
		v.partialCache = map[id.EntityID]partialCacheRow{}
	}
	if v.destinationCache == nil {
		v.destinationCache = map[id.EntityID]*gdpb.Position{}
	}

	_, isScheduled := v.partialCache[eid]

	// log.Println(dest, v.destinationCache[eid])

	if (
		isExternal &&
		// Only schedule an external move if the scheduled position
		// is different from the previously scheduled one. Partial
		// moves may bypass this check because the point is to update
		// the execution tick.
		!proto.Equal(dest, v.destinationCache[eid]) && (
		// External moves may not be scheduled for the future.
		v.dfStatus.Tick() >= tick && (
			// External moves override all scheduled partial moves.
			!v.partialCache[eid].isExternal || (
				// External moves override all external moves,
				// if those moves were scheduled for an earlier
				// tick.
				isScheduled &&
				tick > v.partialCache[eid].scheduledTick))) || (
		!isExternal &&
		// Internal moves only override existing internal move
		// scheduled to execute at an earlier tick.
		tick > v.partialCache[eid].scheduledTick)) {
		v.partialCache[eid] = partialCacheRow{
			scheduledTick: tick,
			isExternal: isExternal,
		}
		v.destinationCache[eid] = dest
	}

	return nil
}

// Schedule adds a move command to the internal schedule.
func (v *Visitor) Schedule(args interface{}) error {
	argsImpl := args.(Args)

	v.cacheMux.Lock()
	defer v.cacheMux.Unlock()

	return v.scheduleUnsafe(argsImpl.Tick, argsImpl.EntityID, argsImpl.Destination, argsImpl.IsExternal)
}

// Visit mutates the specified entity's position curve.
//
// TODO(minkezhang): Observe "spam clicking behavior" and find out why client
// keeps "jumping" the coordinate.
func (v *Visitor) Visit(e visitor.Entity) error {
	if e.Type() != gcpb.EntityType_ENTITY_TYPE_TANK {
		return nil
	}

	tick := v.dfStatus.Tick()

	// TODO(minkezhang): Make this concurrent.
	v.cacheMux.Lock()
	defer v.cacheMux.Unlock()

	partialCache, found := v.partialCache[e.ID()]
	if !found {
		return nil
	}
	if partialCache.scheduledTick > tick {
		return nil
	}

	destination := v.destinationCache[e.ID()]

	c := e.Curve(gcpb.CurveCategory_CURVE_CATEGORY_MOVE)
	if c == nil {
		return nil
	}

	// TODO(minkezhang): proto.Clone the return values in map.astar.Path.
	p, _, err := astar.Path(
		v.tileMap,
		v.abstractGraph,
		// TODO(minkezhang): Investigate / decide if we should use
		// scheduledTick instead.
		utils.MC(coordinate(c.Get(tick).(*gdpb.Position))),
		utils.MC(coordinate(destination)),
		v.minPathLength)
	if err != nil {
		// TODO(minkezhang): Handle error by logging and continuing.
		return err
	}

	if p == nil {
		return nil
	}

	// Add to the existing curve, while smoothing out the existing
	// trajectory.
	prevPos := c.Get(tick).(*gdpb.Position)
	cv := linearmove.New(e.ID(), tick)
	cv.Add(tick, prevPos)
	for i, tile := range p {
		curPos := position(tile.Val.GetCoordinate())
		tickDelta := id.Tick(ticksPerTile.Value() * d(prevPos, curPos))
		cv.Add(tick+id.Tick(i)*ticksPerTile+tickDelta, curPos)
		prevPos = curPos
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
	log.Println("MOVE: ", tick, proto.Equal(lastPosition, destination))
	if proto.Equal(lastPosition, destination) {
		// We need to keep track of current pending destination.
		// Deleting the destination cache here allows spam clicking.
		delete(v.partialCache, e.ID())
	} else {
		if err := v.scheduleUnsafe(
			tick+ticksPerTile*id.Tick(len(p) - 1),
			e.ID(),
			destination,
			false); err != nil {
			// TODO(minkezhang): Handle error by logging and continuing.
			return err
		}
	}

	return nil
}
