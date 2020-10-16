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
	gcpb "github.com/downflux/game/api/constants_go_proto"
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
	Tick() float64

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
		tickLookup:    map[string]float64{},
		entities:      map[string]entity.Entity{},
		commandQueue:  nil,
	}, nil
}

type Executor struct {
	tickMux    sync.RWMutex
	tick       float64
	tickLookup map[string]float64

	tileMap       *tile.Map
	abstractGraph *graph.Graph

	dataMux  sync.RWMutex
	entities map[string]entity.Entity

	commandQueueMux sync.RWMutex
	commandQueue    []Command
}

func Tick(e *Executor) error {
	// TODO(minkezhang): Increment tick counter.

	e.commandQueueMux.Lock()
	commands := e.commandQueue
	e.commandQueue = nil
	e.commandQueueMux.Unlock()

	for _, cmd := range commands {
		if cmd.Type() == sscpb.CommandType_COMMAND_TYPE_MOVE {
			_, err := cmd.Execute()
			// TODO(minkezhang): Only return early if error is very bad -- else, just log.
			if err != nil {
				return err
			}
		}
	}

	return notImplemented
}

func AddEntity(e *Executor, en entity.Entity) error {
	e.dataMux.Lock()
	defer e.dataMux.Unlock()

	if _, found := e.entities[en.ID()]; found {
		return status.Errorf(codes.AlreadyExists, "given entity ID %v already exists in the entity list", en.ID())
	}

	e.entities[en.ID()] = en
	return nil
}

func addCommands(e *Executor, cs []Command) error {
	e.commandQueueMux.Lock()
	defer e.commandQueueMux.Unlock()

	e.commandQueue = append(e.commandQueue, cs...)

	// TODO(minkezhang): Add client validation as per design doc.
	return nil
}

// buildMoveCommands
//
// Is expected to be called concurrently.
//
// TODO(minkezhang): Decide how / when / if we want to deal with click
// spamming (same eids, multiple move commands per tick).
func buildMoveCommands(e *Executor, cid string, t float64, dest *gdpb.Position, eids []string) []*move.Command {
	e.dataMux.RLock()
	defer e.dataMux.RUnlock()

	var res []*move.Command
	for _, eid := range eids {
		en, found := e.entities[eid]
		if found {
			p, err := en.Curve(gcpb.CurveCategory_CURVE_CATEGORY_MOVE).Get(t)
			if err == nil {
				res = append(res, move.New(e.tileMap, e.abstractGraph, cid, t, p.(*gdpb.Position), dest))
			}
		}
	}
	return res
}

// AddMoveCommands
//
// Is expected to be called concurrently.
func AddMoveCommands(e *Executor, req *apipb.MoveRequest) error {
	tick, err := func() (float64, error) {
		e.tickMux.RLock()
		defer e.tickMux.RLock()

		tick, found := e.tickLookup[req.GetTickId()]
		if !found {
			return 0, status.Errorf(codes.NotFound, "invalid tick ID %v", req.GetTickId())
		}
		return tick, nil
	}()
	if err != nil {
		return err
	}

	var cs []Command
	for _, c := range buildMoveCommands(e, req.GetClientId(), tick, req.GetDestination(), req.GetEntityIds()) {
		cs = append(cs, c)
	}
	return addCommands(e, cs)
}
