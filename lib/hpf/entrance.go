// Package entrance provides a way to detect contiguous open segments on
// Cluster borders.
package entrance

import (
	"math"

	rtscpb "github.com/downflux/pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/downflux/pathing/lib/proto/structs_go_proto"

	"github.com/golang/protobuf/proto"
	"github.com/downflux/pathing/lib/hpf/cluster"
	"github.com/downflux/pathing/lib/hpf/tile"
	"github.com/downflux/pathing/lib/hpf/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// MaxSingleGapWidth represents the distance between closed Tile
	// objects which will need a single Transition node -- gaps wider than
	// this width should be represented by two transition nodes instead,
	// per Botea 2004.
	MaxSingleGapWidth = 3
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
func BuildTransitions(tm *tile.Map, cm *cluster.Map, c1, c2 utils.MapCoordinate) ([]*rtsspb.Transition, error) {
	if tm == nil {
		return nil, status.Error(codes.FailedPrecondition, "input tile.Map reference must be non-nil")
	}
	if cm == nil || cm.Val == nil {
		return nil, status.Error(codes.FailedPrecondition, "input cluster.Map reference must be non-nil")
	}
	if !cluster.IsAdjacent(cm, c1, c2) {
		return nil, status.Errorf(codes.FailedPrecondition, "clusters must be immediately adjacent to one another")
	}
	if !proto.Equal(tm.D, cm.Val.GetTileMapDimension()) {
		return nil, status.Errorf(codes.FailedPrecondition, "tile.Map and cluster.Map dimensions do not agree")
	}
	for _, clusterCoord := range []utils.MapCoordinate{c1, c2} {
		if err := cluster.ValidateClusterInRange(cm, clusterCoord); err != nil {
			return nil, err
		}
	}

	d1, err := cluster.GetRelativeDirection(cm, c1, c2)
	if err != nil {
		return nil, err
	}
	d2 := reverseDirection[d1]

	s1, err := buildClusterEdgeCoordinateSlice(cm, c1, d1)
	if err != nil {
		return nil, err
	}
	s2, err := buildClusterEdgeCoordinateSlice(cm, c2, d2)
	if err != nil {
		return nil, err
	}

	return buildTransitionsAux(tm, s1, s2)
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
func sliceContains(s *rtsspb.CoordinateSlice, t utils.MapCoordinate) (bool, error) {
	switch s.GetOrientation() {
	case rtscpb.Orientation_ORIENTATION_HORIZONTAL:
		return (t.Y == s.GetStart().GetY()) && (s.GetStart().GetX() <= t.X) && (t.X < s.GetStart().GetX()+s.GetLength()), nil
	case rtscpb.Orientation_ORIENTATION_VERTICAL:
		return (t.X == s.GetStart().GetX()) && (s.GetStart().GetY() <= t.Y) && (t.Y < s.GetStart().GetY()+s.GetLength()), nil
	default:
		return false, status.Errorf(codes.FailedPrecondition, "invalid slice orientation %v", s.GetOrientation())
	}
}

// OnClusterEdge checks if the given coordinate coord falls on the edge of a
// cluster coordinate.
func OnClusterEdge(m *cluster.Map, clusterCoord utils.MapCoordinate, coord utils.MapCoordinate) bool {
	if err := cluster.ValidateClusterInRange(m, clusterCoord); err != nil {
		return false
	}
	if _, err := cluster.ClusterCoordinateFromTileCoordinate(m, coord); err != nil {
		return false
	}

	for _, d := range []rtscpb.Direction{
		rtscpb.Direction_DIRECTION_NORTH,
		rtscpb.Direction_DIRECTION_SOUTH,
		rtscpb.Direction_DIRECTION_EAST,
		rtscpb.Direction_DIRECTION_WEST,
	} {
		slice, err := buildClusterEdgeCoordinateSlice(m, clusterCoord, d)
		if err != nil {
			return false
		}

		res, err := sliceContains(slice, coord)
		if err != nil {
			return false
		}

		if res {
			return true
		}
	}
	return false
}

// buildClusterEdgeCoordinateSlice constructs a CoordinateSlice instance
// representing the contiguous edge of a Cluster in the specified direction.
// All Tile t on the edge are between the start and end coordinates,
// i.e. start <= t <= end with usual 2D coordinate comparison.
func buildClusterEdgeCoordinateSlice(m *cluster.Map, c utils.MapCoordinate, d rtscpb.Direction) (*rtsspb.CoordinateSlice, error) {
	var start *rtsspb.Coordinate
	var length int32

	tileBoundary, err := cluster.TileBoundary(m, c)
	if err != nil {
		return nil, err
	}
	tileDimension, err := cluster.TileDimension(m, c)
	if err != nil {
		return nil, err
	}

	switch d {
	case rtscpb.Direction_DIRECTION_NORTH:
		start = &rtsspb.Coordinate{
			X: tileBoundary.X,
			Y: tileBoundary.Y + tileDimension.Y - 1,
		}
	case rtscpb.Direction_DIRECTION_SOUTH:
		start = &rtsspb.Coordinate{
			X: tileBoundary.X,
			Y: tileBoundary.Y,
		}
	case rtscpb.Direction_DIRECTION_EAST:
		start = &rtsspb.Coordinate{
			X: tileBoundary.X + tileDimension.X - 1,
			Y: tileBoundary.Y,
		}
	case rtscpb.Direction_DIRECTION_WEST:
		start = &rtsspb.Coordinate{
			X: tileBoundary.X,
			Y: tileBoundary.Y,
		}
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "invalid direction specified %v", d)
	}

	orientation := edgeDirectionToOrientation[d]
	switch orientation {
	case rtscpb.Orientation_ORIENTATION_HORIZONTAL:
		length = tileDimension.X
	case rtscpb.Orientation_ORIENTATION_VERTICAL:
		length = tileDimension.Y
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "invalid orientation specified %v", orientation)
	}

	return &rtsspb.CoordinateSlice{
		Orientation: orientation,
		Start:       start,
		Length:      length,
	}, nil
}

// buildTransitionsFromOpenCoordinateSlice constructs the actual transition
// points between two contiguous open tile slices. We have configured, as per
// Botea 2004, one transition node for a segment of width three tiles or less,
// and two transition nodes for segments longer than three tiles.
//
// In general, the less nodes we have, the faster the hierarchical part of the
// pathing algorithm will take, which would be the case if we increase
// MaxSingleGapWidth. We may also consider a more contextual reworking
// of this function and take into consideration the nearest transition node
// from adjacent slices, e.g. "transition nodes must be N tiles apart".
func buildTransitionsFromOpenCoordinateSlice(s1, s2 *rtsspb.CoordinateSlice) ([]*rtsspb.Transition, error) {
	if err := verifyCoordinateSlices(s1, s2); err != nil {
		return nil, err
	}

	var transitions []*rtsspb.Transition
	var offsets []int32
	if s1.GetLength() <= MaxSingleGapWidth {
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
				TileCoordinate: c1,
			},
			N2: &rtsspb.AbstractNode{
				TileCoordinate: c2,
			},
		})
	}
	return transitions, nil
}

// buildTransitionsAux constructs the list of Transition nodes given the
// corresponding edges of two adjacent Cluster objects.
func buildTransitionsAux(m *tile.Map, s1, s2 *rtsspb.CoordinateSlice) ([]*rtsspb.Transition, error) {
	if err := verifyCoordinateSlices(s1, s2); err != nil {
		return nil, err
	}

	orientation := s1.GetOrientation()
	var res []*rtsspb.Transition

	var tSegment1, tSegment2 *rtsspb.CoordinateSlice
	for o := int32(0); o < s1.GetLength(); o++ {
		t1, err := buildCoordinateWithCoordinateSlice(s1, o)
		if err != nil {
			return nil, err
		}
		t2, err := buildCoordinateWithCoordinateSlice(s2, o)
		if err != nil {
			return nil, err
		}

		if m.C[m.TileFromCoordinate(t1).TerrainType()] < math.Inf(0) && m.C[m.TileFromCoordinate(t2).TerrainType()] < math.Inf(0) {
			if tSegment1 == nil {
				tSegment1 = &rtsspb.CoordinateSlice{
					Orientation: orientation,
					Start:       t1,
				}
			}
			if tSegment2 == nil {
				tSegment2 = &rtsspb.CoordinateSlice{
					Orientation: orientation,
					Start:       t2,
				}
			}
			tSegment1.Length += 1
			tSegment2.Length += 1
		}
		if math.IsInf(m.C[m.TileFromCoordinate(t1).TerrainType()], 0) || math.IsInf(m.C[m.TileFromCoordinate(t2).TerrainType()], 0) {
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
