package entrance

import (
	"testing"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/cluster"
	// "github.com/cripplet/rts-pathing/lib/hpf/tile"
	"github.com/golang/protobuf/proto"
	// "github.com/google/go-cmp/cmp"
)

/*
func tileComparator(t, other *tile.Tile) bool {
	return proto.Equal(t.Coordinate(), other.Coordinate())
}

func (s *ClusterBorderSegment) Equal(other *ClusterBorderSegment) bool {
	return proto.Equal(s.s, other.s)
}
*/

func TestCandidateVectorError(t *testing.T) {
	testConfigs := []struct {
		name string
		c    *cluster.Cluster
		d    rtscpb.Direction
	}{
		{name: "NullClusterTest", c: cluster.ImportCluster(&rtsspb.Cluster{}), d: rtscpb.Direction_DIRECTION_NORTH},
		{name: "NullXDimensionClusterTest", c: cluster.ImportCluster(
			&rtsspb.Cluster{
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 0, Y: 1},
			}), d: rtscpb.Direction_DIRECTION_NORTH,
		},
		{name: "NullYDimensionClusterTest", c: cluster.ImportCluster(
			&rtsspb.Cluster{
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 0},
			}), d: rtscpb.Direction_DIRECTION_NORTH,
		},
		{name: "InvalidDirectionTest", c: cluster.ImportCluster(
			&rtsspb.Cluster{
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			}), d: rtscpb.Direction_DIRECTION_UNKNOWN,
		},
	}
	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := candidateVector(c.c, c.d); err == nil {
				t.Errorf("candidateVector() = %v, %v, want a non-nil error", got, err)
			}
		})
	}
}

func TestCandidateVector(t *testing.T) {
	trivialCluster := cluster.ImportCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
		TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
	})
	smallCluster := cluster.ImportCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
		TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
	})
	embeddedCluster := cluster.ImportCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 1},
		TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
	})
	rectangularCluster := cluster.ImportCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 1},
		TileDimension: &rtsspb.Coordinate{X: 1, Y: 2},
	})
	testConfigs := []struct {
		name string
		c    *cluster.Cluster
		d    rtscpb.Direction
		want *rtsspb.ClusterBorderSegment
	}{
		{name: "TrivialClusterNorthTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, End: &rtsspb.Coordinate{X: 0, Y: 0}}},
		{name: "TrivialClusterSouthTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, End: &rtsspb.Coordinate{X: 0, Y: 0}}},
		{name: "TrivialClusterEastTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, End: &rtsspb.Coordinate{X: 0, Y: 0}}},
		{name: "TrivialClusterWestTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, End: &rtsspb.Coordinate{X: 0, Y: 0}}},
		{name: "SmallClusterNorthTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, End: &rtsspb.Coordinate{X: 1, Y: 1}}},
		{name: "SmallClusterSouthTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, End: &rtsspb.Coordinate{X: 1, Y: 0}}},
		{name: "SmallClusterEastTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, End: &rtsspb.Coordinate{X: 1, Y: 1}}},
		{name: "SmallClusterWestTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, End: &rtsspb.Coordinate{X: 0, Y: 1}}},
		{name: "EmbeddedClusterNorthTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, End: &rtsspb.Coordinate{X: 2, Y: 2}}},
		{name: "EmbeddedClusterSouthTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, End: &rtsspb.Coordinate{X: 2, Y: 1}}},
		{name: "EmbeddedClusterEastTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 2, Y: 1}, End: &rtsspb.Coordinate{X: 2, Y: 2}}},
		{name: "EmbeddedClusterWestTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, End: &rtsspb.Coordinate{X: 1, Y: 2}}},
		{name: "RectangularClusterNorthTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, End: &rtsspb.Coordinate{X: 1, Y: 2}}},
		{name: "RectangularClusterSouthTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, End: &rtsspb.Coordinate{X: 1, Y: 1}}},
		{name: "RectangularClusterEastTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, End: &rtsspb.Coordinate{X: 1, Y: 2}}},
		{name: "RectangularClusterWestTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.ClusterBorderSegment{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, End: &rtsspb.Coordinate{X: 1, Y: 2}}},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := candidateVector(c.c, c.d); err != nil || !proto.Equal(got, c.want) {
				t.Errorf("candidateVector() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestBuildCoordinateSlice(t *testing.T) {
}
