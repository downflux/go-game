// Package abstractgraph constructs and manages the abstract node space corresponding to a TileMap object.
// This package will be used ast he underlying topology for hiearchical A* searching.
package abstractgraph

import (
	"math"
	"sync"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/astar"
	"github.com/cripplet/rts-pathing/lib/hpf/cluster"
	"github.com/cripplet/rts-pathing/lib/hpf/entrance"
	"github.com/cripplet/rts-pathing/lib/hpf/tile"
	"github.com/cripplet/rts-pathing/lib/hpf/utils"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
	infinity = math.Inf(1)
)

type NodeMap map[utils.MapCoordinate]map[utils.MapCoordinate]*AbstractNode
type EdgeMap map[utils.MapCoordinate]map[utils.MapCoordinate]*AbstractEdge

type AbstractGraph struct {
	L int32

	Mu sync.RWMutex
	C  map[int32]*cluster.ClusterMap

	// N is a [L][ClusterCoordinate][TileCoordinate] map from the Cluster coordinate to a list of nodes contained within that Cluster.
	N map[int32]NodeMap

	// E is a [L][SourceTile][DestinationTile] map of Tile coordinates to the list of edges that connect to the Tle.
	E map[int32]EdgeMap
}

type AbstractEdge struct {
	Val *rtsspb.AbstractEdge
}
type AbstractNode struct {
	Val *rtsspb.AbstractNode
}

func ImportAbstractEdge(pb *rtsspb.AbstractEdge) (*AbstractEdge, error) {
	return &AbstractEdge{Val: pb}, nil
}

func ImportAbstractNode(pb *rtsspb.AbstractNode) (*AbstractNode, error) {
	return &AbstractNode{Val: pb}, nil
}

func ImportAbstractGraph(pb *rtsspb.AbstractGraph) (*AbstractGraph, error) {
	g := &AbstractGraph{
		L: pb.GetLevel(),
	}
	for _, cm := range pb.GetClusterMaps() {
		clusterMap, err := cluster.ImportClusterMap(cm)
		if err != nil {
			return nil, err
		}
		g.C[clusterMap.L] = clusterMap
	}

	for _, n := range pb.GetNodes() {
		if err := g.AddNode(&AbstractNode{Val: n}); err != nil {
			return nil, err
		}
	}
	for _, e := range pb.GetEdges() {
		if err := g.AddEdge(&AbstractEdge{Val: e}); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (m NodeMap) AddNode(n *AbstractNode) error {
	c := utils.MC(n.Val.GetClusterCoordinate())
	t := utils.MC(n.Val.GetTileCoordinate())

	if _, found := m[c]; !found {
		m[c] = map[utils.MapCoordinate]*AbstractNode{}
	}
	m[c][t] = n
	return nil
}

func (m EdgeMap) AddEdge(e *AbstractEdge) error {
	s := utils.MC(e.Val.GetSource())
	d := utils.MC(e.Val.GetDestination())
	if _, found := m[s]; !found {
		m[s] = map[utils.MapCoordinate]*AbstractEdge{}
	}
	if _, found := m[d]; !found {
		m[d] = map[utils.MapCoordinate]*AbstractEdge{}
	}
	// Assuming symmetrical bidirectional graph.
	m[s][d] = e
	m[d][s] = e
	return nil
}

func (g *AbstractGraph) AddNode(n *AbstractNode) error {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	if g.N[n.Val.GetLevel()] == nil {
		g.N[n.Val.GetLevel()] = NodeMap{}
	}
	return g.N[n.Val.GetLevel()].AddNode(n)
}

func (g *AbstractGraph) AddEdge(e *AbstractEdge) error {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	if g.E[e.Val.GetLevel()] == nil {
		g.E[e.Val.GetLevel()] = EdgeMap{}
	}
	return g.E[e.Val.GetLevel()].AddEdge(e)
}

func buildTransitions(cm *cluster.ClusterMap, tm *tile.TileMap) ([]*rtsspb.Transition, error) {
	var transitions []*rtsspb.Transition
	for _, c1 := range cm.M {
		neighbors, err := cm.Neighbors(c1.Val.GetCoordinate())
		if err != nil {
			return nil, err
		}

		for _, c2 := range neighbors {
			ts, err := entrance.BuildTransitions(c1, c2, tm)
			if err != nil {
				return nil, err
			}
			transitions = append(transitions, ts...)
		}
	}
	return transitions, nil
}

func buildEdges(tm *tile.TileMap, nodes []*AbstractNode, boundary, dimension *rtsspb.Coordinate, level int32) ([]*AbstractEdge, error) {
	var edges []*AbstractEdge

	for _, n1 := range nodes {
		c1 := n1.Val.GetTileCoordinate()
		for _, n2 := range nodes {
			c2 := n2.Val.GetTileCoordinate()
			if !proto.Equal(c1, c2) {
				path, cost, err := astar.TileMapPath(
					tm,
					tm.TileFromCoordinate(c1),
					tm.TileFromCoordinate(c2),
					boundary,
					dimension,
				)
				if err != nil {
					return nil, err
				}
				if path != nil {
					edges = append(edges, &AbstractEdge{
						Val: &rtsspb.AbstractEdge{
							Level:       level,
							Source:      c1,
							Destination: c2,
							EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
							Weight:      cost,
						},
					})
				}
			}
		}
	}
	return edges, nil
}

func addGraphLevel(g *AbstractGraph, tm *tile.TileMap, l int32) error {
	return notImplemented
}

func buildBaseGraph(g *AbstractGraph, tm *tile.TileMap) error {
	transitions, err := buildTransitions(g.C[1], tm)
	if err != nil {
		return err
	}

	for _, t := range transitions {
		g.AddNode(&AbstractNode{
			Val: &rtsspb.AbstractNode{
				Level:             1,
				ClusterCoordinate: t.GetN1().GetClusterCoordinate(),
				TileCoordinate:    t.GetN1().GetTileCoordinate(),
			},
		})
		g.AddNode(&AbstractNode{
			&rtsspb.AbstractNode{
				Level:             1,
				ClusterCoordinate: t.GetN2().GetClusterCoordinate(),
				TileCoordinate:    t.GetN2().GetTileCoordinate(),
			},
		})
		g.AddEdge(&AbstractEdge{
			&rtsspb.AbstractEdge{
				Level:       1,
				Source:      t.GetN1().GetTileCoordinate(),
				Destination: t.GetN2().GetTileCoordinate(),
				EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
				Weight:      1, // Inter-edges are always of cost 1, per Botea.
			},
		})
	}

	// Add intra-edges for all tiles within the same Cluster.
	for clusterCoordinate, nodesMap := range g.N[1] {
		var nodes []*AbstractNode
		for _, n := range nodesMap {
			nodes = append(nodes, n)
		}
		nodesCluster := g.C[1].Cluster(clusterCoordinate.X, clusterCoordinate.Y)
		edges, err := buildEdges(
			tm,
			nodes,
			nodesCluster.Val.GetTileBoundary(),
			nodesCluster.Val.GetTileDimension(),
			1)
		if err != nil {
			return err
		}

		for _, e := range edges {
			g.AddEdge(e)
		}
	}
	return nil
}

func buildClusterMaps(tm *tile.TileMap, level int32, clusterDimension *rtsspb.Coordinate) (map[int32]*cluster.ClusterMap, error) {
	clusterMaps := map[int32]*cluster.ClusterMap{}
	for i := int32(1); i < level; i++ {
		cm, err := cluster.BuildClusterMap(tm.D, &rtsspb.Coordinate{X: clusterDimension.GetX(), Y: clusterDimension.GetY()}, i)
		if err != nil {
			return nil, err
		}
		clusterMaps[i] = cm

	}
	return clusterMaps, nil
}

func BuildAbstractGraph(tm *tile.TileMap, level int32, clusterDimension *rtsspb.Coordinate) (*AbstractGraph, error) {
	if level < 1 {
		return nil, status.Error(codes.FailedPrecondition, "level must be a positive non-zero integer")
	}
	// Highest level ClusterMap should still have more than one Cluster,
	// otherwise we'll be routing units to the edge first before going back inwards.
	if (int32(math.Pow(float64(level), float64(clusterDimension.GetX()))) >= tm.D.GetX()) || (int32(math.Pow(float64(level), float64(clusterDimension.GetY()))) >= tm.D.GetY()) {
		return nil, status.Error(codes.FailedPrecondition, "given clusterDimension and level will result in too large a cluster map")
	}

	clusterMaps, err := buildClusterMaps(tm, level, clusterDimension)
	if err != nil {
		return nil, err
	}

	g := &AbstractGraph{
		L: level,
		C: clusterMaps,
	}

	if err = buildBaseGraph(g, tm); err != nil {
		return nil, err
	}

	// Add intra-edges for l > 1.
	for i := int32(2); i <= level; i++ {
		if err := addGraphLevel(g, tm, i); err != nil {
			return nil, err
		}
	}

	return g, nil
}
