package entrance

import (
	"math"
	"testing"

	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/pathing/hpf/cluster"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	gdpb "github.com/downflux/game/api/data_go_proto"
	mcpb "github.com/downflux/game/map/api/constants_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	pcpb "github.com/downflux/game/pathing/api/constants_go_proto"
	pdpb "github.com/downflux/game/pathing/api/data_go_proto"
)

var (
	/**
	 * Y = 0 W W
	 *   X = 0
	 */
	trivialClosedMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 2, Y: 1},
		Tiles: []*mdpb.Tile{
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED},
		},
		TerrainCosts: []*mdpb.TerrainCost{
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	}

	/**
	 * Y = 0 - -
	 *   X = 0
	 */
	trivialOpenMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 2, Y: 1},
		Tiles: []*mdpb.Tile{
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
	}

	/**
	 * Y = 0 - W
	 *   X = 0
	 */
	trivialSemiOpenMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 2, Y: 1},
		Tiles: []*mdpb.Tile{
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED},
		},
		TerrainCosts: []*mdpb.TerrainCost{
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	}

	/**
	 *       - -
	 *       - -
	 *       - -
	 * Y = 0 - -
	 *   X = 0
	 */
	longVerticalOpenMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 2, Y: 4},
		Tiles: []*mdpb.Tile{
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 2}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 3}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 2}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 3}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
	}

	/**
	 *       - - - -
	 * Y = 0 - - - -
	 *   X = 0
	 */
	longHorizontalOpenMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 4, Y: 2},
		Tiles: []*mdpb.Tile{
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 3, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 3, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
	}

	/**
	 *       - -
	 *       W W
	 * Y = 0 - -
	 *   X = 0
	 */
	longSemiOpenMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 2, Y: 3},
		Tiles: []*mdpb.Tile{
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 2}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 2}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
		TerrainCosts: []*mdpb.TerrainCost{
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	}
)

func TestBuildClusterEdgeCoordinateSliceError(t *testing.T) {
	testConfigs := []struct {
		name string
		m    *pdpb.ClusterMap
		c    utils.MapCoordinate
		d    pcpb.Direction
	}{
		{
			name: "NullClusterTest",
			m: &pdpb.ClusterMap{
				TileDimension:    &gdpb.Coordinate{X: 0, Y: 0},
				TileMapDimension: &gdpb.Coordinate{X: 0, Y: 0},
			},
			c: utils.MapCoordinate{X: 0, Y: 0},
			d: pcpb.Direction_DIRECTION_NORTH,
		},
		{
			name: "NullXDimensionClusterTest",
			m: &pdpb.ClusterMap{
				TileDimension:    &gdpb.Coordinate{X: 0, Y: 5},
				TileMapDimension: &gdpb.Coordinate{X: 0, Y: 10},
			},
			c: utils.MapCoordinate{X: 0, Y: 1},
			d: pcpb.Direction_DIRECTION_NORTH,
		},
		{
			name: "NullYDimensionClusterTest",
			m: &pdpb.ClusterMap{
				TileDimension:    &gdpb.Coordinate{X: 5, Y: 0},
				TileMapDimension: &gdpb.Coordinate{X: 10, Y: 0},
			},
			c: utils.MapCoordinate{X: 1, Y: 0},
			d: pcpb.Direction_DIRECTION_NORTH,
		},
		{
			name: "InvalidDirectionTest",
			m: &pdpb.ClusterMap{
				TileDimension:    &gdpb.Coordinate{X: 5, Y: 5},
				TileMapDimension: &gdpb.Coordinate{X: 10, Y: 10},
			},
			c: utils.MapCoordinate{X: 1, Y: 1},
			d: pcpb.Direction_DIRECTION_UNKNOWN,
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
	trivialClusterMap := &pdpb.ClusterMap{
		TileDimension:    &gdpb.Coordinate{X: 1, Y: 1},
		TileMapDimension: &gdpb.Coordinate{X: 1, Y: 1},
	}
	smallClusterMap := &pdpb.ClusterMap{
		TileDimension:    &gdpb.Coordinate{X: 2, Y: 2},
		TileMapDimension: &gdpb.Coordinate{X: 2, Y: 2},
	}
	embeddedClusterMap := &pdpb.ClusterMap{
		TileDimension:    &gdpb.Coordinate{X: 2, Y: 2},
		TileMapDimension: &gdpb.Coordinate{X: 4, Y: 4},
	}
	rectangularClusterMap := &pdpb.ClusterMap{
		TileDimension:    &gdpb.Coordinate{X: 1, Y: 2},
		TileMapDimension: &gdpb.Coordinate{X: 2, Y: 4},
	}
	testConfigs := []struct {
		name string
		m    *pdpb.ClusterMap
		c    utils.MapCoordinate
		d    pcpb.Direction
		want coordinateSlice
	}{
		{name: "TrivialClusterNorthTest", m: trivialClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: pcpb.Direction_DIRECTION_NORTH, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterSouthTest", m: trivialClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: pcpb.Direction_DIRECTION_SOUTH, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterEastTest", m: trivialClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: pcpb.Direction_DIRECTION_EAST, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterWestTest", m: trivialClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: pcpb.Direction_DIRECTION_WEST, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "SmallClusterNorthTest", m: smallClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: pcpb.Direction_DIRECTION_NORTH, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 1}, Length: 2}},
		{name: "SmallClusterSouthTest", m: smallClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: pcpb.Direction_DIRECTION_SOUTH, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 2}},
		{name: "SmallClusterEastTest", m: smallClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: pcpb.Direction_DIRECTION_EAST, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 0}, Length: 2}},
		{name: "SmallClusterWestTest", m: smallClusterMap, c: utils.MapCoordinate{X: 0, Y: 0}, d: pcpb.Direction_DIRECTION_WEST, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 2}},
		{name: "EmbeddedClusterNorthTest", m: embeddedClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: pcpb.Direction_DIRECTION_NORTH, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 2, Y: 3}, Length: 2}},
		{name: "EmbeddedClusterSouthTest", m: embeddedClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: pcpb.Direction_DIRECTION_SOUTH, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 2, Y: 2}, Length: 2}},
		{name: "EmbeddedClusterEastTest", m: embeddedClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: pcpb.Direction_DIRECTION_EAST, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 3, Y: 2}, Length: 2}},
		{name: "EmbeddedClusterWestTest", m: embeddedClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: pcpb.Direction_DIRECTION_WEST, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 2, Y: 2}, Length: 2}},
		{name: "RectangularClusterNorthTest", m: rectangularClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: pcpb.Direction_DIRECTION_NORTH, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 1, Y: 3}, Length: 1}},
		{name: "RectangularClusterSouthTest", m: rectangularClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: pcpb.Direction_DIRECTION_SOUTH, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 1, Y: 2}, Length: 1}},
		{name: "RectangularClusterEastTest", m: rectangularClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: pcpb.Direction_DIRECTION_EAST, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 2}, Length: 2}},
		{name: "RectangularClusterWestTest", m: rectangularClusterMap, c: utils.MapCoordinate{X: 1, Y: 1}, d: pcpb.Direction_DIRECTION_WEST, want: coordinateSlice{
			Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 2}, Length: 2}},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			m, err := cluster.ImportMap(c.m)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			if got, err := buildClusterEdgeCoordinateSlice(m, c.c, c.d); err != nil || !cmp.Equal(got, c.want, protocmp.Transform()) {
				t.Errorf("buildClusterEdgeCoordinateSlice() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestBuildCoordinateWithCoordinateSlice(t *testing.T) {
	testConfigs := []struct {
		name   string
		s      coordinateSlice
		offset int32
		want   *gdpb.Coordinate
	}{
		{
			name:   "SingleTileSliceHorizontal",
			s:      coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			offset: 0,
			want:   &gdpb.Coordinate{X: 0, Y: 0},
		},
		{
			name:   "SingleTileSliceVertical",
			s:      coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			offset: 0,
			want:   &gdpb.Coordinate{X: 0, Y: 0},
		},
		{
			name:   "MultiTileTileSliceHorizontal",
			s:      coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 1, Y: 1}, Length: 2},
			offset: 1,
			want:   &gdpb.Coordinate{X: 2, Y: 1},
		},
		{
			name:   "MultiTileTileSliceVertical",
			s:      coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 1}, Length: 2},
			offset: 1,
			want:   &gdpb.Coordinate{X: 1, Y: 2},
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
		s      coordinateSlice
		offset int32
	}{
		{name: "NullTileSlice", s: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 0}, offset: 0},
		{name: "OutOfBoundsTileSliceBefore", s: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1}, offset: -1},
		{name: "OutOfBoundsTileSliceAfter", s: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1}, offset: 2},
		{name: "InvalidOrientationTileSlice", s: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_UNKNOWN, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1}, offset: 0},
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
		s1, s2 coordinateSlice
		want   []Transition
	}{
		{
			name: "SingleWidthEntranceHorizontal",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 1}, Length: 1},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 1}},
				},
			},
		},
		{
			name: "SingleWidthEntranceVertical",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 0}, Length: 1},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 0}},
				},
			},
		},
		{
			name: "DoubleWidthEntranceHorizontal",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 1}, Length: 2},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 1}},
				},
			},
		},
		{
			name: "DoubleWidthEntranceVertical",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 0}, Length: 2},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 1}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 1}},
				},
			},
		},
		{
			name: "TripleWidthEntranceHorizontal",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 3},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 1}, Length: 3},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 1}},
				},
			},
		},
		{
			name: "TripleWidthEntranceVertical",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 3},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 0}, Length: 3},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 1}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 1}},
				},
			},
		},
		{
			name: "QuadrupleWidthEntranceHorizontal",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 1}, Length: 4},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 1}},
				},
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 3, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 3, Y: 1}},
				},
			},
		},
		{
			name: "QuadrupleWidthEntranceVertical",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 0}, Length: 4},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 0}},
				},
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 3}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 3}},
				},
			},
		},
		{
			name: "QuadrupleWidthEmbeddedEntranceHorizontal",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 1, Y: 1}, Length: 4},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 1, Y: 2}, Length: 4},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 1}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 2}},
				},
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 4, Y: 1}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 4, Y: 2}},
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
		s1, s2 coordinateSlice
	}{
		{
			name: "MismatchedLengths",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 1}, Length: 2},
		},
		{
			name: "MismatchedOrientations",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
		},
		{
			name: "NonAdjacentHorizontalSlice",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 2}, Length: 1},
		},
		{
			name: "NonAdjacentVerticalSlice",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 2, Y: 0}, Length: 1},
		},
		{
			name: "NonAlignedHorizontalSlice",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 1, Y: 1}, Length: 2},
		},
		{
			name: "NonAlignedVerticalSlice",
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 1}, Length: 2},
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
	trivialOpenClusterMap := &pdpb.ClusterMap{
		TileDimension:    &gdpb.Coordinate{X: 1, Y: 1},
		TileMapDimension: trivialOpenMap.GetDimension(),
	}
	longVerticalOpenClusterMap := &pdpb.ClusterMap{
		TileDimension:    &gdpb.Coordinate{X: 2, Y: 1},
		TileMapDimension: longVerticalOpenMap.GetDimension(),
	}

	testConfigs := []struct {
		name   string
		m      *mdpb.TileMap
		cm     *pdpb.ClusterMap
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
		m      *mdpb.TileMap
		s1, s2 coordinateSlice
		want   []Transition
	}{
		{name: "TrivialClosedMap", m: trivialClosedMap,
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 0}, Length: 1},
			want: nil,
		},
		{name: "TrivialSemiOpenMap", m: trivialSemiOpenMap,
			s1:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 0}, Length: 1},
			want: nil,
		},
		{name: "TrivialOpenMap", m: trivialOpenMap,
			s1: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 0}, Length: 1},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 0}},
				},
			},
		},
		{name: "LongVerticalOpenMap", m: longVerticalOpenMap,
			s1: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 0}, Length: 4},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 0}},
				},
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 3}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 3}},
				},
			},
		},
		{name: "LongHorizontalOpenMap", m: longHorizontalOpenMap,
			s1: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 1}, Length: 4},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 1}},
				},
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 3, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 3, Y: 1}},
				},
			},
		},
		{name: "LongSemiOpenMap", m: longSemiOpenMap,
			s1: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 3},
			s2: coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 1, Y: 0}, Length: 3},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 0}},
				},
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 2}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 2}},
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
	trivialClusterMap := &pdpb.ClusterMap{TileDimension: &gdpb.Coordinate{X: 1, Y: 1}, TileMapDimension: trivialClosedMap.GetDimension()}
	longVerticalClusterMap := &pdpb.ClusterMap{TileDimension: &gdpb.Coordinate{X: 1, Y: 4}, TileMapDimension: longVerticalOpenMap.GetDimension()}
	longHorizontalClusterMap := &pdpb.ClusterMap{TileDimension: &gdpb.Coordinate{X: 4, Y: 1}, TileMapDimension: longHorizontalOpenMap.GetDimension()}
	longSemiOpenClusterMap := &pdpb.ClusterMap{TileDimension: &gdpb.Coordinate{X: 1, Y: 3}, TileMapDimension: longSemiOpenMap.GetDimension()}

	testConfigs := []struct {
		name   string
		m      *mdpb.TileMap
		cm     *pdpb.ClusterMap
		c1, c2 *gdpb.Coordinate
		want   []Transition
	}{
		{name: "TrivialClosedMap", m: trivialClosedMap, cm: trivialClusterMap, c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 1, Y: 0}, want: nil},
		{name: "TrivialSemiOpenMap", m: trivialSemiOpenMap, cm: trivialClusterMap, c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 1, Y: 0}, want: nil},
		{name: "TrivialOpenMap", m: trivialOpenMap, cm: trivialClusterMap, c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 1, Y: 0},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 0}},
				},
			},
		},
		{name: "LongVerticalOpenMap", m: longVerticalOpenMap, cm: longVerticalClusterMap, c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 1, Y: 0},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 0}},
				},
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 3}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 3}},
				},
			},
		},
		{name: "LongHorizontalOpenMap", m: longHorizontalOpenMap, cm: longHorizontalClusterMap, c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 0, Y: 1},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 1}},
				},
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 3, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 3, Y: 1}},
				},
			},
		},
		{name: "LongSemiOpenMap", m: longSemiOpenMap, cm: longSemiOpenClusterMap, c1: &gdpb.Coordinate{X: 0, Y: 0}, c2: &gdpb.Coordinate{X: 1, Y: 0},
			want: []Transition{
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 0}},
				},
				{
					N1: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 2}},
					N2: &pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 1, Y: 2}},
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
	s := coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_UNKNOWN, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1}
	if _, err := sliceContains(s, utils.MC(&gdpb.Coordinate{X: 0, Y: 0})); err == nil {
		t.Error("sliceContains() = _, nil, want a non-nil error")
	}
}

func TestSliceContains(t *testing.T) {
	testConfigs := []struct {
		name string
		s    coordinateSlice
		c    *gdpb.Coordinate
		want bool
	}{
		{
			name: "TrivialSliceContains",
			s:    coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			c:    &gdpb.Coordinate{X: 0, Y: 0},
			want: true,
		},
		{
			name: "TrivialPreSliceNoContains",
			s:    coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			c:    &gdpb.Coordinate{X: -1, Y: 0},
			want: false,
		},
		{
			name: "TrivialPostSliceNoContains",
			s:    coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			c:    &gdpb.Coordinate{X: 1, Y: 0},
			want: false,
		},
		{
			name: "TrivialBadAxisSliceNoContains",
			s:    coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 1},
			c:    &gdpb.Coordinate{X: 0, Y: -1},
			want: false,
		},
		{
			name: "SimpleSliceContainsHorizontal",
			s:    coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_HORIZONTAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 2},
			c:    &gdpb.Coordinate{X: 1, Y: 0},
			want: true,
		},
		{
			name: "SimpleSliceContainsVertical",
			s:    coordinateSlice{Orientation: pcpb.Orientation_ORIENTATION_VERTICAL, Start: &gdpb.Coordinate{X: 0, Y: 0}, Length: 2},
			c:    &gdpb.Coordinate{X: 0, Y: 1},
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
		m    *pdpb.ClusterMap
		c    *gdpb.Coordinate
		t    *gdpb.Coordinate
		want bool
	}{
		{
			name: "TrivialClusterContains",
			m:    &pdpb.ClusterMap{TileDimension: &gdpb.Coordinate{X: 1, Y: 1}, TileMapDimension: &gdpb.Coordinate{X: 1, Y: 1}},
			c:    &gdpb.Coordinate{X: 0, Y: 0},
			t:    &gdpb.Coordinate{X: 0, Y: 0},
			want: true,
		},
		{
			name: "TrivialClusterNoContains",
			m:    &pdpb.ClusterMap{TileDimension: &gdpb.Coordinate{X: 1, Y: 1}, TileMapDimension: &gdpb.Coordinate{X: 2, Y: 2}},
			c:    &gdpb.Coordinate{X: 0, Y: 0},
			t:    &gdpb.Coordinate{X: 0, Y: 1},
			want: false,
		},
		{
			name: "ClusterInternalNoContains",
			m:    &pdpb.ClusterMap{TileDimension: &gdpb.Coordinate{X: 3, Y: 3}, TileMapDimension: &gdpb.Coordinate{X: 3, Y: 3}},
			c:    &gdpb.Coordinate{X: 0, Y: 0},
			t:    &gdpb.Coordinate{X: 1, Y: 1},
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
