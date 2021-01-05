package move

import (
	"sync"
	"time"

	"github.com/downflux/game/engine/curve/common/linearmove"
	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/pathing/hpf/astar"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move"

	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
	vcpb "github.com/downflux/game/engine/visitor/api/constants_go_proto"
	tile "github.com/downflux/game/map/map"
)

const (
	// visitorType is the registered VisitorType of the move visitor.
	visitorType = vcpb.VisitorType_VISITOR_TYPE_MOVE
)

// coordinate transforms a gdpb.Position instance into a gdpb.Coordinate
// instance. We're assuming the position values are sane and don't overflow
// int32.
func coordinate(p *gdpb.Position) *gdpb.Coordinate {
	return &gdpb.Coordinate{
		X: int32(p.GetX()),
		Y: int32(p.GetY()),
	}
}

// position transforms a gdpb.Coordinate instance into a gdpb.Position
// instance.
func position(c *gdpb.Coordinate) *gdpb.Position {
	return &gdpb.Position{
		X: float64(c.GetX()),
		Y: float64(c.GetY()),
	}
}

type Visitor struct {
	// mux guarantees we're running only one tile.Map astar at a time.
	// TODO(minkezhang): Make this concurrent.
	mux sync.Mutex

	// tileMap is the underlying Map object used for the game.
	tileMap *tile.Map

	// abstractGraph is the underlying abstracted pathing logic data layer
	// for the associated Map.
	abstractGraph *graph.Graph

	// status is a shared object with the game engine and indicates
	// current tick, etc.
	status *status.Status

	// dirties is a shared object between the game engine and the
	// Visitor.
	dirties *dirty.List

	// minPathLength represents the minimum lookahead path length to
	// calculate, where the path is a list of tile.Map coordinates.
	//
	// Longer calculated paths is discouraged, as these paths become
	// outdated once a new move command is issued for the Entity, which
	// may happen frequently in an RTS game.
	minPathLength int
}

// New constructs a new move Visitor instance.
func New(
	tileMap *tile.Map,
	abstractGraph *graph.Graph,
	dfStatus *status.Status,
	dirties *dirty.List,
	minPathLength int) *Visitor {
	return &Visitor{
		tileMap:       tileMap,
		abstractGraph: abstractGraph,
		status:        dfStatus,
		dirties:       dirties,
		minPathLength: minPathLength,
	}
}

// Type returns the registered VisitorType.
func (v *Visitor) Type() vcpb.VisitorType { return visitorType }

// TODO(minkezhang): Remove function.
func (v *Visitor) Schedule(interface{}) error { return nil }

func (v *Visitor) visitFSM(i action.Action) error {
	if i.Type() != fcpb.FSMType_FSM_TYPE_MOVE {
		return nil
	}

	s, err := i.State()
	if err != nil {
		return err
	}

	m := i.(*move.Action)

	tick := v.status.Tick()

	switch s {
	case commonstate.Executing:
		e := m.Component()
		c := e.PositionCurve()
		if c == nil {
			return nil
		}

		ticksPerSecond := float64(time.Second / v.status.TickDuration())
		ticksPerTile := id.Tick(ticksPerSecond / e.Velocity())

		p, _, err := astar.Path(
			v.tileMap,
			v.abstractGraph,
			utils.MC(coordinate(e.Position(tick))),
			utils.MC(coordinate(m.Destination())),
			v.minPathLength,
		)
		if err != nil {
			// TODO(minkezhang): Handle error by logging and continuing.
			return err
		}

		// Add to the existing curve, while smoothing out the existing
		// trajectory.
		prevPos := e.Position(tick)
		cv := linearmove.New(e.ID(), tick)
		cv.Add(tick, prevPos)
		for i, tile := range p {
			curPos := position(tile.Val.GetCoordinate())
			tickDelta := ticksPerTile * id.Tick(utils.Euclidean(prevPos, curPos))
			cv.Add(tick+id.Tick(i)*ticksPerTile+tickDelta, curPos)
			prevPos = curPos
		}
		if err := v.dirties.AddCurve(dirty.Curve{
			EntityID: e.ID(),
			Property: c.Property(),
		}); err != nil {
			return err
		}
		if err := c.Merge(cv); err != nil {
			return err
		}

		// Delay next lookup iteration until a suitable time in the
		// future.
		//
		// TODO(minkezhang): Add test for scheduling here.
		if err := m.SchedulePartialMove(tick + ticksPerTile*id.Tick(len(p))); err != nil {
			// TODO(minkezhang): Handle error by logging and continuing.
			return err
		}
	default:
		return nil
	}

	return nil
}

// Visit mutates the specified entity's position curve.
func (v *Visitor) Visit(a visitor.Agent) error {
	v.mux.Lock()
	defer v.mux.Unlock()

	switch t := a.AgentType(); t {
	case vcpb.AgentType_AGENT_TYPE_FSM:
		return v.visitFSM(a.(action.Action))
	default:
		return nil
	}
}
