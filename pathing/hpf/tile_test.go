package tile

import (
	"testing"

	rtscpb "github.com/downflux/game/pathing/proto/constants_go_proto"
	rtsspb "github.com/downflux/game/pathing/proto/structs_go_proto"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/downflux/game/pathing/hpf/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	/**
	 *       - - -
	 *       - - -
	 * Y = 0 - - -
	 *   X = 0
	 */
	simpleMapProto = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 3, Y: 3},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
	}
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
				&Tile{Val: &rtsspb.Tile{Coordinate: c.c1}},
				&Tile{Val: &rtsspb.Tile{Coordinate: c.c2}}); res != c.want {
				t.Errorf("IsAdjacent((%v, %v), (%v, %v)) = %v, want = %v", c.c1.GetX(), c.c1.GetY(), c.c2.GetX(), c.c2.GetY(), res, c.want)
			}
		})
	}
}

func TestDNotAdjacent(t *testing.T) {
	t1 := &Tile{Val: &rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}}}
	t2 := &Tile{Val: &rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 1, Y: 1}}}
	if res, err := D(nil, t1, t2); err == nil {
		t.Errorf("D(nil, (%v, %v), (%v, %v)) = (%v, nil), want a non-nil error", t1.X(), t1.Y(), t2.X(), t2.Y(), res)
	}
}

func TestD(t *testing.T) {
	cost := map[rtscpb.TerrainType]float64{
		rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED: 1000,
		rtscpb.TerrainType_TERRAIN_TYPE_PLAINS:  1,
	}
	c1 := &rtsspb.Coordinate{X: 0, Y: 0}
	c2 := &rtsspb.Coordinate{X: 0, Y: 1}
	testConfigs := []struct {
		name         string
		terrainType1 rtscpb.TerrainType
		terrainType2 rtscpb.TerrainType
		want         float64
	}{
		{
			name:         "SimpleD",
			terrainType1: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
			terrainType2: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
			want:         1,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if res, _ := D(
				cost,
				&Tile{Val: &rtsspb.Tile{Coordinate: c1, TerrainType: c.terrainType1}},
				&Tile{Val: &rtsspb.Tile{Coordinate: c2, TerrainType: c.terrainType2}}); res != c.want {
				t.Errorf(
					"D((%v, %v, c=%v), (%v, %v, c=%v)) = %v, want = %v",
					c1.GetX(), c1.GetY(), cost[c.terrainType1],
					c2.GetX(), c2.GetY(), cost[c.terrainType2], res, c.want)
			}
		})
	}
}

func TestH(t *testing.T) {
	testConfigs := []struct {
		name string
		c1   *rtsspb.Coordinate
		c2   *rtsspb.Coordinate
		want float64
	}{
		{name: "TrivialH", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 0, Y: 0}, want: 0},
		{name: "SimpleH", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 1, Y: 1}, want: 2},
		{name: "PythagorasH", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 3, Y: 4}, want: 25},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if res, _ := H(
				&Tile{Val: &rtsspb.Tile{Coordinate: c.c1}},
				&Tile{Val: &rtsspb.Tile{Coordinate: c.c2}}); res != c.want {
				t.Errorf("H((%v, %v), (%v, %v)) = %v, want = %v", c.c1.GetX(), c.c1.GetY(), c.c2.GetX(), c.c2.GetY(), res, c.want)
			}
		})
	}
}

func TestGetTile(t *testing.T) {
	simpleMap, err := ImportMap(simpleMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}

	testConfigs := []struct {
		name       string
		coordinate *rtsspb.Coordinate
		want       *Tile
	}{
		{name: "TrivialTile", coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, want: simpleMap.M[utils.MapCoordinate{X: 0, Y: 0}]},
		{name: "DNETile", coordinate: &rtsspb.Coordinate{X: 100, Y: 100}, want: nil},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if res := simpleMap.Tile(c.coordinate.GetX(), c.coordinate.GetY()); res != c.want {
				t.Errorf("Tile((%v, %v)) = %v, want = %v", c.coordinate.GetX(), c.coordinate.GetY(), res, c.want)
			}
		})
	}
}

func tileLess(t1, t2 *Tile) bool {
	return t1.Val.GetCoordinate().GetX() < t2.Val.GetCoordinate().GetX() || (t1.Val.GetCoordinate().GetX() == t2.Val.GetCoordinate().GetX() && t1.Val.GetCoordinate().GetY() < t2.Val.GetCoordinate().GetY())
}

func TestGetNeighbors(t *testing.T) {
	simpleMap, err := ImportMap(simpleMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}

	testConfigs := []struct {
		name       string
		m          *Map
		coordinate *rtsspb.Coordinate
		want       []*Tile
		err        error
	}{
		{name: "NullMapNeighbors", m: &Map{}, coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, want: nil, err: status.Errorf(codes.NotFound, "")},
		{name: "DNETileNeighbors", m: simpleMap, coordinate: &rtsspb.Coordinate{X: 100, Y: 100}, want: nil, err: status.Errorf(codes.NotFound, "")},
		{name: "FullNeighbors", m: simpleMap, coordinate: &rtsspb.Coordinate{X: 1, Y: 1}, want: []*Tile{
			simpleMap.Tile(0, 1),
			simpleMap.Tile(2, 1),
			simpleMap.Tile(1, 0),
			simpleMap.Tile(1, 2),
		}, err: nil},
		{name: "EdgeNeighbors", m: simpleMap, coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, want: []*Tile{
			simpleMap.Tile(0, 2),
			simpleMap.Tile(0, 0),
			simpleMap.Tile(1, 1),
		}, err: nil},
		{name: "CornerNeighbors", m: simpleMap, coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, want: []*Tile{
			simpleMap.Tile(0, 1),
			simpleMap.Tile(1, 0),
		}, err: nil},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			res, err := c.m.Neighbors(c.coordinate)
			resStatus, success := status.FromError(err)
			wantStatus, _ := status.FromError(c.err)
			if !success || resStatus.Code() != wantStatus.Code() {
				t.Errorf("Neighbors((%v, %v)) = (_, %v), want = (_, %v)", c.coordinate.GetX(), c.coordinate.GetY(), err, c.err)
				return
			}
			if !cmp.Equal(res, c.want, cmpopts.SortSlices(tileLess), protocmp.Transform()) {
				t.Errorf("Neighbors((%v, %v) = %v, want = %v", c.coordinate.GetX(), c.coordinate.GetY(), res, c.want)
			}
		})
	}
}
