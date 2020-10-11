package move

import (
	"github.com/downflux/game/pathing/hpf/graph"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	tile "github.com/downflux/game/map/map"
	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

const commandType = sscpb.CommandType_COMMAND_TYPE_MOVE

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

func New(pb *apipb.MoveRequest, m *tile.Map, g *graph.Graph) *Command {
	return &Command{
		tileMap: m,
		abstractGraph: g,
		clientID: pb.GetClientID(),
		tickID:   pb.GetTickID(),
	}
}

type Command struct {
	tileMap *tile.Map
	abstractGraph *graph.Graph
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
