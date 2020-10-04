package astar

import (
	"fmt"
	"math"
	"testing"

	rtscpb "github.com/minkezhang/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	"github.com/google/go-cmp/cmp"
	"github.com/minkezhang/rts-pathing/lib/hpf/graph"
	"github.com/minkezhang/rts-pathing/lib/hpf/tile"
	"github.com/minkezhang/rts-pathing/lib/hpf/utils"
	"google.golang.org/protobuf/testing/protocmp"
)

func buildTileMap(d utils.MapCoordinate, walls []utils.MapCoordinate) (*tile.Map, error) {
	wallHash := map[utils.MapCoordinate]bool{}
	for _, w := range walls {
		wallHash[w] = true
	}

	var tiles []*rtsspb.Tile
	for x := int32(0); x < d.X; x++ {
		for y := int32(0); y < d.Y; y++ {
			c := utils.MC(&rtsspb.Coordinate{X: x, Y: y})
			var t rtscpb.TerrainType
			if _, found := wallHash[c]; found {
				t = rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED
			} else {
				t = rtscpb.TerrainType_TERRAIN_TYPE_PLAINS
			}
			tiles = append(tiles, &rtsspb.Tile{
				Coordinate:  utils.PB(c),
				TerrainType: t,
			})
		}
	}

	return tile.ImportMap(&rtsspb.TileMap{
		Dimension: utils.PB(d),
		Tiles:     tiles,
		TerrainCosts: []*rtsspb.TerrainCost{
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	})
}

type aStarResult struct {
	path []*tile.Tile
	cost float64
}

func TestPath(t *testing.T) {
	trivialMap, err := buildTileMap(utils.MC(&rtsspb.Coordinate{X: 2, Y: 1}), nil)
	if err != nil {
		t.Fatalf("buildTileMap() = _, %v, want = _, nil", err)
	}
	sourceDestinationGraph, err := graph.BuildGraph(trivialMap, &rtsspb.Coordinate{X: 1, Y: 1})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}
	trivialInterClusterGraph, err := graph.BuildGraph(trivialMap, &rtsspb.Coordinate{X: 1, Y: 1})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}

	intraClusterMap, err := buildTileMap(utils.MC(&rtsspb.Coordinate{X: 6, Y: 1}), nil)
	if err != nil {
		t.Fatalf("buildTileMap() = _, %v, want = _, nil", err)
	}

	trivialIntraClusterGraph, err := graph.BuildGraph(intraClusterMap, &rtsspb.Coordinate{X: 2, Y: 1})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}

	testConfigs := []struct {
		name string
		tm   *tile.Map
		g    *graph.Graph
		src  *rtsspb.Coordinate
		dest *rtsspb.Coordinate
		l    int
		want aStarResult
	}{
		{
			name: "SameSourceDestination",
			tm:   trivialMap,
			g:    sourceDestinationGraph,
			src:  &rtsspb.Coordinate{X: 0, Y: 0},
			dest: &rtsspb.Coordinate{X: 0, Y: 0},
			l:    10,
			want: aStarResult{
				path: []*tile.Tile{
					{
						Val: &rtsspb.Tile{
							Coordinate:  &rtsspb.Coordinate{X: 0, Y: 0},
							TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
				},
				cost: 0,
			},
		},
		{
			name: "TrivialInterClusterPath",
			tm:   trivialMap,
			g:    trivialInterClusterGraph,
			src:  &rtsspb.Coordinate{X: 0, Y: 0},
			dest: &rtsspb.Coordinate{X: 1, Y: 0},
			l:    10,
			want: aStarResult{
				path: []*tile.Tile{
					{
						Val: &rtsspb.Tile{
							Coordinate:  &rtsspb.Coordinate{X: 0, Y: 0},
							TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
					{
						Val: &rtsspb.Tile{
							Coordinate:  &rtsspb.Coordinate{X: 1, Y: 0},
							TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
				},
				cost: 1,
			},
		},
		{
			name: "TrivialIntraClusterPath",
			tm:   intraClusterMap,
			g:    trivialIntraClusterGraph,
			src:  &rtsspb.Coordinate{X: 0, Y: 0},
			dest: &rtsspb.Coordinate{X: 1, Y: 0},
			l:    10,
			want: aStarResult{
				path: []*tile.Tile{
					{
						Val: &rtsspb.Tile{
							Coordinate:  &rtsspb.Coordinate{X: 0, Y: 0},
							TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
					{
						Val: &rtsspb.Tile{
							Coordinate:  &rtsspb.Coordinate{X: 1, Y: 0},
							TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
				},
				cost: 1,
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			path, cost, err := Path(c.tm, c.g, utils.MC(c.src), utils.MC(c.dest), c.l)
			if err != nil {
				t.Fatalf("Path() = _, _, %v, want = _, _, nil", err)
			}

			got := aStarResult{
				path: path,
				cost: cost,
			}

			fmt.Println(got)
			if diff := cmp.Diff(c.want, got, cmp.AllowUnexported(aStarResult{}), protocmp.Transform()); diff != "" {
				t.Errorf("Path() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
