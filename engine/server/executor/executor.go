// Package executor contains the logic for the core game loop.
package executor

import (
	"log"
	"time"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/schedule"
	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/gamestate/gamestate"
	"github.com/downflux/game/engine/id/id"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	clientlist "github.com/downflux/game/engine/server/client/list"
	visitorlist "github.com/downflux/game/engine/visitor/list"
)

const (
	// idLen represents the default length of a UUID (e.g. ClientID,
	// EntityID, etc.).
	idLen = 8
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
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

	gamestate *gamestate.GameState

	// dirty is a list of Entity and Curve instances which have been
	// modified during the current game tick. The Executor broadcasts this
	// list to all clients to update the game state.
	dirty *dirty.List

	// clients is an append-only set of connected players / AI.
	clients *clientlist.List

	schedule      *schedule.Schedule
	scheduleCache *schedule.Schedule
}

func New(
	visitors *visitorlist.List,
	state *gamestate.GameState,
	dcs *dirty.List,
	fsmSchedule *schedule.Schedule,
) *Executor {
	return &Executor{
		visitors:      visitors,
		gamestate:     state,
		dirty:         dcs,
		clients:       clientlist.New(idLen),
		schedule:      fsmSchedule,
		scheduleCache: fsmSchedule.Pop(),
	}
}

// Status returns the current Executor status.
func (e *Executor) Status() *gdpb.ServerStatus { return e.gamestate.Status().PB() }

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

// broadcast will send the current game state delta or full game state to
// all connected clients. This is a blocking call.
func (e *Executor) broadcast() error {
	partial := e.gamestate.Export(e.gamestate.Status().Tick()-100, e.dirty.Pop())

	return e.clients.Broadcast(
		// Return the game state update that will need to be broadcast
		// to all valid clients for the current server tick.
		func() *apipb.StreamDataResponse {
			// TODO(minkezhang): Decide if it's okay that the reported tick may not
			// coincide with the ticks of the curve and entities.
			return &apipb.StreamDataResponse{
				Tick:  e.gamestate.Status().Tick().Value(),
				State: partial,
			}
		},
		// Return a list of all Curve and Entity protos as of the
		// current tick. This is used to broadcast the full game state
		// to new or reconnecting clients.
		func() *apipb.StreamDataResponse {
			return &apipb.StreamDataResponse{
				Tick:  e.gamestate.Status().Tick().Value(),
				State: e.gamestate.Export(e.gamestate.Status().Tick(), e.gamestate.NoFilter()),
			}
		},
	)
}

// Stop will teardown the Executor and close all client channels. This is
// called at the end of the game.
func (e *Executor) Stop() error {
	if err := e.gamestate.Status().SetIsStopped(); err != nil {
		return err
	}
	e.clients.StopAll()
	return nil
}

// Run executes the core game loop.
func (e *Executor) Run() error {
	e.gamestate.Status().SetStartTime()
	if err := e.gamestate.Status().SetIsStarted(); err != nil {
		return err
	}
	for !e.gamestate.Status().IsStopped() {
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
	e.gamestate.Status().IncrementTick()

	e.schedule.Clear()
	if err := e.schedule.Merge(e.scheduleCache.Pop()); err != nil {
		return err
	}

	for _, v := range e.visitors.Iter() {
		if err := e.schedule.Get(v.Type()).Accept(v); err != nil {
			return err
		}
	}

	if err := e.broadcast(); err != nil {
		return err
	}

	// TODO(minkezhang): Add metrics collection here for tick
	// distribution.
	tickDuration := e.gamestate.Status().TickDuration()
	u := e.gamestate.Status().StartTime().Add(
		time.Duration(e.gamestate.Status().Tick()) * tickDuration).Sub(t)
	if u < tickDuration {
		time.Sleep(u)
	} else {
		log.Printf(
			"[%.f] took too long: execution time %v > %v",
			e.gamestate.Status().Tick(), u, tickDuration)
	}
	return nil
}

func (e *Executor) Schedule(actions []action.Action) error {
	return e.scheduleCache.Extend(actions)
}
