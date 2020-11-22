package move

import (
	"sync"

	"github.com/downflux/game/curve/linearmove"
	"github.com/downflux/game/entity/entity"
	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/pathing/hpf/astar"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/service/visitor/dirty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	serverstatus "github.com/downflux/game/server/service/status"
)

const (
	pathLength = 0

	// TODO(minkezhang): Make this a property of the entity.
	ticksPerTile = float64(10)
)

func unsupportedError(entityType gcpb.EntityType) error {
	return status.Errorf(codes.Unimplemented, "move not implemented for Entity type %v", entityType)
}

// coordinate transforms a gdpb.Position instance into a gdpb.Coordinate
// instance. We're assuming the position values are sane and don't overflow
// int32.
func coordinate(p *gdpb.Position) *gdpb.Coordinate {
	return &gdpb.Coordinate{
		X: int32(p.GetX()),
		Y: int32(p.GetY()),
	}
}

func position(c *gdpb.Coordinate) *gdpb.Position {
	return &gdpb.Position{
		X: float64(c.GetX()),
		Y: float64(c.GetY()),
	}
}

type cacheRow struct {
	scheduledTick float64
	destination   *gdpb.Position
}

type Args struct {
	EID         string
	Destination *gdpb.Position
}

type Visitor struct {
	// tileMap is the underlying Map object used for the game.
	tileMap *tile.Map

	// abstractGraph is the underlying abstracted pathing logic data layer
	// for the associated Map.
	abstractGraph *graph.Graph

	dfStatus *serverstatus.Status

	dirtyCurves *dirty.List

	mux   sync.Mutex
	cache map[string]cacheRow
}

func New(
	tileMap *tile.Map,
	abstractGraph *graph.Graph,
	dfStatus *serverstatus.Status,
	dirtyCurves *dirty.List) *Visitor {
	return &Visitor{
		tileMap:       tileMap,
		abstractGraph: abstractGraph,
		dfStatus:      dfStatus,
		dirtyCurves:   dirtyCurves,
		cache:         map[string]cacheRow{},
	}
}

func (v *Visitor) scheduleUnsafe(tick float64, eid string, dest *gdpb.Position) error {
	if v.cache == nil {
		v.cache = map[string]cacheRow{}
	}

	if v.cache[eid].scheduledTick < tick {
		v.cache[eid] = cacheRow{
			scheduledTick: tick,
			destination:   dest,
		}
	}

	return nil
}

func (v *Visitor) Schedule(tick float64, args interface{}) error {
	v.mux.Lock()
	defer v.mux.Unlock()

	argsImpl := args.(Args)

	return v.scheduleUnsafe(tick, argsImpl.EID, argsImpl.Destination)
}

func (v *Visitor) Visit(e entity.Entity) error {
	if e.Type() != gcpb.EntityType_ENTITY_TYPE_TANK {
		return unsupportedError(e.Type())
	}

	tick := v.dfStatus.Tick()

	// TODO(minkezhang): Make this concurrent.
	v.mux.Lock()
	defer v.mux.Unlock()

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
	// TODO(minkezhang): Add additional infrastructure necessary to set pathLength > 0.
	p, _, err := astar.Path(
		v.tileMap,
		v.abstractGraph,
		utils.MC(coordinate(c.Get(tick).(*gdpb.Position))),
		utils.MC(coordinate(cRow.destination)),
		pathLength)
	if err != nil {
		// TODO(minkezhang): Handle error by logging and continuing.
		return err
	}

	cv := linearmove.New(e.ID(), tick)
	for i, tile := range p {
		cv.Add(tick+float64(i)*ticksPerTile, position(tile.Val.GetCoordinate()))
	}
	if err := v.dirtyCurves.Add(dirty.Curve{
		EntityID: e.ID(),
		Category: c.Category(),
	}); err != nil {
		return err
	}
	if err := c.ReplaceTail(cv); err != nil {
		return err
	}

	lastPosition := position(p[len(p)-1].Val.GetCoordinate())
	if lastPosition == cRow.destination {
		delete(v.cache, e.ID())
	} else {
		if err := v.scheduleUnsafe(tick+float64(len(p)-1), e.ID(), lastPosition); err != nil {
			// TODO(minkezhang): Handle error by logging and continuing.
			return err
		}
	}
	return nil
}
