package command

import (
	"github.com/downflux/game/curve/curve"

	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

type Command interface {
	Type() sscpb.CommandType
	ClientID() string
	Tick() float64

	Execute() (curve.Curve, error)
}
