package executor

import (
	"sync"

	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/entity/entity"
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
	// TODO(minkezhang): Refactor Curve interface to not be dependent on curve.Curve.
	Execute() ([]curve.Curve, error)
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
		entities: map[string]entity.Entity{},
		curves: map[string]curve.Curve{},
		commandQueue: nil,
	}, nil
}

type Executor struct {
	tileMap       *tile.Map
	abstractGraph *graph.Graph

	entitiesMux sync.RWMutex
	entities    map[string]entity.Entity

	curvesMux sync.RWMutex
	curves    map[string]curve.Curve

	commandQueueMux sync.RWMutex
	commandQueue    []Command
}

func Tick(e *Executor) error {
	e.commandQueueMux.Lock()
	defer e.commandQueueMux.Unlock()

	for _, cmd := range e.commandQueue {
		if cmd.Type() == sscpb.CommandType_COMMAND_TYPE_MOVE {
			curves, err := cmd.Execute()
			// TODO(minkezhang): Only return early if error is very bad -- else, just log.
			if err != nil {
				return err
			}
		}
	}

	return notImplemented
}

func AddEntity(e *Executor, en entity.Entity) error {
	e.entitiesMux.Lock()
	defer e.entitiesMux.Unlock()

	e.entities[en.ID()] = en
	return nil
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
