package command

import (
	"github.com/downflux/game/curve/curve"

	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

type Command interface {
	Type() sscpb.CommandType
	ClientID() string

	// Execute will run subroutine for the command implementation at the
	// specified server tick.
	Execute(tick float64) (curve.Curve, error)
}
