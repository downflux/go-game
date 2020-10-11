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

func Import(pb *apipb.MoveRequest) *Command {
	return &Command{
		clientID: pb.GetClientID(),
		tickID:   pb.GetTickID(),
	}
}

type Command struct {
	clientID string
	tickID   string
}

func (c *Command) Type() sscpb.CommandType {
	return commandType
}

func (c *Command) ClientID() string {
	return c.clientID
}

func (c *Command) TickID() string {
	return c.tickID
}

func (c *Command) Execute() error {
	return notImplemented
}
