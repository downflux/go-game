package cluster

import (
	"testing"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/utils"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestIsAdjacent(t *testing.T) {
	testConfigs := []struct {
		name string
		c1   *rtsspb.Coordinate
		c2   *rtsspb.Coordinate
		want bool
	}{
		{name: "IsAdjacent", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 0, Y: 1}, want: true},
		{name: "IsSame", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 0, Y: 0}, want: false},
		{name: "IsDiagonal", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 1, Y: 1}, want: false},
		{name: "IsNotAdjacent", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 100, Y: 100}, want: false},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if res := IsAdjacent(
				&Cluster{Val: &rtsspb.Cluster{Coordinate: c.c1}},
				&Cluster{Val: &rtsspb.Cluster{Coordinate: c.c2}}); res != c.want {
				t.Errorf("IsAdjacent((%v, %v), (%v, %v)) = %v, want = %v", c.c1.GetX(), c.c1.GetY(), c.c2.GetX(), c.c2.GetY(), res, c.want)
			}
		})
	}

}

func TestAdjacentDirection(t *testing.T) {
	testConfigs := []struct {
		name        string
		c1          *rtsspb.Coordinate
		c2          *rtsspb.Coordinate
		want        rtscpb.Direction
		wantSuccess bool
	}{
		{name: "AdjacentDirectionNorth", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 0, Y: 1}, want: rtscpb.Direction_DIRECTION_NORTH, wantSuccess: true},
		{name: "AdjacentDirectionSouth", c1: &rtsspb.Coordinate{X: 0, Y: 1}, c2: &rtsspb.Coordinate{X: 0, Y: 0}, want: rtscpb.Direction_DIRECTION_SOUTH, wantSuccess: true},
		{name: "AdjacentDirectionEast", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 1, Y: 0}, want: rtscpb.Direction_DIRECTION_EAST, wantSuccess: true},
		{name: "AdjacentDirectionWest", c1: &rtsspb.Coordinate{X: 1, Y: 0}, c2: &rtsspb.Coordinate{X: 0, Y: 0}, want: rtscpb.Direction_DIRECTION_WEST, wantSuccess: true},
		{name: "NonAdjacentDirection", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 1, Y: 1}, want: rtscpb.Direction_DIRECTION_UNKNOWN, wantSuccess: false},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			res, err := GetRelativeDirection(
				&Cluster{Val: &rtsspb.Cluster{Coordinate: c.c1}},
				&Cluster{Val: &rtsspb.Cluster{Coordinate: c.c2}})
			if (err == nil) != c.wantSuccess {
				t.Fatalf("GetRelativeDirection() = _, %v, want wantSuccess = %v", err, c.wantSuccess)
			}
			if err == nil && res != c.want {
				t.Errorf("GetRelativeDirection() = %v, want = %v", res, c.want)
			}
		})
	}
}

func TestPartition(t *testing.T) {
	testConfigs := []struct {
		name             string
		tileMapDimension int32
		tileDimension    int32
		want             []partitionInfo
		wantSuccess      bool
	}{
		{name: "ZeroWidthMapTest", tileMapDimension: 0, tileDimension: 1, want: nil, wantSuccess: true},
		{name: "ZeroWidthMapZeroDimTest", tileMapDimension: 0, tileDimension: 0, want: nil, wantSuccess: false},
		{name: "SimplePartitionTest", tileMapDimension: 1, tileDimension: 1, want: []partitionInfo{
			{TileBoundary: 0, TileDimension: 1},
		}, wantSuccess: true},
		{name: "SimplePartitionMultipleTest", tileMapDimension: 2, tileDimension: 1, want: []partitionInfo{
			{TileBoundary: 0, TileDimension: 1},
			{TileBoundary: 1, TileDimension: 1},
		}, wantSuccess: true},
		{name: "PartialPartitionTest", tileMapDimension: 1, tileDimension: 2, want: []partitionInfo{
			{TileBoundary: 0, TileDimension: 1},
		}, wantSuccess: true},
		{name: "PartialPartitionMultipleTest", tileMapDimension: 3, tileDimension: 2, want: []partitionInfo{
			{TileBoundary: 0, TileDimension: 2},
			{TileBoundary: 2, TileDimension: 1},
		}, wantSuccess: true},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			ps, err := partition(c.tileMapDimension, c.tileDimension)
			if (err == nil) != c.wantSuccess {
				t.Fatalf("partition() = _, %v, want wantSuccess = %v", err, c.wantSuccess)
			}
			if err == nil && !cmp.Equal(ps, c.want) {
				t.Errorf("partition() = %v, want = %v", ps, c.want)
			}
		})
	}
}

func TestBuildCluster(t *testing.T) {
	testConfigs := []struct {
		name             string
		tileMapDimension *rtsspb.Coordinate
		tileDimension    *rtsspb.Coordinate
		want             *ClusterMap
		wantSuccess      bool
	}{
		{name: "ZeroWidthDimTest", tileMapDimension: &rtsspb.Coordinate{X: 1, Y: 1}, tileDimension: &rtsspb.Coordinate{X: 0, Y: 0}, want: nil, wantSuccess: false},
		{name: "ZeroXDimTest", tileMapDimension: &rtsspb.Coordinate{X: 1, Y: 1}, tileDimension: &rtsspb.Coordinate{X: 0, Y: 1}, want: nil, wantSuccess: false},
		{name: "ZeroYDimTest", tileMapDimension: &rtsspb.Coordinate{X: 1, Y: 1}, tileDimension: &rtsspb.Coordinate{X: 1, Y: 0}, want: nil, wantSuccess: false},
		{name: "ZeroWidthMapTest", tileMapDimension: &rtsspb.Coordinate{X: 0, Y: 0}, tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, want: &ClusterMap{
			L: 1, D: &rtsspb.Coordinate{X: 0, Y: 0}, M: nil}, wantSuccess: true},
		{name: "ZeroXMapTest", tileMapDimension: &rtsspb.Coordinate{X: 0, Y: 1}, tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, want: &ClusterMap{
			L: 1, D: &rtsspb.Coordinate{X: 0, Y: 0}, M: nil}, wantSuccess: true},
		{name: "ZeroYMapTest", tileMapDimension: &rtsspb.Coordinate{X: 1, Y: 0}, tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, want: &ClusterMap{
			L: 1, D: &rtsspb.Coordinate{X: 0, Y: 0}, M: nil}, wantSuccess: true},
		{name: "SimpleTest", tileMapDimension: &rtsspb.Coordinate{X: 1, Y: 1}, tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, want: &ClusterMap{
			L: 1, D: &rtsspb.Coordinate{X: 1, Y: 1}, M: map[utils.MapCoordinate]*Cluster{
				{X: 0, Y: 0}: {
					Val: &rtsspb.Cluster{
						Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
						TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
						TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
					},
				},
			}}, wantSuccess: true},
		{name: "MultiplePartitionTest", tileMapDimension: &rtsspb.Coordinate{X: 2, Y: 3}, tileDimension: &rtsspb.Coordinate{X: 2, Y: 2}, want: &ClusterMap{
			L: 1, D: &rtsspb.Coordinate{X: 1, Y: 2}, M: map[utils.MapCoordinate]*Cluster{
				{X: 0, Y: 0}: {
					Val: &rtsspb.Cluster{
						Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
						TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
						TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
					},
				},
				{X: 0, Y: 1}: {
					Val: &rtsspb.Cluster{
						Coordinate:    &rtsspb.Coordinate{X: 0, Y: 1},
						TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 2},
						TileDimension: &rtsspb.Coordinate{X: 2, Y: 1},
					},
				},
			}}, wantSuccess: true},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			m, err := BuildClusterMap(c.tileMapDimension, c.tileDimension, 1)
			if (err == nil) != c.wantSuccess {
				t.Fatalf("BuildClusterMap() = _, %v, want wantSuccess = %v", err, c.wantSuccess)
			}
			if err == nil && !cmp.Equal(m, c.want, protocmp.Transform()) {
				t.Errorf("BuildClusterMap() = %v, want = %v", m, c.want)
			}
		})
	}
}
