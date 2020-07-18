// Package abstractgraph constructs and manages the abstract node space corresponding to a TileMap object.
// This package will be used ast he underlying topology for hiearchical A* searching.
package abstractgraph

import (
	"math"
	"sync"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	// "github.com/cripplet/rts-pathing/lib/hpf/entrance"
	"github.com/cripplet/rts-pathing/lib/hpf/cluster"
	"github.com/cripplet/rts-pathing/lib/hpf/tile"
	"github.com/cripplet/rts-pathing/lib/hpf/utils"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

var (
	infinity = math.Inf(1)
)

type AbstractGraph struct {
	L int32

	Mu sync.RWMutex
	// N is a map from the Cluster coordinate to a list of nodes contained within that Cluster.
	N map[utils.MapCoordinate][]*AbstractNode
	// E is a map of Tile coordinates to the list of edges that connect to the Tle.
	E map[utils.MapCoordinate][]*AbstractEdge
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

	g.N[utils.MC(n.Val.GetTileCoordinate())] = append(g.N[utils.MC(n.Val.GetTileCoordinate())], n)
	return nil
}

func (g *AbstractGraph) AddEdge(pb *rtsspb.AbstractEdge) error {
	e, err := ImportAbstractEdge(pb)
	if err != nil {
		return err
	}

	g.Mu.Lock()
	defer g.Mu.Unlock()

	// Assuming symmetrical bidirectional graph.
	g.E[utils.MC(e.Val.GetSource())] = append(g.E[utils.MC(e.Val.GetSource())], e)
	g.E[utils.MC(e.Val.GetDestination())] = append(g.E[utils.MC(e.Val.GetDestination())], e)
	return nil
}

func BuildAbstractGraph(tm *tile.TileMap, cm *cluster.ClusterMap, ts []*rtsspb.Transition) (*AbstractGraph, error) {
	g := &AbstractGraph{
		L: cm.L,
	}

	for _, t := range ts {
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
	for _, ns := range g.N {
		for i, n1 := range ns {
			for j, n2 := range ns {
				if i != j {
					// TODO(cripplet): Add INTRA edge calculation for L > 1.
					_, cost, found, err := tile.Path(tm.TileFromCoordinate(n1.Val.GetTileCoordinate()), tm.TileFromCoordinate(n2.Val.GetTileCoordinate()))
					if err != nil {
						return nil, err
					}
					if found {
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
