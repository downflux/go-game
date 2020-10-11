package executor

import (
	"sync"

	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/service/commands/move"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tile "github.com/downflux/game/map/map"
	apipb "github.com/downflux/game/api/api_go_proto"
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

func New() *Executor {
	return &Executor{}
}

type Executor struct {
	tileMap *tile.Map
	abstractGraph *graph.Graph

	commandQueueMux sync.RWMutex
	commandQueue    []Command
}

func AddCommand(e *Executor, c Command) error {
	e.commandQueueMux.Lock()
	defer e.commandQueueMux.Lock()

	e.commandQueue = append(e.commandQueue, c)

	// TODO(minkezhang): Add client validation as per design doc.
	return notImplemented
}

func Tick(e *Executor) error {
	return notImplemented
}

func NewMoveCommand(e *Executor, req *apipb.MoveRequest) *move.Command {
	return move.New(req, e.tileMap, e.abstractGraph)
}
