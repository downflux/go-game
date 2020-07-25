// Package entrance provides a way to detect contiguous open segments on
// Cluster borders.
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
	// edgeDirectionToOrientation indicates the orientatino of a Cluster
	// edge.
	edgeDirectionToOrientation = map[rtscpb.Direction]rtscpb.Orientation{
		rtscpb.Direction_DIRECTION_NORTH: rtscpb.Orientation_ORIENTATION_HORIZONTAL,
		rtscpb.Direction_DIRECTION_SOUTH: rtscpb.Orientation_ORIENTATION_HORIZONTAL,
		rtscpb.Direction_DIRECTION_EAST:  rtscpb.Orientation_ORIENTATION_VERTICAL,
		rtscpb.Direction_DIRECTION_WEST:  rtscpb.Orientation_ORIENTATION_VERTICAL,
	}

	// reverseDirection transforms a cardinal direction into its
	// complement.
	reverseDirection = map[rtscpb.Direction]rtscpb.Direction{
		rtscpb.Direction_DIRECTION_NORTH: rtscpb.Direction_DIRECTION_SOUTH,
		rtscpb.Direction_DIRECTION_SOUTH: rtscpb.Direction_DIRECTION_NORTH,
		rtscpb.Direction_DIRECTION_EAST:  rtscpb.Direction_DIRECTION_WEST,
		rtscpb.Direction_DIRECTION_WEST:  rtscpb.Direction_DIRECTION_EAST,
	}
)

// BuildTransitions takes in two adjacent map Cluster objects and returns the
// list of Transition nodes which connect them. Transition nodes are a tuple of
// non-blocking Tiles which are
//   1. immediately adjacent to one another, and
//   2. are in different Clusters.
// See Botea 2004 for more information.
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

	v1, err := buildClusterEdgeCoordinateSlice(c1, d1)
	if err != nil {
		return nil, err
	}
	v2, err := buildClusterEdgeCoordinateSlice(c2, d2)
	if err != nil {
		return nil, err
	}

	return buildTransitionsAux(v1, v2, m)
}

// buildCoordinateWithCoordinateSlice reconstructs the Coordinate object back
// from the given slice info.
func buildCoordinateWithCoordinateSlice(s *rtsspb.CoordinateSlice, offset int32) (*rtsspb.Coordinate, error) {
	if offset < 0 || offset >= s.GetLength() {
		return nil, status.Errorf(codes.FailedPrecondition, "invalid offset specified, end coordinate must be contained within the slice")
	}
	switch s.GetOrientation() {
	case rtscpb.Orientation_ORIENTATION_HORIZONTAL:
		return &rtsspb.Coordinate{
			X: s.GetStart().GetX() + offset,
			Y: s.GetStart().GetY(),
		}, nil
	case rtscpb.Orientation_ORIENTATION_VERTICAL:
		return &rtsspb.Coordinate{
			X: s.GetStart().GetX(),
			Y: s.GetStart().GetY() + offset,
		}, nil
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "invalid orientation specified")
	}
}

// sliceContains checks if the given Coordinate falls within the slice.
func sliceContains(s *rtsspb.CoordinateSlice, coord *rtsspb.Coordinate) (bool, error) {
	switch s.GetOrientation() {
	case rtscpb.Orientation_ORIENTATION_HORIZONTAL:
		return coord.GetY() == s.GetStart().GetY() && s.GetStart().GetX() <= coord.GetX() && coord.GetX() < s.GetStart().GetX()+s.GetLength(), nil
	case rtscpb.Orientation_ORIENTATION_VERTICAL:
		return coord.GetX() == s.GetStart().GetX() && s.GetStart().GetY() <= coord.GetY() && coord.GetY() < s.GetStart().GetY()+s.GetLength(), nil
	default:
		return false, status.Errorf(codes.FailedPrecondition, "invalid slice orientation %v", s.GetOrientation())
	}
}

// OnClusterEdge checks if the given Coordinate falls on the edge of a Cluster.
func OnClusterEdge(c *cluster.Cluster, coord *rtsspb.Coordinate) (bool, error) {
	for _, d := range []rtscpb.Direction{
		rtscpb.Direction_DIRECTION_NORTH,
		rtscpb.Direction_DIRECTION_SOUTH,
		rtscpb.Direction_DIRECTION_EAST,
		rtscpb.Direction_DIRECTION_WEST,
	} {
		slice, err := buildClusterEdgeCoordinateSlice(c, d)
		if err != nil {
			return false, err
		}

		res, err := sliceContains(slice, coord)
		if err != nil {
			return false, err
		}

		if res {
			return true, nil
		}
	}
	return false, nil
}

// buildClusterEdgeCoordinateSlice constructs a CoordinateSlice instance
// representing the contiguous edge of a Cluster in the specified direction.
// All Tile t on the edge are between the start and end coordinates,
// i.e. start <= t <= end with usual 2D coordinate comparison.
func buildClusterEdgeCoordinateSlice(c *cluster.Cluster, d rtscpb.Direction) (*rtsspb.CoordinateSlice, error) {
	if c.Val.GetTileDimension().GetX() == 0 || c.Val.GetTileDimension().GetY() == 0 {
		return nil, status.Error(codes.FailedPrecondition, "input Cluster must have non-zero dimensions")
	}

	var start *rtsspb.Coordinate
	var length int32

	switch d {
	case rtscpb.Direction_DIRECTION_NORTH:
		start = &rtsspb.Coordinate{
			X: c.Val.GetTileBoundary().GetX(),
			Y: c.Val.GetTileBoundary().GetY() + c.Val.GetTileDimension().GetY() - 1,
		}
	case rtscpb.Direction_DIRECTION_SOUTH:
		start = proto.Clone(c.Val.GetTileBoundary()).(*rtsspb.Coordinate)
	case rtscpb.Direction_DIRECTION_EAST:
		start = &rtsspb.Coordinate{
			X: c.Val.GetTileBoundary().GetX() + c.Val.GetTileDimension().GetX() - 1,
			Y: c.Val.GetTileBoundary().GetY(),
		}
	case rtscpb.Direction_DIRECTION_WEST:
		start = proto.Clone(c.Val.GetTileBoundary()).(*rtsspb.Coordinate)
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "invalid direction specified %v", d)
	}

	orientation := edgeDirectionToOrientation[d]
	switch orientation {
	case rtscpb.Orientation_ORIENTATION_HORIZONTAL:
		length = c.Val.GetTileDimension().GetX()
	case rtscpb.Orientation_ORIENTATION_VERTICAL:
		length = c.Val.GetTileDimension().GetY()
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "invalid orientation specified %v", orientation)
	}

	return &rtsspb.CoordinateSlice{
		Orientation:       orientation,
		Start:             start,
		Length:            length,
		ClusterCoordinate: proto.Clone(c.Val.GetCoordinate()).(*rtsspb.Coordinate),
	}, nil
}

// buildTransitionsFromOpenCoordinateSlice constructs the actual transition
// points between two contiguous open tile slices. We have configured, as per
// Botea 2004, one transition node for a segment of width three tiles or less,
// and two transition nodes for segments longer than three tiles.
//
// In general, the less nodes we have, the faster the hierarchical part of the
// pathing algorithm will take, which would be the case if we increase
// minLength. We may also consider a more contextual reworking of this function
// and take into consideration the nearest transition node from adjacent
// slices, e.g. "transition nodes must be N tiles apart".
func buildTransitionsFromOpenCoordinateSlice(s1, s2 *rtsspb.CoordinateSlice) ([]*rtsspb.Transition, error) {
	const minLength = 3

	if err := verifyCoordinateSlices(s1, s2); err != nil {
		return nil, err
	}

	var transitions []*rtsspb.Transition
	var offsets []int32
	if s1.GetLength() <= minLength {
		offsets = append(offsets, s1.GetLength()/2)
	} else {
		offsets = append(offsets, 0, s1.GetLength()-1)
	}

	for _, o := range offsets {
		c1, err := buildCoordinateWithCoordinateSlice(s1, o)
		if err != nil {
			return nil, err
		}
		c2, err := buildCoordinateWithCoordinateSlice(s2, o)
		if err != nil {
			return nil, err
		}

		transitions = append(transitions, &rtsspb.Transition{
			N1: &rtsspb.AbstractNode{
				TileCoordinate:    c1,
				ClusterCoordinate: s1.GetClusterCoordinate(),
			},
			N2: &rtsspb.AbstractNode{
				TileCoordinate:    c2,
				ClusterCoordinate: s2.GetClusterCoordinate(),
			},
		})
	}
	return transitions, nil
}

// buildTransitionsAux constructs the list of Transition nodes given the
// corresponding edges of two adjacent Cluster objects.
func buildTransitionsAux(s1, s2 *rtsspb.CoordinateSlice, m *tile.TileMap) ([]*rtsspb.Transition, error) {
	if err := verifyCoordinateSlices(s1, s2); err != nil {
		return nil, err
	}

	orientation := s1.GetOrientation()
	var res []*rtsspb.Transition

	var tSegment1, tSegment2 *rtsspb.CoordinateSlice
	for o := int32(0); o < s1.GetLength(); o++ {
		c1, err := buildCoordinateWithCoordinateSlice(s1, o)
		if err != nil {
			return nil, err
		}
		c2, err := buildCoordinateWithCoordinateSlice(s2, o)
		if err != nil {
			return nil, err
		}

		if (m.TileFromCoordinate(c1).TerrainType() != rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED) && (m.TileFromCoordinate(c2).TerrainType() != rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED) {
			if tSegment1 == nil {
				tSegment1 = &rtsspb.CoordinateSlice{
					Orientation:       orientation,
					Start:             c1,
					ClusterCoordinate: s1.GetClusterCoordinate(),
				}
			}
			if tSegment2 == nil {
				tSegment2 = &rtsspb.CoordinateSlice{
					Orientation:       orientation,
					Start:             c2,
					ClusterCoordinate: s2.GetClusterCoordinate(),
				}
			}
			tSegment1.Length += 1
			tSegment2.Length += 1
		}
		if (m.TileFromCoordinate(c1).TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED) || (m.TileFromCoordinate(c2).TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED) {
			if tSegment1 != nil && tSegment2 != nil {
				transitions, err := buildTransitionsFromOpenCoordinateSlice(tSegment1, tSegment2)
				if err != nil {
					return nil, err
				}
				res = append(res, transitions...)
			}

			tSegment1 = nil
			tSegment2 = nil
		}
	}
	if tSegment1 != nil && tSegment2 != nil {
		transitions, err := buildTransitionsFromOpenCoordinateSlice(tSegment1, tSegment2)
		if err != nil {
			return nil, err
		}
		res = append(res, transitions...)
	}

	return res, nil
}

// verifyCoordinateSlices ensures our input slices meet some basic criteria,
// e.g. adjacent, same orientation, etc. If we need to optimize, we can skip
// this step, as it's only called in internal functions that are not exposed to
// the end user.
func verifyCoordinateSlices(s1, s2 *rtsspb.CoordinateSlice) error {
	if s1.GetOrientation() != s2.GetOrientation() || s1.GetLength() != s2.GetLength() {
		return status.Error(codes.FailedPrecondition, "input CoordinateSlice instances mismatch")
	}

	switch s1.GetOrientation() {
	case rtscpb.Orientation_ORIENTATION_HORIZONTAL:
		if s1.GetStart().GetX() != s2.GetStart().GetX() || math.Abs(float64(s2.GetStart().GetY()-s1.GetStart().GetY())) != 1 {
			return status.Error(codes.FailedPrecondition, "input CoordinateSlice instances mismatch")
		}
	case rtscpb.Orientation_ORIENTATION_VERTICAL:
		if s1.GetStart().GetY() != s2.GetStart().GetY() || math.Abs(float64(s2.GetStart().GetX()-s1.GetStart().GetX())) != 1 {
			return status.Error(codes.FailedPrecondition, "input CoordinateSlice instances mismatch")
		}
	}
	return nil
}
