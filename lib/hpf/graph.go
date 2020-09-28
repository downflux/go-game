// Package graph constructs and manages the abstract node space corresponding to a tile.Map object.
// This package will be used ast he underlying topology for hiearchical A* searching.
package graph

import (
	"math"
	"math/rand"

	rtscpb "github.com/minkezhang/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	"github.com/golang/protobuf/proto"
	"github.com/minkezhang/rts-pathing/lib/hpf/cluster"
	"github.com/minkezhang/rts-pathing/lib/hpf/edge"
	"github.com/minkezhang/rts-pathing/lib/hpf/entrance"
	"github.com/minkezhang/rts-pathing/lib/hpf/node"
	"github.com/minkezhang/rts-pathing/lib/hpf/tile"
	"github.com/minkezhang/rts-pathing/lib/hpf/tileastar"
	"github.com/minkezhang/rts-pathing/lib/hpf/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

// D gets exact cost between two neighboring AbstractNodes.
func D(g *Graph, src, dst *rtsspb.AbstractNode) (float64, error) {
	edge, err := g.EdgeMap.Get(utils.MC(src.GetTileCoordinate()), utils.MC(dst.GetTileCoordinate()))
	if err != nil {
		return 0, err
	}
	if edge == nil {
		return 0, status.Error(codes.NotFound, "an AbstractEdge does not exist with the given AbstractNode endpoints in the Graph")
	}

	return edge.GetWeight(), nil
}

// H gets the estimated cost of moving between two arbitrary AbstractNodes.
func H(src, dst *rtsspb.AbstractNode) (float64, error) {
	return math.Pow(float64(dst.GetTileCoordinate().GetX()-src.GetTileCoordinate().GetX()), 2) + math.Pow(float64(dst.GetTileCoordinate().GetY()-src.GetTileCoordinate().GetY()), 2), nil
}

// Graph contains the necessary state information to make an efficient
// path planning call on very large maps via hierarchical A* search, as
// described in Botea 2004.
type Graph struct {
	// NodeMap is a hash of AbstractNodes.
	NodeMap *node.Map

	// EdgeMap is a hash of AbstractEdges.
	EdgeMap *edge.Map
}

// BuildGraph build a higher-level representation of a tile.Map
// populated with information about how to travel between different subsections
// between tiles. tileDimension represents a subsection ("cluster") size.
func BuildGraph(tm *tile.Map, tileDimension *rtsspb.Coordinate) (*Graph, error) {
	// Create all node and edge map instances. These will be referenced and
	// mutated later on by passing the Graph object as a function
	// arg.
	cm, err := cluster.BuildMap(tm.D, tileDimension)
	if err != nil {
		return nil, err
	}

	g := &Graph{
		NodeMap: &node.Map{ClusterMap: cm},
		EdgeMap: &edge.Map{},
	}

	// Build the Tile-Tile edges which connect between two adjacent
	// clusters in the cluster.Map object and store this data into the
	// Graph.
	transitions, err := buildTransitions(tm, g.NodeMap.ClusterMap)
	if err != nil {
		return nil, err
	}
	for _, t := range transitions {
		g.NodeMap.Add(t.GetN1())
		g.NodeMap.Add(t.GetN2())
		g.EdgeMap.Add(&rtsspb.AbstractEdge{
			Source:      t.GetN1().GetTileCoordinate(),
			Destination: t.GetN2().GetTileCoordinate(),
			EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
			Weight:      1, // Inter-edges are always of cost 1, per Botea.
		})
	}

	// Build Tile-Tile edges within a cluster of a cluster.Map.
	for _, c := range cluster.Iterator(g.NodeMap.ClusterMap) {
		nodes, err := g.NodeMap.GetByCluster(c)
		if err != nil {
			return nil, err
		}
		for _, n1 := range nodes {
			for _, n2 := range nodes {
				if n1 != n2 {
					e, err := buildIntraEdge(tm, g.NodeMap.ClusterMap, n1, n2)
					if err != nil {
						return nil, err
					}

					if e != nil {
						g.EdgeMap.Add(e)
					}
				}
			}
		}
	}

	return g, nil
}

// AddEphemeralNode adds a temporary AbstractNode to the Graph and connects it
// to the rest of the cluster via AbstractEdge instances.
//
// Function returns a UUID which needs to be tracked by the caller to be passed
// into RemoveEphemeralNode. This function is a no-op if the input coordinates
// match a non-ephemeral AbstractNode.
//
// TODO(minkezhang): Support rollback in case errors happen so that
// AddEphemeralNode is idempotent.
func AddEphemeralNode(tm *tile.Map, g *Graph, t utils.MapCoordinate) (int64, error) {
	n, err := g.NodeMap.Get(t)
	if err != nil {
		return 0, err
	}

	if n == nil {
		n = &rtsspb.AbstractNode{IsEphemeral: true, TileCoordinate: utils.PB(t)}
		g.NodeMap.Add(n)
	}

	var ephemeralKey int64
	if n.GetIsEphemeral() {
		for _, found := n.GetEphemeralKeys()[ephemeralKey]; found || ephemeralKey == 0; ephemeralKey = rand.Int63() {
		}
		n.GetEphemeralKeys()[ephemeralKey] = true
	}

	c, err := cluster.ClusterCoordinateFromTileCoordinate(g.NodeMap.ClusterMap, t)
	if err != nil {
		return 0, err
	}

	borderNodes, err := g.NodeMap.GetByCluster(c)
	if err != nil {
		return 0, err
	}

	for _, borderNode := range borderNodes {
		boundary, err := cluster.TileBoundary(g.NodeMap.ClusterMap, c)
		if err != nil {
			return 0, err
		}
		dimension, err := cluster.TileDimension(g.NodeMap.ClusterMap, c)
		if err != nil {
			return 0, err
		}

		_, cost, err := tileastar.Path(
			tm,
			utils.MC(n.GetTileCoordinate()),
			utils.MC(borderNode.GetTileCoordinate()),
			utils.PB(boundary),
			utils.PB(dimension))
		if err != nil {
			return 0, err
		}
		if err := g.EdgeMap.Add(&rtsspb.AbstractEdge{
			Source:      n.GetTileCoordinate(),
			Destination: borderNode.GetTileCoordinate(),
			EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
			Weight:      cost,
		}); err != nil {
			return 0, err
		}
	}

	return ephemeralKey, nil
}

// RemoveEphemeralNode drops a temporary AbstractNode from the Graph. This is
// a no-op if the input coordinate is a non-ephemeral AbstractNode instance.
//
// TODO(minkezhang): Support rollback in case errors happen so that
// RemoveEphemeralNode is idempotent.
func RemoveEphemeralNode(g *Graph, t utils.MapCoordinate, ephemeralKey int64) error {
	n, err := g.NodeMap.Get(t)
	if err != nil {
		return err
	}

	delete(n.GetEphemeralKeys(), ephemeralKey)
	if n.GetIsEphemeral() && len(n.GetEphemeralKeys()) == 0 {
		g.NodeMap.Pop(t)
		edges, err := g.EdgeMap.GetBySource(t)
		if err != nil {
			return err
		}

		for _, e := range edges {
			if _, err := g.EdgeMap.Pop(utils.MC(e.GetSource()), utils.MC(e.GetDestination())); err != nil {
				return err
			}
		}
	}
	return nil
}

// buildIntraEdge constructs a single AbstractEdge instance with the correct
// traversal cost between two underlying AbstractNode objects. The cost
// function is calculated from the tile.Map entity, which holds information
// on e.g. the terrain information of the map.
func buildIntraEdge(tm *tile.Map, cm *cluster.Map, n1, n2 *rtsspb.AbstractNode) (*rtsspb.AbstractEdge, error) {
	c1, err := cluster.ClusterCoordinateFromTileCoordinate(cm, utils.MC(n1.GetTileCoordinate()))
	if err != nil {
		return nil, err
	}
	c2, err := cluster.ClusterCoordinateFromTileCoordinate(cm, utils.MC(n2.GetTileCoordinate()))
	if err != nil {
		return nil, err
	}
	if c1 != c2 {
		return nil, status.Errorf(codes.FailedPrecondition, "input AbstractNode instances are not bounded by the same cluster")
	}

	tileBoundary, err := cluster.TileBoundary(cm, c1)
	if err != nil {
		return nil, err
	}
	tileDimension, err := cluster.TileDimension(cm, c1)
	if err != nil {
		return nil, err
	}

	p, cost, err := tileastar.Path(
		tm,
		utils.MC(n1.GetTileCoordinate()),
		utils.MC(n2.GetTileCoordinate()),
		utils.PB(tileBoundary),
		utils.PB(tileDimension))
	if err != nil {
		return nil, err
	}

	if p != nil {
		return &rtsspb.AbstractEdge{
			Source:      n1.GetTileCoordinate(),
			Destination: n2.GetTileCoordinate(),
			EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
			Weight:      cost,
		}, nil
	}
	return nil, nil
}

// buildTransitions iterates over the tile.Map for the input cluster.Map overlay
// and look for adjacent, open nodes along cluster-cluster borders.
func buildTransitions(tm *tile.Map, cm *cluster.Map) ([]*rtsspb.Transition, error) {
	var ts []*rtsspb.Transition
	for _, c1 := range cluster.Iterator(cm) {
		neighbors, err := cluster.Neighbors(cm, c1)
		if err != nil {
			return nil, err
		}

		for _, c2 := range neighbors {
			if cluster.IsAdjacent(cm, c1, c2) && utils.LessThan(c1, c2) {
				transitions, err := entrance.BuildTransitions(tm, cm, c1, c2)
				if err != nil {
					return nil, err
				}
				ts = append(ts, transitions...)
			}
		}
	}
	return ts, nil
}

// Neighbors returns all adjacent AbstractNode instances in the AbstractGraph.
// Two AbstractNode instances are considered adjacent if there exists an edge
// defined between the two instances. Note that the instances returned here
// also include ephemeral AbstractNodes (n.GetIsEphemeral() == true) -- DFS
// should take care not to expand these second-order nodes.
func (g *Graph) Neighbors(n *rtsspb.AbstractNode) ([]*rtsspb.AbstractNode, error) {
	node, err := g.NodeMap.Get(utils.MC(n.GetTileCoordinate()))
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, status.Error(codes.FailedPrecondition, "cannot find specified node")
	}

	edges, err := g.EdgeMap.GetBySource(utils.MC(node.GetTileCoordinate()))
	if err != nil {
		return nil, err
	}

	var neighbors []*rtsspb.AbstractNode
	for _, e := range edges {
		var d *rtsspb.Coordinate
		if proto.Equal(node.GetTileCoordinate(), e.GetSource()) {
			d = e.GetDestination()
		} else {
			d = e.GetSource()
		}

		t, err := g.NodeMap.Get(utils.MC(d))
		if err != nil {
			return nil, err
		}
		if t == nil {
			return nil, status.Errorf(codes.NotFound, "invalid node coordinate %v specified for edge %v", d, e)
		}

		neighbors = append(neighbors, t)
	}
	return neighbors, nil
}
