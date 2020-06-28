// Package hpf implements the hiearchical pathfinding algorithm as described in
// Botea and Muller "Near Optimal Hierarchical Path-Finding", 2004.
package hpf

import (
// "math"
// astar "github.com/beefsack/go-astar"
//	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
)

var (
	offset = [][]int{
		{-1, -1},
		{-1, 1},
		{1, -1},
		{1, 1},
	}
)

type TileMap map[int]map[int]*Tile

func (m *TileMap) Neighbors(t *Tile) []*Tile {
	// (X, Y) coordinate candidates
	neighborCoordinateCandidates := [][]int{
		{-1, -1},
		{-1, 1},
		{1, -1},
		{1, 1},
	}
	var neighbors []hpf.TileInterface
	for _, c := range neighborCoordinateCandidates {
		if n := t.Map().Tile(t.x+c[0], t.y+c[1]).(*TerrainTile); n != nil && n.terrainType != rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

type Tile struct {
	x, y int
}

func (t *Tile) X() int {
	return t.x
}

func (t *Tile) Y() int {
	return t.y
}

/*
type ClusterList map[int]ClusterMap
type ClusterMap []Cluster
type Cluster interface {
	Adjacent(other Cluster) bool
}
func LookupCluster(e entrance, l int) Cluster
func BuildEntrances(c1, c2 Cluster) []Entrance
type Entrance interface {
}
type Graph interface {
	Build(l int) error
	AddNode(n *Node, l)
	AddEdge(e *Edge)
	RemoveNode(n *Node)
	RemoveEdge(e *Edge)
}
type Node struct {
	l int
	t *Tile
}
type Edge struct {
	src *Node
	dst *Node
	h float64
	d float64
	t rtscpb.EdgeType
}


func BuildTileMap(l int) TileInterface { return nil }

// Return entrance tiles which exist either in A or B, and let caller deduplicate.
func BuildEntrances(a, b TileInterface) []LGraphNode { return nil }

type TileMap interface {
	// map[int]map[int]TileInterface
	Parent() TileMap  // L + 1
	Child() TileMap  // L - 1
	Tile(x, y int) TileInterface
}

type TileInterface interface {
	Map() TileMap

	// May be dynamically calculated. L + 1.
	// Parent() TileInterface

	// May be dynamically calculated. L - 1.
	// Children() []TileInterface

	// Minimum (x, y) coordinates for the tile.
	X() int
	Y() int

	// May be dynamically calculated from global.
	// Tile is bounded between ([X, X + DimX), [Y, Y + DimY)).
	DimX() int
	DimY() int

	Neighbors() []TileInterface
	Entrances() []TileInterface
}

func BuildLGraph(m TileMap, l int) LGraph { return nil }

type LGraph interface {
	Map(l int) TileMap
	Nodes() []LGraphNode
	Edges() map[LGraphNode]LGraphEdge

	AddNode(n LGraphNode) error
	AddEdge(n LGraphEdge) error  // Directed graph.
	RemoveNode(n LGraphNode) error
	RemoveEdge(n LGraphEdge) error

	Path(source, destination LGraphNode) []LGraphNode
}

type LGraphEdge interface {
	Source() TileInterface
	Destination() TileInterface

	Type() rtscpb.EdgeType
	D() (float64, error)
	H() (float64, error)
}

type LGraphNode interface {
	L() int
	Tile() TileInterface
	Neighbors() []LGraphNode
	D(to LGraphNode) (float64, error)
	H(to LGraphNode) (float64, error)
}
*/
