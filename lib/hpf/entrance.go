// Package entrance provides a way to detect contiguous open segments on Cluster borders.
package entrance

import (
	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/cluster"
	"github.com/cripplet/rts-pathing/lib/hpf/tile"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	reverseDirection := map[rtscpb.Direction]rtscpb.Direction{
		rtscpb.Direction_DIRECTION_NORTH: rtscpb.Direction_DIRECTION_SOUTH,
		rtscpb.Direction_DIRECTION_SOUTH: rtscpb.Direction_DIRECTION_NORTH,
		rtscpb.Direction_DIRECTION_EAST: rtscpb.Direction_DIRECTION_WEST,
		rtscpb.Direction_DIRECTION_WEST: rtscpb.Direction_DIRECTION_EAST,
	}
	directionVector := map[rtscpb.Direction]*rtsspb.Coordinate{
		rtscpb.Direction_DIRECTION_NORTH: &rtsspb.Coordinate{X: 0, Y: 1},
		rtscpb.Direction_DIRECTION_SOUTH: &rtsspb.Coordinate{X: 0,Y: -1},
		rtscpb.Direction_DIRECTION_EAST: &rtsspb.Coordinate{X: 1,Y: 0},
		rtscpb.Direction_DIRECTION_WEST: &rtsspb.Coordinate{X: -1,Y: 0},
	}
)

// ClusterBorderSegment encapsulates a specific contiguous open segment within a cluster border.
// All coordinates within the segment are contained within the cluster itself.
type ClusterBorderSegment struct {
	s *rtsspb.ClusterBorderSegment
}

// BuildClusterBorderSegments constructs a list of ClusterBorderSegment instances based on the input Cluter and TileMap objects.
func BuildClusterBorderSegments(c1, c2 *cluster.Cluster, m *tile.TileMap) ([]*ClusterBorderSegment, []*ClusterBorderSegment, error) {
	d1, err := cluster.GetRelativeDirection(c1, c2)
	if err != nil {
		return nil, nil, err
	}
	d2, found := reverseDirection[d1]
	if !found {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "reverse vector of %v not found", d1)
	}

	start1, end1, err := candidateVector(c1, d1)
	if err != nil {
		return nil, nil, err
	}
	start2, end2, err := candidateVector(c2, d2)
	if err != nil {
		return nil, nil, err
	}

	candidates1, candidates2, err := candidateVectorTiles(start1, end1, start2, end2, m)
	if err != nil {
		return nil, err
	}

	return segments(candidates)
}

// candidateVector constructs a (start, end) Coordinate pair representing the edge of a Cluster in the specified direction.
// All Tile t on the edge are between the start and end coordinates, i.e. start <= t <= end with usual 2D coordinate comparison.
func candidateVector(c *cluster.Cluster, d rtscpb.Direction) (*rtsspb.Coordinate, *rtsspb.Coordinate, error) {
	if c.Cluster().GetTileDimension().GetX() == 0 || c.Cluster().GetTileDimension().GetY() == 0 {
		return nil, nil, status.Error(codes.FailedPrecondition, "input cluster must have non-zero dimensions")
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
		return nil, nil, status.Errorf(codes.FailedPrecondition, "invalid direction specified %v", d)
	}
	return start, end, nil
}

func containerBox(start, end *rtsspb.Coordinate) *rtsspb.Coordinate {
	return &rtsspb.Coordinate{
		X: int32(math.Abs(float64(end.GetX() - start.GetX()))),
		X: int32(math.Abs(float64(end.GetY() - start.GetY()))),
	}
}
func coordinateAdd(c1, c2) *rtsspb.Coordinate{
	return &rtsspb.Coordinate{
		X: c1.GetX() + c2.GetX(),
		Y: c2.GetY() + c2.GetY(),
	}
}

// candidateVectorTiles return the Tile objects caught in between the start and end range Coordinates.
// for all tiles returned, start <= t <= end, per usual 2D coordinate comparison. Tiles returned are
// ordered.
func candidateVectorTiles(start1, end1, start2, end2 *rtsspb.Coordinate, d rtscpb.Direction, m *tile.TileMap) ([]*tile.Tile, []*tile.Tile, error) {
	if box1, box2 := containerBox(start1, end1), containerBox(start2, end2); !proto.Equal(box1, box2) {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "edge tiles shape is different: %v != %v", box1, box2)
	}

	var candidates1, candidates2 []*tile.Tile
	for x := start.GetX(); x <= end.GetX(); x++ {
		for y := start.GetY(); y <= end.GetY(); y++ {
			t := m.Tile(x, y)
			if t == nil {
				return nil, status.Errorf(codes.FailedPrecondition, "Tile (%v, %v) not found in underlying TileMap", x, y)
			}
			candidates = append(candidates, t)
		}
	}
	return candidates, nil
}

// segments constructs the list of contiguous open segments based on the input Tile instances.
// segments expect that the list of Tile are sorted in 2D coordinate comparison order,
// i.e. A < B iff A.Y < B.Y || A.Y == B.Y && A.X < B.X. This is trivially true in this package
// since we're only dealing with a single row or column of tiles for all expected input.
func segments(candidates []*tile.Tile) ([]*ClusterBorderSegment, error) {
	var segments []*ClusterBorderSegment
	var segment *ClusterBorderSegment

	for _, t := range candidates {
		switch t.TerrainType() {
		case rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED:
			if segment != nil {
				segments = append(segments, segment)
				segment = nil
			}
		default:
			if segment == nil {
				segment = &ClusterBorderSegment{
					s: &rtsspb.ClusterBorderSegment{Start: t.Coordinate()},
				}
			}
			segment.s.End = t.Coordinate()
		}
	}
	if segment != nil {
		segments = append(segments, segment)
	}
	return segments, nil
}
