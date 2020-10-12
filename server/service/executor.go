package executor

import (
	"sync"

	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/service/commands/move"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

type Command interface {
	Type() sscpb.CommandType
	ClientID() string
	TickID() string
	Execute() error
}

func New(pb *mdpb.TileMap, d *gdpb.Coordinate) (*Executor, error) {
	tm, err := tile.ImportMap(pb)
	if err != nil {
		return nil, err
	}
	g, err := graph.BuildGraph(tm, d)
	if err != nil {
		return nil, err
	}
	return &Executor{
		tileMap:       tm,
		abstractGraph: g,
	}, nil
}

type Executor struct {
	tileMap       *tile.Map
	abstractGraph *graph.Graph

	commandQueueMux sync.RWMutex
	commandQueue    []Command
}

func Tick(e *Executor) error {
	return notImplemented
}

func AddCommand(e *Executor, c Command) error {
	e.commandQueueMux.Lock()
	defer e.commandQueueMux.Lock()

	e.commandQueue = append(e.commandQueue, c)

	// TODO(minkezhang): Add client validation as per design doc.
	return notImplemented
}

func NewMoveCommand(e *Executor, req *apipb.MoveRequest) *move.Command {
	return move.New(req, e.tileMap, e.abstractGraph)
}
