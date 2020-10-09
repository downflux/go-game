package entrance

import (
	"math"
	"testing"

	rtscpb "github.com/downflux/game/pathing/proto/constants_go_proto"
	rtsspb "github.com/downflux/game/pathing/proto/structs_go_proto"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/downflux/game/pathing/hpf/cluster"
	"github.com/downflux/game/pathing/hpf/tile"
	"github.com/downflux/game/pathing/hpf/utils"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	/**
	 * Y = 0 W W
	 *   X = 0
	 */
	trivialClosedMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 2, Y: 1},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
		},
		TerrainCosts: []*rtsspb.TerrainCost{
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	}

	/**
	 * Y = 0 - -
	 *   X = 0
	 */
	trivialOpenMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 2, Y: 1},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
	}

	/**
	 * Y = 0 - W
	 *   X = 0
	 */
	trivialSemiOpenMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 2, Y: 1},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
		},
		TerrainCosts: []*rtsspb.TerrainCost{
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	}

	/**
	 *       - -
	 *       - -
	 *       - -
	 * Y = 0 - -
	 *   X = 0
	 */
	longVerticalOpenMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 2, Y: 4},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 3}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 3}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
	}

	/**
	 *       - - - -
	 * Y = 0 - - - -
	 *   X = 0
	 */
	longHorizontalOpenMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 4, Y: 2},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 3, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 3, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
	}

	/**
	 *       - -
	 *       W W
	 * Y = 0 - -
	 *   X = 0
	 */
	longSemiOpenMap = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 2, Y: 3},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
		TerrainCosts: []*rtsspb.TerrainCost{
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	}
)

func TestBuildClusterEdgeCoordinateSliceError(t *testing.T) {
	testConfigs := []struct {
		name string
		m    *rtsspb.ClusterMap
		c    utils.MapCoordinate
		d    rtscpb.Direction
	}{
		{
			name: "NullClusterTest",
			m: &rtsspb.ClusterMap{
				TileDimension:    &rtsspb.Coordinate{X: 0, Y: 0},
				TileMapDimension: &rtsspb.Coordinate{X: 0, Y: 0},
			},
			c: utils.MapCoordinate{X: 0, Y: 0},
			d: rtscpb.Direction_DIRECTION_NORTH,
		},
		{
			name: "NullXDimensionClusterTest",
			m: &rtsspb.ClusterMap{
				TileDimension:    &rtsspb.Coordinate{X: 0, Y: 5},
				TileMapDimension: &rtsspb.Coordinate{X: 0, Y: 10},
			},
			c: utils.MapCoordinate{X: 0, Y: 1},
			d: rtscpb.Direction_DIRECTION_NORTH,
		},
		{
			name: "NullYDimensionClusterTest",
			m: &rtsspb.ClusterMap{
				TileDimension:    &rtsspb.Coordinate{X: 5, Y: 0},
				TileMapDimension: &rtsspb.Coordinate{X: 10, Y: 0},
			},
			c: utils.MapCoordinate{X: 1, Y: 0},
			d: rtscpb.Direction_DIRECTION_NORTH,
		},
		{
			name: "InvalidDirectionTest",
			m: &rtsspb.ClusterMap{
				TileDimension:    &rtsspb.Coordinate{X: 5, Y: 5},
				TileMapDimension: &rtsspb.Coordinate{X: 10, Y: 10},
			},
			c: utils.MapCoordinate{X: 1, Y: 1},
			d: rtscpb.Direction_DIRECTION_UNKNOWN,
		},
	}
	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			m, err := cluster.ImportMap(c.m)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			if got, err := buildClusterEdgeCoordinateSlice(m, c.c, c.d); err == nil {
				t.Errorf("buildClusterEdgeCoordinateSlice() = %v, %v, want a non-nil error", got, err)
			}
		})
	}
}

func TestBuildClusterEdgeCoordinateSlice(t *testing.T) {
	trivialClusterMap := &rtsspb.ClusterMap{
		TileDimension:    &rtsspb.Coordinate{X: 1, Y: 1},
		TileMapDimension: &rtsspb.Coordinate{X: 1, Y: 1},
	}
	smallClusterMap := &rtsspb.ClusterMap{
		TileDimension:    &rtsspb.Coordinate{X: 2, Y: 2},
		TileMapDimension: &rtsspb.Coordinate{X: 2, Y: 2},
	}
	embeddedClusterMap := &rtsspb.ClusterMap{
		TileDimension:    &rtsspb.Coordinate{X: 2, Y: 2},
		TileMapDimension: &rtsspb.Coordinate{X: 4, Y: 4},
	}
	rectangularClusterMap := &rtsspb.ClusterMap{
		TileDimension:    &rtsspb.Coordinate{X: 1, Y: 2},
		TileMapDimension: &rtsspb.Coordinate{X: 2, Y: 4},
	}
	testConfigs := []struct {
		name string
		m    *rtsspb.ClusterMap
		c    utils.MapCoordinate
		d    rtscpb.Direction
		want *rtsspb.CoordinateSlice
	}{
		{name: "TrivialClusterNorthTest", m: trivialClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterSouthTest", m: trivialClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterEastTest", m: trivialClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterWestTest", m: trivialClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "SmallClusterNorthTest", m: smallClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 2}},
		{name: "SmallClusterSouthTest", m: smallClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2}},
		{name: "SmallClusterEastTest", m: smallClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 2}},
		{name: "SmallClusterWestTest", m: smallClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2}},
		{name: "EmbeddedClusterNorthTest", m: embeddedClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 2, Y: 3}, Length: 2}},
		{name: "EmbeddedClusterSouthTest", m: embeddedClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 2, Y: 2}, Length: 2}},
		{name: "EmbeddedClusterEastTest", m: embeddedClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 3, Y: 2}, Length: 2}},
		{name: "EmbeddedClusterWestTest", m: embeddedClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 2, Y: 2}, Length: 2}},
		{name: "RectangularClusterNorthTest", m: rectangularClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 3}, Length: 1}},
		{name: "RectangularClusterSouthTest", m: rectangularClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, Length: 1}},
		{name: "RectangularClusterEastTest", m: rectangularClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, Length: 2}},
		{name: "RectangularClusterWestTest", m: rectangularClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, Length: 2}},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			m, err := cluster.ImportMap(c.m)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			if got, err := buildClusterEdgeCoordinateSlice(m, c.c, c.d); err != nil || !proto.Equal(got, c.want) {
				t.Errorf("buildClusterEdgeCoordinateSlice() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestBuildCoordinateWithCoordinateSlice(t *testing.T) {
	testConfigs := []struct {
		name   string
		s      *rtsspb.CoordinateSlice
		offset int32
		want   *rtsspb.Coordinate
	}{
		{
			name:   "SingleTileSliceHorizontal",
			s:      &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			offset: 0,
			want:   &rtsspb.Coordinate{X: 0, Y: 0},
		},
		{
			name:   "SingleTileSliceVertical",
			s:      &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			offset: 0,
			want:   &rtsspb.Coordinate{X: 0, Y: 0},
		},
		{
			name:   "MultiTileTileSliceHorizontal",
			s:      &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2},
			offset: 1,
			want:   &rtsspb.Coordinate{X: 2, Y: 1},
		},
		{
			name:   "MultiTileTileSliceVertical",
			s:      &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2},
			offset: 1,
			want:   &rtsspb.Coordinate{X: 1, Y: 2},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := buildCoordinateWithCoordinateSlice(c.s, c.offset); err != nil || !proto.Equal(got, c.want) {
				t.Errorf("buildCoordinateWithCoordinateSlice() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestBuildCoordinateWithCoordinateSliceError(t *testing.T) {
	testConfigs := []struct {
		name   string
		s      *rtsspb.CoordinateSlice
		offset int32
	}{
		{name: "NullTileSlice", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 0}, offset: 0},
		{name: "OutOfBoundsTileSliceBefore", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}, offset: -1},
		{name: "OutOfBoundsTileSliceAfter", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}, offset: 2},
		{name: "InvalidOrientationTileSlice", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_UNKNOWN, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}, offset: 0},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if _, err := buildCoordinateWithCoordinateSlice(c.s, c.offset); err == nil {
				t.Error("buildCoordinateWithCoordinateSlice() = nil, want a non-nil error")
			}
		})
	}
}

func TestBuildTransitionsFromOpenCoordinateSlice(t *testing.T) {
	testConfigs := []struct {
		name   string
		s1, s2 *rtsspb.CoordinateSlice
		want   []*rtsspb.Transition
	}{
		{
			name: "SingleWidthEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 1},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 1}},
				},
			},
		},
		{
			name: "SingleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 1},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
				},
			},
		},
		{
			name: "DoubleWidthEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 2},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 1}},
				},
			},
		},
		{
			name: "DoubleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 2},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 1}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 1}},
				},
			},
		},
		{
			name: "TripleWidthEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 3},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 3},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 1}},
				},
			},
		},
		{
			name: "TripleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 3},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 3},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 1}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 1}},
				},
			},
		},
		{
			name: "QuadrupleWidthEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 4},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 1}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 3, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 3, Y: 1}},
				},
			},
		},
		{
			name: "QuadrupleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 4},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 3}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 3}},
				},
			},
		},
		{
			name: "QuadrupleWidthEmbeddedEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 4},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, Length: 4},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 1}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 2}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 4, Y: 1}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 4, Y: 2}},
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := buildTransitionsFromOpenCoordinateSlice(c.s1, c.s2); err != nil || !cmp.Equal(got, c.want, protocmp.Transform()) {
				t.Errorf("buildTransitionsFromOpenCoordinateSlice() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestVerifyCoordinateSlicesError(t *testing.T) {
	testConfigs := []struct {
		name   string
		s1, s2 *rtsspb.CoordinateSlice
	}{
		{
			name: "MismatchedLengths",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 2},
		},
		{
			name: "MismatchedOrientations",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
		},
		{
			name: "NonAdjacentHorizontalSlice",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 2}, Length: 1},
		},
		{
			name: "NonAdjacentVerticalSlice",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 2, Y: 0}, Length: 1},
		},
		{
			name: "NonAlignedHorizontalSlice",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2},
		},
		{
			name: "NonAlignedVerticalSlice",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if err := verifyCoordinateSlices(c.s1, c.s2); err == nil {
				t.Error("verifyCoordinateSlices() = nil, want a non-nil error")
			}
		})
	}
}

func TestBuildTransitionsError(t *testing.T) {
	trivialOpenClusterMap := &rtsspb.ClusterMap{
		TileDimension:    &rtsspb.Coordinate{X: 1, Y: 1},
		TileMapDimension: trivialOpenMap.GetDimension(),
	}
	longVerticalOpenClusterMap := &rtsspb.ClusterMap{
		TileDimension:    &rtsspb.Coordinate{X: 2, Y: 1},
		TileMapDimension: longVerticalOpenMap.GetDimension(),
	}

	testConfigs := []struct {
		name   string
		m      *rtsspb.TileMap
		cm     *rtsspb.ClusterMap
		c1, c2 utils.MapCoordinate
	}{
		{name: "NullCluster", m: trivialOpenMap, cm: nil, c1: utils.MapCoordinate{}, c2: utils.MapCoordinate{}},
		{name: "NullMap", m: nil, cm: trivialOpenClusterMap, c1: utils.MapCoordinate{X: 0, Y: 0}, c2: utils.MapCoordinate{X: 1, Y: 0}},
		{name: "NonAdjacentClusters", m: longVerticalOpenMap, cm: longVerticalOpenClusterMap, c1: utils.MapCoordinate{X: 0, Y: 0}, c2: utils.MapCoordinate{X: 1, Y: 1}},
	}
	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			m, err := tile.ImportMap(c.m)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil")
			}
			cm, err := cluster.ImportMap(c.cm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil")
			}

			if got, err := BuildTransitions(m, cm, c.c1, c.c2); err == nil {
				t.Errorf("BuildTransitions() = %v, %v, want a non-nil error", got, err)
			}
		})
	}
}

func TestBuildTransitionsAux(t *testing.T) {
	testConfigs := []struct {
		name   string
		m      *rtsspb.TileMap
		s1, s2 *rtsspb.CoordinateSlice
		want   []*rtsspb.Transition
	}{
		{name: "TrivialClosedMap", m: trivialClosedMap,
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 1},
			want: nil,
		},
		{name: "TrivialSemiOpenMap", m: trivialSemiOpenMap,
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 1},
			want: nil,
		},
		{name: "TrivialOpenMap", m: trivialOpenMap,
			s1: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 1},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
				},
			},
		},
		{name: "LongVerticalOpenMap", m: longVerticalOpenMap,
			s1: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 4},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 3}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 3}},
				},
			},
		},
		{name: "LongHorizontalOpenMap", m: longHorizontalOpenMap,
			s1: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 4},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 1}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 3, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 3, Y: 1}},
				},
			},
		},
		{name: "LongSemiOpenMap", m: longSemiOpenMap,
			s1: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 3},
			s2: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 3},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 2}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 2}},
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tileMap, err := tile.ImportMap(c.m)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			if got, err := buildTransitionsAux(tileMap, c.s1, c.s2); err != nil || !cmp.Equal(got, c.want, protocmp.Transform()) {
				t.Errorf("buildTransitionsAux() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}
func TestBuildTransitions(t *testing.T) {
	trivialClusterMap := &rtsspb.ClusterMap{TileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, TileMapDimension: trivialClosedMap.GetDimension()}
	longVerticalClusterMap := &rtsspb.ClusterMap{TileDimension: &rtsspb.Coordinate{X: 1, Y: 4}, TileMapDimension: longVerticalOpenMap.GetDimension()}
	longHorizontalClusterMap := &rtsspb.ClusterMap{TileDimension: &rtsspb.Coordinate{X: 4, Y: 1}, TileMapDimension: longHorizontalOpenMap.GetDimension()}
	longSemiOpenClusterMap := &rtsspb.ClusterMap{TileDimension: &rtsspb.Coordinate{X: 1, Y: 3}, TileMapDimension: longSemiOpenMap.GetDimension()}

	testConfigs := []struct {
		name   string
		m      *rtsspb.TileMap
		cm     *rtsspb.ClusterMap
		c1, c2 *rtsspb.Coordinate
		want   []*rtsspb.Transition
	}{
		{name: "TrivialClosedMap", m: trivialClosedMap, cm: trivialClusterMap, c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 1, Y: 0}, want: nil},
		{name: "TrivialSemiOpenMap", m: trivialSemiOpenMap, cm: trivialClusterMap, c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 1, Y: 0}, want: nil},
		{name: "TrivialOpenMap", m: trivialOpenMap, cm: trivialClusterMap, c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 1, Y: 0},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
				},
			},
		},
		{name: "LongVerticalOpenMap", m: longVerticalOpenMap, cm: longVerticalClusterMap, c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 1, Y: 0},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 3}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 3}},
				},
			},
		},
		{name: "LongHorizontalOpenMap", m: longHorizontalOpenMap, cm: longHorizontalClusterMap, c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 0, Y: 1},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 1}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 3, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 3, Y: 1}},
				},
			},
		},
		{name: "LongSemiOpenMap", m: longSemiOpenMap, cm: longSemiOpenClusterMap, c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 1, Y: 0},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 2}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 2}},
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			m, err := tile.ImportMap(c.m)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}
			cm, err := cluster.ImportMap(c.cm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			if got, err := BuildTransitions(m, cm, utils.MC(c.c1), utils.MC(c.c2)); err != nil || !cmp.Equal(got, c.want, protocmp.Transform()) {
				t.Errorf("BuildTransitions() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}
func TestSliceContainsError(t *testing.T) {
	s := &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_UNKNOWN, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}
	if _, err := sliceContains(s, utils.MC(&rtsspb.Coordinate{X: 0, Y: 0})); err == nil {
		t.Error("sliceContains() = _, nil, want a non-nil error")
	}
}

func TestSliceContains(t *testing.T) {
	testConfigs := []struct {
		name string
		s    *rtsspb.CoordinateSlice
		c    *rtsspb.Coordinate
		want bool
	}{
		{
			name: "TrivialSliceContains",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			c:    &rtsspb.Coordinate{X: 0, Y: 0},
			want: true,
		},
		{
			name: "TrivialPreSliceNoContains",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			c:    &rtsspb.Coordinate{X: -1, Y: 0},
			want: false,
		},
		{
			name: "TrivialPostSliceNoContains",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			c:    &rtsspb.Coordinate{X: 1, Y: 0},
			want: false,
		},
		{
			name: "TrivialBadAxisSliceNoContains",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			c:    &rtsspb.Coordinate{X: 0, Y: -1},
			want: false,
		},
		{
			name: "SimpleSliceContainsHorizontal",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			c:    &rtsspb.Coordinate{X: 1, Y: 0},
			want: true,
		},
		{
			name: "SimpleSliceContainsVertical",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			c:    &rtsspb.Coordinate{X: 0, Y: 1},
			want: true,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if res, err := sliceContains(c.s, utils.MC(c.c)); err != nil || res != c.want {
				t.Errorf("sliceContains() = %v, %v, want = %v, nil", res, err, c.want)
			}
		})
	}
}

func TestOnClusterEdge(t *testing.T) {
	testConfigs := []struct {
		name string
		m    *rtsspb.ClusterMap
		c    *rtsspb.Coordinate
		t    *rtsspb.Coordinate
		want bool
	}{
		{
			name: "TrivialClusterContains",
			m:    &rtsspb.ClusterMap{TileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, TileMapDimension: &rtsspb.Coordinate{X: 1, Y: 1}},
			c:    &rtsspb.Coordinate{X: 0, Y: 0},
			t:    &rtsspb.Coordinate{X: 0, Y: 0},
			want: true,
		},
		{
			name: "TrivialClusterNoContains",
			m:    &rtsspb.ClusterMap{TileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, TileMapDimension: &rtsspb.Coordinate{X: 2, Y: 2}},
			c:    &rtsspb.Coordinate{X: 0, Y: 0},
			t:    &rtsspb.Coordinate{X: 0, Y: 1},
			want: false,
		},
		{
			name: "ClusterInternalNoContains",
			m:    &rtsspb.ClusterMap{TileDimension: &rtsspb.Coordinate{X: 3, Y: 3}, TileMapDimension: &rtsspb.Coordinate{X: 3, Y: 3}},
			c:    &rtsspb.Coordinate{X: 0, Y: 0},
			t:    &rtsspb.Coordinate{X: 1, Y: 1},
			want: false,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			m, err := cluster.ImportMap(c.m)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			if got := OnClusterEdge(m, utils.MC(c.c), utils.MC(c.t)); got != c.want {
				t.Errorf("OnClusterEdge() = %v, want = %v", got, c.want)
			}
		})
	}
}
