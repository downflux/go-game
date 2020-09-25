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
	 * Y = 0 -
	 *   X = 0
	 */
	trivialOpenMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 1, Y: 1},
		Tiles: []*rtsspb.Tile{
			{
				Coordinate:  &rtsspb.Coordinate{X: 0, Y: 0},
				TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
			},
		},
		TerrainCosts: []*rtsspb.TerrainCost{},
	}
)

type aStarResult struct {
	path []*rtsspb.AbstractNode
	cost float64
}

func TestFoo(t *testing.T) {
	testConfigs := []struct {
		name      string
		tm        *rtsspb.TileMap
		g         *graph.Graph
		src, dest *rtsspb.AbstractNode
		want      aStarResult
	}{}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tm, err := tile.ImportMap(c.tm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			nodes, cost, err := Path(tm, c.g, c.src, c.dest)
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
