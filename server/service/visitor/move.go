package move

import (
	"sync"

	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/curve/linearmove"
	"github.com/downflux/game/entity/entity"
	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/pathing/hpf/astar"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/service/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
)

const (
	pathLength = 0

	// TODO(minkezhang): Make this a property of the entity.
	ticksPerTile = float64(10)
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

	serverStatus *status.Status

	// entities is an append-only set of game entities.
	entities map[string]entity.Entity

	mux   sync.Mutex
	cache map[string]cacheRow
}

func New(tileMap *tile.Map, abstractGraph *graph.Graph, serverStatus *status.Status, entities map[string]entity.Entity) *Visitor {
	if entities == nil {
		entities = map[string]entity.Entity{}
	}
	return &Visitor{
		tileMap:       tileMap,
		abstractGraph: abstractGraph,
		serverStatus:  serverStatus,
		entities:      entities,
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

func (v *Visitor) Executor(tick float64) ([]curve.Curve, error) {
	v.mux.Lock()
	defer v.mux.Unlock()

	var curves []curve.Curve

	// TODO(minkezhang): Make this concurrent.
	for eid, cRow := range v.cache {
		if cRow.scheduledTick > tick {
			continue
		}
		e := v.entities[eid]
		c := e.Curve(gcpb.CurveCategory_CURVE_CATEGORY_MOVE)
		if c == nil {
			continue
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
			return nil, err
		}

		cv := linearmove.New(eid, tick)
		for i, tile := range p {
			cv.Add(tick+float64(i)*ticksPerTile, position(tile.Val.GetCoordinate()))
		}
		curves = append(curves, cv)

		lastPosition := position(p[len(p)-1].Val.GetCoordinate())
		if lastPosition == cRow.destination {
			delete(v.cache, eid)
		} else {
			if err := v.scheduleUnsafe(tick+float64(len(p)-1), eid, lastPosition); err != nil {
				// TODO(minkezhang): Handle error by logging and continuing.
				return nil, err
			}
		}
	}

	return curves, nil
}
