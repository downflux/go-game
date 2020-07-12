package entrance

import (
	"testing"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/cluster"
	"github.com/cripplet/rts-pathing/lib/hpf/tile"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
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
		Dimension: &rtsspb.Coordinate{X: 4, Y: 4},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
	}
)

func TestbuildClusterEdgeCoordinateSliceError(t *testing.T) {
	testConfigs := []struct {
		name string
		c    *rtsspb.Cluster
		d    rtscpb.Direction
	}{
		{name: "NullClusterTest", c: &rtsspb.Cluster{}, d: rtscpb.Direction_DIRECTION_NORTH},
		{name: "NullXDimensionClusterTest", c: &rtsspb.Cluster{
			TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
			TileDimension: &rtsspb.Coordinate{X: 0, Y: 1},
		}, d: rtscpb.Direction_DIRECTION_NORTH,
		},
		{name: "NullYDimensionClusterTest", c: &rtsspb.Cluster{
			TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
			TileDimension: &rtsspb.Coordinate{X: 1, Y: 0},
		}, d: rtscpb.Direction_DIRECTION_NORTH,
		},
		{name: "InvalidDirectionTest", c: &rtsspb.Cluster{
			TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
			TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
		}, d: rtscpb.Direction_DIRECTION_UNKNOWN,
		},
	}
	for _, c := range testConfigs {
		tmpCluster, err := cluster.ImportCluster(c.c)
		if err != nil {
			t.Fatalf("ImportCluster() = _, %v, want = _, nil", err)
		}

		t.Run(c.name, func(t *testing.T) {
			if got, err := buildClusterEdgeCoordinateSlice(tmpCluster, c.d); err == nil {
				t.Errorf("buildClusterEdgeCoordinateSlice() = %v, %v, want a non-nil error", got, err)
			}
		})
	}
}

func TestBuildClusterEdgeCoordinateSlice(t *testing.T) {
	trivialCluster := &rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
		TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
	}
	smallCluster := &rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
		TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
	}
	embeddedCluster := &rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 1},
		TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
	}
	rectangularCluster := &rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 1},
		TileDimension: &rtsspb.Coordinate{X: 1, Y: 2},
	}
	testConfigs := []struct {
		name string
		c    *rtsspb.Cluster
		d    rtscpb.Direction
		want *rtsspb.CoordinateSlice
	}{
		{name: "TrivialClusterNorthTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterSouthTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterEastTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterWestTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "SmallClusterNorthTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 2}},
		{name: "SmallClusterSouthTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2}},
		{name: "SmallClusterEastTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 2}},
		{name: "SmallClusterWestTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2}},
		{name: "EmbeddedClusterNorthTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, Length: 2}},
		{name: "EmbeddedClusterSouthTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2}},
		{name: "EmbeddedClusterEastTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 2, Y: 1}, Length: 2}},
		{name: "EmbeddedClusterWestTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2}},
		{name: "RectangularClusterNorthTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, Length: 1}},
		{name: "RectangularClusterSouthTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 1}},
		{name: "RectangularClusterEastTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2}},
		{name: "RectangularClusterWestTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2}},
	}

	for _, c := range testConfigs {
		tmpCluster, err := cluster.ImportCluster(c.c)
		if err != nil {
			t.Fatalf("ImportCluster() = _, %v, want = _, nil", err)
		}

		t.Run(c.name, func(t *testing.T) {
			if got, err := buildClusterEdgeCoordinateSlice(tmpCluster, c.d); err != nil || !proto.Equal(got, c.want) {
				t.Errorf("buildClusterEdgeCoordinateSlice() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestBuildCoordinateWithCoordinateSlice(t *testing.T) {
	testConfigs := []struct {
		name string
		s    *rtsspb.CoordinateSlice
		o    int32
		want *rtsspb.Coordinate
	}{
		{
			name: "SingleTileSliceHorizontal",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			o:    0,
			want: &rtsspb.Coordinate{X: 0, Y: 0},
		},
		{
			name: "SingleTileSliceVertical",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			o:    0,
			want: &rtsspb.Coordinate{X: 0, Y: 0},
		},
		{
			name: "MultiTileTileSliceHorizontal",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2},
			o:    1,
			want: &rtsspb.Coordinate{X: 2, Y: 1},
		},
		{
			name: "MultiTileTileSliceVertical",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2},
			o:    1,
			want: &rtsspb.Coordinate{X: 1, Y: 2},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := buildCoordinateWithCoordinateSlice(c.s, c.o); err != nil || !proto.Equal(got, c.want) {
				t.Errorf("buildCoordinateWithCoordinateSlice() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestBuildCoordinateWithCoordinateSliceError(t *testing.T) {
	testConfigs := []struct {
		name string
		s    *rtsspb.CoordinateSlice
		o    int32
	}{
		{name: "NullTileSlice", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 0}, o: 0},
		{name: "OutOfBoundsTileSliceBefore", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}, o: -1},
		{name: "OutOfBoundsTileSliceAfter", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}, o: 2},
		{name: "InvalidOrientationTileSlice", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_UNKNOWN, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}, o: 0},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if _, err := buildCoordinateWithCoordinateSlice(c.s, c.o); err == nil {
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
				{C1: &rtsspb.Coordinate{X: 0, Y: 0}, C2: &rtsspb.Coordinate{X: 0, Y: 1}},
			},
		},
		{
			name: "SingleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 1},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 0, Y: 0}, C2: &rtsspb.Coordinate{X: 1, Y: 0}},
			},
		},
		{
			name: "DoubleWidthEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 2},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 1, Y: 0}, C2: &rtsspb.Coordinate{X: 1, Y: 1}},
			},
		},
		{
			name: "DoubleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 2},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 0, Y: 1}, C2: &rtsspb.Coordinate{X: 1, Y: 1}},
			},
		},
		{
			name: "TripleWidthEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 3},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 3},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 1, Y: 0}, C2: &rtsspb.Coordinate{X: 1, Y: 1}},
			},
		},
		{
			name: "TripleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 3},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 3},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 0, Y: 1}, C2: &rtsspb.Coordinate{X: 1, Y: 1}},
			},
		},
		{
			name: "QuadrupleWidthEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 4},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 0, Y: 0}, C2: &rtsspb.Coordinate{X: 0, Y: 1}},
				{C1: &rtsspb.Coordinate{X: 3, Y: 0}, C2: &rtsspb.Coordinate{X: 3, Y: 1}},
			},
		},
		{
			name: "QuadrupleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 4},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 0, Y: 0}, C2: &rtsspb.Coordinate{X: 1, Y: 0}},
				{C1: &rtsspb.Coordinate{X: 0, Y: 3}, C2: &rtsspb.Coordinate{X: 1, Y: 3}},
			},
		},
		{
			name: "QuadrupleWidthEmbeddedEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 4},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, Length: 4},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 1, Y: 1}, C2: &rtsspb.Coordinate{X: 1, Y: 2}},
				{C1: &rtsspb.Coordinate{X: 4, Y: 1}, C2: &rtsspb.Coordinate{X: 4, Y: 2}},
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
	m, err := tile.ImportTileMap(trivialOpenMap)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil")
	}
	longM, err := tile.ImportTileMap(longVerticalOpenMap)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil")
	}

	testConfigs := []struct {
		name   string
		m      *tile.TileMap
		c1, c2 *cluster.Cluster
	}{
		{name: "NullCluster", m: m, c1: nil, c2: nil},
		{name: "NullMap", m: nil, c1: &cluster.Cluster{
			Val: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			}}, c2: &cluster.Cluster{
			Val: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 1, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			}},
		},
		{name: "NonAdjacentClusters", m: longM,
			c1: &cluster.Cluster{
				Val: &rtsspb.Cluster{
					Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
					TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
					TileDimension: &rtsspb.Coordinate{X: 2, Y: 1},
				},
			},
			c2: &cluster.Cluster{
				Val: &rtsspb.Cluster{
					Coordinate:    &rtsspb.Coordinate{X: 0, Y: 2},
					TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 2},
					TileDimension: &rtsspb.Coordinate{X: 2, Y: 1},
				},
			},
		},
	}
	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := BuildTransitions(c.c1, c.c2, c.m); err == nil {
				t.Errorf("BuildTransitions() = %v, %v, want a non-nil error", got, err)
			}
		})
	}
}

func TestBuildTransitions(t *testing.T) {
	testConfigs := []struct {
		name   string
		m      *rtsspb.TileMap
		c1, c2 *rtsspb.Cluster
		want   []*rtsspb.Transition
	}{
		{name: "TrivialClosedMap", m: trivialClosedMap,
			c1: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			},
			c2: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 1, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			},
			want: nil,
		},
		{name: "TrivialSemiOpenMap", m: trivialSemiOpenMap,
			c1: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			},
			c2: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 1, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			},
			want: nil,
		},
		{name: "TrivialOpenMap", m: trivialOpenMap,
			c1: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			},
			c2: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 1, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			},
			want: []*rtsspb.Transition{
				{
					C1: &rtsspb.Coordinate{X: 0, Y: 0},
					C2: &rtsspb.Coordinate{X: 1, Y: 0},
				},
			},
		},
		{name: "LongVerticalOpenMap", m: longVerticalOpenMap,
			c1: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 4},
			},
			c2: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 1, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 4},
			},
			want: []*rtsspb.Transition{
				{
					C1: &rtsspb.Coordinate{X: 0, Y: 0},
					C2: &rtsspb.Coordinate{X: 1, Y: 0},
				},
				{
					C1: &rtsspb.Coordinate{X: 0, Y: 3},
					C2: &rtsspb.Coordinate{X: 1, Y: 3},
				},
			},
		},
		{name: "LongHorizontalOpenMap", m: longHorizontalOpenMap,
			c1: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 4, Y: 1},
			},
			c2: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 0, Y: 1},
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 1},
				TileDimension: &rtsspb.Coordinate{X: 4, Y: 1},
			},
			want: []*rtsspb.Transition{
				{
					C1: &rtsspb.Coordinate{X: 0, Y: 0},
					C2: &rtsspb.Coordinate{X: 0, Y: 1},
				},
				{
					C1: &rtsspb.Coordinate{X: 3, Y: 0},
					C2: &rtsspb.Coordinate{X: 3, Y: 1},
				},
			},
		},
		{name: "LongSemiOpenMap", m: longSemiOpenMap,
			c1: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 3},
			},
			c2: &rtsspb.Cluster{
				Coordinate:    &rtsspb.Coordinate{X: 1, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 3},
			},
			want: []*rtsspb.Transition{
				{
					C1: &rtsspb.Coordinate{X: 0, Y: 0},
					C2: &rtsspb.Coordinate{X: 1, Y: 0},
				},
				{
					C1: &rtsspb.Coordinate{X: 0, Y: 2},
					C2: &rtsspb.Coordinate{X: 1, Y: 2},
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			cluster1, err := cluster.ImportCluster(c.c1)
			if err != nil {
				t.Fatalf("ImportCluster() = _, %v, want = _, nil", err)
			}
			cluster2, err := cluster.ImportCluster(c.c2)
			if err != nil {
				t.Fatalf("ImportCluster() = _, %v, want = _, nil", err)
			}
			tileMap, err := tile.ImportTileMap(c.m)
			if err != nil {
				t.Fatalf("ImportTileMap() = _, %v, want = _, nil", err)
			}

			if got, err := BuildTransitions(cluster1, cluster2, tileMap); err != nil || !cmp.Equal(got, c.want, protocmp.Transform()) {
				t.Errorf("BuildTransitions() = %v, %v, want = %v, nil", got, err, c.want)
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
					C1: &rtsspb.Coordinate{X: 0, Y: 0},
					C2: &rtsspb.Coordinate{X: 1, Y: 0},
				},
			},
		},
		{name: "LongVerticalOpenMap", m: longVerticalOpenMap,
			s1: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 4},
			want: []*rtsspb.Transition{
				{
					C1: &rtsspb.Coordinate{X: 0, Y: 0},
					C2: &rtsspb.Coordinate{X: 1, Y: 0},
				},
				{
					C1: &rtsspb.Coordinate{X: 0, Y: 3},
					C2: &rtsspb.Coordinate{X: 1, Y: 3},
				},
			},
		},
		{name: "LongHorizontalOpenMap", m: longHorizontalOpenMap,
			s1: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 4},
			want: []*rtsspb.Transition{
				{
					C1: &rtsspb.Coordinate{X: 0, Y: 0},
					C2: &rtsspb.Coordinate{X: 0, Y: 1},
				},
				{
					C1: &rtsspb.Coordinate{X: 3, Y: 0},
					C2: &rtsspb.Coordinate{X: 3, Y: 1},
				},
			},
		},
		{name: "LongSemiOpenMap", m: longSemiOpenMap,
			s1: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 3},
			s2: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 3},
			want: []*rtsspb.Transition{
				{
					C1: &rtsspb.Coordinate{X: 0, Y: 0},
					C2: &rtsspb.Coordinate{X: 1, Y: 0},
				},
				{
					C1: &rtsspb.Coordinate{X: 0, Y: 2},
					C2: &rtsspb.Coordinate{X: 1, Y: 2},
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tileMap, err := tile.ImportTileMap(c.m)
			if err != nil {
				t.Fatalf("ImportTileMap() = _, %v, want = _, nil", err)
			}

			if got, err := buildTransitionsAux(c.s1, c.s2, tileMap); err != nil || !cmp.Equal(got, c.want, protocmp.Transform()) {
				t.Errorf("buildTransitionsAux() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}
