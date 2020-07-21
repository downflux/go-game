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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	infinity = math.Inf(1)
)

type AbstractGraph struct {
	L int32

	Mu sync.RWMutex
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

	g := &AbstractGraph{
		L: level,
	}

	var transitions []*rtsspb.Transition

	cm, err := cluster.BuildClusterMap(tm.D, &rtsspb.Coordinate{X: clusterDimension.GetX(), Y: clusterDimension.GetY()})
	if err != nil {
		return nil, err
	}

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

	for _, t := range transitions {
		g.AddNode(&rtsspb.AbstractNode{
			Level:             g.L,
			ClusterCoordinate: t.GetN1().GetClusterCoordinate(),
			TileCoordinate:    t.GetN1().GetTileCoordinate(),
		})
		g.AddNode(&rtsspb.AbstractNode{
			Level:             g.L,
			ClusterCoordinate: t.GetN2().GetClusterCoordinate(),
			TileCoordinate:    t.GetN2().GetTileCoordinate(),
		})
		g.AddEdge(&rtsspb.AbstractEdge{
			Level:       g.L,
			Source:      t.GetN1().GetTileCoordinate(),
			Destination: t.GetN2().GetTileCoordinate(),
			EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
			Weight:      1, // Inter-edges are always of cost 1, per Botea.
		})
	}

	// Add intra-edges for all tiles within the same Cluster.
	for _, clusterNodes := range g.N[1] {
		for c1, n1 := range clusterNodes {
			for c2, n2 := range clusterNodes {
				if c1 != c2 {
					// TODO(cripplet): Add INTRA edge calculation for L > 1.
					path, cost, err := astar.TileMapPath(
						tm,
						tm.TileFromCoordinate(n1.Val.GetTileCoordinate()),
						tm.TileFromCoordinate(n2.Val.GetTileCoordinate()))
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

	return g, nil
}
