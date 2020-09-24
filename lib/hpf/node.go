// Package node constructs and manages the abstract node space
// corresponding to a TileMap object.
package node

import (
	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	"github.com/minkezhang/rts-pathing/lib/hpf/cluster"
	"github.com/minkezhang/rts-pathing/lib/hpf/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Map contains a collection of AbstractNode instances, which
// represent an AbstractGraph node used for hierarchical A* search.
//
// AbstractNodes are indexed by cluster coordinate and then Tile coordinate.
type Map struct {
	ClusterMap *cluster.Map
	nodes      map[utils.MapCoordinate]map[utils.MapCoordinate]*rtsspb.AbstractNode
}

// GetByCluster filters the Map by the input cluster coordinate
// and returns all AbstractNode objects that are bounded by the input.
func (nm Map) GetByCluster(c utils.MapCoordinate) ([]*rtsspb.AbstractNode, error) {
	if err := cluster.ValidateClusterInRange(nm.ClusterMap, c); err != nil {
		return nil, err
	}

	if nm.nodes == nil {
		return nil, nil
	}

	var nodes []*rtsspb.AbstractNode
	for _, n := range nm.nodes[c] {
		if cluster.CoordinateInCluster(nm.ClusterMap, c, utils.MC(n.GetTileCoordinate())) {
			nodes = append(nodes, n)
		}
	}
	return nodes, nil
}

func (nm *Map) Get(t utils.MapCoordinate) (*rtsspb.AbstractNode, error) {
	c, err := cluster.ClusterCoordinateFromTileCoordinate(nm.ClusterMap, t)
	if err != nil {
		return nil, err
	}

	if nm.nodes == nil || nm.nodes[c] == nil {
		return nil, nil
	}

	return nm.nodes[c][t], nil
}

func (nm *Map) Pop(t utils.MapCoordinate) (*rtsspb.AbstractNode, error) {
	n, err := nm.Get(t)
	if err != nil {
		return nil, err
	}

	if n != nil {
		c, _ := cluster.ClusterCoordinateFromTileCoordinate(nm.ClusterMap, t)
		delete(nm.nodes[c], t)
	}

	return n, nil
}

func (nm *Map) Add(n *rtsspb.AbstractNode) error {
	if n.GetLevel() != nm.ClusterMap.Val.GetLevel() {
		return status.Error(codes.FailedPrecondition, "input mismatch, given AbstractNode does not have the same hierarchy level as the cluster.Map bound to the tile.Map")
	}

	t := utils.MC(n.GetTileCoordinate())

	existingNode, err := nm.Get(t)
	if err != nil {
		return err
	}

	if existingNode != nil {
		return status.Errorf(codes.AlreadyExists, "an AbstractNode already exists for Map with the given Coordinate %v", n.GetTileCoordinate())
	}

	c, err := cluster.ClusterCoordinateFromTileCoordinate(nm.ClusterMap, utils.MC(n.GetTileCoordinate()))
	if err != nil {
		return err
	}

	if nm.nodes == nil {
		nm.nodes = map[utils.MapCoordinate]map[utils.MapCoordinate]*rtsspb.AbstractNode{}
	}
	if nm.nodes[c] == nil {
		nm.nodes[c] = map[utils.MapCoordinate]*rtsspb.AbstractNode{}
	}

	nm.nodes[c][t] = n

	return nil
}
