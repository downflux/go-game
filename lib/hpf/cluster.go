// Package cluster implements the clustering logic necessary to build and
// operate on logical tile.Map subsets.
//
// See Botea 2004 for more details.
package cluster

import (
	"math"

	rtscpb "github.com/downflux/pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/downflux/pathing/lib/proto/structs_go_proto"

	"github.com/downflux/pathing/lib/hpf/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")

	// neighborCoordinates provides the Coordinate deltas between a
	// specific Coordinate and adjacent Coordinates to expand to in
	// a graph search.
	neighborCoordinates = []*rtsspb.Coordinate{
		{X: 0, Y: 1},
		{X: 0, Y: -1},
		{X: 1, Y: 0},
		{X: -1, Y: 0},
	}
)

// Map is a logical abstraction of an underlying tile.Map. A tile.Map may
// be broken up into separate rectangular partitions, where the cost of a
// partition-partition move is known. This will save cycles when iterating over
// large maps.
type Map struct {
	Val *rtsspb.ClusterMap
}

// Dimension returns the X, Y dimension of a given Map.
func Dimension(m *Map) *rtsspb.Coordinate {
	return &rtsspb.Coordinate{
		X: int32(math.Ceil(
			float64(m.Val.GetTileMapDimension().GetX()) / float64(m.Val.GetTileDimension().GetX()))),
		Y: int32(math.Ceil(
			float64(m.Val.GetTileMapDimension().GetY()) / float64(m.Val.GetTileDimension().GetY()))),
	}
}

// ValidateClusterInRange checks if the given MapCoordinate is a valid cluster
// coordinate for the given Map.
func ValidateClusterInRange(m *Map, c utils.MapCoordinate) error {
	dim := utils.MC(Dimension(m))
	if 0 > c.X || c.X >= dim.X || 0 > c.Y || c.Y >= dim.Y {
		return status.Errorf(codes.OutOfRange, "invalid cluster coordinate %v for Map", c)
	}
	return nil
}

// ClusterCoordinateFromTileCoordinate gets the corresponding cluster to which
// an input Tile coordinate belongs.
func ClusterCoordinateFromTileCoordinate(m *Map, t utils.MapCoordinate) (utils.MapCoordinate, error) {
	if t.X >= m.Val.GetTileMapDimension().GetX() || t.Y >= m.Val.GetTileMapDimension().GetY() {
		return utils.MapCoordinate{}, status.Errorf(codes.OutOfRange, "input Tile coordinate %v is outside the map boundary %v", t, m.Val.GetTileMapDimension())
	}

	return utils.MapCoordinate{
		X: t.X / m.Val.GetTileDimension().GetX(),
		Y: t.Y / m.Val.GetTileDimension().GetY(),
	}, nil
}

// ImportMap constructs a Map object from the given protobuf.
func ImportMap(pb *rtsspb.ClusterMap) (*Map, error) {
	return &Map{
		Val: pb,
	}, nil
}

// ExportMap constructs a protobuf from the given Map object.
func ExportMap(m *Map) (*rtsspb.ClusterMap, error) {
	return nil, notImplemented
}

// IsAdjacent checks if two clusters are next to each other in the same
// Map.
func IsAdjacent(m *Map, c1, c2 utils.MapCoordinate) bool {
	return math.Abs(float64(c2.X-c1.X))+math.Abs(float64(c2.Y-c1.Y)) == 1
}

// TileBoundary returns the starting X, Y coordinates of the specified
// cluster coordinate.
func TileBoundary(m *Map, c utils.MapCoordinate) (utils.MapCoordinate, error) {
	if err := ValidateClusterInRange(m, c); err != nil {
		return utils.MapCoordinate{}, err
	}
	return utils.MapCoordinate{
		X: c.X * m.Val.GetTileDimension().GetX(),
		Y: c.Y * m.Val.GetTileDimension().GetY(),
	}, nil
}

// TileDimension calculates the length of the specified cluster coordinate.
func TileDimension(m *Map, c utils.MapCoordinate) (utils.MapCoordinate, error) {
	if err := ValidateClusterInRange(m, c); err != nil {
		return utils.MapCoordinate{}, err
	}
	return utils.MapCoordinate{
		X: int32(math.Min(
			float64(m.Val.GetTileDimension().GetX()),
			float64(m.Val.GetTileMapDimension().GetX()-c.X*m.Val.GetTileDimension().GetX()))),
		Y: int32(math.Min(
			float64(m.Val.GetTileDimension().GetY()),
			float64(m.Val.GetTileMapDimension().GetY()-c.Y*m.Val.GetTileDimension().GetY()))),
	}, nil
}

// Iterator provides a flattened list of MapCoordinate instances representing
// cluster coordinates.
func Iterator(m *Map) []utils.MapCoordinate {
	var cs []utils.MapCoordinate
	for x := int32(0); x < Dimension(m).GetX(); x++ {
		for y := int32(0); y < Dimension(m).GetY(); y++ {
			cs = append(cs, utils.MC(&rtsspb.Coordinate{X: x, Y: y}))
		}
	}
	return cs
}

// CoordinateInCluster checks if the given coordinate is bounded by the input
// cluster coordinate c.
func CoordinateInCluster(m *Map, c, t utils.MapCoordinate) bool {
	tileBoundary, err := TileBoundary(m, c)
	if err != nil {
		return false
	}

	tileDimension, err := TileDimension(m, c)
	if err != nil {
		return false
	}

	return (tileBoundary.X <= t.X && t.X < tileBoundary.X+tileDimension.X) && (tileBoundary.Y <= t.Y && t.Y < tileBoundary.Y+tileDimension.Y)
}

// Neighbors returns the adjacent cluster coordinates given a cluster
// coordinate.
func Neighbors(m *Map, c utils.MapCoordinate) ([]utils.MapCoordinate, error) {
	if err := ValidateClusterInRange(m, c); err != nil {
		return nil, err
	}

	var neighbors []utils.MapCoordinate
	for _, coord := range neighborCoordinates {
		dest := utils.AddMapCoordinate(c, utils.MC(coord))
		if err := ValidateClusterInRange(m, dest); err == nil {
			neighbors = append(neighbors, dest)
		}
	}
	return neighbors, nil
}

// GetRelativeDirection will return the direction of travel from c to other.
// c and other must be immediately adjacent to one another.
func GetRelativeDirection(m *Map, c, other utils.MapCoordinate) (rtscpb.Direction, error) {
	if !IsAdjacent(m, c, other) {
		return rtscpb.Direction_DIRECTION_UNKNOWN, status.Errorf(codes.FailedPrecondition, "input clusters are not immediately adjacent to one another")
	}

	if c.X == other.X && c.Y < other.Y {
		return rtscpb.Direction_DIRECTION_NORTH, nil
	}
	if c.X == other.X && c.Y > other.Y {
		return rtscpb.Direction_DIRECTION_SOUTH, nil
	}
	if c.X < other.X && c.Y == other.Y {
		return rtscpb.Direction_DIRECTION_EAST, nil
	}
	if c.X > other.X && c.Y == other.Y {
		return rtscpb.Direction_DIRECTION_WEST, nil
	}
	return rtscpb.Direction_DIRECTION_UNKNOWN, status.Errorf(codes.FailedPrecondition, "clusters which are immediately adjacent are somehow not traversible via cardinal directions")
}

// BuildMap constructs a Map instance which will be used to
// organize and group Tile objects in the underlying tile.Map. Map does
// not link to the actual Tile -- we need to manually pass the tile.Map object
// along when looking up the Tile by a given coordinate.
func BuildMap(tileMapDimension *rtsspb.Coordinate, tileDimension *rtsspb.Coordinate) (*Map, error) {
	return &Map{
		Val: &rtsspb.ClusterMap{
			TileDimension:    tileDimension,
			TileMapDimension: tileMapDimension,
		},
	}, nil
}
