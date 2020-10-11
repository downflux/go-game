// Package edge constructs and manages the abstract edge space
// corresponding to a TileMap object.
package edge

import (
	pdpb "github.com/downflux/game/pathing/api/data_go_proto"

	"github.com/downflux/game/map/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Map contains a collection of AbstractEdge instances, which
// represent an AbstractGraph edge; these edges represent the cost to move
// between different AbstractNode instances.
type Map struct {
	// edges contain a (Tile, Tile) coordinate tuple-indexed list of
	// AbstractEdges, where the coordinates represent WLOG the source and
	// destination Tile instances in the tile.Map.
	//
	// We want to explicitly disallow direct access.
	edges map[utils.MapCoordinate]map[utils.MapCoordinate]*pdpb.AbstractEdge
}

// Iterator provides a flattened list of AbstractEdge instances constructed
// from the 2D map Map.edges object.
func (em *Map) Iterator() []*pdpb.AbstractEdge {
	var edges []*pdpb.AbstractEdge
	for _, row := range em.edges {
		for _, e := range row {
			edges = append(edges, e)
		}
	}
	return edges
}

// GetBySource returns a list of AbstractEdge instances which WLOG originate
// from the specified Tile coordinate.
func (em *Map) GetBySource(t utils.MapCoordinate) ([]*pdpb.AbstractEdge, error) {
	if em.edges == nil {
		return nil, nil
	}

	var res []*pdpb.AbstractEdge
	for _, e := range em.edges[t] {
		res = append(res, e)
	}
	return res, nil
}

// Add appends an AbstractEdge instance into the Map collection.
//
// We're assuming the graph is undirected -- that is, for nodes A, B, if
// A --> B, then B --> A with the same cost.
func (em *Map) Add(e *pdpb.AbstractEdge) error {
	t1 := utils.MC(e.GetSource())
	t2 := utils.MC(e.GetDestination())

	if t1 == t2 {
		return status.Errorf(codes.FailedPrecondition, "AbstractEdge may not specify the same source and destination")
	}

	edge, err := em.Get(t1, t2)
	if err != nil {
		return err
	}
	if edge != nil {
		return status.Errorf(codes.AlreadyExists, "AbstractEdge unexpectedly found at %v, %v", t1, t2)
	}

	if em.edges == nil {
		em.edges = map[utils.MapCoordinate]map[utils.MapCoordinate]*pdpb.AbstractEdge{}
	}
	if _, found := em.edges[t1]; !found {
		em.edges[t1] = map[utils.MapCoordinate]*pdpb.AbstractEdge{}
	}
	if _, found := em.edges[t2]; !found {
		em.edges[t2] = map[utils.MapCoordinate]*pdpb.AbstractEdge{}
	}

	em.edges[t1][t2] = e
	em.edges[t2][t1] = e

	return nil
}

// Get queries the Map for an AbstractEdge instance which connects
// two Tile coordinates.
func (em *Map) Get(t1, t2 utils.MapCoordinate) (*pdpb.AbstractEdge, error) {
	if em.edges == nil {
		return nil, nil
	}

	if _, found := em.edges[t1]; found {
		if e, found := em.edges[t1][t2]; found {
			return e, nil
		}
	}

	return nil, nil
}

// Pop deletes the specified AbstractEdge from the Map.
func (em *Map) Pop(t1, t2 utils.MapCoordinate) (*pdpb.AbstractEdge, error) {
	e, err := em.Get(t1, t2)
	if err != nil {
		return nil, err
	}

	if e != nil {
		delete(em.edges[t1], t2)
		delete(em.edges[t2], t1)
	}

	return e, nil
}
