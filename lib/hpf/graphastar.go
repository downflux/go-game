// Package graphastar defines fzipp.astar.Graph implementations for
// graph.Graph.
package graphastar

import (
	"math"

	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	fastar "github.com/fzipp/astar"
	"github.com/minkezhang/rts-pathing/lib/hpf/graph"
	"github.com/minkezhang/rts-pathing/lib/hpf/tile"
	"github.com/minkezhang/rts-pathing/lib/hpf/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

// We cannot do direct protobuf equality checking because of unexported
// field differences. We need to export the public fields into a controllable
// data struct.
type nodeImpl struct {
	Level          int32
	EphemeralKey   int32
	TileCoordinate utils.MapCoordinate
}

func importNodeImpl(n *rtsspb.AbstractNode) nodeImpl {
	return nodeImpl{
		Level:          n.GetLevel(),
		EphemeralKey:   n.GetEphemeralKey(),
		TileCoordinate: utils.MC(n.GetTileCoordinate()),
	}
}

func exportNodeImpl(n nodeImpl) *rtsspb.AbstractNode {
	return &rtsspb.AbstractNode{
		Level:          n.Level,
		EphemeralKey:   n.EphemeralKey,
		TileCoordinate: utils.PB(n.TileCoordinate),
	}
}

// dFunc provides a shim for the graph.Graph neighbor distance
// function.
func dFunc(g *graph.Graph, src, dest fastar.Node) float64 {
	cost, err := graph.D(g, exportNodeImpl(src.(nodeImpl)), exportNodeImpl(dest.(nodeImpl)))
	if err != nil {
		return math.Inf(0)
	}

	return cost
}

// hFunc provides a shim for the graph.Graph heuristic function.
func hFunc(src, dest fastar.Node) float64 {
	cost, err := graph.H(exportNodeImpl(src.(nodeImpl)), exportNodeImpl(dest.(nodeImpl)))
	if err != nil {
		return math.Inf(0)
	}

	return cost
}

// graphImpl implements fzipp.astar.Graph for the graph.Graph struct.
type graphImpl struct {
	// g holds information on how different AbstractNode objects are
	// connected via AbstractEdge links.
	g *graph.Graph
}

// Neighbours returns neighboring AbstractNode objects from a
// graph.Graph.
func (g graphImpl) Neighbours(n fastar.Node) []fastar.Node {
	neighbors, _ := g.g.Neighbors(exportNodeImpl(n.(nodeImpl)))
	var res []fastar.Node
	for _, n := range neighbors {
		res = append(res, importNodeImpl(n))
	}
	return res
}

// Path returns pathing information for two AbstractNode objects embedded in a
// graph.Graph. This function returns a (path, cost, error) tuple,
// where path is a list of AbstractNode objects and cost is the actual cost as
// calculated by calling D over the returned path. An empty path indicates
// there is no path found between the two AbstractNode objects.
func Path(tm *tile.Map, g *graph.Graph, src, dest *rtsspb.AbstractNode) ([]*rtsspb.AbstractNode, float64, error) {
	if tm == nil {
		return nil, 0, status.Error(codes.FailedPrecondition, "cannot have nil tile.Map input")
	}
	if g == nil {
		return nil, 0, status.Error(codes.FailedPrecondition, "cannot have nil graph.Graph input")
	}

	// TODO(minkezhang): Implement logic when multi-level AbstractGraph
	// objects are a thing.
	if src.GetLevel() != 1 || dest.GetLevel() != 1 {
		return nil, 0, notImplemented
	}

	d := func(a, b fastar.Node) float64 {
		return dFunc(g, a, b)
	}
	nodes := fastar.FindPath(graphImpl{g: g}, importNodeImpl(src), importNodeImpl(dest), d, hFunc)

	var res []*rtsspb.AbstractNode
	for _, node := range nodes {
		res = append(res, exportNodeImpl(node.(nodeImpl)))
	}

	return res, nodes.Cost(d), nil
}
