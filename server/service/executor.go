// Package executor contains the logic for the core game loop.
package executor

import (
	"log"
	"sync"
	"time"

	"github.com/downflux/game/entity/entity"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/service/clientlist"
	"github.com/downflux/game/server/service/command/command"
	"github.com/downflux/game/server/service/command/move"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
	serverstatus "github.com/downflux/game/server/service/status"
)

const (
	idLen        = 8
	tickDuration = 100 * time.Millisecond
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

// TODO(minkezhang): Add ClientID string type.

// dirtyCurve represents a Curve instance which was altered in the current
// tick and will need to be broadcast to all clients.
//
// The Entity UUID and CurveCategory uniquely identifies a curve.
type dirtyCurve struct {
	// eid is the parent Entity UUID.
	eid string

	// category is the Entity property for which this Curve instance
	// represents.
	category gcpb.CurveCategory
}

// New creates a new instance of the Executor.
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
		entities:      map[string]entity.Entity{},
		commandQueue:  nil,
		clients:       clientlist.New(idLen),
		statusImpl:    serverstatus.New(tickDuration),
	}, nil
}

// Executor encapsulates logic for executing the core game loop.
type Executor struct {
	// tileMap is the underlying Map object used for the game.
	tileMap *tile.Map

	// abstractGraph is the underlying abstracted pathing logic data layer
	// for the associated Map.
	abstractGraph *graph.Graph

	// statusImpl represents the current Executor state metadata.
	statusImpl *serverstatus.Status

	// clients is an append-only set of connected players / AI.
	clients *clientlist.List

	// commandQueueMux guards the commandQueue property.
	commandQueueMux sync.Mutex

	// commandQueue is a FIFO list of commands to be run per tick. This
	// list is reset per tick.
	//
	// TODO(minkezhang): Refactor this into a CommandQueue object, that
	// hashes a command the tick it is scheduled at to run, and by a UUID
	// for cancelling the command.
	commandQueue []command.Command

	// dataMux guards the entities property.
	//
	// This lock must be acquired first.
	dataMux sync.RWMutex

	// entities is an append-only set of game entities.
	entities map[string]entity.Entity

	// entityQueueMux guards the entityQueue property.
	// This lock must be acquired second.
	entityQueueMux sync.RWMutex

	// entityQueue is an unordered list of new entites created during the
	// current game tick. This list is reset per tick.
	//
	// TODO(minkezhang): Determine if we can isolate this property and not
	// rely on the entities property.
	entityQueue []string

	// curveQueueMux protects the curveQueue property.
	//
	// This lock must be acquired third.
	curveQueueMux sync.RWMutex

	// curveQueue is an unordered list of curves that have been mutated
	// during the current game tick. This list is reset per tick.
	//
	// TODO(minkezhang): Determine if we can isolate this property and not
	// rely on the entities property.
	curveQueue []dirtyCurve
}

// Status returns the current Executor status.
func (e *Executor) Status() *gdpb.ServerStatus { return e.statusImpl.PB() }

// ClientExists tests for if the specified Client UUID is currently being
// tracked by the Executor.
func (e *Executor) ClientExists(cid string) bool { return e.clients.In(cid) }

// AddClient creates a new Client to be tracked by the Executor.
func (e *Executor) AddClient() (string, error) { return e.clients.Add() }

// StartClientStream instructs the Executor to mark the associated client
// ready for game state updates.
func (e *Executor) StartClientStream(cid string) error { return e.clients.Start(cid) }

// StopClientStreamError instructs the Executor to mark the associated client
// as having been disconnected, and stop broadcasting future game states to the
// linked channel.
func (e *Executor) StopClientStreamError(cid string) error { return e.clients.Stop(cid, false) }

// ClientChannel returns a read-only game state channel. This is consumed by
// the gRPC server and forwarded to the end-user.
func (e *Executor) ClientChannel(cid string) (<-chan *apipb.StreamDataResponse, error) {
	return e.clients.Channel(cid)
}

// popCommandQueue returns the list of commands for the current server tick and
// resets the internal list.
func (e *Executor) popCommandQueue() ([]command.Command, error) {
	e.commandQueueMux.Lock()
	defer e.commandQueueMux.Unlock()

	commands := e.commandQueue
	e.commandQueue = nil
	return commands, nil
}

// processCommand executes a single command with side-effects.
func (e *Executor) processCommand(cmd command.Command) error {
	if cmd.Type() == sscpb.CommandType_COMMAND_TYPE_MOVE {
		c, err := cmd.Execute(move.Args{
			Tick:   e.statusImpl.Tick(),
			Source: e.entities[cmd.(*move.Command).EntityID()].Curve(gcpb.CurveCategory_CURVE_CATEGORY_MOVE).Get(e.statusImpl.Tick()).(*gdpb.Position)})
		if err != nil {
			return err
		}

		if err := func() error {
			e.dataMux.RLock()
			defer e.dataMux.RUnlock()

			if err := e.entities[c.EntityID()].Curve(c.Category()).ReplaceTail(c); err != nil {
				return err
			}

			e.curveQueueMux.Lock()
			e.curveQueue = append(e.curveQueue, dirtyCurve{
				eid:      c.EntityID(),
				category: c.Category(),
			})
			e.curveQueueMux.Unlock()

			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}

// popTickQueue returns the list of curves and entity protos that will need to
// be broadcast to all valid cliens for the current server tick.
func (e *Executor) popTickQueue() ([]*gdpb.Curve, []*gdpb.Entity) {
	e.dataMux.RLock()
	e.entityQueueMux.Lock()
	e.curveQueueMux.Lock()
	defer e.curveQueueMux.Unlock()
	defer e.entityQueueMux.Unlock()
	defer e.dataMux.RUnlock()

	var processedCurves = map[string]map[gcpb.CurveCategory]bool{}

	var curves []*gdpb.Curve
	var entities []*gdpb.Entity

	// TODO(minkezhang): Make concurrent.
	for _, eid := range e.entityQueue {
		entities = append(entities, &gdpb.Entity{
			EntityId: eid,
			Type:     e.entities[eid].Type(),
		})
	}
	for _, dc := range e.curveQueue {
		// Do not broadcast curve twice.
		if _, found := processedCurves[dc.eid]; !found {
			processedCurves[dc.eid] = map[gcpb.CurveCategory]bool{}
		}
		if _, found := processedCurves[dc.eid][dc.category]; !found {
			processedCurves[dc.eid][dc.category] = true
			curves = append(curves, e.entities[dc.eid].Curve(dc.category).ExportTail(e.statusImpl.Tick()))
		}
	}

	e.curveQueue = nil
	e.entityQueue = nil

	return curves, entities
}

// allCurvesAndEntities returns a list of all Curve and Entity protos as of the
// current tick. This is used to broadcast the full game state to new or
// reconnecting clients.
func (e *Executor) allCurvesAndEntities() ([]*gdpb.Curve, []*gdpb.Entity) {
	e.dataMux.RLock()
	defer e.dataMux.RUnlock()

	var curves []*gdpb.Curve
	var entities []*gdpb.Entity

	// TODO(minkezhang): Give some leeway here, broadcast a bit in the
	// past.
	beginningTick := e.statusImpl.Tick()

	for _, en := range e.entities {
		entities = append(entities, &gdpb.Entity{
			EntityId: en.ID(),
			Type:     en.Type(),
		})
		for _, cat := range en.CurveCategories() {
			curves = append(curves, e.entities[en.ID()].Curve(cat).ExportTail(beginningTick))
		}
	}
	return curves, entities
}

// broadcastCurves will send the current game state delta or full game state to
// all connected clients. This is a blocking call.
func (e *Executor) broadcastCurves() error {
	log.Printf("[%.f]: broadcasting curves", e.statusImpl.Tick())

	curves, entities := e.popTickQueue()

	return e.clients.Broadcast(
		func() *apipb.StreamDataResponse {
			// TODO(minkezhang): Decide if it's okay that the reported tick may not
			// coincide with the ticks of the curve and entities.
			return &apipb.StreamDataResponse{
				Tick:     e.statusImpl.Tick(),
				Curves:   curves,
				Entities: entities,
			}
		},
		func() *apipb.StreamDataResponse {
			full := &apipb.StreamDataResponse{
				Tick: e.statusImpl.Tick(),
			}
			allCurves, allEntities := e.allCurvesAndEntities()
			full.Curves = allCurves
			full.Entities = allEntities
			return full
		},
	)
}

// Stop will teardown the Executor and close all client channels. This is
// called at the end of the game.
func (e *Executor) Stop() {
	e.statusImpl.SetIsStopped()
	e.clients.StopAll()
}

// Run executes the core game loop.
func (e *Executor) Run() error {
	e.statusImpl.SetStartTime()
	e.statusImpl.SetIsStarted()
	for !e.statusImpl.IsStopped() {
		t := time.Now()
		e.statusImpl.IncrementTick()

		if err := e.doTick(); err != nil {
			return err
		}

		// TODO(minkezhang): Add metrics collection here for tick
		// distribution.
		if d := time.Now().Sub(t); d < tickDuration {
			time.Sleep(tickDuration - d)
		}
	}
	return nil
}

// doTick executes a single iteration of the core game loop.
func (e *Executor) doTick() error {
	commands, err := e.popCommandQueue()
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		// TODO(minkezhang): Add actual error handling here -- only
		// Only return early if error is very bad.
		if err := e.processCommand(cmd); err != nil {
			return err
		}
	}

	if err := e.broadcastCurves(); err != nil {
		return err
	}

	return nil
}

// AddEntity creates a new entity.
//
// TODO(minkezhang): Make this method private -- this is currently public for
// debugging purposes.
//
// TODO(minkezhang): Make this schedule a generated Command instead for the
// next tick.
func (e *Executor) AddEntity(en entity.Entity) error {
	e.dataMux.Lock()
	e.entityQueueMux.Lock()
	e.curveQueueMux.Lock()
	defer e.curveQueueMux.Unlock()
	defer e.entityQueueMux.Unlock()
	defer e.dataMux.Unlock()

	if _, found := e.entities[en.ID()]; found {
		return status.Errorf(codes.AlreadyExists, "given entity ID %v already exists in the entity list", en.ID())
	}

	e.entities[en.ID()] = en

	e.entityQueue = append(e.entityQueue, en.ID())
	for _, cat := range en.CurveCategories() {
		e.curveQueue = append(e.curveQueue, dirtyCurve{
			eid:      en.ID(),
			category: cat,
		})
	}

	return nil
}

// addCommands extends the current tick-specific command queue with the input.
func (e *Executor) addCommands(cs []command.Command) error {
	e.commandQueueMux.Lock()
	defer e.commandQueueMux.Unlock()

	e.commandQueue = append(e.commandQueue, cs...)

	// TODO(minkezhang): Add client validation as per design doc.
	return nil
}

// buildMoveCommands constructs a list of command.Command instances.
//
// TODO(minkezhang): Decide how / when / if we want to deal with click
// spamming (same eids, multiple move commands per tick).
func (e *Executor) buildMoveCommands(cid string, dest *gdpb.Position, eids []string) []*move.Command {
	e.dataMux.RLock()
	defer e.dataMux.RUnlock()

	var res []*move.Command
	for _, eid := range eids {
		_, found := e.entities[eid]
		if found {
			res = append(
				res,
				move.New(
					e.tileMap,
					e.abstractGraph,
					cid,
					eid,
					dest))
		} else {
			log.Printf("entity ID %s not found in server entity lookup, could not build Move command", eid)
		}
	}
	return res
}

// AddMoveCommands transforms the player MoveRequest input into a list of
// Command instances, and schedules them to be executed in the next tick.
func (e *Executor) AddMoveCommands(req *apipb.MoveRequest) error {
	// TODO(minkezhang): If tick outside window, return error.
	var cs []command.Command

	for _, c := range e.buildMoveCommands(req.GetClientId(), req.GetDestination(), req.GetEntityIds()) {
		cs = append(cs, c)
	}

	return e.addCommands(cs)
}
