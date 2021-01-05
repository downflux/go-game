// Package utils contains some shared logic between different HPF packages.
package utils

import (
	"math"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

// MapCoordinate is a hashable convenience struct used as map keys.
type MapCoordinate struct {
	X, Y int32
}

// MC constructs a MapCoordinate struct from the associated protobuf.
func MC(c *gdpb.Coordinate) MapCoordinate {
	return MapCoordinate{
		X: c.GetX(),
		Y: c.GetY(),
	}
}

// PB constructs a protobuf from the associated MapCoordinate struct.
func PB(c MapCoordinate) *gdpb.Coordinate {
	return &gdpb.Coordinate{
		X: c.X,
		Y: c.Y,
	}
}

// AddMapCoordinate computes the pair-wise sum of two MapCoordinate instances.
func AddMapCoordinate(a, b MapCoordinate) MapCoordinate {
	return MapCoordinate{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

// LessThan applies the canonical lexicographical ordering to the input
// MapCoordinates.
func LessThan(a, b MapCoordinate) bool {
	return a.X < b.X || a.X == b.X && a.Y < b.Y
}

func Manhattan(a, b *gdpb.Position) float64 {
	return math.Abs(a.GetX() - b.GetX()) + math.Abs(a.GetY() - b.GetY())
}

func Euclidean(a, b *gdpb.Position) float64 {
	return math.Sqrt(math.Pow(a.GetX() - b.GetX(), 2) + math.Pow(a.GetY()-b.GetY(), 2))
}
