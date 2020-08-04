// Package abstractgraph constructs and manages the abstract node space corresponding to a TileMap object.
// This package will be used ast he underlying topology for hiearchical A* searching.
package abstractgraph

import (
	"math"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/cluster"
	"github.com/cripplet/rts-pathing/lib/hpf/entrance"
	"github.com/cripplet/rts-pathing/lib/hpf/tile"
	"github.com/cripplet/rts-pathing/lib/hpf/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

// AbstractNodeMap contains a collection of AbstractNode instances, which
// represent an AbstractGraph node used for hierarchical A* search.
type AbstractNodeMap map[utils.MapCoordinate]*rtsspb.AbstractNode

// AbstractEdgeMap contains a collection of AbstractEdge instances, which
// represent an AbstractGraph edge; these edges represent the cost to move
// between different AbstractNode instances.
type AbstractEdgeMap map[utils.MapCoordinate]map[utils.MapCoordinate]*rtsspb.AbstractEdge

// Remove deletes the specified AbstractNode from the AbstractNodeMap.
func (m AbstractNodeMap) Remove(c utils.MapCoordinate) error {
	delete(m, c)
	return nil
}

func (m AbstractNodeMap) GetByClusterEdge(c *cluster.Cluster) ([]*rtsspb.AbstractNode, error) {
	var res []*rtsspb.AbstractNode
	nodes, err := m.GetByCluster(c)
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		if entrance.OnClusterEdge(c, n.GetTileCoordinate()) {
			res = append(res, n)
		}
	}

	return res, nil
}

// GetByCluster filters the AbstractNodeMap by the input Cluster and returns
// all AbstractNode objects that are bounded by the input.
func (m AbstractNodeMap) GetByCluster(c *cluster.Cluster) ([]*rtsspb.AbstractNode, error) {
	var nodes []*rtsspb.AbstractNode
	for _, n := range m {
		if cluster.CoordinateInCluster(n.GetTileCoordinate(), c) {
			nodes = append(nodes, n)
		}
	}
	return nodes, nil
}

// Add appends an AbstractNode instance into the AbstractNodeMap collection.
func (m AbstractNodeMap) Add(n *rtsspb.AbstractNode) error {
	m[utils.MC(n.GetTileCoordinate())] = n
	return nil
}

// Get queries the AbstractNodeMap for the AbstractNode instance at a specific
// TileMap Coordinate.
func (m AbstractNodeMap) Get(c utils.MapCoordinate) (*rtsspb.AbstractNode, error) {
	return m[c], nil
}

// Remove deletes the specified AbstractEdge from the AbstractEdgeMap.
func (m AbstractEdgeMap) Remove(s, d utils.MapCoordinate) error {
	if _, found := m[s]; found {
		delete(m[s], d)
	}
	if _, found := m[d]; found {
		delete(m[d], s)
	}
	return nil
}

// Add appends an AbstractEdge instance into the AbstractEdgeMap collection.
func (m AbstractEdgeMap) Add(e *rtsspb.AbstractEdge) error {
	s := utils.MC(e.GetSource())
	d := utils.MC(e.GetDestination())

	edge, err := m.Get(s, d)
	if err != nil {
		return err
	}
	if edge != nil {
		return status.Errorf(codes.AlreadyExists, "AbstractEdge unexpectedly found at %v, %v", s, d)
	}

	if _, found := m[s]; !found {
		m[s] = map[utils.MapCoordinate]*rtsspb.AbstractEdge{}
	}

	m[s][d] = e
	return nil
}

// Get queries the AbstractEdgeMap for an AbstractEdge instance which connects
// two TileMap Coordinate instances.
func (m AbstractEdgeMap) Get(s, d utils.MapCoordinate) (*rtsspb.AbstractEdge, error) {
	if _, found := m[s]; found {
		return m[s][d], nil
	}
	if _, found := m[d]; found {
		return m[d][s], nil
	}
	return nil, nil
}

type AbstractGraph struct {
	Level int32

	ClusterMap map[int32]*cluster.ClusterMap
	NodeMap    map[int32]AbstractNodeMap
	EdgeMap    map[int32]AbstractEdgeMap
}

func BuildAbstractGraph(tm *tile.TileMap, level int32, clusterDimension *rtsspb.Coordinate) (*AbstractGraph, error) {
	if level < 1 {
		return nil, status.Error(codes.FailedPrecondition, "level must be a positive non-zero integer")
	}

	// TODO(cripplet): Add higher-level node generation.
	if level >= 2 {
		return nil, notImplemented
	}

	// Highest level ClusterMap should still have more than one Cluster,
	// otherwise we'll be routing units to the edge first before going back
	// inwards.
	if (int32(math.Pow(float64(level), float64(clusterDimension.GetX()))) >= tm.D.GetX()) || (int32(math.Pow(float64(level), float64(clusterDimension.GetY()))) >= tm.D.GetY()) {
		return nil, status.Error(codes.FailedPrecondition, "given clusterDimension and level will result in too large a cluster map")
	}

	clusterMaps, err := buildTieredClusterMaps(tm, level, clusterDimension)
	if err != nil {
		return nil, err
	}

	transitions, err := buildTransitions(clusterMaps[1], tm)
	if err != nil {
		return nil, err
	}

	nm, err := buildBaseNodes(transitions)
	if err != nil {
		return nil, err
	}

	em, err := buildBaseEdges(transitions, tm)
	if err != nil {
		return nil, err
	}

	g := &AbstractGraph{
		Level:      level,
		ClusterMap: clusterMaps,
		NodeMap:    map[int32]AbstractNodeMap{1: nm},
		EdgeMap:    map[int32]AbstractEdgeMap{1: em},
	}

	return g, nil
}

func buildBaseNodes(transitions []*rtsspb.Transition) (AbstractNodeMap, error) {
	nm := AbstractNodeMap{}
	for _, t := range transitions {
		nm.Add(t.GetN1())
		nm.Add(t.GetN2())
	}
	return nm, nil
}

func buildBaseEdges(transitions []*rtsspb.Transition, tm *tile.TileMap) (AbstractEdgeMap, error) {
	em := AbstractEdgeMap{}
	for _, t := range transitions {
		em.Add(&rtsspb.AbstractEdge{
			Level:       1,
			Source:      t.GetN1().GetTileCoordinate(),
			Destination: t.GetN2().GetTileCoordinate(),
			EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
			Weight:      1, // Inter-edges are always of cost 1, per Botea.
		})
	}
	return em, nil
}

func buildTieredClusterMaps(tm *tile.TileMap, level int32, clusterDimension *rtsspb.Coordinate) (map[int32]*cluster.ClusterMap, error) {
	cms := map[int32]*cluster.ClusterMap{}
	for i := int32(1); i <= level; i++ {
		cm, err := cluster.BuildClusterMap(tm.D, clusterDimension, i)
		if err != nil {
			return nil, err
		}
		cms[i] = cm
	}
	return cms, nil
}

func buildTransitions(cm *cluster.ClusterMap, tm *tile.TileMap) ([]*rtsspb.Transition, error) {
	var ts []*rtsspb.Transition
	for _, c1 := range cm.M {
		for _, c2 := range cm.M {
			if cluster.IsAdjacent(c1, c2) {
				transitions, err := entrance.BuildTransitions(c1, c2, tm)
				if err != nil {
					return nil, err
				}
				ts = append(ts, transitions...)
			}
		}
	}
	return ts, nil
}
