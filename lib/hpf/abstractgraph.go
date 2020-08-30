// Package abstractgraph constructs and manages the abstract node space corresponding to a TileMap object.
// This package will be used ast he underlying topology for hiearchical A* searching.
package abstractgraph

import (
	"math"

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

// GetByClusterEdge filters the AbstractNodeMap by the input cluster coordinate
// and returns all AbstractNode objects that are bounded by the edges of the
// input.
func (m AbstractNodeMap) GetByClusterEdge(cm *cluster.ClusterMap, c utils.MapCoordinate) ([]*rtsspb.AbstractNode, error) {
	var res []*rtsspb.AbstractNode
	nodes, err := m.GetByCluster(cm, c)
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		if entrance.OnClusterEdge(cm, c, utils.MC(n.GetTileCoordinate())) {
			res = append(res, n)
		}
	}

	return res, nil
}

// GetByCluster filters the AbstractNodeMap by the input cluster coordinate
//  and returns all AbstractNode objects that are bounded by the input.
func (m AbstractNodeMap) GetByCluster(cm *cluster.ClusterMap, c utils.MapCoordinate) ([]*rtsspb.AbstractNode, error) {
	var nodes []*rtsspb.AbstractNode
	for _, n := range m {
		if cluster.CoordinateInCluster(cm, c, utils.MC(n.GetTileCoordinate())) {
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
//
// We're assuming the graph is undirected -- that is, for nodes A, B, if
// A --> B, then B --> A with the same cost.
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
	if _, found := m[d]; !found {
		m[d] = map[utils.MapCoordinate]*rtsspb.AbstractEdge{}
	}

	m[s][d] = e
	m[d][s] = e
	return nil
}

// Get queries the AbstractEdgeMap for an AbstractEdge instance which connects
// two TileMap Coordinate instances.
func (m AbstractEdgeMap) Get(s, d utils.MapCoordinate) (*rtsspb.AbstractEdge, error) {
	if _, found := m[s]; found {
		if e, found := m[s][d]; found {
			return e, nil
		}
	}
	return nil, nil
}

// GetBySource returns all edges in an AbstractEdgeMap which originate from the
// given source coordinate.
func (m AbstractEdgeMap) GetBySource(s utils.MapCoordinate) ([]*rtsspb.AbstractEdge, error) {
	var edges []*rtsspb.AbstractEdge
	for _, e := range m[s] {
		edges = append(edges, e)
	}
	return edges, nil
}

// AbstractGraph contains the necessary state information to make an efficient
// path planning call on very large maps via hierarchical A* search, as
// described in Botea 2004.
type AbstractGraph struct {
	// Level is the maximum hierarchy of AbstractNodes in this graph;
	// this is a positive, non-zero integer. The 0th level here loosely
	// refers to the underlying base map.
	Level int32

	// ClusterMap is a Level: ClusterMap dict representing all generated
	// tile boundaries for an AbstractGraph of the given Level. There
	// is a corresponding ClusterMap object per AbstractGraph.Level.
	ClusterMap map[int32]*cluster.ClusterMap

	// NodeMap contains a Level: AbstractNodeMap dict representing the
	// AbstractNodes per Level. As per AbstractGraph.ClusterMap, there
	// is a corresponding AbstractNodeMap object per level. Nodes within
	// a specific AbstractNodeMap may move between levels, and may be
	// deleted when the underlying terrain changes.
	NodeMap map[int32]AbstractNodeMap

	// EdgeMap contains a Level: AbstractEdgeMap dict representing the
	// AbstractEdges per Level. Edges may move between levels and may
	// be deleted when the underlying terrain changes.
	EdgeMap map[int32]AbstractEdgeMap
}

// BuildAbstractGraph build a higher-level representation of a TileMap
// populated with information about how to travel between different subsections
// between tiles. The level specified in input represents the number of
// abstractions that should be built for this map (l > 1 is useful for very,
// very large maps), and the clusterDimension represents a subsection size.
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
	if (int32(math.Pow(float64(clusterDimension.GetX()), float64(level))) >= tm.D.GetX()) || (int32(math.Pow(float64(clusterDimension.GetY()), float64(level))) >= tm.D.GetY()) {
		return nil, status.Error(codes.FailedPrecondition, "given clusterDimension and level will result in too large a cluster map")
	}

	// This does not add any value for an AbstractGraph.
	if clusterDimension.GetX() <= 1 && clusterDimension.GetY() <= 1 {
		return nil, status.Error(codes.FailedPrecondition, "invalid clusterDimension")
	}

	clusterMaps, err := buildTieredClusterMaps(tm, level, clusterDimension)
	if err != nil {
		return nil, err
	}

	transitions, err := buildTransitions(tm, clusterMaps[1])
	if err != nil {
		return nil, err
	}

	nm := AbstractNodeMap{}
	em := AbstractEdgeMap{}

	nodes, err := buildBaseNodes(transitions)
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		nm.Add(n)
	}

	edges, err := buildBaseInterEdges(tm, transitions)
	if err != nil {
		return nil, err
	}

	for _, e := range edges {
		em.Add(e)
	}

	edges, err = buildBaseIntraEdges(tm, clusterMaps[1], nm)
	if err != nil {
		return nil, err
	}

	for _, e := range edges {
		em.Add(e)
	}

	g := &AbstractGraph{
		Level:      level,
		ClusterMap: clusterMaps,
		NodeMap:    map[int32]AbstractNodeMap{1: nm},
		EdgeMap:    map[int32]AbstractEdgeMap{1: em},
	}

	return g, nil
}

// Neighbors returns all adjacent AbstractNode instances of the same hierarchy
// level of the input. Two AbstractNode instances are considered adjacent if
// there exists an edge defined between the two instances. Note that the
// instances returned here also include ephemeral AbstractNodes
// (n.GetEphemeralKey() > 0) -- DFS should take care not to expand these
// secondary neighbors.
func (g *AbstractGraph) Neighbors(n *rtsspb.AbstractNode) ([]*rtsspb.AbstractNode, error) {
	nm, found := g.NodeMap[n.GetLevel()]
	if !found {
		return nil, status.Error(codes.FailedPrecondition, "invalid level specified for input")
	}

	node, err := nm.Get(utils.MC(n.GetTileCoordinate()))
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, status.Error(codes.FailedPrecondition, "cannot find specified node")
	}

	em, found := g.EdgeMap[n.GetLevel()]
	if !found {
		return nil, status.Error(codes.NotFound, "EdgeMap of specified level not found")
	}

	edges, err := em.GetBySource(utils.MC(node.GetTileCoordinate()))
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

		t, err := nm.Get(utils.MC(d))
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

func buildBaseNodes(transitions []*rtsspb.Transition) ([]*rtsspb.AbstractNode, error) {
	var res []*rtsspb.AbstractNode
	for _, t := range transitions {
		res = append(res, t.GetN1())
		res = append(res, t.GetN2())
	}
	return res, nil
}

func buildBaseInterEdges(tm *tile.TileMap, transitions []*rtsspb.Transition) ([]*rtsspb.AbstractEdge, error) {
	var res []*rtsspb.AbstractEdge
	for _, t := range transitions {
		res = append(res, &rtsspb.AbstractEdge{
			Level:       1,
			Source:      t.GetN1().GetTileCoordinate(),
			Destination: t.GetN2().GetTileCoordinate(),
			EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTER,
			Weight:      1, // Inter-edges are always of cost 1, per Botea.
		})
	}
	return res, nil
}

// buildBaseIntraEdges generates a list of AbstractEdges corresponding to a
// totally connected graph of the AbstractNodes for each Cluster in a
// ClusterMap object.
func buildBaseIntraEdges(tm *tile.TileMap, cm *cluster.ClusterMap, nm AbstractNodeMap) ([]*rtsspb.AbstractEdge, error) {
	if cm.Val.GetLevel() > 1 {
		return nil, notImplemented
	}

	var edges []*rtsspb.AbstractEdge
	for _, c := range cluster.Iterator(cm) {
		// TODO(cripplet): Determine if we only need GetByClusterEdge
		// instead here.
		nodes, err := nm.GetByClusterEdge(cm, c)
		if err != nil {
			return nil, err
		}

		for _, n1 := range nodes {
			for _, n2 := range nodes {
				if n1 != n2 && cm.Val.GetLevel() == n1.GetLevel() && cm.Val.GetLevel() == n2.GetLevel() {
					tileBoundary, err := cluster.TileBoundary(cm, c)
					if err != nil {
						return nil, err
					}
					tileDimension, err := cluster.TileDimension(cm, c)
					if err != nil {
						return nil, err
					}

					p, cost, err := astar.TileMapPath(
						tm,
						tm.TileFromCoordinate(n1.GetTileCoordinate()),
						tm.TileFromCoordinate(n2.GetTileCoordinate()),
						utils.PB(tileBoundary),
						utils.PB(tileDimension))
					if err != nil {
						return nil, err
					}
					if p != nil {
						edges = append(edges, &rtsspb.AbstractEdge{
							Level:       cm.Val.GetLevel(),
							Source:      n1.GetTileCoordinate(),
							Destination: n2.GetTileCoordinate(),
							EdgeType:    rtscpb.EdgeType_EDGE_TYPE_INTRA,
							Weight:      cost,
						})
					}
				}
			}
		}
	}
	return edges, nil
}

// buildTieredClusterMaps constructs a list of ClusterMap objects. Each set of
// Cluster instances inside a ClusterMap are completely nested inside the
// corresponding Cluster instances in a higher-level ClusterMap.
func buildTieredClusterMaps(tm *tile.TileMap, level int32, clusterDimension *rtsspb.Coordinate) (map[int32]*cluster.ClusterMap, error) {
	cms := map[int32]*cluster.ClusterMap{}
	for i := int32(1); i <= level; i++ {
		cm, err := cluster.BuildClusterMap(
			tm.D,
			&rtsspb.Coordinate{
				X: int32(math.Pow(float64(clusterDimension.GetX()), float64(i))),
				Y: int32(math.Pow(float64(clusterDimension.GetY()), float64(i))),
			},
			i,
		)
		if err != nil {
			return nil, err
		}
		cms[i] = cm
	}
	return cms, nil
}

// buildTransitions iterates over the TileMap for the input ClusterMap overlay
// and look for adjacent, open nodes along Cluster-Cluster borders.
func buildTransitions(tm *tile.TileMap, cm *cluster.ClusterMap) ([]*rtsspb.Transition, error) {
	var ts []*rtsspb.Transition
	for _, c1 := range cluster.Iterator(cm) {
		neighbors, err := cluster.Neighbors(cm, c1)
		if err != nil {
			return nil, err
		}

		for _, c2 := range neighbors {
			if cluster.IsAdjacent(cm, c1, c2) {
				transitions, err := entrance.BuildTransitions(tm, cm, c1, c2)
				if err != nil {
					return nil, err
				}
				ts = append(ts, transitions...)
			}
		}
	}
	for _, t := range ts {
		t.GetN1().Level = cm.Val.GetLevel()
		t.GetN2().Level = cm.Val.GetLevel()
	}
	return ts, nil
}
