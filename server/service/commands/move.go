package move

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

const commandType = sscpb.CommandType_COMMAND_TYPE_MOVE

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

func Import(pb *apipb.MoveRequest) *MoveCommand {
	return &MoveCommand{
		clientID: pb.GetClientID(),
		tickID:   pb.GetTickID(),
	}
}

type MoveCommand struct {
	clientID string
	tickID   string
}

func (c *MoveCommand) Type() sscpb.CommandType {
	return commandType
}

func (c *MoveCommand) ClientID() string {
	return c.clientID
}

func (c *MoveCommand) TickID() string {
	return c.tickID
}

func (c *MoveCommand) Execute() error {
	return notImplemented
}
