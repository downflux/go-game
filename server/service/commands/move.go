package move

import (
	"github.com/downflux/game/curve/curve"
	// "github.com/downflux/game/pathing/hpf/astar"
	"github.com/downflux/game/pathing/hpf/graph"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
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
		tileMap:       m,
		abstractGraph: g,
		clientID:      pb.GetClientId(),
		entityIDs:     pb.GetEntityIds(),
		tickID:        pb.GetTickId(),
		destination:   pb.GetDestination(),
	}
}

type Command struct {
	tileMap       *tile.Map
	abstractGraph *graph.Graph
	entityIDs     []string
	clientID      string
	tickID        string
	destination   *gdpb.Coordinate
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

func (c *Command) Execute() ([]curve.Curve, error) {
	/*
	// TODO(minkezhang): Make entity singular, tie source at creation time.
	for _, e := c.entityIDs {
		p, _, err := astar.Path(c.tm, c.g, utils.MC(nil), utils.MC(c.destination), 10)
		p = 
	}
	 */
	return nil, notImplemented
}
