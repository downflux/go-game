package graphastar

import (
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
			{
				Coordinate:  &rtsspb.Coordinate{X: 0, Y: 0},
				TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
			},
			{
				Coordinate:  &rtsspb.Coordinate{X: 1, Y: 0},
				TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
			},
		},
		TerrainCosts: []*rtsspb.TerrainCost{},
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

func TestFoo(t *testing.T) {
	testConfigs := []struct {
		name      string
		tm        *rtsspb.TileMap
		g         buildGraphInput
		src, dest *rtsspb.AbstractNode
		want      aStarResult
	}{
		{
			// TODO(minkezhang): Decide if src and dest should
			// be passed in by reference instead (i.e. by
			// referencing an actual underlying AbstractNode)
			// instead of constructing a new reference.
			//
			// We probably should.
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
