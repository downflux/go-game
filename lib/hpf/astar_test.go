package astar

import (
	"fmt"
	"math"
	"testing"

	rtscpb "github.com/minkezhang/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	"github.com/minkezhang/rts-pathing/lib/hpf/graph"
	"github.com/minkezhang/rts-pathing/lib/hpf/tile"
	"github.com/minkezhang/rts-pathing/lib/hpf/utils"
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
				Coordinate: utils.PB(c),
				TerrainType: t,
			})
		}
	}

	return tile.ImportMap(&rtsspb.TileMap{
		Dimension: utils.PB(d),
		Tiles: tiles,
		TerrainCosts: []*rtsspb.TerrainCost{
			&rtsspb.TerrainCost{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
			&rtsspb.TerrainCost{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	})
}

func TestPath(t *testing.T) {
	tm, err := buildTileMap(utils.MC(&rtsspb.Coordinate{X: 1, Y: 2}), nil)
	if err != nil {
		t.Fatalf("buildTileMap() = _, %v, want = _, nil", err)
	}

	g, err := graph.BuildGraph(tm, &rtsspb.Coordinate{X: 1, Y: 1})
	path, cost, err := Path(tm, g, utils.MC(&rtsspb.Coordinate{X: 0, Y: 0}), utils.MC(&rtsspb.Coordinate{X: 0, Y: 0}), 10)
	fmt.Println(path, cost, err)
}
