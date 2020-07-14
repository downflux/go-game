package utils

import (
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"
)

type MapCoordinate struct {
	x, y int32
}

func MC(c *rtsspb.Coordinate) MapCoordinate {
	return MapCoordinate{
		x: c.GetX(),
		y: c.GetY(),
	}
}
