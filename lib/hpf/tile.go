// Package tile implements the TileInterface for the underlying map, i.e. for l = 0.
// Higher l-levels abstract the map away into logical nodes instead.
package tile

import (
	"math"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"
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
)

// IsAdjacent tests if two Tile objects are adjacent to one another.
func IsAdjacent(src, dst *Tile) bool {
	return math.Abs(float64(dst.X()-src.X()))+math.Abs(float64(dst.Y()-src.Y())) == 1
}

// D gets exact cost of two neighboring Tiles.
func D(c map[rtscpb.TerrainType]float64, src, dst *Tile) (float64, error) {
	if !IsAdjacent(src, dst) {
		return 0, status.Error(codes.InvalidArgument, "input tiles are not adjacent to one another")
	}
	return c[src.TerrainType()] + c[dst.TerrainType()], nil
}

// H gets the estimated cost of moving between two arbitrary Tiles.
func H(src, dst *Tile) (float64, error) {
	return math.Pow(float64(dst.X()-src.X()), 2) + math.Pow(float64(dst.Y()-src.Y()), 2), nil
}

// TileMap is a 2D hashmap of the terrain map file.
// Coordinates are expected to be accessed in (x, y) order.
type TileMap struct {
	d *rtsspb.Coordinate
	m map[int32]map[int32]*Tile
	c map[rtscpb.TerrainType]float64 // terrain cost
}

// NewTileMap constructs a new TileMap object from the input protobuf.
// List of TileMap.Tiles may be sparse.
func NewTileMap(pb *rtsspb.TileMap) (*TileMap, error) {
	return nil, notImplemented
}

// TileMapProto converts an internal TileMap object into an exportable
// protobuf. Certain tiles may be ignored to be reconstructed later.
func TileMapProto(m *TileMap) (*rtsspb.TileMap, error) {
	return nil, notImplemented
}

// Tile returns the Tile object from the input coordinates.
func (m TileMap) Tile(x, y int32) *Tile {
	if _, found := m.m[x]; !found {
		return nil
	}
	return m.m[x][y]
}

// Neighbors returns the adjacent Tiles of an input Tile object.
func (m TileMap) Neighbors(coordinate *rtsspb.Coordinate) ([]*Tile, error) {
	if m.Tile(coordinate.GetX(), coordinate.GetY()) == nil {
		return nil, status.Error(
			codes.NotFound, "tile not found in the map")
	}

	var neighbors []*Tile
	for _, c := range neighborCoordinates {
		if t := m.Tile(coordinate.GetX()+c.GetX(), coordinate.GetY()+c.GetY()); t != nil && t.TerrainType() != rtscpb.TerrainType_TERRAIN_TYPE_BLOCKED {
			neighbors = append(neighbors, t)
		}
	}
	return neighbors, nil
}

// Tile represents a physical map node.
type Tile struct {
	// t is the underlying representation of the map node. It may be
	// mutated, e.g. changing TerrainType to / from TERRAIN_TYPE_BLOCKED
	t *rtsspb.Tile
}

// X returns the X-coordinate of the Tile.
func (t *Tile) X() int32 {
	return t.t.GetCoordinate().GetX()
}

// X returns the Y-coordinate of the Tile.
func (t *Tile) Y() int32 {
	return t.t.GetCoordinate().GetY()
}

// TerrainType returns the TerrainType enum of the Tile.
func (t *Tile) TerrainType() rtscpb.TerrainType {
	return t.t.GetTerrainType()
}

// SetTerrainType sets the TerrainType enum of the Tile.
func (t *Tile) SetTerrainType(terrainType rtscpb.TerrainType) {
	t.t.TerrainType = terrainType
}

// Equal checks if two Tile objects are equal to one another.
// This is useful in tests; we should generally avoid this in actual code.
func (t *Tile) Equal(other *Tile) bool {
	return proto.Equal(t.t, other.t)
}
