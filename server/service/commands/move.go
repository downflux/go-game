package move

import (
	"github.com/downflux/game/curve/curve"
	// "github.com/downflux/game/pathing/hpf/astar"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gdpb "github.com/downflux/game/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

const commandType = sscpb.CommandType_COMMAND_TYPE_MOVE

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

func New(m *tile.Map, g *graph.Graph, cid string, t float64, source *gdpb.Position, destination *gdpb.Position) *Command {
	return &Command{
		tileMap:       m,
		abstractGraph: g,
		clientID:      cid,
		tick:          t,
		source:        proto.Clone(source).(*gdpb.Position),
		destination:   proto.Clone(destination).(*gdpb.Position),
	}
}

type Command struct {
	tileMap       *tile.Map
	abstractGraph *graph.Graph
	clientID      string
	tick          float64
	source        *gdpb.Position
	destination   *gdpb.Position
}

func (c *Command) Type() sscpb.CommandType {
	return commandType
}

func (c *Command) ClientID() string {
	return c.clientID
}

func (c *Command) Tick() float64 {
	return c.tick
}

func (c *Command) Execute() ([]curve.Curve, error) {
	/*
		p, _, err := astar.Path(c.tm, c.g, utils.MC(nil), utils.MC(c.destination), 10)
		p =
	*/
	return nil, notImplemented
}
