// Package graph constructs and manages the abstract node space corresponding to a tile.Map object.
// This package will be used ast he underlying topology for hiearchical A* searching.
package graph

import (
	"math"
	"math/rand"

	gdpb "github.com/downflux/game/api/data_go_proto"
	rtscpb "github.com/downflux/game/pathing/proto/constants_go_proto"
	rtsspb "github.com/downflux/game/pathing/proto/structs_go_proto"

	"github.com/golang/protobuf/proto"
	"github.com/downflux/game/pathing/hpf/cluster"
	"github.com/downflux/game/pathing/hpf/edge"
	"github.com/downflux/game/pathing/hpf/entrance"
	"github.com/downflux/game/pathing/hpf/node"
	"github.com/downflux/game/pathing/hpf/tile"
	"github.com/downflux/game/pathing/hpf/tileastar"
	"github.com/downflux/game/pathing/hpf/utils"
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
func BuildGraph(tm *tile.Map, tileDimension *gdpb.Coordinate) (*Graph, error) {
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
		for _, n := range nodes {
			if err := connect(tm, g, utils.MC(n.GetTileCoordinate())); err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}

// connect takes as input an AbstractNode, builds all possible INTRA_EDGE
// AbstractEdge instances within the same cluster, and inserts them into
// the Graph.
//
// connect does not rebuild edges if they already exist between two nodes, and
// does not create an edge between two ephemeral AbstractNode instances.
func connect(tm *tile.Map, g *Graph, t utils.MapCoordinate) error {
	n1, err := g.NodeMap.Get(t)
	if err != nil {
		return err
	}

	c, err := cluster.ClusterCoordinateFromTileCoordinate(g.NodeMap.ClusterMap, t)
	if err != nil {
		return err
	}

	borderNodes, err := g.NodeMap.GetByCluster(c)
	if err != nil {
		return err
	}

	for _, n2 := range borderNodes {
		e, err := g.EdgeMap.Get(utils.MC(n1.GetTileCoordinate()), utils.MC(n2.GetTileCoordinate()))
		if err != nil {
			return err
		}
		// Ephemeral nodes are invisible to one another and should not be
		// considered neighbors.
		//
		// Also don't re-run pathing if we don't have to -- we're
		// assuming the TileMap at this point is static.
		if e != nil || (n1.GetIsEphemeral() && n2.GetIsEphemeral()) {
			continue
		}

		tileBoundary, err := cluster.TileBoundary(g.NodeMap.ClusterMap, c)
		if err != nil {
			return err
		}
		tileDimension, err := cluster.TileDimension(g.NodeMap.ClusterMap, c)
		if err != nil {
			return err
		}

		p, cost, err := tileastar.Path(
			tm,
			utils.MC(n1.GetTileCoordinate()),
			utils.MC(n2.GetTileCoordinate()),
			utils.PB(tileBoundary),
			utils.PB(tileDimension))
		if err != nil {
			return err
		}

		if p != nil {
			g.EdgeMap.Add(&rtsspb.AbstractEdge{
				Source:      n1.GetTileCoordinate(),
				Destination: n2.GetTileCoordinate(),
				EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
				Weight:      cost,
			})
		}
	}

	return nil
}

// AddEphemeralNode adds a temporary AbstractNode to the Graph and connects it
// to the rest of the cluster via AbstractEdge instances.
//
// Function returns a UUID which needs to be tracked by the caller to be passed
// into RemoveEphemeralNode. This function is a no-op if the input coordinates
// match a non-ephemeral AbstractNode.
//
// TODO(minkezhang): Support rollback in case errors happen so that
// InsertEphemeralNode is idempotent.
func InsertEphemeralNode(tm *tile.Map, g *Graph, t utils.MapCoordinate) (int64, error) {
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
		if n.GetEphemeralKeys() == nil {
			n.EphemeralKeys = map[int64]bool{}
		}
		n.GetEphemeralKeys()[ephemeralKey] = true
	}

	if err := connect(tm, g, t); err != nil {
		return 0, err
	}
	return ephemeralKey, nil
}

// disconnect takes as input an AbstractNode and removes all edges connected
// to it, as well as remove the node itself from the Graph.
func disconnect(g *Graph, t utils.MapCoordinate) error {
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
	return nil
}

// RemoveEphemeralNode drops a temporary AbstractNode from the Graph. This is
// a no-op if the input coordinate is a non-ephemeral AbstractNode instance.
//
// No error will be raised if the input key is non-existent.
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
		return disconnect(g, t)
	}
	return nil
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
		var d *gdpb.Coordinate
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
