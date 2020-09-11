package abstractgraph

import (
	"testing"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/abstractedgemap"
	// "github.com/cripplet/rts-pathing/lib/hpf/abstractnodemap"
	"github.com/cripplet/rts-pathing/lib/hpf/cluster"
	"github.com/cripplet/rts-pathing/lib/hpf/tile"
	"github.com/cripplet/rts-pathing/lib/hpf/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

	/**
	 *       - - - - - -
	 *       - - - - - -
	 *       - - - - - -
	 *       - - - - - -
	 *       - - - - - -
	 * Y = 0 - - - - - -
	 *   X = 0
	 */
	largeMapProto = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 6, Y: 6},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 3}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 4}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 5}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 3}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 4}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 1, Y: 5}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 3}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 4}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 2, Y: 5}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 3, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 3, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 3, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 3, Y: 3}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 3, Y: 4}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 3, Y: 5}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 4, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 4, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 4, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 4, Y: 3}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 4, Y: 4}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 4, Y: 5}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 5, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 5, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 5, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 5, Y: 3}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 5, Y: 4}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 5, Y: 5}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
	}
)

func coordLess(c1, c2 *rtsspb.Coordinate) bool {
	return c1.GetX() < c2.GetX() || (c1.GetX() == c2.GetX() && c1.GetY() < c2.GetY())
}

func nodeLess(n1, n2 *rtsspb.AbstractNode) bool {
	return coordLess(n1.GetTileCoordinate(), n2.GetTileCoordinate())
}

func transitionLess(t1, t2 *rtsspb.Transition) bool {
	return nodeLess(t1.GetN1(), t2.GetN1())
}

func edgeLess(e1, e2 *rtsspb.AbstractEdge) bool {
	return coordLess(e1.GetSource(), e2.GetSource()) || cmp.Equal(
		e1.GetSource(),
		e2.GetSource(),
		protocmp.Transform()) && coordLess(e1.GetDestination(), e2.GetDestination())
}

func abstractEdgeEqual(e1, e2 *rtsspb.AbstractEdge) bool {
	if cmp.Equal(e1, e2, protocmp.Transform()) {
		return true
	}

	return cmp.Equal(
		e1.GetSource(),
		e2.GetDestination(),
		protocmp.Transform(),
	) && cmp.Equal(
		e1.GetDestination(),
		e2.GetSource(),
		protocmp.Transform(),
	) && cmp.Equal(
		e1,
		e2,
		protocmp.Transform(),
		protocmp.IgnoreFields(&rtsspb.AbstractEdge{}, "source", "destination"),
	)
}

func abstractEdgeMapEqual(em1, em2 abstractedgemap.Map) bool {
	for _, e1 := range em1.Iterator() {
		e2, err := em2.Get(utils.MC(e1.GetSource()), utils.MC(e1.GetDestination()))
		if err != nil || e2 == nil {
			return false
		}
		if !cmp.Equal(e1, e2, cmp.Comparer(abstractEdgeEqual)) {
			return false
		}
	}

	for _, e2 := range em2.Iterator() {
		e1, err := em1.Get(utils.MC(e2.GetSource()), utils.MC(e2.GetDestination()))
		if err != nil || e1 == nil {
			return false
		}
		if !cmp.Equal(e1, e2, cmp.Comparer(abstractEdgeEqual)) {
			return false
		}
	}
	return true
}

/*
func TestBuildTieredClusterMapsError(t *testing.T) {
	testConfigs := []struct {
		name string
		tm   *rtsspb.TileMap
		l    int32
		dim  *rtsspb.Coordinate
	}{
		{name: "NilMapError", tm: nil, l: 1, dim: &rtsspb.Coordinate{X: 2, Y: 2}},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tm, err := tile.ImportTileMap(c.tm)
			if err != nil {
				t.Fatalf("ImportTileMap() = _, %v, want = _, nil", err)
			}

			if _, err := buildTieredClusterMaps(tm, c.l, c.dim); err == nil {
				t.Fatal("buildTieredClusterMaps() = _, nil, want a non-nil error")
			}
		})
	}
}

func TestBuildTieredClusterMaps(t *testing.T) {
	tm := &rtsspb.TileMap{Dimension: &rtsspb.Coordinate{X: 8, Y: 8}}
	testConfigs := []struct {
		name string
		tm   *rtsspb.TileMap
		l    int32
		dim  *rtsspb.Coordinate
		want map[int32]*rtsspb.ClusterMap
	}{
		{
			name: "SingleTier",
			tm:   tm,
			l:    1,
			dim:  &rtsspb.Coordinate{X: 4, Y: 4},
			want: map[int32]*rtsspb.ClusterMap{
				1: {
					Level:     1,
					Dimension: &rtsspb.Coordinate{X: 2, Y: 2},
					Clusters: []*rtsspb.Cluster{
						{
							Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
							TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 4, Y: 4},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 0, Y: 1},
							TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 4},
							TileDimension: &rtsspb.Coordinate{X: 4, Y: 4},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 1, Y: 0},
							TileBoundary:  &rtsspb.Coordinate{X: 4, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 4, Y: 4},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 1, Y: 1},
							TileBoundary:  &rtsspb.Coordinate{X: 4, Y: 4},
							TileDimension: &rtsspb.Coordinate{X: 4, Y: 4},
						},
					},
				},
			},
		},
		{
			name: "ImperfectBorderScaling",
			tm:   tm,
			l:    1,
			dim:  &rtsspb.Coordinate{X: 5, Y: 5},
			want: map[int32]*rtsspb.ClusterMap{
				1: {
					Level:     1,
					Dimension: &rtsspb.Coordinate{X: 2, Y: 2},
					Clusters: []*rtsspb.Cluster{
						{
							Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
							TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 5, Y: 5},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 0, Y: 1},
							TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 5},
							TileDimension: &rtsspb.Coordinate{X: 5, Y: 3},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 1, Y: 0},
							TileBoundary:  &rtsspb.Coordinate{X: 5, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 3, Y: 5},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 1, Y: 1},
							TileBoundary:  &rtsspb.Coordinate{X: 5, Y: 5},
							TileDimension: &rtsspb.Coordinate{X: 3, Y: 3},
						},
					},
				},
			},
		},
		{
			name: "DoubleTier",
			tm:   tm,
			l:    2,
			dim:  &rtsspb.Coordinate{X: 2, Y: 2},
			want: map[int32]*rtsspb.ClusterMap{
				1: {
					Level:     1,
					Dimension: &rtsspb.Coordinate{X: 4, Y: 4},
					Clusters: []*rtsspb.Cluster{
						{
							Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
							TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 0, Y: 1},
							TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 2},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 0, Y: 2},
							TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 4},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 0, Y: 3},
							TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 6},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 1, Y: 0},
							TileBoundary:  &rtsspb.Coordinate{X: 2, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 1, Y: 1},
							TileBoundary:  &rtsspb.Coordinate{X: 2, Y: 2},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 1, Y: 2},
							TileBoundary:  &rtsspb.Coordinate{X: 2, Y: 4},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 1, Y: 3},
							TileBoundary:  &rtsspb.Coordinate{X: 2, Y: 6},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 2, Y: 0},
							TileBoundary:  &rtsspb.Coordinate{X: 4, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 2, Y: 1},
							TileBoundary:  &rtsspb.Coordinate{X: 4, Y: 2},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 2, Y: 2},
							TileBoundary:  &rtsspb.Coordinate{X: 4, Y: 4},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 2, Y: 3},
							TileBoundary:  &rtsspb.Coordinate{X: 4, Y: 6},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 3, Y: 0},
							TileBoundary:  &rtsspb.Coordinate{X: 6, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 3, Y: 1},
							TileBoundary:  &rtsspb.Coordinate{X: 6, Y: 2},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 3, Y: 2},
							TileBoundary:  &rtsspb.Coordinate{X: 6, Y: 4},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 3, Y: 3},
							TileBoundary:  &rtsspb.Coordinate{X: 6, Y: 6},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
					},
				},
				2: {
					Level:     2,
					Dimension: &rtsspb.Coordinate{X: 2, Y: 2},
					Clusters: []*rtsspb.Cluster{
						{
							Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
							TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 4, Y: 4},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 0, Y: 1},
							TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 4},
							TileDimension: &rtsspb.Coordinate{X: 4, Y: 4},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 1, Y: 0},
							TileBoundary:  &rtsspb.Coordinate{X: 4, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 4, Y: 4},
						},
						{
							Coordinate:    &rtsspb.Coordinate{X: 1, Y: 1},
							TileBoundary:  &rtsspb.Coordinate{X: 4, Y: 4},
							TileDimension: &rtsspb.Coordinate{X: 4, Y: 4},
						},
					},
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tm, err := tile.ImportTileMap(c.tm)
			if err != nil {
				t.Fatalf("ImportTileMap() = _, %v, want = _, nil", err)
			}

			want := map[int32]*cluster.ClusterMap{}
			for l, clPB := range c.want {
				cl, err := cluster.ImportClusterMap(clPB)
				if err != nil {
					t.Fatalf("ImportClusterMap() = _, %v, want = _, nil", err)
				}
				want[l] = cl
			}

			got, err := buildTieredClusterMaps(tm, c.l, c.dim)
			if err != nil {
				t.Fatalf("buildTieredClusterMaps() = _, %v, want = _, nil", err)
			}

			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("buildTieredClusterMaps() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
*/
func TestBuildTransitions(t *testing.T) {
	testConfigs := []struct {
		name string
		cm   *rtsspb.ClusterMap
		tm   *rtsspb.TileMap
		want []*rtsspb.Transition
	}{
		{
			name: "TrivialOpenMap",
			cm: &rtsspb.ClusterMap{
				Level:            1,
				TileDimension:    &rtsspb.Coordinate{X: 1, Y: 3},
				TileMapDimension: &rtsspb.Coordinate{X: 2, Y: 6},
			},
			tm: &rtsspb.TileMap{
				Dimension: &rtsspb.Coordinate{X: 2, Y: 6},
				Tiles: []*rtsspb.Tile{
					{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
					{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}},
					{Coordinate: &rtsspb.Coordinate{X: 0, Y: 2}},
					{Coordinate: &rtsspb.Coordinate{X: 0, Y: 3}},
					{Coordinate: &rtsspb.Coordinate{X: 0, Y: 4}},
					{Coordinate: &rtsspb.Coordinate{X: 0, Y: 5}},
					{Coordinate: &rtsspb.Coordinate{X: 1, Y: 0}},
					{Coordinate: &rtsspb.Coordinate{X: 1, Y: 1}},
					{Coordinate: &rtsspb.Coordinate{X: 1, Y: 2}},
					{Coordinate: &rtsspb.Coordinate{X: 1, Y: 3}},
					{Coordinate: &rtsspb.Coordinate{X: 1, Y: 4}},
					{Coordinate: &rtsspb.Coordinate{X: 1, Y: 5}},
				},
			},
			want: []*rtsspb.Transition{
				{
					N1: &rtsspb.AbstractNode{
						Level:          1,
						TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 1},
					},
					N2: &rtsspb.AbstractNode{
						Level:          1,
						TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 1},
					},
				},
				{
					N1: &rtsspb.AbstractNode{
						Level:          1,
						TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 2},
					},
					N2: &rtsspb.AbstractNode{
						Level:          1,
						TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 3},
					},
				},
				{
					N1: &rtsspb.AbstractNode{
						Level:          1,
						TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 2},
					},
					N2: &rtsspb.AbstractNode{
						Level:          1,
						TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 3},
					},
				},
				{
					N1: &rtsspb.AbstractNode{
						Level:          1,
						TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 4},
					},
					N2: &rtsspb.AbstractNode{
						Level:          1,
						TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 4},
					},
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			cm, err := cluster.ImportClusterMap(c.cm)
			if err != nil {
				t.Fatalf("ImportClusterMap() = _, %v, want = _, nil", err)
			}

			tm, err := tile.ImportTileMap(c.tm)
			if err != nil {
				t.Fatalf("ImportTileMap() = _, %v, want = _, nil", err)
			}

			got, err := buildTransitions(tm, cm)
			if diff := cmp.Diff(c.want, got, protocmp.Transform(), cmpopts.SortSlices(transitionLess)); diff != "" {
				t.Errorf("buildTranactions() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

/*
func TestBuildBaseIntraEdges(t *testing.T) {
	testConfigs := []struct {
		name string
		cm   *rtsspb.ClusterMap
		tm   *rtsspb.TileMap
		nm   Map
		want []*rtsspb.AbstractEdge
	}{
		{name: "NilCase", cm: nil, nm: nil, tm: nil, want: nil},
		{
			name: "SingleAbstractNode",
			cm: &rtsspb.ClusterMap{
				Level:     1,
				Dimension: &rtsspb.Coordinate{X: 1, Y: 1},
				Clusters: []*rtsspb.Cluster{
					{
						Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
						TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
						TileDimension: &rtsspb.Coordinate{X: 3, Y: 3},
					},
				},
			},
			nm: map[utils.MapCoordinate]*rtsspb.AbstractNode{
				{X: 0, Y: 0}: {
					Level:             1,
					ClusterCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
					TileCoordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
				},
			},
			tm:   simpleMapProto,
			want: nil,
		},
		{
			name: "MultiAbstractNode",
			cm: &rtsspb.ClusterMap{
				Level:     1,
				Dimension: &rtsspb.Coordinate{X: 1, Y: 1},
				Clusters: []*rtsspb.Cluster{
					{
						Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
						TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
						TileDimension: &rtsspb.Coordinate{X: 3, Y: 3},
					},
				},
			},
			nm: map[utils.MapCoordinate]*rtsspb.AbstractNode{
				{X: 0, Y: 0}: {
					Level:             1,
					ClusterCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
					TileCoordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
				},

				{X: 0, Y: 1}: {
					Level:             1,
					ClusterCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
					TileCoordinate:    &rtsspb.Coordinate{X: 0, Y: 1},
				},
			},
			tm: simpleMapProto,
			want: []*rtsspb.AbstractEdge{
				{
					Level:       1,
					Source:      &rtsspb.Coordinate{X: 0, Y: 0},
					Destination: &rtsspb.Coordinate{X: 0, Y: 1},
					EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
				},
				{
					Level:       1,
					Source:      &rtsspb.Coordinate{X: 0, Y: 1},
					Destination: &rtsspb.Coordinate{X: 0, Y: 0},
					EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			cm, err := cluster.ImportClusterMap(c.cm)
			if err != nil {
				t.Fatalf("ImportClusterMap() = _, %v, want = _, nil", err)
			}

			tm, err := tile.ImportTileMap(c.tm)
			if err != nil {
				t.Fatalf("ImportTileMap() = _, %v, want = _, nil", err)
			}

			got, err := buildBaseIntraEdges(cm, tm, c.nm)
			if err != nil {
				t.Fatalf("buildBaseIntraEdges() = _, %v, want = _, nil", err)
			}

			if diff := cmp.Diff(c.want, got, protocmp.Transform(), cmpopts.SortSlices(edgeLess)); diff != "" {
				t.Errorf("buildBaseIntraEdges() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestBuildAbstractGraphError(t *testing.T) {
	testConfigs := []struct {
		name             string
		tm               *rtsspb.TileMap
		level            int32
		clusterDimension *rtsspb.Coordinate
	}{
		{
			name:             "UnimplementedHigherLevelError",
			tm:               simpleMapProto,
			level:            2,
			clusterDimension: &rtsspb.Coordinate{X: 2, Y: 2},
		},
		{
			name:             "ClusterDimensionTooLargeError",
			tm:               simpleMapProto,
			level:            1,
			clusterDimension: &rtsspb.Coordinate{X: 100, Y: 100},
		},
		{
			name:             "ClusterDimensionTooSmall",
			tm:               simpleMapProto,
			level:            1,
			clusterDimension: &rtsspb.Coordinate{X: 1, Y: 1},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tm, err := tile.ImportTileMap(c.tm)
			if err != nil {
				t.Fatalf("ImportTileMap() = _, %v, want = _, nil", err)
			}

			if _, err := BuildAbstractGraph(tm, c.level, c.clusterDimension); err == nil {
				t.Error("BuildAbstractGraph() = _, nil, want a non-nil error")
			}
		})
	}
}

func TestBuildAbstractGraph(t *testing.T) {
	simpleMapClusterMapProto := &rtsspb.ClusterMap{
		Level:     1,
		Dimension: &rtsspb.Coordinate{X: 2, Y: 2},
		Clusters: []*rtsspb.Cluster{
			{
				Coordinate:    &rtsspb.Coordinate{X: 0, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
			},
			{
				Coordinate:    &rtsspb.Coordinate{X: 0, Y: 1},
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 2},
				TileDimension: &rtsspb.Coordinate{X: 2, Y: 1},
			},
			{
				Coordinate:    &rtsspb.Coordinate{X: 1, Y: 0},
				TileBoundary:  &rtsspb.Coordinate{X: 2, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 2},
			},
			{
				Coordinate:    &rtsspb.Coordinate{X: 1, Y: 1},
				TileBoundary:  &rtsspb.Coordinate{X: 2, Y: 2},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			},
		},
	}
	simpleMapClusterMap, err := cluster.ImportClusterMap(simpleMapClusterMapProto)
	if err != nil {
		t.Fatalf("ImportClusterMap() = _, %v, want = _, nil", err)
	}

	testConfigs := []struct {
		name             string
		tm               *rtsspb.TileMap
		level            int32
		clusterDimension *rtsspb.Coordinate
		want             *AbstractGraph
	}{
		{
			name:             "SimpleMap",
			tm:               simpleMapProto,
			level:            1,
			clusterDimension: simpleMapClusterMap.D,
			want: &AbstractGraph{
				Level: 1,
				ClusterMap: map[int32]*cluster.ClusterMap{
					1: simpleMapClusterMap,
				},
				NodeMap: map[int32]Map{
					1: {
						utils.MapCoordinate{X: 1, Y: 1}: &rtsspb.AbstractNode{
							Level:             1,
							ClusterCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
							TileCoordinate:    &rtsspb.Coordinate{X: 1, Y: 1},
						},
						utils.MapCoordinate{X: 1, Y: 2}: &rtsspb.AbstractNode{
							Level:             1,
							ClusterCoordinate: &rtsspb.Coordinate{X: 0, Y: 1},
							TileCoordinate:    &rtsspb.Coordinate{X: 1, Y: 2},
						},
						utils.MapCoordinate{X: 2, Y: 1}: &rtsspb.AbstractNode{
							Level:             1,
							ClusterCoordinate: &rtsspb.Coordinate{X: 1, Y: 0},
							TileCoordinate:    &rtsspb.Coordinate{X: 2, Y: 1},
						},
						utils.MapCoordinate{X: 2, Y: 2}: &rtsspb.AbstractNode{
							Level:             1,
							ClusterCoordinate: &rtsspb.Coordinate{X: 1, Y: 1},
							TileCoordinate:    &rtsspb.Coordinate{X: 2, Y: 2},
						},
					},
				},
				EdgeMap: map[int32]Map{
					1: {
						utils.MapCoordinate{X: 1, Y: 1}: map[utils.MapCoordinate]*rtsspb.AbstractEdge{
							{X: 1, Y: 2}: {
								Level:       1,
								Source:      &rtsspb.Coordinate{X: 1, Y: 1},
								Destination: &rtsspb.Coordinate{X: 1, Y: 2},
								EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
								Weight:      1,
							},
							{X: 2, Y: 1}: {
								Level:       1,
								Source:      &rtsspb.Coordinate{X: 1, Y: 1},
								Destination: &rtsspb.Coordinate{X: 2, Y: 1},
								EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
								Weight:      1,
							},
						},
						utils.MapCoordinate{X: 1, Y: 2}: map[utils.MapCoordinate]*rtsspb.AbstractEdge{
							{X: 2, Y: 2}: {
								Level:       1,
								Source:      &rtsspb.Coordinate{X: 1, Y: 2},
								Destination: &rtsspb.Coordinate{X: 2, Y: 2},
								EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
								Weight:      1,
							},
							{X: 1, Y: 1}: {
								Level:       1,
								Source:      &rtsspb.Coordinate{X: 1, Y: 2},
								Destination: &rtsspb.Coordinate{X: 1, Y: 1},
								EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
								Weight:      1,
							},
						},
						utils.MapCoordinate{X: 2, Y: 1}: map[utils.MapCoordinate]*rtsspb.AbstractEdge{
							{X: 2, Y: 2}: {
								Level:       1,
								Source:      &rtsspb.Coordinate{X: 2, Y: 1},
								Destination: &rtsspb.Coordinate{X: 2, Y: 2},
								EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
								Weight:      1,
							},
							{X: 1, Y: 1}: {
								Level:       1,
								Source:      &rtsspb.Coordinate{X: 2, Y: 1},
								Destination: &rtsspb.Coordinate{X: 1, Y: 1},
								EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
								Weight:      1,
							},
						},
						utils.MapCoordinate{X: 2, Y: 2}: map[utils.MapCoordinate]*rtsspb.AbstractEdge{
							{X: 2, Y: 1}: {
								Level:       1,
								Source:      &rtsspb.Coordinate{X: 2, Y: 2},
								Destination: &rtsspb.Coordinate{X: 2, Y: 1},
								EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
								Weight:      1,
							},
							{X: 1, Y: 2}: {
								Level:       1,
								Source:      &rtsspb.Coordinate{X: 2, Y: 2},
								Destination: &rtsspb.Coordinate{X: 1, Y: 2},
								EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
								Weight:      1,
							},
						},
					},
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tm, err := tile.ImportTileMap(c.tm)
			if err != nil {
				t.Fatalf("ImportTileMap() = _, %v, want = _, nil", err)
			}

			got, err := BuildAbstractGraph(tm, c.level, c.clusterDimension)
			if err != nil {
				t.Fatalf("BuildAbstractGraph() = _, %v, want = _, nil", err)
			}

			if diff := cmp.Diff(
				c.want,
				got,
				cmp.Comparer(abstractEdgeMapEqual),
				protocmp.Transform(),
			); diff != "" {
				t.Errorf("BuildAbstractGraph() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAbstractEdgeGraphGetBySource(t *testing.T) {
	source := &rtsspb.Coordinate{X: 0, Y: 0}
	want := []*rtsspb.AbstractEdge{
		{
			Source:      source,
			Destination: &rtsspb.Coordinate{X: 1, Y: 1},
		},
		{
			Source:      source,
			Destination: &rtsspb.Coordinate{X: 2, Y: 2},
		},
	}

	em := Map{}
	for _, e := range want {
		em.Add(e)
	}
	em.Add(&rtsspb.AbstractEdge{
		Source:      &rtsspb.Coordinate{X: 1, Y: 1},
		Destination: &rtsspb.Coordinate{X: 2, Y: 2},
	})

	got, err := em.GetBySource(utils.MC(source))
	if err != nil {
		t.Fatalf("GetBySource() = _, %v, want = _, nil", err)
	}

	if diff := cmp.Diff(want, got, cmp.Comparer(abstractEdgeEqual), cmpopts.SortSlices(edgeLess)); diff != "" {
		t.Errorf("GetBySource() mismatch (-want +got):\n%s", diff)
	}
}

func TestAbstractGraphGetNeighbors(t *testing.T) {
	const level = 1
	clusterDimension := &rtsspb.Coordinate{X: 3, Y: 3}
	nodeCoordinate := &rtsspb.Coordinate{X: 2, Y: 1}
	want := []*rtsspb.AbstractNode{
		{
			Level:             1,
			ClusterCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
			TileCoordinate:    &rtsspb.Coordinate{X: 1, Y: 2},
		},
		{
			Level:             1,
			ClusterCoordinate: &rtsspb.Coordinate{X: 1, Y: 0},
			TileCoordinate:    &rtsspb.Coordinate{X: 3, Y: 1},
		},
	}

	tm, err := tile.ImportTileMap(largeMapProto)
	if err != nil {
		t.Fatalf("ImportTileMap() = _, %v, want = _, nil", err)
	}

	g, err := BuildAbstractGraph(tm, level, clusterDimension)
	if err != nil {
		t.Fatalf("BuildAbstractGraph() = _, %v, want = _, nil", err)
	}

	n, err := g.NodeMap[level].Get(utils.MC(nodeCoordinate))
	if err != nil {
		t.Fatalf("Get() = _, %v, want = _, nil", err)
	}
	if n == nil {
		t.Fatal("Get() = nil, want a non-nil result")
	}

	got, err := g.Neighbors(n)
	if err != nil {
		t.Fatalf("Neighbors() = _, %v, want = _, nil", err)
	}

	if diff := cmp.Diff(
		want,
		got,
		protocmp.Transform(),
		cmpopts.SortSlices(nodeLess)); diff != "" {
		t.Errorf("Neighbors() mismatch (-want +got):\n%s", diff)
	}
}
*/
