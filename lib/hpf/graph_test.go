package graph

import (
	"math"
	"testing"

	rtscpb "github.com/minkezhang/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/minkezhang/rts-pathing/lib/hpf/cluster"
	"github.com/minkezhang/rts-pathing/lib/hpf/edge"
	"github.com/minkezhang/rts-pathing/lib/hpf/node"
	"github.com/minkezhang/rts-pathing/lib/hpf/tile"
	"github.com/minkezhang/rts-pathing/lib/hpf/utils"
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
			{TerrainType: rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
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

func edgeMapEqual(em1, em2 edge.Map) bool {
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
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 1}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 1}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 2}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 3}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 2}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 3}},
				},
				{
					N1: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 4}},
					N2: &rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 4}},
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

func TestBuildGraphError(t *testing.T) {
	testConfigs := []struct {
		name             string
		tm               *rtsspb.TileMap
		clusterDimension *rtsspb.Coordinate
	}{}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tm, err := tile.ImportMap(c.tm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			if _, err := BuildGraph(tm, c.clusterDimension); err == nil {
				t.Error("BuildGraph() = _, nil, want a non-nil error")
			}
		})
	}
}

func newAbstractNode(cm *cluster.Map, nodes []*rtsspb.AbstractNode) *node.Map {
	nm := &node.Map{
		ClusterMap: cm,
	}
	for _, n := range nodes {
		nm.Add(n)
	}

	return nm
}

func newAbstractEdge(edges []*rtsspb.AbstractEdge) *edge.Map {
	em := &edge.Map{}
	for _, e := range edges {
		em.Add(e)
	}

	return em
}

func TestBuildGraph(t *testing.T) {
	simpleMapClusterMapProto := &rtsspb.ClusterMap{
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
		clusterDimension *rtsspb.Coordinate
		want             *Graph
	}{
		{
			name:             "SimpleMap",
			tm:               simpleMapProto,
			clusterDimension: simpleMapClusterMap.Val.GetTileDimension(),
			want: &Graph{
				NodeMap: newAbstractNode(simpleMapClusterMap, []*rtsspb.AbstractNode{
					{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 1}},
					{TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 2}},
					{TileCoordinate: &rtsspb.Coordinate{X: 2, Y: 1}},
					{TileCoordinate: &rtsspb.Coordinate{X: 2, Y: 2}},
				}),
				EdgeMap: newAbstractEdge([]*rtsspb.AbstractEdge{
					{
						Source:      &rtsspb.Coordinate{X: 1, Y: 1},
						Destination: &rtsspb.Coordinate{X: 1, Y: 2},
						EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
						Weight:      1,
					}, {
						Source:      &rtsspb.Coordinate{X: 1, Y: 1},
						Destination: &rtsspb.Coordinate{X: 2, Y: 1},
						EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
						Weight:      1,
					}, {
						Source:      &rtsspb.Coordinate{X: 1, Y: 2},
						Destination: &rtsspb.Coordinate{X: 2, Y: 2},
						EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
						Weight:      1,
					}, {
						Source:      &rtsspb.Coordinate{X: 1, Y: 2},
						Destination: &rtsspb.Coordinate{X: 1, Y: 1},
						EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
						Weight:      1,
					}, {
						Source:      &rtsspb.Coordinate{X: 2, Y: 1},
						Destination: &rtsspb.Coordinate{X: 2, Y: 2},
						EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
						Weight:      1,
					}, {
						Source:      &rtsspb.Coordinate{X: 2, Y: 1},
						Destination: &rtsspb.Coordinate{X: 1, Y: 1},
						EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
						Weight:      1,
					}, {
						Source:      &rtsspb.Coordinate{X: 2, Y: 2},
						Destination: &rtsspb.Coordinate{X: 2, Y: 1},
						EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
						Weight:      1,
					}, {
						Source:      &rtsspb.Coordinate{X: 2, Y: 2},
						Destination: &rtsspb.Coordinate{X: 1, Y: 2},
						EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
						Weight:      1,
					},
				}),
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tm, err := tile.ImportMap(c.tm)
			if err != nil {
				t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
			}

			got, err := BuildGraph(tm, c.clusterDimension)
			if err != nil {
				t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
			}

			if diff := cmp.Diff(
				c.want,
				got,
				cmp.Comparer(edgeMapEqual),
				cmp.AllowUnexported(edge.Map{}, node.Map{}),
				protocmp.Transform(),
			); diff != "" {
				t.Errorf("BuildGraph() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGraphGetNeighbors(t *testing.T) {
	clusterDimension := &rtsspb.Coordinate{X: 3, Y: 3}
	nodeCoordinate := &rtsspb.Coordinate{X: 2, Y: 1}
	want := []*rtsspb.AbstractNode{
		{
			TileCoordinate: &rtsspb.Coordinate{X: 1, Y: 2},
		},
		{
			TileCoordinate: &rtsspb.Coordinate{X: 3, Y: 1},
		},
	}

	tm, err := tile.ImportMap(largeMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}

	g, err := BuildGraph(tm, clusterDimension)
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}

	n, err := g.NodeMap.Get(utils.MC(nodeCoordinate))
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

func TestConnect(t *testing.T) {
	connectEdgeNodeCoordinate := &rtsspb.Coordinate{X: 0, Y: 0}
	connectEdgeNodeMap, err := tile.ImportMap(simpleMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}
	connectEdgeNodeGraph, err := BuildGraph(connectEdgeNodeMap, simpleMapProto.GetDimension())
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}
	for _, tc := range []*rtsspb.Coordinate{
		{X: 0, Y: 1},
		connectEdgeNodeCoordinate,
	} {
		if err := connectEdgeNodeGraph.NodeMap.Add(&rtsspb.AbstractNode{TileCoordinate: tc}); err != nil {
			t.Fatalf("Add() = %v, want = nil", err)
		}
	}

	connectEphemeralNodeCoordinate := &rtsspb.Coordinate{X: 0, Y: 0}
	connectEphemeralNodeMap, err := tile.ImportMap(simpleMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}
	connectEphemeralNodeGraph, err := BuildGraph(connectEphemeralNodeMap, simpleMapProto.GetDimension())
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}
	if err := connectEphemeralNodeGraph.NodeMap.Add(&rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 1}, IsEphemeral: true}); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}
	if err := connectEphemeralNodeGraph.NodeMap.Add(&rtsspb.AbstractNode{TileCoordinate: connectEphemeralNodeCoordinate}); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	noConnectEphemeralNodeCoordinate := &rtsspb.Coordinate{X: 0, Y: 0}
	noConnectEphemeralNodeMap, err := tile.ImportMap(simpleMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}
	noConnectEphemeralNodeGraph, err := BuildGraph(noConnectEphemeralNodeMap, simpleMapProto.GetDimension())
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}
	for _, tc := range []*rtsspb.Coordinate{
		{X: 0, Y: 1},
		noConnectEphemeralNodeCoordinate,
	} {
		if err := noConnectEphemeralNodeGraph.NodeMap.Add(&rtsspb.AbstractNode{TileCoordinate: tc, IsEphemeral: true}); err != nil {
			t.Fatalf("Add() = %v, want = nil", err)
		}
	}

	testConfigs := []struct {
		name string
		tm   *tile.Map
		g    *Graph
		t    *rtsspb.Coordinate
		want []*rtsspb.AbstractEdge
	}{
		{
			name: "ConnectEdgeNode",
			tm:   connectEdgeNodeMap,
			g:    connectEdgeNodeGraph,
			t:    connectEdgeNodeCoordinate,
			want: []*rtsspb.AbstractEdge{
				{
					Source:      connectEdgeNodeCoordinate,
					Destination: &rtsspb.Coordinate{X: 0, Y: 1},
					EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
					Weight:      1,
				},
			},
		},
		{
			name: "ConnectEphemeralNode",
			tm:   connectEphemeralNodeMap,
			g:    connectEphemeralNodeGraph,
			t:    connectEphemeralNodeCoordinate,
			want: []*rtsspb.AbstractEdge{
				{
					Source:      connectEphemeralNodeCoordinate,
					Destination: &rtsspb.Coordinate{X: 0, Y: 1},
					EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
					Weight:      1,
				},
			},
		},
		{
			name: "NoConnectEphemeralNode",
			tm:   noConnectEphemeralNodeMap,
			g:    noConnectEphemeralNodeGraph,
			t:    noConnectEphemeralNodeCoordinate,
			want: nil,
		},
	}
	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if err := connect(c.tm, c.g, utils.MC(c.t)); err != nil {
				t.Fatalf("connect() = %v, want = nil", err)
			}

			got, err := c.g.EdgeMap.GetBySource(utils.MC(c.t))
			if err != nil {
				t.Fatalf("GetBySource() = _, %v, want = _, nil", err)
			}

			if diff := cmp.Diff(
				c.want,
				got,
				cmpopts.SortSlices(edgeLess),
				cmp.Comparer(abstractEdgeEqual),
			); diff != "" {
				t.Errorf("GetBySource() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestAddNonEphemeralNodeNoOp(t *testing.T) {
	tileCoordinate := utils.MapCoordinate{X: 1, Y: 1}

	tm, err := tile.ImportMap(simpleMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}
	g, err := BuildGraph(tm, &rtsspb.Coordinate{X: 2, Y: 2})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}

	uuid, err := InsertEphemeralNode(tm, g, tileCoordinate)
	if err != nil || uuid != 0 {
		t.Fatalf("InsertEphemeralNode() = %v, %v, want = 0, nil", uuid, err)
	}

	n, err := g.NodeMap.Get(tileCoordinate)
	if err != nil || n == nil {
		t.Fatalf("Get() = %v, %v, want = _, nil", n, err)
	}

	if n.GetIsEphemeral() {
		t.Fatalf("GetIsEphemeral() = %v, want = false", n.GetIsEphemeral())
	}
}

func TestSimpleAddEphemeralNode(t *testing.T) {
	const nInserts = 1000
	tileCoordinate := utils.MapCoordinate{X: 0, Y: 0}

	tm, err := tile.ImportMap(simpleMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}
	g, err := BuildGraph(tm, &rtsspb.Coordinate{X: 2, Y: 2})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}

	uuids := map[int64]bool{}
	for len(uuids) < nInserts {
		u, err := InsertEphemeralNode(tm, g, tileCoordinate)
		if err != nil || u == 0 {
			t.Fatalf("InsertEphemeralNode() = %v, %v, want = _, nil", u, err)
		}
		if _, found := uuids[u]; found {
			t.Fatalf("uuids[u] = %v, want = false", found)
		}
		uuids[u] = true
	}

	n, err := g.NodeMap.Get(tileCoordinate)
	if err != nil || n == nil {
		t.Fatalf("Get() = %v, %v, want = _, nil", n, err)
	}

	if !n.GetIsEphemeral() {
		t.Fatalf("GetIsEphemeral() = %v, want = true", n.GetIsEphemeral())
	}

	for u := range uuids {
		if got, found := n.GetEphemeralKeys()[u]; !found || got != true {
			t.Errorf("GetEphemeralKeys()[u] = %v, %v, want = true, true", got, found)
		}
	}
}

func TestDeleteNonEphemeralNodeNoOp(t *testing.T) {
	tileCoordinate := utils.MapCoordinate{X: 1, Y: 1}

	tm, err := tile.ImportMap(simpleMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}
	g, err := BuildGraph(tm, &rtsspb.Coordinate{X: 2, Y: 2})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}

	if err := RemoveEphemeralNode(g, tileCoordinate, 0); err != nil {
		t.Errorf("RemoveEphemeralNode() = %v, want = nil", err)
	}

	if n, err := g.NodeMap.Get(tileCoordinate); err != nil || n == nil {
		t.Errorf("Get() = %v, %v, want = _, nil", n, err)
	}
}

func TestDeleteEphemeralNode(t *testing.T) {
	const nInserts = 1000
	tileCoordinate := utils.MapCoordinate{X: 0, Y: 0}

	tm, err := tile.ImportMap(simpleMapProto)
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}
	g, err := BuildGraph(tm, &rtsspb.Coordinate{X: 2, Y: 2})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}

	uuids := map[int64]bool{}
	for len(uuids) < nInserts {
		u, err := InsertEphemeralNode(tm, g, tileCoordinate)
		if err != nil || u == 0 {
			t.Fatalf("InsertEphemeralNode() = %v, %v, want = _, nil", u, err)
		}
		if _, found := uuids[u]; found {
			t.Fatalf("uuids[u] = %v, want = false", found)
		}
		uuids[u] = true
	}

	for u := range uuids {
		// The node returned by Get is a reference; we don't need to keep
		// querying for this.
		//
		// TODO(minkezhang): Decide if this is what we actually want.
		n, err := g.NodeMap.Get(tileCoordinate)
		if err != nil || n == nil {
			t.Fatalf("Get() = %v, %v, want = _, nil", n, err)
		}

		if err := RemoveEphemeralNode(g, tileCoordinate, u); err != nil {
			t.Errorf("RemoveEphemeralNode() = %v, want = nil", err)
		}

		if _, found := n.GetEphemeralKeys()[u]; found {
			t.Fatalf("GetEphemeralKeys()[u] = _, %v, want = _, false", found)
		}
	}

	if n, err := g.NodeMap.Get(tileCoordinate); err != nil || n != nil {
		t.Errorf("Get() = %v, %v, want = nil, nil", n, err)
	}
}
