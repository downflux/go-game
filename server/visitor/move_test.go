package move

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/status"
	"github.com/downflux/game/server/visitor/dirty"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/testing/protocmp"

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

type cacheRow struct {
	scheduledTick id.Tick
	destination *gdpb.Position
	isExternal bool
}

func getCache(v *Visitor, eid id.EntityID) cacheRow {
	v.cacheMux.Lock()
	defer v.cacheMux.Unlock()

	return cacheRow{
		scheduledTick: v.partialCache[eid].scheduledTick,
		destination: v.destinationCache[eid],
		isExternal: v.partialCache[eid].isExternal,
	}
}

func TestScheduleIdempotency(t *testing.T) {
	const eid = "entity-id"
	p0 := &gdpb.Position{X: 0, Y: 0}
	p1 := &gdpb.Position{X: 1, Y: 1}
	p2 := &gdpb.Position{X: 2, Y: 2}

	inOrderExternalMoveUpdate := newVisitor(t)
	inOrderExternalMoveUpdate.dfStatus.IncrementTick()

	outOfOrderExternalMoveIdempotence := newVisitor(t)
	outOfOrderExternalMoveIdempotence.dfStatus.IncrementTick()

	testConfigs := []struct{
		name string
		moves []Args
		want cacheRow
		v *Visitor
	}{
		{
			name: "TestTrivialIdempotence",
			moves: []Args{
				Args{ EntityID: eid, Tick: 0, Destination: p1, IsExternal: true },
				Args{ EntityID: eid, Tick: 0, Destination: p1, IsExternal: true },
			},
			want: cacheRow{ scheduledTick: 0, destination: p1, isExternal: true },
			v: newVisitor(t),
		},
		{
			name: "TestScheduleDefaultPositionOverride",
			moves: []Args{
				Args{ EntityID: eid, Tick: 0, Destination: p0, IsExternal: true },
			},
			want: cacheRow{ scheduledTick: 0, destination: p0, isExternal: true },
			v: newVisitor(t),
		},
		{
			name: "TestSpamClickIdempotence",
			moves: []Args{
				Args{ EntityID: eid, Tick: 0, Destination: p1, IsExternal: true },
				Args{ EntityID: eid, Tick: 1, Destination: p1, IsExternal: true },
			},
			want: cacheRow{ scheduledTick: 0, destination: p1, isExternal: true },
			v: newVisitor(t),
		},
		{
			name: "TestUpdatePartialMove",
			moves: []Args{
				Args{ EntityID: eid, Tick: 0, Destination: p1, IsExternal: true },
				Args{ EntityID: eid, Tick: 1, Destination: p1, IsExternal: false },
			},
			want: cacheRow{ scheduledTick: 1, destination: p1, isExternal: false },
			v: newVisitor(t),
		},
		{
			name: "TestExternalFutureScheduleNoOp",
			moves: []Args{
				Args{ EntityID: eid, Tick: 1, Destination: p1, IsExternal: true },
			},
			want: cacheRow{ scheduledTick: 0, destination: nil, isExternal: false },
			v: newVisitor(t),
		},
		{
			name: "TestInOrderExternalMoveUpdate",
			moves: []Args{
				Args{ EntityID: eid, Tick: 0, Destination: p1, IsExternal: true },
				Args{ EntityID: eid, Tick: 1, Destination: p2, IsExternal: true },
			},
			want: cacheRow{ scheduledTick: 1, destination: p2, isExternal: true },
			v: inOrderExternalMoveUpdate,
		},
		{
			name: "TestInternalMovePastNoOp",
			moves: []Args{
				Args{ EntityID: eid, Tick: 1, Destination: p1, IsExternal: false },
				Args{ EntityID: eid, Tick: 0, Destination: p2, IsExternal: false },
			},
			want: cacheRow{ scheduledTick: 1, destination: p1, isExternal: false },
			v: newVisitor(t),
		},
		{
			name: "TestInternalMove",
			moves: []Args{
				Args{ EntityID: eid, Tick: 0, Destination: p1, IsExternal: false },
				Args{ EntityID: eid, Tick: 1, Destination: p2, IsExternal: false },
			},
			want: cacheRow{ scheduledTick: 1, destination: p2, isExternal: false },
			v: newVisitor(t),
		},
		{
			name: "TestExternalMovePastPrecedence",
			moves: []Args{
				Args{ EntityID: eid, Tick: 1, Destination: p1, IsExternal: false },
				Args{ EntityID: eid, Tick: 0, Destination: p2, IsExternal: true },
			},
			want: cacheRow{ scheduledTick: 0, destination: p2, isExternal: true },
			v: outOfOrderExternalMoveIdempotence,
		},
		{
			name: "TestExternalMovePastNoOp",
			moves: []Args{
				Args{ EntityID: eid, Tick: 1, Destination: p1, IsExternal: true },
				Args{ EntityID: eid, Tick: 0, Destination: p2, IsExternal: true },
			},
			want: cacheRow{ scheduledTick: 1, destination: p1, isExternal: true },
			v: outOfOrderExternalMoveIdempotence,
		},
		{
			name: "TestExternalMoveSameTickNoOp",
			moves: []Args{
				Args{ EntityID: eid, Tick: 0, Destination: p1, IsExternal: true },
				Args{ EntityID: eid, Tick: 0, Destination: p2, IsExternal: true },
			},
			want: cacheRow{ scheduledTick: 0, destination: p1, isExternal: true },
			v: newVisitor(t),
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			for _, m := range c.moves {
				if err := c.v.Schedule(m); err != nil {
					t.Fatalf("Schedule() = %v, want = nil", err)
				}
			}

			got := getCache(c.v, c.moves[0].EntityID)
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(cacheRow{}),
				protocmp.Transform(),
			); diff != "" {
				t.Errorf("getCache() mismatch(-want +got):\n%v", diff)
			}
		})
	}
}

func TestSchedule(t *testing.T) {
	const nClients = 1000
	v := newVisitor(t)

	var eg errgroup.Group
	for i := 0; i < nClients; i++ {
		i := i
		eg.Go(func() error {
			return v.Schedule(Args{
				Tick: 0,
				EntityID: id.EntityID(fmt.Sprintf("entity-%d", i)),
				Destination: &gdpb.Position{X: 1, Y: 1},
				IsExternal: true,
			})
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("Wait() = %v, want = nil", err)
	}

	if got := len(v.partialCache); got != nClients {
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
	v.Schedule(Args{Tick: t0, EntityID: eid, Destination: dest, IsExternal: true})

	if err := v.Visit(e); err != nil {
		t.Fatalf("Visit() = %v, want = nil", err)
	}

	func(t *testing.T) {
		v.cacheMux.Lock()
		defer v.cacheMux.Unlock()

		got := v.destinationCache[eid]
		if diff := cmp.Diff(
			dest,
			got,
			cmp.AllowUnexported(cacheRow{}),
			cmpopts.IgnoreFields(cacheRow{}, "scheduledTick"),
			protocmp.Transform()); diff != "" {
			t.Fatalf("cache[] mismatch (-want +got):\n%v", diff)
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
