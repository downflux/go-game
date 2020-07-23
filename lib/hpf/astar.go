package astar

import (
	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

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
	m                   *tile.TileMap
	boundary, dimension *rtsspb.Coordinate
}

// boundedBy returns true if a <= b < c.
func boundedBy(a, b, c *rtsspb.Coordinate) bool {
	return (a.GetX() <= b.GetX() && a.GetY() <= b.GetY()) && (b.GetX() < c.GetX() && b.GetY() < c.GetY())
}

func (t tileMapGraph) Neighbours(n fastar.Node) []fastar.Node {
	neighbors, _ := t.m.Neighbors(n.(*tile.Tile).Val.GetCoordinate())
	var res []fastar.Node
	for _, n := range neighbors {
		if boundedBy(t.boundary, n.Val.GetCoordinate(), &rtsspb.Coordinate{
			X: t.boundary.GetX() + t.dimension.GetX(),
			Y: t.boundary.GetY() + t.dimension.GetY(),
		},
		) {
			res = append(res, n)
		}
	}
	return res
}

func TileMapPath(m *tile.TileMap, src, dest *tile.Tile, boundary, dimension *rtsspb.Coordinate) ([]*tile.Tile, float64, error) {
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
	nodes := fastar.FindPath(tileMapGraph{m: m, boundary: boundary, dimension: dimension}, src, dest, d, tileMapH)

	var tiles []*tile.Tile
	for _, node := range nodes {
		tiles = append(tiles, node.(*tile.Tile))
	}
	return tiles, nodes.Cost(d), nil
}
