// Package tile implements the TileInterface for the underlying map, i.e. for l = 0.
// Higher l-levels abstract the map away into logical nodes instead.
package tile

import (
	"math"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/fzipp/astar"
	"github.com/cripplet/rts-pathing/lib/hpf/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
	neighborCoordinates = []*rtsspb.Coordinate{
		{X: 0, Y: 1},
		{X: 0, Y: -1},
		{X: 1, Y: 0},
		{X: -1, Y: 0},
	}
	infinity = math.Inf(1)
)

// IsAdjacent tests if two Tile objects are adjacent to one another.
func IsAdjacent(src, dst *Tile) bool {
	return math.Abs(float64(dst.X()-src.X()))+math.Abs(float64(dst.Y()-src.Y())) == 1
}

// D gets exact cost between two neighboring Tiles.
//
// We're only taking the "difficulty" metric of the target Tile here; moving into and out of
// a Tile will essentially just double its cost, and the cost doesn't matter when
// the Tile is a source or target (since moving there is mandatory).
func D(c map[rtscpb.TerrainType]float64, src, dst *Tile) (float64, error) {
	if !IsAdjacent(src, dst) {
		return 0, status.Error(codes.InvalidArgument, "input tiles are not adjacent to one another")
	}
	return c[dst.TerrainType()], nil
}

// H gets the estimated cost of moving between two arbitrary Tiles.
func H(src, dst *Tile) (float64, error) {
	return math.Pow(float64(dst.X()-src.X()), 2) + math.Pow(float64(dst.Y()-src.Y()), 2), nil
}

// TileMap is a 2D hashmap of the terrain map file.
// Coordinates are expected to be accessed in (x, y) order.
//
// TileMap implements astar.Graph.
type TileMap struct {
	D *rtsspb.Coordinate
	M map[utils.MapCoordinate]*Tile
	C map[rtscpb.TerrainType]float64 // terrain cost
}

// ImportTileMap constructs a new TileMap object from the input protobuf.
// List of TileMap.Tiles may be sparse.
func ImportTileMap(pb *rtsspb.TileMap) (*TileMap, error) {
	m := make(map[utils.MapCoordinate]*Tile)
	tc := make(map[rtscpb.TerrainType]float64)

	for _, c := range pb.GetTerrainCosts() {
		tc[c.GetTerrainType()] = c.GetCost()
	}
	tm := &TileMap{
		D: pb.GetDimension(),
		M: m,
		C: tc,
	}

	for _, pbt := range pb.GetTiles() {
		t, err := ImportTile(pbt)
		if err != nil {
			return nil, err
		}
		m[utils.MC(t.Val.GetCoordinate())] = t
	}

	return tm, nil
}

// ExportTileMap converts an internal TileMap object into an exportable
// protobuf. Certain tiles may be ignored to be reconstructed later.
func ExportTileMap(m *TileMap) (*rtsspb.TileMap, error) {
	return nil, notImplemented
}

// Tile returns the Tile object from the input coordinates.
func (m *TileMap) Tile(x, y int32) *Tile {
	return m.M[utils.MapCoordinate{X: x, Y: y}]
}

func (m *TileMap) TileFromCoordinate(c *rtsspb.Coordinate) *Tile {
	return m.Tile(c.GetX(), c.GetY())
}

// Returns a list of adjacent Tiles.
//
// Implements astar.Graph.Neighbours.
func (m TileMap) Neighbours(n astar.Node) []astar.Node {
	ts, _ := m.Neighbors(n.(*Tile).Val.GetCoordinate())

	var res []astar.Node
	for _, t := range ts {
		res = append(res, t)
	}
	return res
}

// Neighbors returns the adjacent Tiles of an input Tile object.
func (m TileMap) Neighbors(coordinate *rtsspb.Coordinate) ([]*Tile, error) {
	if m.TileFromCoordinate(coordinate) == nil {
		return nil, status.Error(
			codes.NotFound, "tile not found in the map")

	}

	src := m.TileFromCoordinate(coordinate)
	var neighbors []*Tile
	for _, c := range neighborCoordinates {
		if t := m.Tile(coordinate.GetX()+c.GetX(), coordinate.GetY()+c.GetY()); src.TerrainType() != rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED && t != nil && t.TerrainType() != rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED {
			neighbors = append(neighbors, t)
		}
	}
	return neighbors, nil
}

// Tile represents a physical map node.
type Tile struct {
	// Val is the underlying representation of the map node. It may be
	// mutated, e.g. changing TerrainType to / from TERRAIN_TYPE_BLOCKED
	Val *rtsspb.Tile
}

func Path(m *TileMap, src, dest *Tile) ([]*Tile, error) {
	if src == nil || dest == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "cannot have nil Tile inputs")
	}
	if src.TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED || dest.TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED {
		return nil, nil
	}

	var res []*Tile
	path := astar.FindPath(m, src, dest, func(a, b astar.Node) float64{
		cost, _ := D(m.C, a.(*Tile), b.(*Tile))
		return cost
	}, func(a, b astar.Node) float64{
		cost, _ := H(a.(*Tile), b.(*Tile))
		return cost
	})
	for _, p := range path {
		res = append(res, p.(*Tile))
	}
	return res, nil
}

func ImportTile(pb *rtsspb.Tile) (*Tile, error) {
	return &Tile{
		Val: pb,
	}, nil
}

func ExportTile(t *Tile) (*rtsspb.Tile, error) {
	return nil, notImplemented
}

func (t *Tile) Coordinate() *rtsspb.Coordinate {
	return t.Val.GetCoordinate()
}

// X returns the X-coordinate of the Tile.
func (t *Tile) X() int32 {
	return t.Val.GetCoordinate().GetX()
}

// X returns the Y-coordinate of the Tile.
func (t *Tile) Y() int32 {
	return t.Val.GetCoordinate().GetY()
}

// TerrainType returns the TerrainType enum of the Tile.
func (t *Tile) TerrainType() rtscpb.TerrainType {
	return t.Val.GetTerrainType()
}

// SetTerrainType sets the TerrainType enum of the Tile.
func (t *Tile) SetTerrainType(terrainType rtscpb.TerrainType) {
	t.Val.TerrainType = terrainType
}
