// Package utils contains some shared logic between different HPF packages.
package utils

import (
	rtsspb "github.com/downflux/pathing/lib/proto/structs_go_proto"
)

// MapCoordinate is a hashable convenience struct used as map keys.
type MapCoordinate struct {
	X, Y int32
}

// MC constructs a MapCoordinate struct from the associated protobuf.
func MC(c *rtsspb.Coordinate) MapCoordinate {
	return MapCoordinate{
		X: c.GetX(),
		Y: c.GetY(),
	}
}

// PB constructs a protobuf from the associated MapCoordinate struct.
func PB(c MapCoordinate) *rtsspb.Coordinate {
	return &rtsspb.Coordinate{
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
