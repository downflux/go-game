package move

import (
	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/curve/linearmove"
	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/pathing/hpf/astar"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/id"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gdpb "github.com/downflux/game/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

const (
	commandType  = sscpb.CommandType_COMMAND_TYPE_MOVE
	pathLength   = 0
	ticksPerTile = float64(10)
	idLen        = 8
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

func New(m *tile.Map, g *graph.Graph, cid string, eid string, t float64, source *gdpb.Position, destination *gdpb.Position) *Command {
	return &Command{
		tileMap:       m,
		abstractGraph: g,
		clientID:      cid,
		entityID:      eid,
		tick:          t,
		source:        proto.Clone(source).(*gdpb.Position),
		destination:   proto.Clone(destination).(*gdpb.Position),
	}
}

type Command struct {
	tileMap       *tile.Map
	abstractGraph *graph.Graph
	entityID      string
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

// Start time.
func (c *Command) Tick() float64 {
	return c.tick
}

// We're assuming the position values are sane and doesn't overflow int32.
func coordinate(p *gdpb.Position) *gdpb.Coordinate {
	return &gdpb.Coordinate{
		X: int32(p.GetX()),
		Y: int32(p.GetY()),
	}
}

func position(c *gdpb.Coordinate) *gdpb.Position {
	return &gdpb.Position{
		X: float64(c.GetX()),
		Y: float64(c.GetY()),
	}
}

func (c *Command) Execute() (curve.Curve, error) {
	// Called concurrently (across multiple commands).
	// TODO(minkezhang): proto.Clone the return values in map.astar.Path.
	// TODO(minkezhang): Add additional infrastructure necessary to set pathLength > 0.
	p, _, err := astar.Path(c.tileMap, c.abstractGraph, utils.MC(coordinate(c.source)), utils.MC(coordinate(c.destination)), pathLength)
	if err != nil {
		return nil, err
	}

	cv := linearmove.New(id.RandomString(idLen), c.entityID)
	for i, tile := range p {
		cv.Add(c.tick+float64(i)*ticksPerTile, position(tile.Val.GetCoordinate()))
	}

	return cv, nil
}
