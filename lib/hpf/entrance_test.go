package entrance

import (
	"testing"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/cluster"
	"github.com/cripplet/rts-pathing/lib/hpf/tile"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
)

func tileComparator(t, other *tile.Tile) bool {
	return proto.Equal(t.Coordinate(), other.Coordinate())
}

func (s *ClusterBorderSegment) Equal(other *ClusterBorderSegment) bool {
	return proto.Equal(s.s, other.s)
}

func TestCandidateVectorError(t *testing.T) {
	testConfigs := []struct {
		name string
		c    *cluster.Cluster
		d    rtscpb.Direction
	}{
		{name: "NullClusterTest", c: cluster.NewCluster(&rtsspb.Cluster{}), d: rtscpb.Direction_DIRECTION_NORTH},
		{name: "NullXDimensionClusterTest", c: cluster.NewCluster(
			&rtsspb.Cluster{
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 0, Y: 1},
			}), d: rtscpb.Direction_DIRECTION_NORTH,
		},
		{name: "NullYDimensionClusterTest", c: cluster.NewCluster(
			&rtsspb.Cluster{
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 0},
			}), d: rtscpb.Direction_DIRECTION_NORTH,
		},
		{name: "InvalidDirectionTest", c: cluster.NewCluster(
			&rtsspb.Cluster{
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			}), d: rtscpb.Direction_DIRECTION_UNKNOWN,
		},
	}
	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if start, end, err := candidateVector(c.c, c.d); err == nil {
				t.Errorf("candidateVector() = %v, %v, %v, want a non-nil error", start, end, err)
			}
		})
	}
}

func TestCandidateVector(t *testing.T) {
	trivialCluster := cluster.NewCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
		TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
	})
	smallCluster := cluster.NewCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
		TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
	})
	embeddedCluster := cluster.NewCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 1},
		TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
	})
	rectangularCluster := cluster.NewCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 1},
		TileDimension: &rtsspb.Coordinate{X: 1, Y: 2},
	})
	testConfigs := []struct {
		name               string
		c                  *cluster.Cluster
		d                  rtscpb.Direction
		wantStart, wantEnd *rtsspb.Coordinate
	}{
		{name: "TrivialClusterNorthTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_NORTH, wantStart: &rtsspb.Coordinate{X: 0, Y: 0}, wantEnd: &rtsspb.Coordinate{X: 0, Y: 0}},
		{name: "TrivialClusterSouthTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_SOUTH, wantStart: &rtsspb.Coordinate{X: 0, Y: 0}, wantEnd: &rtsspb.Coordinate{X: 0, Y: 0}},
		{name: "TrivialClusterEastTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_EAST, wantStart: &rtsspb.Coordinate{X: 0, Y: 0}, wantEnd: &rtsspb.Coordinate{X: 0, Y: 0}},
		{name: "TrivialClusterWestTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_WEST, wantStart: &rtsspb.Coordinate{X: 0, Y: 0}, wantEnd: &rtsspb.Coordinate{X: 0, Y: 0}},
		{name: "SmallClusterNorthTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_NORTH, wantStart: &rtsspb.Coordinate{X: 0, Y: 1}, wantEnd: &rtsspb.Coordinate{X: 1, Y: 1}},
		{name: "SmallClusterSouthTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_SOUTH, wantStart: &rtsspb.Coordinate{X: 0, Y: 0}, wantEnd: &rtsspb.Coordinate{X: 1, Y: 0}},
		{name: "SmallClusterEastTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_EAST, wantStart: &rtsspb.Coordinate{X: 1, Y: 0}, wantEnd: &rtsspb.Coordinate{X: 1, Y: 1}},
		{name: "SmallClusterWestTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_WEST, wantStart: &rtsspb.Coordinate{X: 0, Y: 0}, wantEnd: &rtsspb.Coordinate{X: 0, Y: 1}},
		{name: "EmbeddedClusterNorthTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_NORTH, wantStart: &rtsspb.Coordinate{X: 1, Y: 2}, wantEnd: &rtsspb.Coordinate{X: 2, Y: 2}},
		{name: "EmbeddedClusterSouthTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_SOUTH, wantStart: &rtsspb.Coordinate{X: 1, Y: 1}, wantEnd: &rtsspb.Coordinate{X: 2, Y: 1}},
		{name: "EmbeddedClusterEastTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_EAST, wantStart: &rtsspb.Coordinate{X: 2, Y: 1}, wantEnd: &rtsspb.Coordinate{X: 2, Y: 2}},
		{name: "EmbeddedClusterWestTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_WEST, wantStart: &rtsspb.Coordinate{X: 1, Y: 1}, wantEnd: &rtsspb.Coordinate{X: 1, Y: 2}},
		{name: "RectangularClusterNorthTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_NORTH, wantStart: &rtsspb.Coordinate{X: 1, Y: 2}, wantEnd: &rtsspb.Coordinate{X: 1, Y: 2}},
		{name: "RectangularClusterSouthTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_SOUTH, wantStart: &rtsspb.Coordinate{X: 1, Y: 1}, wantEnd: &rtsspb.Coordinate{X: 1, Y: 1}},
		{name: "RectangularClusterEastTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_EAST, wantStart: &rtsspb.Coordinate{X: 1, Y: 1}, wantEnd: &rtsspb.Coordinate{X: 1, Y: 2}},
		{name: "RectangularClusterWestTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_WEST, wantStart: &rtsspb.Coordinate{X: 1, Y: 1}, wantEnd: &rtsspb.Coordinate{X: 1, Y: 2}},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if start, end, err := candidateVector(c.c, c.d); err != nil || !proto.Equal(start, c.wantStart) || !proto.Equal(end, c.wantEnd) {
				t.Errorf("candidateVector() = %v, %v, %v, want = %v, %v, nil", start, end, err, c.wantStart, c.wantEnd)
			}
		})
	}
}

func TestCandidateVectorTilesError(t *testing.T) {
	m, err := tile.ImportTileMap(
		&rtsspb.TileMap{
			Dimension: &rtsspb.Coordinate{X: 1, Y: 1},
		},
	)
	if err != nil {
		t.Fatalf("ImportTileMap() = %v, want = nil", err)
	}
	if res, err := candidateVectorTiles(&rtsspb.Coordinate{X: 0, Y: 0}, &rtsspb.Coordinate{X: 0, Y: 0}, m); err == nil {
		t.Errorf("candidateVectorTiles() = %v, nil, want a non-nil error", res)
	}
}

func TestCandidateVectorTiles(t *testing.T) {
	trivialMap := &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 1, Y: 1},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
		},
	}
	squareMap := &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 2, Y: 2},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 1}},
		},
	}

	testConfigs := []struct {
		name       string
		m          *rtsspb.TileMap
		start, end *rtsspb.Coordinate
		want       []*rtsspb.Coordinate
	}{
		{name: "TrivialMapTest", m: trivialMap,
			start: &rtsspb.Coordinate{X: 0, Y: 0},
			end:   &rtsspb.Coordinate{X: 0, Y: 0},
			want: []*rtsspb.Coordinate{
				{X: 0, Y: 0},
			},
		},
		{name: "SquareMapHorizontalTest", m: squareMap,
			start: &rtsspb.Coordinate{X: 0, Y: 0},
			end:   &rtsspb.Coordinate{X: 1, Y: 0},
			want: []*rtsspb.Coordinate{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
			},
		},
		{name: "SquareMapVerticalTest", m: squareMap,
			start: &rtsspb.Coordinate{X: 0, Y: 0},
			end:   &rtsspb.Coordinate{X: 0, Y: 1},
			want: []*rtsspb.Coordinate{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			m, err := tile.ImportTileMap(c.m)
			if err != nil {
				t.Errorf("ImportTileMap() = %v, want = nil", err)
				return
			}

			var wantTiles []*tile.Tile
			for _, coord := range c.want {
				wt, err := tile.ImportTile(&rtsspb.Tile{Coordinate: coord})
				if err != nil {
					t.Errorf("ImportTile() = %v, want = nil", err)
					return
				}
				wantTiles = append(wantTiles, wt)
			}

			if res, err := candidateVectorTiles(c.start, c.end, m); err != nil || !cmp.Equal(res, wantTiles, cmp.Comparer(tileComparator)) {
				t.Errorf("candidateVectorTiles() = %v, %v, want = %v, nil", res, err, wantTiles)
			}
		})
	}
}

func TestSegments(t *testing.T) {
	testConfigs := []struct {
		name string
		ts   []*rtsspb.Tile
		want []*ClusterBorderSegment
	}{
		{name: "EmptySegmentsTest", ts: []*rtsspb.Tile{}, want: nil},
		{name: "TrivialSegmentsTest", ts: []*rtsspb.Tile{
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		}, want: []*ClusterBorderSegment{
			&ClusterBorderSegment{s: &rtsspb.ClusterBorderSegment{Start: &rtsspb.Coordinate{X: 0, Y: 0}, End: &rtsspb.Coordinate{X: 0, Y: 0}}},
		}},
		{name: "AllBlockedSegmentsTest", ts: []*rtsspb.Tile{
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
		}, want: nil},
		{name: "LongerContinuousYSegmentsTest", ts: []*rtsspb.Tile{
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		}, want: []*ClusterBorderSegment{
			&ClusterBorderSegment{s: &rtsspb.ClusterBorderSegment{Start: &rtsspb.Coordinate{X: 0, Y: 0}, End: &rtsspb.Coordinate{X: 0, Y: 2}}},
		}},
		{name: "LongerContinuousXSegmentsTest", ts: []*rtsspb.Tile{
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 2, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		}, want: []*ClusterBorderSegment{
			&ClusterBorderSegment{s: &rtsspb.ClusterBorderSegment{Start: &rtsspb.Coordinate{X: 0, Y: 0}, End: &rtsspb.Coordinate{X: 2, Y: 0}}},
		}},
		{name: "BlockedSegmentsTest", ts: []*rtsspb.Tile{
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		}, want: []*ClusterBorderSegment{
			&ClusterBorderSegment{s: &rtsspb.ClusterBorderSegment{Start: &rtsspb.Coordinate{X: 0, Y: 0}, End: &rtsspb.Coordinate{X: 0, Y: 0}}},
			&ClusterBorderSegment{s: &rtsspb.ClusterBorderSegment{Start: &rtsspb.Coordinate{X: 0, Y: 2}, End: &rtsspb.Coordinate{X: 0, Y: 2}}},
		}},
		{name: "BlockedStartSegmentsTest", ts: []*rtsspb.Tile{
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		}, want: []*ClusterBorderSegment{
			&ClusterBorderSegment{s: &rtsspb.ClusterBorderSegment{Start: &rtsspb.Coordinate{X: 0, Y: 1}, End: &rtsspb.Coordinate{X: 0, Y: 2}}},
		}},
		{name: "BlockedEndSegmentsTest", ts: []*rtsspb.Tile{
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			&rtsspb.Tile{Coordinate: &rtsspb.Coordinate{X: 0, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
		}, want: []*ClusterBorderSegment{
			&ClusterBorderSegment{s: &rtsspb.ClusterBorderSegment{Start: &rtsspb.Coordinate{X: 0, Y: 0}, End: &rtsspb.Coordinate{X: 0, Y: 1}}},
		}},
	}

	for _, c := range testConfigs {
		var tiles []*tile.Tile
		for _, pb := range c.ts {
			obj, err := tile.ImportTile(pb)
			if err != nil {
				t.Errorf("ImportTile() = %v, want = nil", err)
				return
			}
			tiles = append(tiles, obj)
		}

		t.Run(c.name, func(t *testing.T) {
			if res, err := segments(tiles); err != nil || !cmp.Equal(res, c.want) {
				t.Errorf("segments() = %v, %v, want = %v, nil", res, err, c.want)
			}
		})
	}
}
