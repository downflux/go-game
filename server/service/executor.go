// Package executor contains the logic for the core game loop.
package executor

import (
	"log"
	"time"

	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/fsm/schedule"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/clientlist"
	"github.com/downflux/game/server/visitor/dirty"
	"github.com/downflux/game/server/visitor/move"
	"github.com/downflux/game/server/visitor/produce"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	entitylist "github.com/downflux/game/engine/entity/list"
	vcpb "github.com/downflux/game/engine/visitor/api/constants_go_proto"
	visitorlist "github.com/downflux/game/engine/visitor/list"
	fcpb "github.com/downflux/game/fsm/api/constants_go_proto"
	moveinstance "github.com/downflux/game/fsm/move"
	produceinstance "github.com/downflux/game/fsm/produce"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	serverstatus "github.com/downflux/game/server/service/status"
)

const (
	// idLen represents the default length of a UUID (e.g. ClientID,
	// EntityID, etc.).
	idLen = 8

	// tickDuration is the targeted loop iteration time delta. If a tick
	// loop exceeds this time, it should delay commands until the next
	// cycle and ensure the dirty curves are being broadcasted instead.
	//
	// TODO(minkezhang): Ensure tick timeout actually occurs.
	tickDuration = 100 * time.Millisecond

	// entityListID is a preset ID for the global EntityList Entity
	// instance.
	entityListID = "entity-list"

	// minPathLength represents the minimum lookahead path length to
	// calculate, where the path is a list of tile.Map coordinates.
	minPathLength = 8
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")

	fsmVisitorTypeLookup = map[vcpb.VisitorType]fcpb.FSMType{
		vcpb.VisitorType_VISITOR_TYPE_MOVE:    fcpb.FSMType_FSM_TYPE_MOVE,
		vcpb.VisitorType_VISITOR_TYPE_PRODUCE: fcpb.FSMType_FSM_TYPE_PRODUCE,
	}
)

// Executor encapsulates logic for executing the core game loop.
type Executor struct {
	// visitors is a list of all Visitor instances used by the Executor.
	// A Visitor instance takes as state input an arbitrary subset of the
	// game state and mutates some or all entities every tick.
	//
	// The Executor uses the Visitor pattern for the central tick loop.
	// See https://en.wikipedia.org/wiki/Visitor_pattern.
	visitors *visitorlist.List

	// entities is a list of all Entity instances for the current game.
	// An Entity is an arbitrary stateful object -- it may not be a
	// physical game object like a tank; the entitylist.List object
	// itself is implements the Entity interface.
	//
	// Entity object states are mutated by Visitor instances.
	entities *entitylist.List

	// dirties is a list of Entity and Curve instances which have been
	// modified during the current game tick. The Executor broadcasts this
	// list to all clients to update the game state.
	dirties *dirty.List

	// statusImpl represents the current Executor state metadata.
	statusImpl *serverstatus.Status

	// clients is an append-only set of connected players / AI.
	clients *clientlist.List

	sot     *schedule.Schedule
	cache   *schedule.Schedule
	move    *move.Visitor
	produce *produce.Visitor
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

	dirties := dirty.New()
	statusImpl := serverstatus.New(tickDuration)

	entities := entitylist.New(entityListID)

	visitors, err := visitorlist.New([]visitor.Visitor{
		produce.New(statusImpl, entities, dirties),
		move.New(tm, g, statusImpl, dirties, minPathLength),
	})
	if err != nil {
		return nil, err
	}

	return &Executor{
		visitors:   visitors,
		entities:   entities,
		dirties:    dirties,
		clients:    clientlist.New(idLen),
		statusImpl: statusImpl,
		sot:        schedule.New(schedule.FSMTypes),
		cache:      schedule.New(schedule.FSMTypes),
	}, nil
}

// Status returns the current Executor status.
func (e *Executor) Status() *gdpb.ServerStatus { return e.statusImpl.PB() }

// ClientExists tests for if the specified Client UUID is currently being
// tracked by the Executor.
func (e *Executor) ClientExists(cid id.ClientID) bool { return e.clients.In(cid) }

// AddClient creates a new Client to be tracked by the Executor.
func (e *Executor) AddClient() (id.ClientID, error) { return e.clients.Add() }

// StartClientStream instructs the Executor to mark the associated client
// ready for game state updates.
func (e *Executor) StartClientStream(cid id.ClientID) error { return e.clients.Start(cid) }

// StopClientStreamError instructs the Executor to mark the associated client
// as having been disconnected, and stop broadcasting future game states to the
// linked channel.
func (e *Executor) StopClientStreamError(cid id.ClientID) error { return e.clients.Stop(cid, false) }

// ClientChannel returns a read-only game state channel. This is consumed by
// the gRPC server and forwarded to the end-user.
func (e *Executor) ClientChannel(cid id.ClientID) (<-chan *apipb.StreamDataResponse, error) {
	return e.clients.Channel(cid)
}

// popTickQueue returns the list of curves and entity protos that will need to
// be broadcast to all valid cliens for the current server tick.
func (e *Executor) popTickQueue() ([]*gdpb.Curve, []*gdpb.Entity) {
	var curves []*gdpb.Curve
	var entities []*gdpb.Entity

	tailTick := e.statusImpl.Tick() - 100
	if tailTick < 0 {
		tailTick = 0
	}

	// TODO(minkezhang): Make concurrent.
	for _, de := range e.dirties.PopEntities() {
		entities = append(entities, &gdpb.Entity{
			EntityId: de.ID.Value(),
			Type:     e.entities.Get(de.ID).Type(),
		})
	}
	for _, dc := range e.dirties.Pop() {
		curves = append(
			curves,
			e.entities.Get(
				dc.EntityID).Curve(dc.Property).ExportTail(tailTick))
	}

	return curves, entities
}

// allCurvesAndEntities returns a list of all Curve and Entity protos as of the
// current tick. This is used to broadcast the full game state to new or
// reconnecting clients.
func (e *Executor) allCurvesAndEntities() ([]*gdpb.Curve, []*gdpb.Entity) {
	var curves []*gdpb.Curve
	var entities []*gdpb.Entity

	// TODO(minkezhang): Give some leeway here, broadcast a bit in the
	// past.
	beginningTick := e.statusImpl.Tick()

	for _, en := range e.entities.Iter() {
		entities = append(entities, &gdpb.Entity{
			EntityId: en.ID().Value(),
			Type:     en.Type(),
		})
		for _, p := range en.Properties() {
			curves = append(curves, en.Curve(p).ExportTail(beginningTick))
		}
	}
	return curves, entities
}

// broadcastCurves will send the current game state delta or full game state to
// all connected clients. This is a blocking call.
func (e *Executor) broadcastCurves() error {
	curves, entities := e.popTickQueue()

	return e.clients.Broadcast(
		func() *apipb.StreamDataResponse {
			// TODO(minkezhang): Decide if it's okay that the reported tick may not
			// coincide with the ticks of the curve and entities.
			return &apipb.StreamDataResponse{
				Tick:     e.statusImpl.Tick().Value(),
				Curves:   curves,
				Entities: entities,
			}
		},
		func() *apipb.StreamDataResponse {
			full := &apipb.StreamDataResponse{
				Tick: e.statusImpl.Tick().Value(),
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
func (e *Executor) Stop() error {
	if err := e.statusImpl.SetIsStopped(); err != nil {
		return err
	}
	e.clients.StopAll()
	return nil
}

// Run executes the core game loop.
func (e *Executor) Run() error {
	e.statusImpl.SetStartTime()
	if err := e.statusImpl.SetIsStarted(); err != nil {
		return err
	}
	for !e.statusImpl.IsStopped() {
		if err := e.doTick(); err != nil {
			// TODO(minkezhang): Only return if error is fatal.
			return err
		}
	}
	return nil
}

// doTick executes a single iteration of the core game loop.
func (e *Executor) doTick() error {
	t := time.Now()
	e.statusImpl.IncrementTick()

	x := e.cache.Pop()

	if err := e.sot.Merge(x); err != nil {
		return err
	}

	// TODO(minkezhang): Clear CANCELED or FINISHED instances in a Visitor.
	e.sot.Clear()
	for _, v := range e.visitors.Iter() {
		if fsmType, found := fsmVisitorTypeLookup[v.Type()]; found {
			if err := e.sot.Get(fsmType).Accept(v); err != nil {
				return err
			}
		}
	}

	if err := e.broadcastCurves(); err != nil {
		return err
	}

	// TODO(minkezhang): Add metrics collection here for tick
	// distribution.
	u := e.statusImpl.StartTime().Add(
		time.Duration(e.statusImpl.Tick()) * tickDuration).Sub(t)
	if u < tickDuration {
		time.Sleep(u)
	} else {
		log.Printf(
			"[%.f] took too long: execution time %v > %v",
			e.statusImpl.Tick(), u, tickDuration)
	}
	return nil
}

// AddEntity schedules adding a new entity in the next game tick.
//
// TODO(minkezhang): Delete this method -- this is currently public for
// debugging purposes.
func (e *Executor) AddEntity(entityType gcpb.EntityType, p *gdpb.Position) error {
	return e.cache.Add(
		produceinstance.New(e.statusImpl, e.statusImpl.Tick(), entityType, p),
	)
}

// AddMoveCommands transforms the player MoveRequest input into a list of
// Command instances, and schedules them to be executed in the next tick.
func (e *Executor) AddMoveCommands(req *apipb.MoveRequest) error {
	// TODO(minkezhang): If tick outside window, return error.

	for _, eid := range req.GetEntityIds() {
		if err := e.cache.Add(
			moveinstance.New(e.entities.Get(id.EntityID(eid)), e.statusImpl, req.GetDestination()),
		); err != nil {
			return err
		}
	}

	return nil
}
