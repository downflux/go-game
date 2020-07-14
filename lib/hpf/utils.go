package utils

import (
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"
)

type MapCoordinate struct {
	X, Y int32
}

func MC(c *rtsspb.Coordinate) MapCoordinate {
	return MapCoordinate{
		X: c.GetX(),
		Y: c.GetY(),
	}
}
