package tile

import (
	"testing"

	gdpb "github.com/downflux/game/api/data_go_proto"
	pcpb "github.com/downflux/game/pathing/api/constants_go_proto"
	pdpb "github.com/downflux/game/pathing/api/data_go_proto"

	"github.com/downflux/game/pathing/hpf/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
	simpleMapProto = &pdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 3, Y: 3},
		Tiles: []*pdpb.Tile{
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}, TerrainType: pcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 1}, TerrainType: pcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 2}, TerrainType: pcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 0}, TerrainType: pcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 1}, TerrainType: pcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 2}, TerrainType: pcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 0}, TerrainType: pcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 1}, TerrainType: pcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 2}, TerrainType: pcpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
	}
)

func TestIsAdjacent(t *testing.T) {
	testConfigs := []struct {
		name string
		c1   *gdpb.Coordinate
		c2   *gdpb.Coordinate
		want bool
	}{
		{name: "IsAdjacent", c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 0, Y: 1}, want: true},
		{name: "IsSame", c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 0, Y: 0}, want: false},
		{name: "IsDiagonal", c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 1, Y: 1}, want: false},
		{name: "IsNotAdjacent", c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 100, Y: 100}, want: false},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if res := IsAdjacent(
				&Tile{Val: &pdpb.Tile{Coordinate: c.c1}},
				&Tile{Val: &pdpb.Tile{Coordinate: c.c2}}); res != c.want {
				t.Errorf("IsAdjacent((%v, %v), (%v, %v)) = %v, want = %v", c.c1.GetX(), c.c1.GetY(), c.c2.GetX(), c.c2.GetY(), res, c.want)
			}
		})
	}
}

func TestDNotAdjacent(t *testing.T) {
	t1 := &Tile{Val: &pdpb.Tile{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}}}
	t2 := &Tile{Val: &pdpb.Tile{Coordinate: &gdpb.Coordinate{X: 1, Y: 1}}}
	if res, err := D(nil, t1, t2); err == nil {
		t.Errorf("D(nil, (%v, %v), (%v, %v)) = (%v, nil), want a non-nil error", t1.X(), t1.Y(), t2.X(), t2.Y(), res)
	}
}

func TestD(t *testing.T) {
	cost := map[pcpb.TerrainType]float64{
		pcpb.TerrainType_TERRAIN_TYPE_BLOCKED: 1000,
		pcpb.TerrainType_TERRAIN_TYPE_PLAINS:  1,
	}
	c1 := &gdpb.Coordinate{X: 0, Y: 0}
	c2 := &gdpb.Coordinate{X: 0, Y: 1}
	testConfigs := []struct {
		name         string
		terrainType1 pcpb.TerrainType
		terrainType2 pcpb.TerrainType
		want         float64
	}{
		{
			name:         "SimpleD",
			terrainType1: pcpb.TerrainType_TERRAIN_TYPE_PLAINS,
			terrainType2: pcpb.TerrainType_TERRAIN_TYPE_PLAINS,
			want:         1,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if res, _ := D(
				cost,
				&Tile{Val: &pdpb.Tile{Coordinate: c1, TerrainType: c.terrainType1}},
				&Tile{Val: &pdpb.Tile{Coordinate: c2, TerrainType: c.terrainType2}}); res != c.want {
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
		c1   *gdpb.Coordinate
		c2   *gdpb.Coordinate
		want float64
	}{
		{name: "TrivialH", c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 0, Y: 0}, want: 0},
		{name: "SimpleH", c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 1, Y: 1}, want: 2},
		{name: "PythagorasH", c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 3, Y: 4}, want: 25},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if res, _ := H(
				&Tile{Val: &pdpb.Tile{Coordinate: c.c1}},
				&Tile{Val: &pdpb.Tile{Coordinate: c.c2}}); res != c.want {
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
		coordinate *gdpb.Coordinate
		want       *Tile
	}{
		{name: "TrivialTile", coordinate: &gdpb.Coordinate{X: 0, Y: 0}, want: simpleMap.M[utils.MapCoordinate{X: 0, Y: 0}]},
		{name: "DNETile", coordinate: &gdpb.Coordinate{X: 100, Y: 100}, want: nil},
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
		coordinate *gdpb.Coordinate
		want       []*Tile
		err        error
	}{
		{name: "NullMapNeighbors", m: &Map{}, coordinate: &gdpb.Coordinate{X: 0, Y: 0}, want: nil, err: status.Errorf(codes.NotFound, "")},
		{name: "DNETileNeighbors", m: simpleMap, coordinate: &gdpb.Coordinate{X: 100, Y: 100}, want: nil, err: status.Errorf(codes.NotFound, "")},
		{name: "FullNeighbors", m: simpleMap, coordinate: &gdpb.Coordinate{X: 1, Y: 1}, want: []*Tile{
			simpleMap.Tile(0, 1),
			simpleMap.Tile(2, 1),
			simpleMap.Tile(1, 0),
			simpleMap.Tile(1, 2),
		}, err: nil},
		{name: "EdgeNeighbors", m: simpleMap, coordinate: &gdpb.Coordinate{X: 0, Y: 1}, want: []*Tile{
			simpleMap.Tile(0, 2),
			simpleMap.Tile(0, 0),
			simpleMap.Tile(1, 1),
		}, err: nil},
		{name: "CornerNeighbors", m: simpleMap, coordinate: &gdpb.Coordinate{X: 0, Y: 0}, want: []*Tile{
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
