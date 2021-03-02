package executorutils

import (
	"time"

	"github.com/downflux/game/engine/fsm/schedule"
	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/gamestate/gamestate"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/server/executor/executor"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/entity/component/attackable"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/entity/component/targetable"
	"github.com/downflux/game/server/visitor/attack"
	"github.com/downflux/game/server/visitor/chase"
	"github.com/downflux/game/server/visitor/move"
	"github.com/downflux/game/server/visitor/produce"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	entitylist "github.com/downflux/game/engine/entity/list"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
	serverstatus "github.com/downflux/game/engine/status/status"
	visitorlist "github.com/downflux/game/engine/visitor/list"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	attackaction "github.com/downflux/game/server/fsm/attack"
	chaseaction "github.com/downflux/game/server/fsm/chase"
	moveaction "github.com/downflux/game/server/fsm/move"
	produceaction "github.com/downflux/game/server/fsm/produce"
)

type Utils struct {
	executor *executor.Executor

	// gamestate links to the internal Executor.gamestate property. This is
	// for linking to FSM action constructors and must not be mutated here.
	// Making read-only calls is okay.
	gamestate *gamestate.GameState
}

func New(pb *mdpb.TileMap, d *gdpb.Coordinate, tickDuration time.Duration, minPathLength int) (*Utils, error) {
	tm, err := tile.ImportMap(pb)
	if err != nil {
		return nil, err
	}
	g, err := graph.BuildGraph(tm, d)
	if err != nil {
		return nil, err
	}

	fsmSchedule := schedule.New([]fcpb.FSMType{
		fcpb.FSMType_FSM_TYPE_CHASE,
		fcpb.FSMType_FSM_TYPE_MOVE,
		fcpb.FSMType_FSM_TYPE_PRODUCE,
		fcpb.FSMType_FSM_TYPE_ATTACK,
	})

	state := gamestate.New(serverstatus.New(tickDuration), entitylist.New())
	dirtystate := dirty.New()
	visitors, err := visitorlist.New([]visitor.Visitor{
		chase.New(state.Status(), fsmSchedule),
		produce.New(state.Status(), state.Entities(), dirtystate),
		move.New(tm, g, state.Status(), dirtystate, minPathLength),
		attack.New(state.Status(), dirtystate),
	})
	if err != nil {
		return nil, err
	}

	return &Utils{
		executor:  executor.New(visitors, state, dirtystate, fsmSchedule),
		gamestate: state,
	}, nil
}

func (u *Utils) Executor() *executor.Executor        { return u.executor }
func (u *Utils) Status() serverstatus.ReadOnlyStatus { return u.gamestate.Status() }

// Move transforms the player MoveRequest input into a list of move actions
// and schedules them to be executed in the next tick.
func (u *Utils) Move(pb *apipb.MoveRequest) error {
	// TODO(minkezhang): If tick outside window, return error.

	// TODO(minkezhang): Return list of errors instead.
	for _, eid := range pb.GetEntityIds() {
		m, ok := u.gamestate.Entities().Get(id.EntityID(eid)).(moveable.Component)
		if !ok {
			return status.Error(codes.FailedPrecondition, "specified entity is not moveable")
		}

		if err := u.executor.Schedule(
			moveaction.New(m, u.Status(), pb.GetDestination())); err != nil {
			return err
		}
	}

	return nil
}

func (u *Utils) Attack(pb *apipb.AttackRequest) error {
	t, ok := u.gamestate.Entities().Get(id.EntityID(pb.GetTargetEntityId())).(targetable.Component)
	if !ok {
		return status.Error(codes.FailedPrecondition, "specified entity is not targetable")
	}

	// TODO(minkezhang): Return list of errors instead.
	for _, eid := range pb.GetEntityIds() {
		e := u.gamestate.Entities().Get(id.EntityID(eid))
		a, ok := e.(attackable.Component)
		if !ok {
			return status.Error(codes.FailedPrecondition, "specified entity is not attackable")
		}
		m, ok := e.(moveable.Component)
		if !ok {
			return status.Error(codes.FailedPrecondition, "specified entity is not moveable")
		}

		chaseAction := chaseaction.New(u.Status(), m, t)
		attackAction := attackaction.New(u.Status(), a, t, chaseAction)

		if err := u.executor.Schedule(chaseAction); err != nil {
			return err
		}
		if err := u.executor.Schedule(attackAction); err != nil {
			return err
		}
	}
	return nil
}

// ProduceDebug schedules adding a new entity in the next game tick.
func (u *Utils) ProduceDebug(entityType gcpb.EntityType, spawnPosition *gdpb.Position) error {
	// TODO(minkezhang): Use arbitrary client-id after implementing
	// per-instance ACLs and setting to PublicWritable here.
	return u.executor.Schedule(
		produceaction.New(
			u.Status(),
			u.Status().Tick(),
			entityType,
			spawnPosition,
			id.ClientID("")))
}
