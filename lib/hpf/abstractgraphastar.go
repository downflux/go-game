// Package abstractgraphastar defines fzipp.astar.Graph implementations for
// abstractgraph.Graph.
package abstractgraphastar

import (
	"math"

	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	fastar "github.com/fzipp/astar"
	"github.com/minkezhang/rts-pathing/lib/hpf/abstractgraph"
	"github.com/minkezhang/rts-pathing/lib/hpf/tile"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

// dFunc provides a shim for the abstractgraph.Graph neighbor distance
// function.
func dFunc(g *abstractgraph.Graph, src, dest fastar.Node) float64 {
	cost, err := abstractgraph.D(g, src.(*rtsspb.AbstractNode), dest.(*rtsspb.AbstractNode))
	if err != nil {
		return math.Inf(0)
	}

	return cost
}

// hFunc provides a shim for the abstractgraph.Graph heuristic function.
func hFunc(src, dest fastar.Node) float64 {
	cost, err := abstractgraph.H(src.(*rtsspb.AbstractNode), dest.(*rtsspb.AbstractNode))
	if err != nil {
		return math.Inf(0)
	}

	return cost
}

// graph implements fzipp.astar.Graph for the abstractgraph.Graph struct.
type graph struct {
	g *abstractgraph.Graph
}

// Neighbours returns neighboring AbstractNode objects from a
// abstractgraph.Graph.
func (g graph) Neighbours(n fastar.Node) []fastar.Node {
	neighbors, _ := g.g.Neighbors(n.(*rtsspb.AbstractNode))
	var res []fastar.Node
	for _, n := range neighbors {
		res = append(res, n)
	}
	return res
}

// Path returns pathing information for two AbstractNode objects embedded in a
// abstractgraph.Graph. This function returns a (path, cost, error) tuple,
// where path is a list of AbstractNode objects and cost is the actual cost as
// calculated by calling D over the returned path. An empty path indicates
// there is no path found between the two AbstractNode objects.
func Path(tm *tile.Map, g *abstractgraph.Graph, src, dest *rtsspb.AbstractNode) ([]*rtsspb.AbstractNode, float64, error) {
	if tm == nil {
		return nil, 0, status.Error(codes.FailedPrecondition, "cannot have nil tile.Map input")
	}
	if g == nil {
		return nil, 0, status.Error(codes.FailedPrecondition, "cannot have nil abstractgraph.Graph input")
	}

	// TODO(minkezhang): Implement logic when multi-level AbstractGraph
	// objects are a thing.
	if src.GetLevel() != 1 || dest.GetLevel() != 1 {
		return nil, 0, notImplemented
	}

	d := func(a, b fastar.Node) float64 {
		return dFunc(g, a, b)
	}
	nodes := fastar.FindPath(graph{g: g}, src, dest, d, hFunc)

	var res []*rtsspb.AbstractNode
	for _, node := range nodes {
		res = append(res, node.(*rtsspb.AbstractNode))
	}

	return res, nodes.Cost(d), nil
}
