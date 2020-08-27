// Package cluster implements the clustering logic necessary to build and
// operate on logical MapTile subsets.
//
// See Botea 2004 for more details.
package cluster

import (
	"math"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/utils"
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

// ClusterMap is a logical abstraction of an underlying TileMap. A TileMap may
// be broken up into separate rectangular partitions, where the cost of a
// partition-partition move is known. This will save cycles when iterating over
// large maps.
type ClusterMap struct {
	Val *rtsspb.ClusterMap
}

func validateClusterInRange(m *ClusterMap, c utils.MapCoordinate) error {
	dim := utils.MapCoordinate{
		X: int32(math.Ceil(
			float64(m.Val.GetTileMapDimension().GetX()) / float64(m.Val.GetTileDimension().GetX()))),
		Y: int32(math.Ceil(
			float64(m.Val.GetTileMapDimension().GetY()) / float64(m.Val.GetTileDimension().GetY()))),
	}

	if 0 < c.X || c.X >= dim.X || 0 < c.Y || c.Y >= dim.X {
		return status.Errorf(codes.OutOfRange, "invalid cluster coordinate %v for ClusterMap", c)
	}
	return nil
}

// ImportClusterMap constructs a ClusterMap object from the given protobuf.
func ImportClusterMap(pb *rtsspb.ClusterMap) (*ClusterMap, error) {
	return &ClusterMap{
		Val: pb,
	}, nil
}

// ExportClusterMap constructs a protobuf from the given ClusterMap object.
func ExportClusterMap(m *ClusterMap) (*rtsspb.ClusterMap, error) {
	return nil, notImplemented
}

// IsAdjacent checks if two clusters are next to each other in the same
// ClusterMap.
func IsAdjacent(m *ClusterMap, c1, c2 utils.MapCoordinate) bool {
	return math.Abs(float64(c2.X-c1.X))+math.Abs(float64(c2.Y-c1.Y)) == 1
}

// TileBoundary returns the starting X, Y coordinates of the specified
// cluster coordinate.
func TileBoundary(m *ClusterMap, c utils.MapCoordinate) (utils.MapCoordinate, error) {
	if err := validateClusterInRange(m, c); err != nil {
		return utils.MapCoordinate{}, err
	}
	return utils.MapCoordinate{
		X: c.X * m.Val.GetTileDimension().GetX(),
		Y: c.Y * m.Val.GetTileDimension().GetY(),
	}, nil
}

// TileDimension calculates the length of the specified cluster coordinate.
func TileDimension(m *ClusterMap, c utils.MapCoordinate) (utils.MapCoordinate, error) {
	if err := validateClusterInRange(m, c); err != nil {
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

// CoordinateInCluster checks if the given coordinate is bounded by the input
// cluster coordinate c.
func CoordinateInCluster(m *ClusterMap, c, t utils.MapCoordinate) bool {
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

// Neighbors returns the adjacent Cluster objects given a Cluster Coordinate.
func Neighbors(m *ClusterMap, c utils.MapCoordinate) ([]utils.MapCoordinate, error) {
	if err := validateClusterInRange(m, c); err != nil {
		return nil, err
	}

	var neighbors []utils.MapCoordinate
	for _, coord := range neighborCoordinates {
		dest := utils.AddMapCoordinate(c, utils.MC(coord))
		if err := validateClusterInRange(m, dest); err == nil {
			neighbors = append(neighbors, dest)
		}
	}
	return neighbors, nil
}

// GetRelativeDirection will return the direction of travel from c to other.
// c and other must be immediately adjacent to one another.
func GetRelativeDirection(m *ClusterMap, c, other utils.MapCoordinate) (rtscpb.Direction, error) {
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

// BuildClusterMap constructs a ClusterMap instance which will be used to
// organize and group Tile objects in the underlying TileMap. ClusterMap does
// not link to the actual Tile -- we need to manually pass the TileMap object
// along when looking up the Tile by a given coordinate.
func BuildClusterMap(tileMapDimension *rtsspb.Coordinate, tileDimension *rtsspb.Coordinate, level int32) (*ClusterMap, error) {
	if level < 1 {
		return nil, status.Error(codes.FailedPrecondition, "level must be a positive non-zero integer")
	}

	return &ClusterMap{
		Val: &rtsspb.ClusterMap{
			Level:            level,
			TileDimension:    tileDimension,
			TileMapDimension: tileMapDimension,
		},
	}, nil
}
