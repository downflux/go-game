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
	infinity = math.Inf(1)
)

type AbstractGraph struct {
	L int32

	Mu sync.RWMutex
	C map[int32]*cluster.ClusterMap
	// N is a [L][ClusterCoordinate][TileCoordinate] map from the Cluster coordinate to a list of nodes contained within that Cluster.
	N map[int32]map[utils.MapCoordinate]map[utils.MapCoordinate]*AbstractNode
	// E is a [L][SourceTile][DestinationTile] map of Tile coordinates to the list of edges that connect to the Tle.
	E map[int32]map[utils.MapCoordinate]map[utils.MapCoordinate]*AbstractEdge
}

type AbstractEdge struct {
	Val *rtsspb.AbstractEdge
}
type AbstractNode struct {
	Val *rtsspb.AbstractNode
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
		if err := g.AddNode(n); err != nil {
			return nil, err
		}
	}
	for _, e := range pb.GetEdges() {
		if err := g.AddEdge(e); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func ImportAbstractEdge(pb *rtsspb.AbstractEdge) (*AbstractEdge, error) {
	return &AbstractEdge{Val: pb}, nil
}

func ImportAbstractNode(pb *rtsspb.AbstractNode) (*AbstractNode, error) {
	return &AbstractNode{Val: pb}, nil
}

func (g *AbstractGraph) AddNode(pb *rtsspb.AbstractNode) error {
	n, err := ImportAbstractNode(pb)
	if err != nil {
		return err
	}

	g.Mu.Lock()
	defer g.Mu.Unlock()

	if g.N[n.Val.GetLevel()] == nil {
		g.N[n.Val.GetLevel()] = map[utils.MapCoordinate]map[utils.MapCoordinate]*AbstractNode{}
	}
	g.N[n.Val.GetLevel()][utils.MC(n.Val.GetClusterCoordinate())][utils.MC(n.Val.GetTileCoordinate())] = n
	return nil
}

func (g *AbstractGraph) AddEdge(pb *rtsspb.AbstractEdge) error {
	e, err := ImportAbstractEdge(pb)
	if err != nil {
		return err
	}

	g.Mu.Lock()
	defer g.Mu.Unlock()

	if g.E[e.Val.GetLevel()] == nil {
		g.E[e.Val.GetLevel()] = map[utils.MapCoordinate]map[utils.MapCoordinate]*AbstractEdge{}
	}
	// Assuming symmetrical bidirectional graph.
	g.E[e.Val.GetLevel()][utils.MC(e.Val.GetSource())][utils.MC(e.Val.GetDestination())] = e
	g.E[e.Val.GetLevel()][utils.MC(e.Val.GetDestination())][utils.MC(e.Val.GetSource())] = e
	return nil
}

func BuildAbstractGraph(tm *tile.TileMap, level int32, clusterDimension *rtsspb.Coordinate) (*AbstractGraph, error) {
	if level < 1 {
		return nil, status.Error(codes.FailedPrecondition, "level must be a positive non-zero integer")
	}
	// Highest level ClusterMap should still have more than one Cluster,
	// otherwise we'll be routing units to the edge first before going back inwards.
	if (
		int32(math.Pow(float64(level), float64(clusterDimension.GetX()))) >= tm.D.GetX()) || (
		int32(math.Pow(float64(level), float64(clusterDimension.GetY()))) >= tm.D.GetY()) {
		return nil, status.Error(codes.FailedPrecondition, "given clusterDimension and level will result in too large a cluster map")
	}

	cm, err := cluster.BuildClusterMap(tm.D, &rtsspb.Coordinate{X: clusterDimension.GetX(), Y: clusterDimension.GetY()}, 1)
	if err != nil {
		return nil, err
	}

	g := &AbstractGraph{
		L: level,
		C: map[int32]*cluster.ClusterMap{
			cm.L: cm,
		},
	}

	var transitions []*rtsspb.Transition
	for _, cluster1 := range g.C[1].M {
		neighbors, err := g.C[1].Neighbors(cluster1.Val.GetCoordinate())
		if err != nil {
			return nil, err
		}

		for _, cluster2 := range neighbors {
			ts, err := entrance.BuildTransitions(cluster1, cluster2, tm)
			if err != nil {
				return nil, err
			}
			transitions = append(transitions, ts...)
		}
	}

	for _, t := range transitions {
		g.AddNode(&rtsspb.AbstractNode{
			Level:             1,
			ClusterCoordinate: t.GetN1().GetClusterCoordinate(),
			TileCoordinate:    t.GetN1().GetTileCoordinate(),
		})
		g.AddNode(&rtsspb.AbstractNode{
			Level:             1,
			ClusterCoordinate: t.GetN2().GetClusterCoordinate(),
			TileCoordinate:    t.GetN2().GetTileCoordinate(),
		})
		g.AddEdge(&rtsspb.AbstractEdge{
			Level:       1,
			Source:      t.GetN1().GetTileCoordinate(),
			Destination: t.GetN2().GetTileCoordinate(),
			EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
			Weight:      1, // Inter-edges are always of cost 1, per Botea.
		})
	}

	// Add intra-edges for all tiles within the same Cluster.
	for clusterCoordinate, nodes := range g.N[1] {
		nodesCluster := g.C[1].Cluster(clusterCoordinate.X, clusterCoordinate.Y)
		for _, n1 := range nodes {
			for _, n2 := range nodes {
				if !proto.Equal(n1.Val.GetTileCoordinate(), n2.Val.GetTileCoordinate()) {
					// TODO(cripplet): Add INTRA edge calculation for L > 1.
					path, cost, err := astar.TileMapPath(
						tm,
						tm.TileFromCoordinate(n1.Val.GetTileCoordinate()),
						tm.TileFromCoordinate(n2.Val.GetTileCoordinate()),
						nodesCluster.Val.GetTileBoundary(),
						nodesCluster.Val.GetTileDimension(),
					)
					if err != nil {
						return nil, err
					}
					if path != nil {
						g.AddEdge(&rtsspb.AbstractEdge{
							Level:       g.L,
							Source:      n1.Val.GetTileCoordinate(),
							Destination: n2.Val.GetTileCoordinate(),
							EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
							Weight:      cost,
						})
					}
				}
			}
		}
	}

	/*
	for i := 2; i <= level; i++ {
		...
	}
	*/

	return g, nil
}
