package abstractedge

import (
	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	"github.com/minkezhang/rts-pathing/lib/hpf/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Map contains a collection of AbstractEdge instances, which
// represent an AbstractGraph edge; these edges represent the cost to move
// between different AbstractNode instances.
type Map struct {
	// We want to explicitly disallow direct access.
	edges map[utils.MapCoordinate]map[utils.MapCoordinate]*rtsspb.AbstractEdge
}

func (em *Map) Iterator() []*rtsspb.AbstractEdge {
	var edges []*rtsspb.AbstractEdge
	for _, row := range em.edges {
		for _, e := range row {
			edges = append(edges, e)
		}
	}
	return edges
}

func (em *Map) GetBySource(t utils.MapCoordinate) ([]*rtsspb.AbstractEdge, error) {
	if em.edges == nil {
		return nil, nil
	}

	var res []*rtsspb.AbstractEdge
	for _, e := range em.edges[t] {
		res = append(res, e)
	}
	return res, nil
}

// Add appends an AbstractEdge instance into the Map collection.
//
// We're assuming the graph is undirected -- that is, for nodes A, B, if
// A --> B, then B --> A with the same cost.
func (em *Map) Add(e *rtsspb.AbstractEdge) error {
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
		em.edges = map[utils.MapCoordinate]map[utils.MapCoordinate]*rtsspb.AbstractEdge{}
	}
	if _, found := em.edges[t1]; !found {
		em.edges[t1] = map[utils.MapCoordinate]*rtsspb.AbstractEdge{}
	}
	if _, found := em.edges[t2]; !found {
		em.edges[t2] = map[utils.MapCoordinate]*rtsspb.AbstractEdge{}
	}

	em.edges[t1][t2] = e
	em.edges[t2][t1] = e

	return nil
}

// Get queries the Map for an AbstractEdge instance which connects
// two TileMap Coordinate instances.
func (em *Map) Get(t1, t2 utils.MapCoordinate) (*rtsspb.AbstractEdge, error) {
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
func (em *Map) Pop(t1, t2 utils.MapCoordinate) (*rtsspb.AbstractEdge, error) {
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
