// Package tileastar defines fzipp.astar.Graph implementations for tile.Map.
package tileastar

import (
	"math"

	"github.com/downflux/game/map/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gdpb "github.com/downflux/game/api/data_go_proto"
	mcpb "github.com/downflux/game/map/api/constants_go_proto"
	tile "github.com/downflux/game/map/map"
	fastar "github.com/fzipp/astar"
)

// dFunc provides a shim for the tile.Map neighbor distance function.
func dFunc(c map[mcpb.TerrainType]float64, src, dest fastar.Node) float64 {
	cost, err := tile.D(c, src.(*tile.Tile), dest.(*tile.Tile))
	if err != nil {
		return math.Inf(0)
	}
	return cost
}

// hFunc provides a shim for the tile.Map heuristic function.
func hFunc(src, dest fastar.Node) float64 {
	cost, err := tile.H(src.(*tile.Tile), dest.(*tile.Tile))
	if err != nil {
		return math.Inf(0)
	}
	return cost
}

// graphImpl implements fzipp.astar.Graph for the tile.Map struct.
type graphImpl struct {
	// m holds a reference to the underlying terrain map.
	m *tile.Map

	// boundary represents the (inclusive) lower-bound of the bounding box
	// used to constrain the path search.
	boundary *gdpb.Coordinate

	// dimension represents the (exclusive) upper-bound of the bounding box
	// used to constrain the path search.
	dimension *gdpb.Coordinate
}

// boundedBy returns true if points a <= b < c. Here, a < b is true in a 2D
// graph if point a is down and to the left of b, and is only a partial
// ordering. Specifically, it is not normal lexicographical order.
func boundedBy(a, b, c *gdpb.Coordinate) bool {
	return (a.GetX() <= b.GetX() && a.GetY() <= b.GetY()) && (b.GetX() < c.GetX() && b.GetY() < c.GetY())
}

// Neighbours returns neighboring Tile objects from a tile.Map.
func (t graphImpl) Neighbours(n fastar.Node) []fastar.Node {
	neighbors, _ := t.m.Neighbors(n.(*tile.Tile).Val.GetCoordinate())
	var res []fastar.Node
	for _, n := range neighbors {
		if boundedBy(t.boundary, n.Val.GetCoordinate(), &gdpb.Coordinate{
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
func Path(m *tile.Map, src, dest utils.MapCoordinate, boundary, dimension *gdpb.Coordinate) ([]*tile.Tile, float64, error) {
	if m == nil {
		return nil, 0, status.Errorf(codes.FailedPrecondition, "cannot have nil tile.Map input")
	}

	tSrc := m.TileFromCoordinate(utils.PB(src))
	tDest := m.TileFromCoordinate(utils.PB(dest))
	if tSrc == nil || tDest == nil {
		return nil, 0, status.Errorf(codes.NotFound, "a Tile cannot be found with the input coordinates")
	}
	if math.IsInf(m.C[tSrc.TerrainType()], 0) || math.IsInf(m.C[tDest.TerrainType()], 0) {
		return nil, math.Inf(0), nil
	}

	d := func(a, b fastar.Node) float64 {
		return dFunc(m.C, a, b)
	}
	nodes := fastar.FindPath(graphImpl{m: m, boundary: boundary, dimension: dimension}, tSrc, tDest, d, hFunc)

	var tiles []*tile.Tile
	for _, node := range nodes {
		tiles = append(tiles, node.(*tile.Tile))
	}

	var cost float64
	if tiles == nil {
		cost = math.Inf(0)
	} else {
		cost = nodes.Cost(d)
	}
	return tiles, cost, nil
}
