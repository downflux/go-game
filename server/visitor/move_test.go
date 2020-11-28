package move

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/status"
	"github.com/downflux/game/server/service/visitor/dirty"
	"github.com/downflux/game/server/service/visitor/entity/tank"
	"golang.org/x/sync/errgroup"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	mcpb "github.com/downflux/game/map/api/constants_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
)

var (
	/**
		 *       - - -
		 *       - - -
	         * Y = 0 - - -
	         *   X = 0
	*/
	simpleMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 3, Y: 3},
		Tiles: []*mdpb.Tile{
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 2}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 2}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 2}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
		TerrainCosts: []*mdpb.TerrainCost{
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
		},
	}
)

func TestSchedule(t *testing.T) {
	const nClients = 1000
	v := New(nil, nil, nil, nil, 0)

	var eg errgroup.Group
	for i := 0; i < nClients; i++ {
		i := i
		eg.Go(func() error {
			return v.Schedule(Args{Tick: 0, EntityID: id.NewEntityID(fmt.Sprintf("entity-%d", i))})
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("Wait() = %v, want = nil", err)
	}

	if got := len(v.cache); got != nClients {
		t.Errorf("len() = %v, want = %v", got, nClients)
	}
}

func newVisitor(t *testing.T) *Visitor {
	tm, err := tile.ImportMap(simpleMap)
	if err != nil {
		t.Fatalf("Import() = _, %v, want = nil", err)
	}

	g, err := graph.BuildGraph(tm, &gdpb.Coordinate{X: 1, Y: 1})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = nil", err)
	}

	s := status.New(time.Millisecond)
	d := dirty.New()

	return New(tm, g, s, d, 1)
}

func TestVisitNoSchedule(t *testing.T) {
	const eid = "entity-id"
	const t0 = 0

	v := newVisitor(t)
	e := tank.New(eid, t0, &gdpb.Position{X: 0, Y: 0})

	if err := v.Visit(e); err != nil {
		t.Fatalf("Visit() = %v, want = nil", err)
	}

	if got := v.dirties.Pop(); got != nil {
		t.Errorf("Pop() = %v, want = nil", got)
	}
}

func TestVisitFutureSchedule(t *testing.T) {
	const eid = "entity-id"
	dest := &gdpb.Position{X: 2, Y: 0}
	const t0 = 0

	v := newVisitor(t)
	e := tank.New(eid, t0, &gdpb.Position{X: 0, Y: 0})
	v.Schedule(Args{Tick: t0 + 1, EntityID: eid, Destination: dest})

	if err := v.Visit(e); err != nil {
		t.Fatalf("Visit() = %v, want = nil", err)
	}

	if got := v.dirties.Pop(); got != nil {
		t.Errorf("Pop() = %v, want = nil", got)
	}
}

func TestVisit(t *testing.T) {
	const eid = "entity-id"
	dest := &gdpb.Position{X: 2, Y: 0}
	const t0 = 0

	v := newVisitor(t)
	e := tank.New(eid, t0, &gdpb.Position{X: 0, Y: 0})
	v.Schedule(Args{Tick: t0, EntityID: eid, Destination: dest})

	if err := v.Visit(e); err != nil {
		t.Fatalf("Visit() = %v, want = nil", err)
	}

	func(t *testing.T) {
		want := cacheRow{
			scheduledTick: t0 + ticksPerTile,
			destination:   dest,
		}

		v.cacheMux.Lock()
		defer v.cacheMux.Unlock()
		if got := v.cache[eid]; got != want {
			t.Fatalf("cache[] = %v, want = %v", got, dest)
		}
	}(t)

	want := dirty.Curve{
		Category: gcpb.CurveCategory_CURVE_CATEGORY_MOVE,
		EntityID: eid,
	}
	if got := v.dirties.Pop(); got[0] != want {
		t.Errorf("Pop() = %v, want = %v", got[0], want)
	}
}
