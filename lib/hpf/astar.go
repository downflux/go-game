package astar

import (
	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/tile"
	fastar "github.com/fzipp/astar"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func tileMapD(c map[rtscpb.TerrainType]float64, src, dest fastar.Node) float64 {
	cost, _ := tile.D(c, src.(*tile.Tile), dest.(*tile.Tile))
	return cost
}

func tileMapH(src, dest fastar.Node) float64 {
	cost, _ := tile.H(src.(*tile.Tile), dest.(*tile.Tile))
	return cost
}

type tileMapGraph struct {
	m *tile.TileMap
}

func (t tileMapGraph) Neighbours(n fastar.Node) []fastar.Node {
	neighbors, _ := t.m.Neighbors(n.(*tile.Tile).Val.GetCoordinate())

	var res []fastar.Node
	for _, t := range neighbors {
		res = append(res, t)
	}
	return res
}

func TileMapPath(m *tile.TileMap, src, dest *tile.Tile) ([]*tile.Tile, float64, error) {
	if m == nil {
		return nil, 0, status.Errorf(codes.FailedPrecondition, "cannot have nil TileMap input")
	}
	if src == nil || dest == nil {
		return nil, 0, status.Errorf(codes.FailedPrecondition, "cannot have nil Tile inputs")
	}
	if src.TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED || dest.TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED {
		return nil, 0, nil
	}

	d := func(a, b fastar.Node) float64 {
		return tileMapD(m.C, a, b)
	}
	nodes := fastar.FindPath(tileMapGraph{m: m}, src, dest, d, tileMapH)

	var tiles []*tile.Tile
	for _, node := range nodes {
		tiles = append(tiles, node.(*tile.Tile))
	}
	return tiles, nodes.Cost(d), nil
}
