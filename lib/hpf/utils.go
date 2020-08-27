// Package utils contains some shared logic between different HPF packages.
package utils

import (
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"
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

func AddMapCoordinate(a, b MapCoordinate) MapCoordinate {
	return MapCoordinate{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}
