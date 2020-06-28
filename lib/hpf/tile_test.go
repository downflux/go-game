package tile

import (
	"testing"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"
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
				&Tile{t: &rtsspb.Tile{Coordinate: c.c1}},
				&Tile{t: &rtsspb.Tile{Coordinate: c.c2}}); res != c.want {
				t.Errorf("IsAdjacent((%v, %v), (%v, %v)) = %v, want = %v", c.c1.GetX(), c.c1.GetY(), c.c2.GetX(), c.c2.GetY(), res, c.want)
			}
		})
	}
}

func TestDNotAdjacent(t *testing.T) {
	t1 := &Tile{t: &rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}}}
	t2 := &Tile{t: &rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 1, Y: 1}}}
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
			want:         2,
		},
		{
			name:         "CommutativeA",
			terrainType1: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED,
			terrainType2: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
			want:         1001,
		},
		{
			name:         "CommutativeB",
			terrainType1: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS,
			terrainType2: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED,
			want:         1001,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if res, _ := D(
				cost,
				&Tile{t: &rtsspb.Tile{Coordinate: c1, TerrainType: c.terrainType1}},
				&Tile{t: &rtsspb.Tile{Coordinate: c2, TerrainType: c.terrainType2}}); res != c.want {
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
				&Tile{t: &rtsspb.Tile{Coordinate: c.c1}},
				&Tile{t: &rtsspb.Tile{Coordinate: c.c2}}); res != c.want {
				t.Errorf("H((%v, %v), (%v, %v)) = %v, want = %v", c.c1.GetX(), c.c1.GetY(), c.c2.GetX(), c.c2.GetY(), res, c.want)
			}
		})
	}
}
