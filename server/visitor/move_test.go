package move

import (
	"math"
	"testing"
	"time"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/fsm/move"
	"github.com/google/go-cmp/cmp"

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

	_ visitor.Visitor = &Visitor{}
)

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

func newTank(t *testing.T, eid id.EntityID, tick id.Tick, p *gdpb.Position) *tank.Entity {
	cid := id.ClientID("client-id")
	tankEntity, err := tank.New(eid, tick, p, cid)
	if err != nil {
		t.Fatalf("New() = %v, want = nil", err)
	}
	return tankEntity
}

func TestVisit(t *testing.T) {
	const eid = "entity-id"
	const t0 = 0
	p0 := &gdpb.Position{X: 0, Y: 0}
	p1 := &gdpb.Position{X: 0, Y: 1}

	testNoMoveVisitor := newVisitor(t)
	testSimpleMoveVisitor := newVisitor(t)

	testConfigs := []struct {
		name string
		v    *Visitor
		i    action.Action
		want []dirty.Curve
	}{
		{
			name: "TestNoMove",
			v:    testNoMoveVisitor,
			i: move.New(
				newTank(t, eid, t0, p0),
				testNoMoveVisitor.status, p0),
			want: nil,
		},
		{
			name: "TestSimpleMove",
			v:    testSimpleMoveVisitor,
			i: move.New(
				newTank(t, eid, t0, p0),
				testSimpleMoveVisitor.status, p1),
			want: []dirty.Curve{
				{EntityID: eid, Property: gcpb.EntityProperty_ENTITY_PROPERTY_POSITION},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if err := c.v.Visit(c.i); err != nil {
				t.Fatalf("Visit() = %v, want = nil", err)
			}
			got := c.v.dirty.Pop().Curves()
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("Pop() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
