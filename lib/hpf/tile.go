// Package tile implements the TileInterface for the underlying map, i.e. for l = 0.
// Higher l-levels abstract the map away into logical nodes instead.
package tile

import (
	"math"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/beefsack/go-astar"
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
		t, err := ImportTile(pbt, tm)
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
//
// Tile implements astar.Pather.
type Tile struct {
	// M is a backreference to the originating map. We need this
	// to implement astar.Pather (i.e. if we expect to call A*
	// on this Tile object.
	//
	// M is not automatically exported when calling ExportTile.
	M *TileMap

	// Val is the underlying representation of the map node. It may be
	// mutated, e.g. changing TerrainType to / from TERRAIN_TYPE_BLOCKED
	Val *rtsspb.Tile
}

func (t *Tile) PathNeighbors() []astar.Pather {
	var ps []astar.Pather
	if t.M == nil {
		return ps
	}
	tiles, err := t.M.Neighbors(t.Val.GetCoordinate())
	if err != nil {
		return ps
	}
	for _, t := range tiles {
		ps = append(ps, t)
	}
	return ps
}

func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	if t.M == nil {
		return infinity
	}
	cost, err := D(t.M.C, t, to.(*Tile))
	if err != nil {
		return infinity
	}
	return cost
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	cost, err := H(t, to.(*Tile))
	if err != nil {
		return infinity
	}
	return cost
}

func Path(src, dest *Tile) ([]*Tile, bool, error) {
	if src == nil || dest == nil {
		return nil, false, status.Errorf(codes.FailedPrecondition, "cannot have nil Tile inputs")
	}
	if src.TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED || dest.TerrainType() == rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED {
		return nil, false, nil
	}

	path, _, found := astar.Path(src, dest)

	// astar.Path returns the path starting from destination, which is not useful for a majority
	// of our cases. Reversing here for convenience.
	var tiles []*Tile
	for i := range path {
		tiles = append(tiles, path[len(path)-1-i].(*Tile))
	}

	return tiles, found, nil
}

func ImportTile(pb *rtsspb.Tile, m *TileMap) (*Tile, error) {
	return &Tile{
		M:   m,
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
