// Package executor contains the logic for the core game loop.
package executor

import (
	"log"
	"time"

	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/clientlist"
	"github.com/downflux/game/server/service/visitor/dirty"
	"github.com/downflux/game/server/service/visitor/entity/entitylist"
	"github.com/downflux/game/server/service/visitor/move"
	"github.com/downflux/game/server/service/visitor/produce"
	"github.com/downflux/game/server/service/visitor/visitor"
	"github.com/downflux/game/server/service/visitor/visitorlist"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	serverstatus "github.com/downflux/game/server/service/status"
	vcpb "github.com/downflux/game/server/service/visitor/api/constants_go_proto"
)

const (
	idLen        = 8
	tickDuration = 100 * time.Millisecond
	entityListID = "entity-list"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

// dirtyCurve represents a Curve instance which was altered in the current
// tick and will need to be broadcast to all clients.
//
// The Entity UUID and CurveCategory uniquely identifies a curve.
type dirtyCurve struct {
	// eid is the parent Entity UUID.
	eid id.EntityID

	// category is the Entity property for which this Curve instance
	// represents.
	category gcpb.CurveCategory
}

// Executor encapsulates logic for executing the core game loop.
type Executor struct {
	visitors *visitorlist.List
	entities *entitylist.List
	dirties  *dirty.List

	// statusImpl represents the current Executor state metadata.
	statusImpl *serverstatus.Status

	// clients is an append-only set of connected players / AI.
	clients *clientlist.List
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

	visitors, err := visitorlist.New(
		[]visitor.Visitor{
			produce.New(statusImpl, dirties),
			move.New(tm, g, statusImpl, dirties, 10),
		},
	)
	if err != nil {
		return nil, err
	}

	return &Executor{
		visitors:   visitors,
		entities:   entitylist.New(entityListID),
		dirties:    dirties,
		clients:    clientlist.New(idLen),
		statusImpl: statusImpl,
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
				dc.EntityID).Curve(dc.Category).ExportTail(e.statusImpl.Tick()))
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
		for _, cat := range en.CurveCategories() {
			curves = append(curves, en.Curve(cat).ExportTail(beginningTick))
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

	for _, v := range e.visitors.Iter() {
		if err := e.entities.Accept(v); err != nil {
			return err
		}
	}

	if err := e.broadcastCurves(); err != nil {
		return err
	}

	// TODO(minkezhang): Add metrics collection here for tick
	// distribution.
	if d := time.Now().Sub(t); d < tickDuration {
		time.Sleep(tickDuration - d)
	}
	return nil
}

// AddEntity schedules a new entity.
//
// TODO(minkezhang): Delete this method -- this is currently public for
// debugging purposes.
func (e *Executor) AddEntity(entityType gcpb.EntityType, p *gdpb.Position) error {
	return e.visitors.Get(vcpb.VisitorType_VISITOR_TYPE_PRODUCE).Schedule(
		produce.Args{
			ScheduledTick: e.statusImpl.Tick(),
			EntityType:    entityType,
			SpawnPosition: p,
		},
	)
}

// AddMoveCommands transforms the player MoveRequest input into a list of
// Command instances, and schedules them to be executed in the next tick.
func (e *Executor) AddMoveCommands(req *apipb.MoveRequest) error {
	// TODO(minkezhang): If tick outside window, return error.

	for _, eid := range req.GetEntityIds() {
		if err := e.visitors.Get(vcpb.VisitorType_VISITOR_TYPE_MOVE).Schedule(
			move.Args{
				Tick:        e.statusImpl.Tick(),
				EntityID:    id.EntityID(eid),
				Destination: req.GetDestination(),
			},
		); err != nil {
			return err
		}
	}

	return nil
}
