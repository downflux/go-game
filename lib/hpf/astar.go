package astar

import (
	"math"

	"github.com/minkezhang/rts-pathing/lib/hpf/cluster"
	"github.com/minkezhang/rts-pathing/lib/hpf/graph"
	"github.com/minkezhang/rts-pathing/lib/hpf/graphastar"
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

func Path(tm *tile.Map, g *graph.Graph, src, dest utils.MapCoordinate, l int) ([]*tile.Tile, float64, error) {
	srcID, err := graph.InsertEphemeralNode(tm, g, src)
	if err != nil {
		return nil, 0, err
	}
	defer graph.RemoveEphemeralNode(g, src, srcID)

	destID, err := graph.InsertEphemeralNode(tm, g, dest)
	if err != nil {
		return nil, 0, err
	}
	defer graph.RemoveEphemeralNode(g, dest, destID)

	nPath, cost, err := graphastar.Path(tm, g, src, dest)
	if err != nil {
		return nil, 0, err
	}
	if math.IsInf(cost, 0) {
		return nil, cost, nil
	}

	var path []*tile.Tile
	for i, n1 := range nPath[:len(nPath)-1] {
		n2 := nPath[i+1]

		t1 := utils.MC(n1.GetTileCoordinate())
		t2 := utils.MC(n2.GetTileCoordinate())

		c1, err := cluster.ClusterCoordinateFromTileCoordinate(g.NodeMap.ClusterMap, t1)
		if err != nil {
			return nil, 0, err
		}
		c2, err := cluster.ClusterCoordinateFromTileCoordinate(g.NodeMap.ClusterMap, t2)
		if err != nil {
			return nil, 0, err
		}

		var p []*tile.Tile
		if c1 == c2 {
			tileBoundary, err := cluster.TileBoundary(g.NodeMap.ClusterMap, c1)
			if err != nil {
				return nil, 0, err
			}
			tileDimension, err := cluster.TileDimension(g.NodeMap.ClusterMap, c1)
			if err != nil {
				return nil, 0, err
			}

			p, _, err = tileastar.Path(tm, t1, t2, utils.PB(tileBoundary), utils.PB(tileDimension))
			if err != nil {
				return nil, 0, err
			}
		} else {
			p, _, err = tileastar.Path(tm, t1, t2, utils.PB(t1), utils.PB(t2))
			if err != nil {
				return nil, 0, err
			}
		}

		for _, n := range p {
			if len(p) < l || l == 0 {
				path = append(path, n)
			}
		}
	}

	return path, cost, err
}
