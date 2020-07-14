// Package abstractgraph constructs and manages the abstract node space corresponding to a TileMap object.
// This package will be used ast he underlying topology for hiearchical A* searching.
package abstractgraph

import (
	"sync"

	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/utils"
)

type AbstractGraph struct {
	L int32

	Mu sync.RWMutex
	N map[utils.MapCoordinate][]*AbstractNode
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

	g.E[utils.MC(e.Val.GetSource())] = append(g.E[utils.MC(e.Val.GetSource())], e)
	g.E[utils.MC(e.Val.GetDestination())] = append(g.E[utils.MC(e.Val.GetDestination())], e)
	return nil
}
