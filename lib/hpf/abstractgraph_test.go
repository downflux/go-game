package abstractgraph

import (
	"testing"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/abstractedgemap"
	"github.com/cripplet/rts-pathing/lib/hpf/abstractnodemap"
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
		TerrainCosts: []*rtsspb.TerrainCost{
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
		},
	}

	/**
	 *	 -
	 *       X
	 * Y = 0 -
	 *   X = 0
	 */
	closedMapProto = &rtsspb.TileMap{
		Dimension: &rtsspb.Coordinate{X: 1, Y: 3},
		Tiles: []*rtsspb.Tile{
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 0}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 1}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &rtsspb.Coordinate{X: 0, Y: 2}, TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
		TerrainCosts: []*rtsspb.TerrainCost{
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
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
			cm, err := cluster.ImportMap(c.cm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			tm, err := tile.ImportMap(c.tm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			got, err := buildTransitions(tm, cm)
			if diff := cmp.Diff(c.want, got, protocmp.Transform(), cmpopts.SortSlices(transitionLess)); diff != "" {
				t.Errorf("buildTranactions() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestBuildIntraEdgeError(t *testing.T) {
	testConfigs := []struct {
		name   string
		tm     *rtsspb.TileMap
		cm     *rtsspb.ClusterMap
		n1, n2 *rtsspb.AbstractNode
	}{
		{name: "NilCaseError", tm: nil, cm: nil, n1: nil, n2: nil},
		{
			name: "NonAdjacentClusters",
			tm:   simpleMapProto,
			cm: &rtsspb.ClusterMap{
				Level:            1,
				TileDimension:    &rtsspb.Coordinate{X: 1, Y: 2},
				TileMapDimension: simpleMapProto.GetDimension(),
			},
			n1: &rtsspb.AbstractNode{
				Level:          1,
				TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
			},
			n2: &rtsspb.AbstractNode{
				Level:          1,
				TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 3},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			cm, err := cluster.ImportMap(c.cm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			tm, err := tile.ImportMap(c.tm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			if got, err := buildIntraEdge(tm, cm, c.n1, c.n2); err == nil {
				t.Fatalf("buildBaseIntraEdges() = %v, nil, want a non-nil error", got)
			}
		})
	}
}

func TestBuildIntraEdge(t *testing.T) {
	testConfigs := []struct {
		name   string
		tm     *rtsspb.TileMap
		cm     *rtsspb.ClusterMap
		n1, n2 *rtsspb.AbstractNode
		want   *rtsspb.AbstractEdge
	}{
		{
			name: "ZeroLengthEdge",
			tm:   simpleMapProto,
			cm: &rtsspb.ClusterMap{
				Level:            1,
				TileDimension:    &rtsspb.Coordinate{X: 2, Y: 2},
				TileMapDimension: simpleMapProto.GetDimension(),
			},
			n1: &rtsspb.AbstractNode{
				Level:          1,
				TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
			},
			n2: &rtsspb.AbstractNode{
				Level:          1,
				TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
			},
			want: &rtsspb.AbstractEdge{
				Level:       1,
				Source:      &rtsspb.Coordinate{X: 0, Y: 0},
				Destination: &rtsspb.Coordinate{X: 0, Y: 0},
				Weight:      0,
				EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
			},
		},
		{
			name: "AdjacentNode",
			tm:   simpleMapProto,
			cm: &rtsspb.ClusterMap{
				Level:            1,
				TileDimension:    &rtsspb.Coordinate{X: 2, Y: 2},
				TileMapDimension: simpleMapProto.GetDimension(),
			},
			n1: &rtsspb.AbstractNode{
				Level:          1,
				TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
			},
			n2: &rtsspb.AbstractNode{
				Level:          1,
				TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 1},
			},
			want: &rtsspb.AbstractEdge{
				Level:       1,
				Source:      &rtsspb.Coordinate{X: 0, Y: 0},
				Destination: &rtsspb.Coordinate{X: 0, Y: 1},
				Weight:      1,
				EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
			},
		},
		{
			name: "BlockedNodes",
			tm:   closedMapProto,
			cm: &rtsspb.ClusterMap{
				Level:            1,
				TileDimension:    &rtsspb.Coordinate{X: 1, Y: 3},
				TileMapDimension: closedMapProto.GetDimension(),
			},
			n1: &rtsspb.AbstractNode{
				Level:          1,
				TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
			},
			n2: &rtsspb.AbstractNode{
				Level:          1,
				TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 2},
			},
			want: nil,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			cm, err := cluster.ImportMap(c.cm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			tm, err := tile.ImportMap(c.tm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			got, err := buildIntraEdge(tm, cm, c.n1, c.n2)
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
			tm, err := tile.ImportMap(c.tm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			if _, err := BuildAbstractGraph(tm, c.clusterDimension, c.level); err == nil {
				t.Error("BuildAbstractGraph() = _, nil, want a non-nil error")
			}
		})
	}
}

func newAbstractNodeMap(cm *cluster.Map, nodes []*rtsspb.AbstractNode) *abstractnodemap.Map {
	nm := &abstractnodemap.Map{
		ClusterMap: cm,
	}
	for _, n := range nodes {
		nm.Add(n)
	}

	return nm
}

func newAbstractEdgeMap(edges []*rtsspb.AbstractEdge) *abstractedgemap.Map {
	em := &abstractedgemap.Map{}
	for _, e := range edges {
		em.Add(e)
	}

	return em
}

func TestBuildAbstractGraph(t *testing.T) {
	simpleMapClusterMapProto := &rtsspb.ClusterMap{
		Level:            1,
		TileDimension:    &rtsspb.Coordinate{X: 2, Y: 2},
		TileMapDimension: simpleMapProto.GetDimension(),
	}
	simpleMapClusterMap, err := cluster.ImportMap(simpleMapClusterMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
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
			clusterDimension: simpleMapClusterMap.Val.GetTileDimension(),
			want: &AbstractGraph{
				Level: 1,
				NodeMap: []*abstractnodemap.Map{
					newAbstractNodeMap(simpleMapClusterMap, []*rtsspb.AbstractNode{
						{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 1}},
						{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 2}},
						{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 2, Y: 1}},
						{Level: 1, TileCoordinate: &rtsspb.Coordinate{X: 2, Y: 2}},
					}),
				},
				EdgeMap: []*abstractedgemap.Map{
					newAbstractEdgeMap([]*rtsspb.AbstractEdge{
						{
							Level:       1,
							Source:      &rtsspb.Coordinate{X: 1, Y: 1},
							Destination: &rtsspb.Coordinate{X: 1, Y: 2},
							EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
							Weight:      1,
						}, {
							Level:       1,
							Source:      &rtsspb.Coordinate{X: 1, Y: 1},
							Destination: &rtsspb.Coordinate{X: 2, Y: 1},
							EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
							Weight:      1,
						}, {
							Level:       1,
							Source:      &rtsspb.Coordinate{X: 1, Y: 2},
							Destination: &rtsspb.Coordinate{X: 2, Y: 2},
							EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
							Weight:      1,
						}, {
							Level:       1,
							Source:      &rtsspb.Coordinate{X: 1, Y: 2},
							Destination: &rtsspb.Coordinate{X: 1, Y: 1},
							EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
							Weight:      1,
						}, {
							Level:       1,
							Source:      &rtsspb.Coordinate{X: 2, Y: 1},
							Destination: &rtsspb.Coordinate{X: 2, Y: 2},
							EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
							Weight:      1,
						}, {
							Level:       1,
							Source:      &rtsspb.Coordinate{X: 2, Y: 1},
							Destination: &rtsspb.Coordinate{X: 1, Y: 1},
							EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
							Weight:      1,
						}, {
							Level:       1,
							Source:      &rtsspb.Coordinate{X: 2, Y: 2},
							Destination: &rtsspb.Coordinate{X: 2, Y: 1},
							EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
							Weight:      1,
						}, {
							Level:       1,
							Source:      &rtsspb.Coordinate{X: 2, Y: 2},
							Destination: &rtsspb.Coordinate{X: 1, Y: 2},
							EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
							Weight:      1,
						},
					}),
				},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tm, err := tile.ImportMap(c.tm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			got, err := BuildAbstractGraph(tm, c.clusterDimension, c.level)
			if err != nil {
				t.Fatalf("BuildAbstractGraph() = _, %v, want = _, nil", err)
			}

			if diff := cmp.Diff(
				c.want,
				got,
				cmp.Comparer(abstractEdgeMapEqual),
				cmp.AllowUnexported(abstractedgemap.Map{}, abstractnodemap.Map{}),
				protocmp.Transform(),
			); diff != "" {
				t.Errorf("BuildAbstractGraph() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAbstractGraphGetNeighbors(t *testing.T) {
	const level = 1
	clusterDimension := &rtsspb.Coordinate{X: 3, Y: 3}
	nodeCoordinate := &rtsspb.Coordinate{X: 2, Y: 1}
	want := []*rtsspb.AbstractNode{
		{
			Level:          level,
			TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 2},
		},
		{
			Level:          level,
			TileCoordinate: &rtsspb.Coordinate{X: 3, Y: 1},
		},
	}

	tm, err := tile.ImportMap(largeMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}

	g, err := BuildAbstractGraph(tm, clusterDimension, level)
	if err != nil {
		t.Fatalf("BuildAbstractGraph() = _, %v, want = _, nil", err)
	}

	n, err := g.NodeMap[listIndex(level)].Get(utils.MC(nodeCoordinate))
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
