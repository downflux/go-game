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

type ClusterBorderSegment struct {
        s *rtsspb.ClusterBorderSegment
}

func BuildClusterBorderSegments(c *cluster.Cluster, m *tile.TileMap, d rtscpb.Direction) ([]*ClusterBorderSegment, error) {
	start, end, err := candidateVector(c, d)
	if err != nil {
		return nil, err
	}

	candidates, err := candidateVectorTiles(start, end, m)
	if err != nil {
		return nil, err
	}

	return segments(candidates)
}

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

// candidateVectorTiles return the Tile objects caught in between the start and end range Coordinates.
// for all tiles returned, start <= t <= end, per usual 2D coordinate comparison. Tiles returned are
// ordered.
func candidateVectorTiles(start, end *rtsspb.Coordinate, m *tile.TileMap) ([]*tile.Tile, error) {
	var candidates []*tile.Tile
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

func segments(candidates []*tile.Tile) ([]*ClusterBorderSegment, error) {
	var segments []*ClusterBorderSegment
	var segment *ClusterBorderSegment

	for _, t := range candidates {
		if t.TerrainType() != rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED {
			if segment == nil {
				segment = &ClusterBorderSegment{
					s: &rtsspb.ClusterBorderSegment{
						Source: t.Coordinate(),
					},
				}
			}
			segment.s.Destination = t.Coordinate()
		} else {
			segments = append(segments, segment)
			segment = nil
		}
	}
	return segments, nil
}
