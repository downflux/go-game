package move

import (
	"math"
	_ "sync"

	"github.com/downflux/game/curve/linearmove"
	"github.com/downflux/game/fsm/fsm"
	"github.com/downflux/game/fsm/instance"
	"github.com/downflux/game/fsm/list"
	"github.com/downflux/game/fsm/move"
	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/pathing/hpf/astar"
	"github.com/downflux/game/pathing/hpf/graph"
	_ "github.com/downflux/game/server/entity/entity"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/status"
	"github.com/downflux/game/server/visitor/dirty"
	"github.com/downflux/game/server/visitor/visitor"
	_ "google.golang.org/protobuf/proto"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/fsm/api/constants_go_proto"
	tile "github.com/downflux/game/map/map"
	vcpb "github.com/downflux/game/server/visitor/api/constants_go_proto"
)

const (
	// ticksPerTile represents the number of ticks necessary to move an
	// Entity instance one full tile width.
	//
	// TODO(minkezhang): Make this a property of the entity.
	ticksPerTile = id.Tick(10)

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

func d(a, b *gdpb.Position) float64 {
	return math.Sqrt(math.Pow(a.GetX()-b.GetX(), 2) + math.Pow(a.GetY()-b.GetY(), 2))
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
	visitor.Base
	visitor.Leaf

	// tileMap is the underlying Map object used for the game.
	tileMap *tile.Map

	// abstractGraph is the underlying abstracted pathing logic data layer
	// for the associated Map.
	abstractGraph *graph.Graph

	// dfStatus is a shared object with the game engine and indicates
	// current tick, etc.
	dfStatus *status.Status

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
		dfStatus:      dfStatus,
		dirties:       dirties,
		minPathLength: minPathLength,
	}
}

// Type returns the registered VisitorType.
func (v *Visitor) Type() vcpb.VisitorType { return visitorType }

// Schedule adds a move command to the internal schedule.
func (v *Visitor) Schedule(args interface{}) error { return nil }

func (v *Visitor) visitFSMList(l *list.List) error { return nil }

func (v *Visitor) visitFSM(i instance.Instance) error {
	if i.Type() != fcpb.FSMType_FSM_TYPE_MOVE {
		return nil
	}

	s, err := i.State()
	if err != nil {
		return err
	}

	m := i.(*move.Instance)

	tick := v.dfStatus.Tick()

	switch s {
	case fsm.State(fcpb.MoveState_MOVE_STATE_EXECUTING):
		e := m.Entity()
		c := e.Curve(gcpb.EntityProperty_ENTITY_PROPERTY_POSITION)
		if c == nil {
			return nil
		}

		p, _, err := astar.Path(
			v.tileMap,
			v.abstractGraph,
			utils.MC(coordinate(c.Get(tick).(*gdpb.Position))),
			utils.MC(coordinate(m.Destination())),
			5,
		)
		if err != nil {
			// TODO(minkezhang): Handle error by logging and continuing.
			return err
		}

		// Add to the existing curve, while smoothing out the existing
		// trajectory.
		prevPos := c.Get(tick).(*gdpb.Position)
		cv := linearmove.New(e.ID(), tick)
		cv.Add(tick, prevPos)
		for i, tile := range p {
			curPos := position(tile.Val.GetCoordinate())
			tickDelta := id.Tick(ticksPerTile.Value() * d(prevPos, curPos))
			cv.Add(tick+id.Tick(i)*ticksPerTile+tickDelta, curPos)
			prevPos = curPos
		}
		if err := v.dirties.Add(dirty.Curve{
			EntityID: e.ID(),
			Property: c.Property(),
		}); err != nil {
			return err
		}
		if err := c.ReplaceTail(cv); err != nil {
			return err
		}

		// Delay next lookup iteration until a suitable time in the
		// future.
		if err := m.Schedule(tick + ticksPerTile*id.Tick(len(p)-1)); err != nil {
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
	switch t := a.AgentType(); t {
	case vcpb.AgentType_AGENT_TYPE_FSM:
		return v.visitFSM(a.(instance.Instance))
	case vcpb.AgentType_AGENT_TYPE_FSM_LIST:
		return v.visitFSMList(a.(*list.List))
	default:
		return nil
	}
}
