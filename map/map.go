// Package tile implements the TileInterface for the underlying map.
package tile

import (
	"math"

	gdpb "github.com/downflux/game/api/data_go_proto"
	mcpb "github.com/downflux/game/map/api/constants_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"

	"github.com/downflux/game/map/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")

	// neighborCoordinates provides the Coordinate deltas between a
	// specific Coordinate and adjacent Coordinates to expand to in
	// a graph search.
	neighborCoordinates = []*gdpb.Coordinate{
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
func D(c map[mcpb.TerrainType]float64, src, dst *Tile) (float64, error) {
	if !IsAdjacent(src, dst) {
		return 0, status.Error(codes.InvalidArgument, "input tiles are not adjacent to one another")
	}
	return c[dst.TerrainType()], nil
}

// H gets the estimated cost of moving between two arbitrary Tiles.
func H(src, dst *Tile) (float64, error) {
	return math.Pow(float64(dst.X()-src.X()), 2) + math.Pow(float64(dst.Y()-src.Y()), 2), nil
}

// Map is a 2D hashmap of the terrain map file.
// Coordinates are expected to be accessed in (x, y) order.
//
// Map implements astar.Graph.
type Map struct {
	// D is the dimension of the Map, i.e. the number of Tile objects
	// in each direction.
	D *gdpb.Coordinate

	// M is the actual list of Tile objects in the map; the MapCoordinate
	// here is the Coordinate of the actual Tile.
	M map[utils.MapCoordinate]*Tile

	// C is an embedded lookup table of terrain costs.
	C map[mcpb.TerrainType]float64
}

// ImportMap constructs a new Map object from the input protobuf.
// List of Map.Tiles may be sparse.
func ImportMap(pb *mdpb.TileMap) (*Map, error) {
	m := make(map[utils.MapCoordinate]*Tile)
	tc := make(map[mcpb.TerrainType]float64)

	for _, c := range pb.GetTerrainCosts() {
		tc[c.GetTerrainType()] = c.GetCost()
	}
	tm := &Map{
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

	for x := int32(0); x < tm.D.GetX(); x++ {
		for y := int32(0); y < tm.D.GetY(); y++ {
			if tm.Tile(x, y) == nil {
				return nil, status.Errorf(codes.InvalidArgument, "Map does not fully specify all tiles within the given map dimensions")
			}
		}
	}

	return tm, nil
}

// ExportMap converts an internal Map object into an exportable
// protobuf. Certain tiles may be ignored to be reconstructed later.
func ExportMap(m *Map) (*mdpb.TileMap, error) {
	return nil, notImplemented
}

// Tile returns the Tile object from the input coordinates.
func (m *Map) Tile(x, y int32) *Tile {
	return m.M[utils.MapCoordinate{X: x, Y: y}]
}

// TileFromCoordinate returns the Tile object from the input Coordinate
// protobuf.
func (m *Map) TileFromCoordinate(c *gdpb.Coordinate) *Tile {
	return m.Tile(c.GetX(), c.GetY())
}

// Neighbors returns the adjacent Tiles of an input Tile object.
func (m *Map) Neighbors(coordinate *gdpb.Coordinate) ([]*Tile, error) {
	src := m.TileFromCoordinate(coordinate)
	if src == nil {
		return nil, status.Error(
			codes.NotFound, "tile not found in the map")

	}

	var neighbors []*Tile
	for _, c := range neighborCoordinates {
		if t := m.Tile(coordinate.GetX()+c.GetX(), coordinate.GetY()+c.GetY()); m.C[src.TerrainType()] < math.Inf(0) && t != nil && m.C[t.TerrainType()] < math.Inf(0) {
			neighbors = append(neighbors, t)
		}
	}
	return neighbors, nil
}

// Tile represents a physical map node.
type Tile struct {
	// Val is the underlying representation of the map node. It may be
	// mutated, e.g. changing TerrainType to / from TERRAIN_TYPE_BLOCKED
	Val *mdpb.Tile
}

// ImportTile constructs the Tile object from the specified protobuf.
func ImportTile(pb *mdpb.Tile) (*Tile, error) {
	return &Tile{
		Val: pb,
	}, nil
}

// ExportTile constrcts a protobuf based on the specified Tile object.
func ExportTile(t *Tile) (*mdpb.Tile, error) {
	return nil, notImplemented
}

// Coordinate returns the Tile Coordinate.
func (t *Tile) Coordinate() *gdpb.Coordinate {
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
func (t *Tile) TerrainType() mcpb.TerrainType {
	return t.Val.GetTerrainType()
}

// SetTerrainType sets the TerrainType enum of the Tile.
func (t *Tile) SetTerrainType(terrainType mcpb.TerrainType) {
	t.Val.TerrainType = terrainType
}
