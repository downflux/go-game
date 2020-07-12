// Package cluster implements the clustering logic necessary to build and operate on logical MapTile subsets.
package cluster

import (
	"math"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
	neighborCoordinates = []*rtsspb.Coordinate{
		{X: 0, Y: 1},
		{X: 0, Y: -1},
		{X: 1, Y: 0},
		{X: -1, Y: 0},
	}
)

type ClusterMap struct {
	l int32
	d *rtsspb.Coordinate
	m map[int32]map[int32]*Cluster
}

func ImportClusterMap(pb *rtsspb.ClusterMap) (*ClusterMap, error) {
	return nil, notImplemented
}

func ExportClusterMap(m *ClusterMap) (*rtsspb.ClusterMap, error) {
	return nil, notImplemented
}

type Cluster struct {
	Val *rtsspb.Cluster
}

func ImportCluster(pb *rtsspb.Cluster) (*Cluster, error) {
	return &Cluster{
		Val: pb,
	}, nil
}

func IsAdjacent(c1, c2 *Cluster) bool {
	return math.Abs(float64(c2.Val.GetCoordinate().GetX()-c1.Val.GetCoordinate().GetX()))+math.Abs(float64(c2.Val.GetCoordinate().GetY()-c1.Val.GetCoordinate().GetY())) == 1
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

func BuildClusterMap(tileMapDimension *rtsspb.Coordinate, tileDimension *rtsspb.Coordinate, l int32) (*ClusterMap, error) {
	if l < 1 {
		return nil, status.Errorf(codes.FailedPrecondition, "specified l-level must be a non-zero positive integer")
	}
	m := &ClusterMap{
		l: l,
		d: &rtsspb.Coordinate{},
		m: nil,
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

	m.m = make(map[int32]map[int32]*Cluster)
	m.d.X = int32(math.Ceil(float64(tileMapDimension.GetX()) / float64(tileDimension.GetX())))
	m.d.Y = int32(math.Ceil(float64(tileMapDimension.GetY()) / float64(tileDimension.GetY())))

	for _, xp := range xPartitions {
		x := xp.TileBoundary / tileDimension.GetX()
		m.m[x] = map[int32]*Cluster{}
		for _, yp := range yPartitions {
			y := yp.TileBoundary / tileDimension.GetY()
			m.m[x][y] = &Cluster{
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

type partitionInfo struct {
	TileBoundary  int32
	TileDimension int32
}

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
