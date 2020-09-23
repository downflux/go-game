package cluster

import (
	"testing"

	rtscpb "github.com/minkezhang/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	"github.com/minkezhang/rts-pathing/lib/hpf/utils"
)

var (
	largeClusterMap = &Map{
		Val: &rtsspb.ClusterMap{
			Level:            1,
			TileDimension:    &rtsspb.Coordinate{X: 10, Y: 10},
			TileMapDimension: &rtsspb.Coordinate{X: 10000, Y: 10000},
		},
	}
)

func TestIsAdjacent(t *testing.T) {
	testConfigs := []struct {
		name string
		c1   utils.MapCoordinate
		c2   utils.MapCoordinate
		want bool
	}{
		{name: "IsAdjacent", c1: utils.MapCoordinate{X: 0, Y: 0}, c2: utils.MapCoordinate{X: 0, Y: 1}, want: true},
		{name: "IsSame", c1: utils.MapCoordinate{X: 0, Y: 0}, c2: utils.MapCoordinate{X: 0, Y: 0}, want: false},
		{name: "IsDiagonal", c1: utils.MapCoordinate{X: 0, Y: 0}, c2: utils.MapCoordinate{X: 1, Y: 1}, want: false},
		{name: "IsNotAdjacent", c1: utils.MapCoordinate{X: 0, Y: 0}, c2: utils.MapCoordinate{X: 100, Y: 100}, want: false},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if res := IsAdjacent(largeClusterMap, c.c1, c.c2); res != c.want {
				t.Errorf("IsAdjacent(m, (%v, %v), (%v, %v)) = %v, want = %v", c.c1.X, c.c1.Y, c.c2.X, c.c2.Y, res, c.want)
			}
		})
	}

}

func TestAdjacentDirection(t *testing.T) {
	testConfigs := []struct {
		name        string
		c1          utils.MapCoordinate
		c2          utils.MapCoordinate
		want        rtscpb.Direction
		wantSuccess bool
	}{
		{name: "AdjacentDirectionNorth", c1: utils.MapCoordinate{X: 0, Y: 0}, c2: utils.MapCoordinate{X: 0, Y: 1}, want: rtscpb.Direction_DIRECTION_NORTH, wantSuccess: true},
		{name: "AdjacentDirectionSouth", c1: utils.MapCoordinate{X: 0, Y: 1}, c2: utils.MapCoordinate{X: 0, Y: 0}, want: rtscpb.Direction_DIRECTION_SOUTH, wantSuccess: true},
		{name: "AdjacentDirectionEast", c1: utils.MapCoordinate{X: 0, Y: 0}, c2: utils.MapCoordinate{X: 1, Y: 0}, want: rtscpb.Direction_DIRECTION_EAST, wantSuccess: true},
		{name: "AdjacentDirectionWest", c1: utils.MapCoordinate{X: 1, Y: 0}, c2: utils.MapCoordinate{X: 0, Y: 0}, want: rtscpb.Direction_DIRECTION_WEST, wantSuccess: true},
		{name: "NonAdjacentDirection", c1: utils.MapCoordinate{X: 0, Y: 0}, c2: utils.MapCoordinate{X: 1, Y: 1}, want: rtscpb.Direction_DIRECTION_UNKNOWN, wantSuccess: false},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			res, err := GetRelativeDirection(largeClusterMap, c.c1, c.c2)
			if (err == nil) != c.wantSuccess {
				t.Fatalf("GetRelativeDirection() = _, %v, want wantSuccess = %v", err, c.wantSuccess)
			}
			if err == nil && res != c.want {
				t.Errorf("GetRelativeDirection() = %v, want = %v", res, c.want)
			}
		})
	}
}

func TestCoordinateInCluster(t *testing.T) {
	testConfigs := []struct {
		name string
		co   utils.MapCoordinate
		cl   utils.MapCoordinate
		want bool
	}{
		{
			name: "TrivialClusterIn",
			co:   utils.MapCoordinate{X: 0, Y: 0},
			cl:   utils.MapCoordinate{X: 0, Y: 0},
			want: true,
		},
		{
			name: "TrivialClusterNotIn",
			co:   utils.MapCoordinate{X: 0, Y: 0},
			cl:   utils.MapCoordinate{X: 1, Y: 0},
			want: false,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := CoordinateInCluster(largeClusterMap, c.cl, c.co); c.want != got {
				t.Errorf("CoordinateInCluster() = %v, want = %v", got, c.want)
			}
		})
	}
}
