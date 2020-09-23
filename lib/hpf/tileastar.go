// Package tileastar defines fzipp.astar.Graph implementations for tile.Map.
package tileastar

import (
	rtscpb "github.com/minkezhang/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	"github.com/minkezhang/rts-pathing/lib/hpf/tile"
	fastar "github.com/fzipp/astar"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// dFunc provides a shim for the tile.Map neighbor distance function.
func dFunc(c map[rtscpb.TerrainType]float64, src, dest fastar.Node) float64 {
	cost, _ := tile.D(c, src.(*tile.Tile), dest.(*tile.Tile))
	return cost
}

// hFunc provides a shim for the tile.Map heuristic function.
func hFunc(src, dest fastar.Node) float64 {
	cost, _ := tile.H(src.(*tile.Tile), dest.(*tile.Tile))
	return cost
}

// graph implements fzipp.astar.Graph for the tile.Map struct.
type graph struct {
	m                   *tile.Map
	boundary, dimension *rtsspb.Coordinate
}

// boundedBy returns true if points a <= b < c. Here, a < b is true in a 2D
// graph if point a is down and to the left of b, and is only a partial
// ordering. Specifically, it is not normal lexicographical order.
func boundedBy(a, b, c *rtsspb.Coordinate) bool {
	return (a.GetX() <= b.GetX() && a.GetY() <= b.GetY()) && (b.GetX() < c.GetX() && b.GetY() < c.GetY())
}

// Neighbours returns neighboring Tile objects from a tile.Map.
func (t graph) Neighbours(n fastar.Node) []fastar.Node {
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

// Path returns pathing information for two Tile objects embedded in a
// tile.Map. This function returns a (path, cost, error) tuple, where path is a
// list of Tile objects and cost is the actual cost as calculated by calling D
// over the returned path. An empty path indicates there is no path found
// between the two Tile objects.
//
// The user needs to additionally supply the bounding box within the tile.Map in
// which to search for a path; if the entire tile.Map should be considered, the
// bounding box as defined by the tile.Map should be used here. The lower bound
// of the bounding box is defined as the boundary Coordinate, and the size of
// the box is specified by the dimension Coordinate.
func Path(m *tile.Map, src, dest *tile.Tile, boundary, dimension *rtsspb.Coordinate) ([]*tile.Tile, float64, error) {
	if m == nil {
		return nil, 0, status.Errorf(codes.FailedPrecondition, "cannot have nil Map input")
	}
	if src == nil || dest == nil {
		return nil, 0, status.Errorf(codes.FailedPrecondition, "cannot have nil Tile inputs")
	}
	if src.TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED || dest.TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED {
		return nil, 0, nil
	}

	d := func(a, b fastar.Node) float64 {
		return dFunc(m.C, a, b)
	}
	nodes := fastar.FindPath(graph{m: m, boundary: boundary, dimension: dimension}, src, dest, d, hFunc)

	var tiles []*tile.Tile
	for _, node := range nodes {
		tiles = append(tiles, node.(*tile.Tile))
	}
	return tiles, nodes.Cost(d), nil
}
