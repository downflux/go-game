// Package entrance provides a way to detect contiguous open segments on Cluster borders.
package entrance

import (
	"math"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/cluster"
	"github.com/cripplet/rts-pathing/lib/hpf/tile"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	edgeDirectionToOrientation = map[rtscpb.Direction]rtscpb.Orientation{
		rtscpb.Direction_DIRECTION_NORTH: rtscpb.Orientation_ORIENTATION_HORIZONTAL,
		rtscpb.Direction_DIRECTION_SOUTH: rtscpb.Orientation_ORIENTATION_HORIZONTAL,
		rtscpb.Direction_DIRECTION_EAST: rtscpb.Orientation_ORIENTATION_VERTICAL,
		rtscpb.Direction_DIRECTION_WEST: rtscpb.Orientation_ORIENTATION_VERTICAL,
	}
	reverseDirection = map[rtscpb.Direction]rtscpb.Direction{
		rtscpb.Direction_DIRECTION_NORTH: rtscpb.Direction_DIRECTION_SOUTH,
		rtscpb.Direction_DIRECTION_SOUTH: rtscpb.Direction_DIRECTION_NORTH,
		rtscpb.Direction_DIRECTION_EAST: rtscpb.Direction_DIRECTION_WEST,
		rtscpb.Direction_DIRECTION_WEST: rtscpb.Direction_DIRECTION_EAST,
	}
)

func BuildTransitions(c1, c2 *cluster.Cluster, m *tile.TileMap) ([]*rtsspb.Transition, error) {
	if c1 == nil || c2 == nil {
		return nil, status.Error(codes.FailedPrecondition, "input Cluster references must be non-nil")
	}
	if m == nil {
		return nil, status.Error(codes.FailedPrecondition, "input TileMap reference must be non-nil")
	}

	if !cluster.IsAdjacent(c1, c2) {
		return nil, status.Errorf(codes.FailedPrecondition, "clusters must be immediately adjacent to one another")
	}

	d1, err := cluster.GetRelativeDirection(c1, c2)
	if err != nil {
		return nil, err
	}
	d2 := reverseDirection[d1]

	v1, err := candidateVector(c1, d1)
	if err != nil {
		return nil, err
	}
	v2, err := candidateVector(c2, d2)
	if err != nil {
		return nil, err
	}

	return buildTransitions(v1, v2, m)
}

type orientationInfo struct {
	rank, file, width int32
	orientation rtscpb.Orientation
}

func segmentCoordinateInfo(s *rtsspb.ClusterBorderSegment) (orientationInfo, error) {
	switch s.GetOrientation() {
		case rtscpb.Orientation_ORIENTATION_HORIZONTAL:
			return orientationInfo{
				orientation: s.GetOrientation(),
				rank: s.GetStart().GetY(),
				file: s.GetStart().GetX(),
				width: int32(math.Abs(float64(s.GetEnd().GetX() - s.GetStart().GetX()))),
			}, nil
		case rtscpb.Orientation_ORIENTATION_VERTICAL:
			return orientationInfo{
				orientation: s.GetOrientation(),
				rank: s.GetStart().GetX(),
				file: s.GetStart().GetY(),
				width: int32(math.Abs(float64(s.GetEnd().GetY() - s.GetStart().GetY()))),
			}, nil
		default:
			return orientationInfo{}, status.Errorf(codes.FailedPrecondition, "invalid orientation specified")
	}
}

func buildCoordinateWithSegmentCoordinateInfo(sInfo orientationInfo, offset int32) (*rtsspb.Coordinate, error) {
	switch sInfo.orientation {
		case rtscpb.Orientation_ORIENTATION_HORIZONTAL:
			return &rtsspb.Coordinate{
				X: sInfo.rank,
				Y: sInfo.file + offset,
			}, nil
		case rtscpb.Orientation_ORIENTATION_VERTICAL:
			return &rtsspb.Coordinate{
				X: sInfo.file + offset,
				Y: sInfo.rank,
			}, nil
		default:
			return nil, status.Errorf(codes.FailedPrecondition, "invalid orientation specified")
	}
}

// candidateVector constructs a ClusterBorderSegment instance representing the contiguous edge of a Cluster in the specified direction.
// All Tile t on the edge are between the start and end coordinates, i.e. start <= t <= end with usual 2D coordinate comparison.
func candidateVector(c *cluster.Cluster, d rtscpb.Direction) (*rtsspb.ClusterBorderSegment, error) {
	if c.Cluster().GetTileDimension().GetX() == 0 || c.Cluster().GetTileDimension().GetY() == 0 {
		return nil, status.Error(codes.FailedPrecondition, "input Cluster must have non-zero dimensions")
	}

	var start, end *rtsspb.Coordinate
	switch d {
	case rtscpb.Direction_DIRECTION_NORTH:
		start = &rtsspb.Coordinate{
			X: c.Cluster().GetTileBoundary().GetX(),
			Y: c.Cluster().GetTileBoundary().GetY() + c.Cluster().GetTileDimension().GetY() - 1,
		}
		end = &rtsspb.Coordinate{
			X: c.Cluster().GetTileBoundary().GetX() + c.Cluster().GetTileDimension().GetX() - 1,
			Y: c.Cluster().GetTileBoundary().GetY() + c.Cluster().GetTileDimension().GetY() - 1,
		}
	case rtscpb.Direction_DIRECTION_SOUTH:
		start = proto.Clone(c.Cluster().GetTileBoundary()).(*rtsspb.Coordinate)
		end = &rtsspb.Coordinate{
			X: c.Cluster().GetTileBoundary().GetX() + c.Cluster().GetTileDimension().GetX() - 1,
			Y: c.Cluster().GetTileBoundary().GetY(),
		}
	case rtscpb.Direction_DIRECTION_EAST:
		start = &rtsspb.Coordinate{
			X: c.Cluster().GetTileBoundary().GetX() + c.Cluster().GetTileDimension().GetX() - 1,
			Y: c.Cluster().GetTileBoundary().GetY(),
		}
		end = &rtsspb.Coordinate{
			X: c.Cluster().GetTileBoundary().GetX() + c.Cluster().GetTileDimension().GetX() - 1,
			Y: c.Cluster().GetTileBoundary().GetY() + c.Cluster().GetTileDimension().GetY() - 1,
		}
	case rtscpb.Direction_DIRECTION_WEST:
		start = proto.Clone(c.Cluster().GetTileBoundary()).(*rtsspb.Coordinate)
		end = &rtsspb.Coordinate{
			X: c.Cluster().GetTileBoundary().GetX(),
			Y: c.Cluster().GetTileBoundary().GetY() + c.Cluster().GetTileDimension().GetY() - 1,
		}
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "invalid direction specified %v", d)
	}

	orientation, found := edgeDirectionToOrientation[d]
	if !found {
		return nil, status.Errorf(codes.FailedPrecondition, "invalid direction specified %v", d)
	}

	return &rtsspb.ClusterBorderSegment{
		Orientation: orientation,
		Start: start,
		End: end,
	}, nil
}

func buildTransitionsFromContiguousOpenSegment(s1, s2 *rtsspb.ClusterBorderSegment) ([]*rtsspb.Transition, error) {
	sInfo1, err := segmentCoordinateInfo(s1)
	if err != nil {
		return nil, err
	}
	sInfo2, err := segmentCoordinateInfo(s2)
	if err != nil {
		return nil, err
	}

	var transitions []*rtsspb.Transition
	var offsets []int32
	if sInfo1.width > 3 {
		offsets = append(offsets, 0, sInfo1.width - 1)
	} else {
		offsets = append(offsets, sInfo1.width / 2)
	}

	for _, o := range offsets {
		c1, err := buildCoordinateWithSegmentCoordinateInfo(sInfo1, o)
		if err != nil {
			return nil, err
		}
		c2, err := buildCoordinateWithSegmentCoordinateInfo(sInfo2, o)
		if err != nil {
			return nil, err
		}

		transitions = append(transitions, &rtsspb.Transition{C1: c1, C2: c2})
	}
	return transitions, nil
}

func buildTransitions(v1, v2 *rtsspb.ClusterBorderSegment, m *tile.TileMap) ([]*rtsspb.Transition, error) {
	sInfo1, err := segmentCoordinateInfo(v1)
	if err != nil {
		return nil, err
	}
	sInfo2, err := segmentCoordinateInfo(v2)
	if err != nil {
		return nil, err
	}

	orientation := v1.GetOrientation()
	var res []*rtsspb.Transition

	var s1, s2 *rtsspb.ClusterBorderSegment
	for o := int32(0); o < sInfo1.width; o++ {
		c1, err := buildCoordinateWithSegmentCoordinateInfo(sInfo1, o)
		if err != nil {
			return nil, err
		}
		c2, err := buildCoordinateWithSegmentCoordinateInfo(sInfo2, o)
		if err != nil {
			return nil, err
		}

		if (
			m.TileFromCoordinate(c1).TerrainType() != rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED) && (
			m.TileFromCoordinate(c2).TerrainType() != rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED) {
			if s1 == nil {
				s1 = &rtsspb.ClusterBorderSegment{
					Orientation: orientation,
					Start: c1,
				}
			}
			if s2 == nil {
				s2 = &rtsspb.ClusterBorderSegment{
					Orientation: orientation,
					Start: c2,
				}
			}
			s1.End = c1
			s2.End = c2
		}
		if (
			m.TileFromCoordinate(c1).TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED) || (
			m.TileFromCoordinate(c2).TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED) {
			if s1 != nil && s2 != nil {
				transitions, err := buildTransitionsFromContiguousOpenSegment(s1, s2)
				if err != nil {
					return nil, err
				}
				res = append(res, transitions...)
			}
			s1 = nil
			s2 = nil
		}
	}
	if s1 != nil && s2 != nil {
		transitions, err := buildTransitionsFromContiguousOpenSegment(s1, s2)
		if err != nil {
			return nil, err
		}
		res = append(res, transitions...)
	}

	return res, nil
}
