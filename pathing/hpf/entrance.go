// Package entrance provides a way to detect contiguous open segments on
// Cluster borders.
package entrance

import (
	"math"
	"reflect"

	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/pathing/hpf/cluster"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gdpb "github.com/downflux/game/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	pcpb "github.com/downflux/game/pathing/api/constants_go_proto"
	pdpb "github.com/downflux/game/pathing/api/data_go_proto"
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
	edgeDirectionToOrientation = map[pcpb.Direction]pcpb.Orientation{
		pcpb.Direction_DIRECTION_NORTH: pcpb.Orientation_ORIENTATION_HORIZONTAL,
		pcpb.Direction_DIRECTION_SOUTH: pcpb.Orientation_ORIENTATION_HORIZONTAL,
		pcpb.Direction_DIRECTION_EAST:  pcpb.Orientation_ORIENTATION_VERTICAL,
		pcpb.Direction_DIRECTION_WEST:  pcpb.Orientation_ORIENTATION_VERTICAL,
	}

	// reverseDirection transforms a cardinal direction into its
	// complement.
	reverseDirection = map[pcpb.Direction]pcpb.Direction{
		pcpb.Direction_DIRECTION_NORTH: pcpb.Direction_DIRECTION_SOUTH,
		pcpb.Direction_DIRECTION_SOUTH: pcpb.Direction_DIRECTION_NORTH,
		pcpb.Direction_DIRECTION_EAST:  pcpb.Direction_DIRECTION_WEST,
		pcpb.Direction_DIRECTION_WEST:  pcpb.Direction_DIRECTION_EAST,
	}
)

type Transition struct {
	N1, N2 *pdpb.AbstractNode
}

type coordinateSlice struct {
	Orientation pcpb.Orientation
	Start       *gdpb.Coordinate
	Length      int32
}

// BuildTransitions takes in two adjacent map Cluster objects and returns the
// list of Transition nodes which connect them. Transition nodes are a tuple of
// non-blocking Tiles which are
//   1. immediately adjacent to one another, and
//   2. are in different Clusters.
// See Botea 2004 for more information.
func BuildTransitions(tm *tile.Map, cm *cluster.Map, c1, c2 utils.MapCoordinate) ([]Transition, error) {
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
func buildCoordinateWithCoordinateSlice(s coordinateSlice, offset int32) (*gdpb.Coordinate, error) {
	if offset < 0 || offset >= s.Length {
		return nil, status.Errorf(codes.FailedPrecondition, "invalid offset specified, end coordinate must be contained within the slice")
	}
	switch s.Orientation {
	case pcpb.Orientation_ORIENTATION_HORIZONTAL:
		return &gdpb.Coordinate{
			X: s.Start.GetX() + offset,
			Y: s.Start.GetY(),
		}, nil
	case pcpb.Orientation_ORIENTATION_VERTICAL:
		return &gdpb.Coordinate{
			X: s.Start.GetX(),
			Y: s.Start.GetY() + offset,
		}, nil
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "invalid orientation specified")
	}
}

// sliceContains checks if the given Coordinate falls within the slice.
func sliceContains(s coordinateSlice, t utils.MapCoordinate) (bool, error) {
	switch s.Orientation {
	case pcpb.Orientation_ORIENTATION_HORIZONTAL:
		return (t.Y == s.Start.GetY()) && (s.Start.GetX() <= t.X) && (t.X < s.Start.GetX()+s.Length), nil
	case pcpb.Orientation_ORIENTATION_VERTICAL:
		return (t.X == s.Start.GetX()) && (s.Start.GetY() <= t.Y) && (t.Y < s.Start.GetY()+s.Length), nil
	default:
		return false, status.Errorf(codes.FailedPrecondition, "invalid slice orientation %v", s.Orientation)
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

	for _, d := range []pcpb.Direction{
		pcpb.Direction_DIRECTION_NORTH,
		pcpb.Direction_DIRECTION_SOUTH,
		pcpb.Direction_DIRECTION_EAST,
		pcpb.Direction_DIRECTION_WEST,
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
func buildClusterEdgeCoordinateSlice(m *cluster.Map, c utils.MapCoordinate, d pcpb.Direction) (coordinateSlice, error) {
	var start *gdpb.Coordinate
	var length int32

	tileBoundary, err := cluster.TileBoundary(m, c)
	if err != nil {
		return coordinateSlice{}, err
	}
	tileDimension, err := cluster.TileDimension(m, c)
	if err != nil {
		return coordinateSlice{}, err
	}

	switch d {
	case pcpb.Direction_DIRECTION_NORTH:
		start = &gdpb.Coordinate{
			X: tileBoundary.X,
			Y: tileBoundary.Y + tileDimension.Y - 1,
		}
	case pcpb.Direction_DIRECTION_SOUTH:
		start = &gdpb.Coordinate{
			X: tileBoundary.X,
			Y: tileBoundary.Y,
		}
	case pcpb.Direction_DIRECTION_EAST:
		start = &gdpb.Coordinate{
			X: tileBoundary.X + tileDimension.X - 1,
			Y: tileBoundary.Y,
		}
	case pcpb.Direction_DIRECTION_WEST:
		start = &gdpb.Coordinate{
			X: tileBoundary.X,
			Y: tileBoundary.Y,
		}
	default:
		return coordinateSlice{}, status.Errorf(codes.FailedPrecondition, "invalid direction specified %v", d)
	}

	orientation := edgeDirectionToOrientation[d]
	switch orientation {
	case pcpb.Orientation_ORIENTATION_HORIZONTAL:
		length = tileDimension.X
	case pcpb.Orientation_ORIENTATION_VERTICAL:
		length = tileDimension.Y
	default:
		return coordinateSlice{}, status.Errorf(codes.FailedPrecondition, "invalid orientation specified %v", orientation)
	}

	return coordinateSlice{
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
func buildTransitionsFromOpenCoordinateSlice(s1, s2 coordinateSlice) ([]Transition, error) {
	if err := verifyCoordinateSlices(s1, s2); err != nil {
		return nil, err
	}

	var transitions []Transition
	var offsets []int32
	if s1.Length <= MaxSingleGapWidth {
		offsets = append(offsets, s1.Length/2)
	} else {
		offsets = append(offsets, 0, s1.Length-1)
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

		transitions = append(transitions, Transition{
			N1: &pdpb.AbstractNode{
				TileCoordinate: c1,
			},
			N2: &pdpb.AbstractNode{
				TileCoordinate: c2,
			},
		})
	}
	return transitions, nil
}

// buildTransitionsAux constructs the list of Transition nodes given the
// corresponding edges of two adjacent Cluster objects.
func buildTransitionsAux(m *tile.Map, s1, s2 coordinateSlice) ([]Transition, error) {
	if err := verifyCoordinateSlices(s1, s2); err != nil {
		return nil, err
	}

	orientation := s1.Orientation
	var res []Transition

	var tSegment1, tSegment2 coordinateSlice
	for o := int32(0); o < s1.Length; o++ {
		t1, err := buildCoordinateWithCoordinateSlice(s1, o)
		if err != nil {
			return nil, err
		}
		t2, err := buildCoordinateWithCoordinateSlice(s2, o)
		if err != nil {
			return nil, err
		}

		if m.C[m.TileFromCoordinate(t1).TerrainType()] < math.Inf(0) && m.C[m.TileFromCoordinate(t2).TerrainType()] < math.Inf(0) {
			if reflect.ValueOf(tSegment1).IsZero() {
				tSegment1 = coordinateSlice{
					Orientation: orientation,
					Start:       t1,
				}
			}
			if reflect.ValueOf(tSegment2).IsZero() {
				tSegment2 = coordinateSlice{
					Orientation: orientation,
					Start:       t2,
				}
			}
			tSegment1.Length += 1
			tSegment2.Length += 1
		}
		if math.IsInf(m.C[m.TileFromCoordinate(t1).TerrainType()], 0) || math.IsInf(m.C[m.TileFromCoordinate(t2).TerrainType()], 0) {
			if !reflect.ValueOf(tSegment1).IsZero() && !reflect.ValueOf(tSegment2).IsZero() {
				transitions, err := buildTransitionsFromOpenCoordinateSlice(tSegment1, tSegment2)
				if err != nil {
					return nil, err
				}
				res = append(res, transitions...)
			}

			tSegment1 = coordinateSlice{}
			tSegment2 = coordinateSlice{}
		}
	}
	if !reflect.ValueOf(tSegment1).IsZero() && !reflect.ValueOf(tSegment2).IsZero() {
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
func verifyCoordinateSlices(s1, s2 coordinateSlice) error {
	if s1.Orientation != s2.Orientation || s1.Length != s2.Length {
		return status.Error(codes.FailedPrecondition, "input CoordinateSlice instances mismatch")
	}

	switch s1.Orientation {
	case pcpb.Orientation_ORIENTATION_HORIZONTAL:
		if s1.Start.GetX() != s2.Start.GetX() || math.Abs(float64(s2.Start.GetY()-s1.Start.GetY())) != 1 {
			return status.Error(codes.FailedPrecondition, "input CoordinateSlice instances mismatch")
		}
	case pcpb.Orientation_ORIENTATION_VERTICAL:
		if s1.Start.GetY() != s2.Start.GetY() || math.Abs(float64(s2.Start.GetX()-s1.Start.GetX())) != 1 {
			return status.Error(codes.FailedPrecondition, "input CoordinateSlice instances mismatch")
		}
	}
	return nil
}
