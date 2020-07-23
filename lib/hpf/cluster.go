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
	"github.com/golang/protobuf/proto"
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
	// L specifies the abstaction level of the ClusterMap -- higher
	// abstraction Clusters contain groups of lower level Clusters.
	L int32

	// D specifies the number of Cluster objects extending in each spatial
	// dimension.
	D *rtsspb.Coordinate

	// M contains the list of Clusters in a ClusterMap. The key here
	// is relative to other Clusters -- it does not represent any TileMap
	// Coordinates.
	M map[utils.MapCoordinate]*Cluster
}

// ImportClusterMap constructs a ClusterMap object from the given protobuf.
func ImportClusterMap(pb *rtsspb.ClusterMap) (*ClusterMap, error) {
	cm := &ClusterMap{
		L: pb.GetLevel(),
		D: pb.GetDimension(),
	}
	for _, c := range pb.GetClusters() {
		cluster, err := ImportCluster(c)
		if err != nil {
			return nil, err
		}
		cm.M[utils.MC(c.GetCoordinate())] = cluster
	}

	return cm, nil
}

// ExportClusterMap constructs a protobuf from the given ClusterMap object.
func ExportClusterMap(m *ClusterMap) (*rtsspb.ClusterMap, error) {
	return nil, notImplemented
}

// Cluster encapsulates a group of Tile objects.
type Cluster struct {
	// Val is the underlying
	Val *rtsspb.Cluster
}

// ImportCluster constructs a Cluster object from the given protobuf.
func ImportCluster(pb *rtsspb.Cluster) (*Cluster, error) {
	return &Cluster{
		Val: proto.Clone(pb).(*rtsspb.Cluster),
	}, nil
}

// ExportCluster constructs a protobuf from the given Cluster object.
func ExportCluster(c *Cluster) (*rtsspb.Cluster, error) {
	return proto.Clone(c.Val).(*rtsspb.Cluster), nil
}

// IsAdjacent checks if two Cluster objects are next to each other in the same ClusterMap.
// TODO(cripplet): Check if we need an l-level in Cluster proto -- if so, we should check that here as well.
func IsAdjacent(c1, c2 *Cluster) bool {
	return math.Abs(float64(c2.Val.GetCoordinate().GetX()-c1.Val.GetCoordinate().GetX()))+math.Abs(float64(c2.Val.GetCoordinate().GetY()-c1.Val.GetCoordinate().GetY())) == 1
}

// Cluster returns the Cluster object from the input coordinates.
func (m *ClusterMap) Cluster(x, y int32) *Cluster {
	return m.M[utils.MapCoordinate{X: x, Y: y}]
}

// Neighbors returns the adjacent Cluster objects given a Cluster Coordinate.
func (m *ClusterMap) Neighbors(coordinate *rtsspb.Coordinate) ([]*Cluster, error) {
	src, found := m.M[utils.MC(coordinate)]
	if !found {
		return nil, status.Error(codes.NotFound, "no Cluster exists with given coordinates in the ClusterMap")
	}

	var neighbors []*Cluster
	for _, c := range neighborCoordinates {
		if dest := m.M[utils.MC(&rtsspb.Coordinate{
			X: src.Val.GetCoordinate().GetX() + c.GetX(),
			Y: src.Val.GetCoordinate().GetY() + c.GetY(),
		})]; dest != nil {
			neighbors = append(neighbors, dest)
		}
	}
	return neighbors, nil
}

// GetRelativeDirection will return the direction of travel from c to other.
// c and other must be immediately adjacent to one another.
func GetRelativeDirection(c, other *Cluster) (rtscpb.Direction, error) {
	if !IsAdjacent(c, other) {
		return rtscpb.Direction_DIRECTION_UNKNOWN, status.Errorf(codes.FailedPrecondition, "input clusters are not immediately adjacent to one another")
	}

	if c.Val.GetCoordinate().GetX() == other.Val.GetCoordinate().GetX() && c.Val.GetCoordinate().GetY() < other.Val.GetCoordinate().GetY() {
		return rtscpb.Direction_DIRECTION_NORTH, nil
	}
	if c.Val.GetCoordinate().GetX() == other.Val.GetCoordinate().GetX() && c.Val.GetCoordinate().GetY() > other.Val.GetCoordinate().GetY() {
		return rtscpb.Direction_DIRECTION_SOUTH, nil
	}
	if c.Val.GetCoordinate().GetX() < other.Val.GetCoordinate().GetX() && c.Val.GetCoordinate().GetY() == other.Val.GetCoordinate().GetY() {
		return rtscpb.Direction_DIRECTION_EAST, nil
	}
	if c.Val.GetCoordinate().GetX() > other.Val.GetCoordinate().GetX() && c.Val.GetCoordinate().GetY() == other.Val.GetCoordinate().GetY() {
		return rtscpb.Direction_DIRECTION_WEST, nil
	}
	return rtscpb.Direction_DIRECTION_UNKNOWN, status.Errorf(codes.FailedPrecondition, "clusters which are immediately adjacent are somehow not traversible via cardinal directions")
}

// BuildClusterMap constructs a ClusterMap instance which will be used to organize and group Tile objects in the
// underlying TileMap. ClusterMap does not link to the actual Tile -- we need to manually pass the TileMap object along
// when looking up the Tile by a given coordinate.
func BuildClusterMap(tileMapDimension *rtsspb.Coordinate, tileDimension *rtsspb.Coordinate, level int32) (*ClusterMap, error) {
	if level < 1 {
		return nil, status.Error(codes.FailedPrecondition, "level must be a positive non-zero integer")
	}

	m := &ClusterMap{
		L: level,
		D: &rtsspb.Coordinate{},
		M: nil,
	}

	xPartitions, err := partition(tileMapDimension.GetX(), tileDimension.GetX())
	if err != nil {
		return nil, err
	}
	yPartitions, err := partition(tileMapDimension.GetY(), tileDimension.GetY())
	if err != nil {
		return nil, err
	}

	if xPartitions == nil || yPartitions == nil {
		return m, nil
	}

	m.M = make(map[utils.MapCoordinate]*Cluster)
	m.D.X = int32(math.Ceil(float64(tileMapDimension.GetX()) / float64(tileDimension.GetX())))
	m.D.Y = int32(math.Ceil(float64(tileMapDimension.GetY()) / float64(tileDimension.GetY())))

	for _, xp := range xPartitions {
		x := xp.TileBoundary / tileDimension.GetX()
		for _, yp := range yPartitions {
			y := yp.TileBoundary / tileDimension.GetY()
			m.M[utils.MC(&rtsspb.Coordinate{X: x, Y: y})] = &Cluster{
				Val: &rtsspb.Cluster{
					Coordinate:    &rtsspb.Coordinate{X: x, Y: y},
					TileBoundary:  &rtsspb.Coordinate{X: xp.TileBoundary, Y: yp.TileBoundary},
					TileDimension: &rtsspb.Coordinate{X: xp.TileDimension, Y: yp.TileDimension},
				},
			}
		}
	}

	return m, nil
}

// partitionInfo is a 1D partition data struct, representing a semi-open
// interval of Tile Coordinates, projected onto a specific axis.
type partitionInfo struct {
	// TileBoundary is the acceptable lower bound of a Tile; Tile objects
	// may include this projected Coordinate.
	TileBoundary int32

	// TileDimension is the length of the partition interval. The upper
	// bound as defined by the TileDimension is open (exclusive).
	TileDimension int32
}

// partition builds a 1D list of partitions -- we will combine the X-specific and Y-specific partitions into
// a 2D partition array.
func partition(tileMapDimension int32, tileDimension int32) ([]partitionInfo, error) {
	if tileDimension == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "invalid tileDimension value %v", tileDimension)
	}
	var partitions []partitionInfo

	for x := int32(0); x*tileDimension < tileMapDimension; x++ {
		minX := x * tileDimension
		maxX := int32(math.Min(
			float64((x+1)*tileDimension-1), float64(tileMapDimension-1)))

		partitions = append(partitions, partitionInfo{
			TileBoundary:  minX,
			TileDimension: maxX - minX + 1,
		})
	}
	return partitions, nil
}
