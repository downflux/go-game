// Package tile implements the TileInterface for the underlying map,
// i.e. for l = 0. Higher l-levels abstract the map away into logical nodes
// instead.
package tile

import (
	"math"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")

	// neighborCoordinates provides the Coordinate deltas between a
	// specific Coordinate and adjacent Coordinates to expand to in
	// a graph search.
	neighborCoordinates = []*rtsspb.Coordinate{
		{X: 0, Y: 1},
		{X: 0, Y: -1},
		{X: 1, Y: 0},
		{X: -1, Y: 0},
	}
)

// IsAdjacent tests if two Tile objects are adjacent to one another.
func IsAdjacent(src, dst *Tile) bool {
	return math.Abs(float64(dst.X()-src.X()))+math.Abs(float64(dst.Y()-src.Y())) == 1
}

// D gets exact cost between two neighboring Tiles.
//
// We're only taking the "difficulty" metric of the target Tile here; moving
// into and out of a Tile will essentially just double its cost, and the cost
// doesn't matter when the Tile is a source or target (since moving there is
// mandatory).
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
	// D is the dimension of the TileMap, i.e. the number of Tile objects
	// in each direction.
	D *rtsspb.Coordinate

	// M is the actual list of Tile objects in the map; the MapCoordinate
	// here is the Coordinate of the actual Tile.
	M map[utils.MapCoordinate]*Tile

	// C is an embedded lookup table of terrain costs.
	C map[rtscpb.TerrainType]float64
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

// TileFromCoordinate returns the Tile object from the input Coordinate
// protobuf.
func (m *TileMap) TileFromCoordinate(c *rtsspb.Coordinate) *Tile {
	return m.Tile(c.GetX(), c.GetY())
}

// Neighbors returns the adjacent Tiles of an input Tile object.
func (m *TileMap) Neighbors(coordinate *rtsspb.Coordinate) ([]*Tile, error) {
	src := m.TileFromCoordinate(coordinate)
	if src == nil {
		return nil, status.Error(
			codes.NotFound, "tile not found in the map")

	}

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

// ImportTile constructs the Tile object from the specified protobuf.
func ImportTile(pb *rtsspb.Tile) (*Tile, error) {
	return &Tile{
		Val: pb,
	}, nil
}

// ExportTile constrcts a protobuf based on the specified Tile object.
func ExportTile(t *Tile) (*rtsspb.Tile, error) {
	return nil, notImplemented
}

// Coordinate returns the Tile Coordinate.
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
