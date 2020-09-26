package graphastar

import (
	"math"
	"testing"

	rtscpb "github.com/minkezhang/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	"github.com/google/go-cmp/cmp"
	"github.com/minkezhang/rts-pathing/lib/hpf/graph"
	"github.com/minkezhang/rts-pathing/lib/hpf/tile"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	/**
	 * A graph.Graph object relies on building Transition objects, which
	 * relies on more than one cluster being formed by BuildGraph.
	 *
	 * Y = 0 - -
	 *   X = 0
	 */
	trivialOpenMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 2, Y: 1},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
		TerrainCosts: []*rtsspb.TerrainCost{
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
		},
	}

	/*
	 * Y = 0 W W
	 *   X = 0
	 */
	trivialClosedMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 2, Y: 1},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
		},
		TerrainCosts: []*rtsspb.TerrainCost{
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
		},
	}

	/*
	 * Y = 0 W -
	 *   X = 0
	 */
	trivialSemiOpenMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 2, Y: 1},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
		TerrainCosts: []*rtsspb.TerrainCost{
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
		},
	}

	/*
	 * Y = 0 - W -
	 *   X = 0
	 */
	simpleBlockedMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 3, Y: 1},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
		TerrainCosts: []*rtsspb.TerrainCost{
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
		},
	}

	/*
	 * Y = 0 - - - - - -
	 *   X = 0
	 */
	simpleLongOpenMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 6, Y: 1},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 3, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 4, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 5, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
		TerrainCosts: []*rtsspb.TerrainCost{
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
		},
	}
)

type buildGraphInput struct {
	tileDimension *rtsspb.Coordinate
	level         int32
}

type aStarResult struct {
	path []*rtsspb.AbstractNode
	cost float64
}

func TestPath(t *testing.T) {
	testConfigs := []struct {
		name      string
		tm        *rtsspb.TileMap
		g         buildGraphInput
		src, dest *rtsspb.AbstractNode
		want      aStarResult
	}{
		{
			name: "TrivialReachablePath",
			tm:   trivialOpenMap,
			g:    buildGraphInput{tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, level: 1},
			src:  &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
			dest: &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
			want: aStarResult{
				path: []*rtsspb.AbstractNode{
					{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
				},
			},
		},
		{
			name: "TrivialIntraClusterPath",
			tm:   simpleLongOpenMap,
			g:    buildGraphInput{tileDimension: &rtsspb.Coordinate{X: 2, Y: 1}, level: 1},
			src:  &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 2, Y: 0}},
			dest: &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 3, Y: 0}},
			want: aStarResult{
				path: []*rtsspb.AbstractNode{
					{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 2, Y: 0}},
					{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 3, Y: 0}},
				},
				cost: 1,
			},
		},
		{
			name: "TrivialInterClusterPath",
			tm:   trivialOpenMap,
			g:    buildGraphInput{tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, level: 1},
			src:  &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
			dest: &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
			want: aStarResult{
				path: []*rtsspb.AbstractNode{
					{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
				},
				cost: 1,
			},
		},
		{
			name: "TrivialClosedPath",
			tm:   trivialClosedMap,
			g:    buildGraphInput{tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, level: 1},
			src:  &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
			dest: &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
			want: aStarResult{
				path: nil,
				cost: math.Inf(0),
			},
		},
		{
			name: "TrivialSemiOpenPath",
			tm:   trivialClosedMap,
			g:    buildGraphInput{tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, level: 1},
			src:  &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
			dest: &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
			want: aStarResult{
				path: nil,
				cost: math.Inf(0),
			},
		},
		{
			name: "SimpleBlockedPath",
			tm:   simpleBlockedMap,
			g:    buildGraphInput{tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, level: 1},
			src:  &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
			dest: &rtsspb.AbstractNode{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 2, Y: 0}},
			want: aStarResult{
				path: nil,
				cost: math.Inf(0),
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tm, err := tile.ImportMap(c.tm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			g, err := graph.BuildGraph(tm, c.g.tileDimension, c.g.level)
			if err != nil {
				t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
			}

			nodes, cost, err := Path(tm, g, c.src, c.dest)
			if err != nil {
				t.Fatalf("Path() = %v, want = nil", err)
			}

			got := aStarResult{
				path: nodes,
				cost: cost,
			}
			if diff := cmp.Diff(c.want, got, cmp.AllowUnexported(aStarResult{}), protocmp.Transform()); diff != "" {
				t.Errorf("Path() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
