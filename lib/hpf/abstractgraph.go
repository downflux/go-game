// Package abstractgraph constructs and manages the abstract node space corresponding to a TileMap object.
// This package will be used ast he underlying topology for hiearchical A* searching.
package abstractgraph

import (
	"math"

	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/cluster"
	"github.com/cripplet/rts-pathing/lib/hpf/entrance"
	"github.com/cripplet/rts-pathing/lib/hpf/utils"
	"github.com/cripplet/rts-pathing/lib/hpf/tile"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AbstractNodeMap map[utils.MapCoordinate]*rtsspb.AbstractNode
type AbstractEdgeMap map[utils.MapCoordinate]map[utils.MapCoordinate]*rtsspb.AbstractEdge

func (m AbstractNodeMap) Remove(c *rtsspb.Coordinate) error {
	delete(m, utils.MC(c))
	return nil
}

func (m AbstractNodeMap) GetByClusterEdge(c *cluster.Cluster) ([]*rtsspb.AbstractNode, error) {
	var res []*rtsspb.AbstractNode
	nodes, err := m.GetByCluster(c)
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		onClusterEdge, err := entrance.OnClusterEdge(c, n.GetTileCoordinate())
		if err != nil {
			return nil, err
		}
		if onClusterEdge {
			res = append(res, n)
		}
	}

	return res, nil
}

func (m AbstractNodeMap) GetByCluster(c *cluster.Cluster) ([]*rtsspb.AbstractNode, error) {
	var nodes []*rtsspb.AbstractNode
	for _, n := range m {
		if cluster.CoordinateInCluster(n.GetTileCoordinate(), c) {
			nodes = append(nodes, n)
		}
	}
	return nodes, nil
}

func (m AbstractNodeMap) Add(n *rtsspb.AbstractNode) error {
	m[utils.MC(n.GetTileCoordinate())] = n
	return nil
}

func (m AbstractEdgeMap) Remove(s *rtsspb.Coordinate) error {
	if destinations, found := m[utils.MC(s)]; found {
		for d := range destinations {
			delete(m[utils.MC(s)], d)
		}
	}
	delete(m, utils.MC(s))
	return nil
}

func (m AbstractEdgeMap) Add(e *rtsspb.AbstractEdge) error {
	s := utils.MC(e.GetSource())
	d := utils.MC(e.GetDestination())

	if _, found := m[s]; !found {
		m[s] = map[utils.MapCoordinate]*rtsspb.AbstractEdge{}
	}
	if _, found := m[d]; !found {
		m[d] = map[utils.MapCoordinate]*rtsspb.AbstractEdge{}
	}

	// Assuming symmetrical bidirectional graph.
	m[s][d] = e
	m[d][s] = e
	return nil
}

type AbstractGraph struct {
	Level int32

	ClusterMap map[int32]*cluster.ClusterMap
	NodeMap    map[int32]map[utils.MapCoordinate]*rtsspb.AbstractNode
	EdgeMap    map[int32]map[utils.MapCoordinate]map[utils.MapCoordinate]*rtsspb.AbstractEdge
}

func BuildAbstractGraph(tm *tile.TileMap, level int32, clusterDimension *rtsspb.Coordinate) (*AbstractGraph, error) {
	if level < 1 {
		return nil, status.Error(codes.FailedPrecondition, "level must be a positive non-zero integer")
	}
	// Highest level ClusterMap should still have more than one Cluster,
	// otherwise we'll be routing units to the edge first before going back
	// inwards.
	if (
		int32(math.Pow(float64(level), float64(clusterDimension.GetX()))) >= tm.D.GetX()) || (
		int32(math.Pow(float64(level), float64(clusterDimension.GetY()))) >= tm.D.GetY()) {
		return nil, status.Error(codes.FailedPrecondition, "given clusterDimension and level will result in too large a cluster map")
	}

	clusterMaps, err := buildTieredClusterMaps(tm, level, clusterDimension)
	if err != nil {
		return nil, err
	}

	g := &AbstractGraph{
		Level: level,
		ClusterMap: clusterMaps,
	}
	/*

	if err = buildBaseGraph(g, tm); err != nil {
		return nil, err
	}

	// Add intra-edges for l > 1.
	for i := int32(2); i <= level; i++ {
		if err := addGraphLevel(g, tm, i); err != nil {
			return nil, err
		}
	} */

	return g, nil
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
