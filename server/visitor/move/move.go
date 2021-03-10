package move

import (
	"sync"
	"time"

	"github.com/downflux/game/engine/curve/common/linearmove"
	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/pathing/hpf/astar"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move/move"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
	serverstatus "github.com/downflux/game/engine/status/status"
	tile "github.com/downflux/game/map/map"
)

const (
	// fsmType is the registered FSMType of the move visitor.
	fsmType = fcpb.FSMType_FSM_TYPE_MOVE
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
	visitor.Base

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
	status serverstatus.ReadOnlyStatus

	// dirty is a shared object between the game engine and the Visitor.
	dirty *dirty.List

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
	dfStatus serverstatus.ReadOnlyStatus,
	dirtystate *dirty.List,
	minPathLength int) *Visitor {
	return &Visitor{
		Base:          *visitor.New(fsmType),
		tileMap:       tileMap,
		abstractGraph: abstractGraph,
		status:        dfStatus,
		dirty:         dirtystate,
		minPathLength: minPathLength,
	}
}

// TODO(minkezhang): Remove function.
func (v *Visitor) Schedule(interface{}) error { return nil }

func (v *Visitor) generatePath(node *move.Action) ([]*tile.Tile, error) {
	t := node.MoveType()
	switch t {
	case move.Default:
		p, _, err := astar.Path(
			v.tileMap,
			v.abstractGraph,
			utils.MC(coordinate(node.Component().Position(v.status.Tick()))),
			utils.MC(coordinate(node.Destination())),
			v.minPathLength,
		)
		// TODO(minkezhang): Handle error by logging and continuing.
		return p, err
	case move.Direct:
		if node.Destination().GetX() >= float64(v.tileMap.D.X) ||
			node.Destination().GetY() >= float64(v.tileMap.D.Y) {
			return nil, status.Errorf(
				codes.OutOfRange,
				"input Tile coordinate %v is outside the map boundary %v",
				node.Destination(),
				v.tileMap.D,
			)
		}

		// TODO(minkezhang): Check for collisions, e.g. walls.
		return []*tile.Tile{v.tileMap.TileFromCoordinate(coordinate(node.Destination()))}, nil
	default:
		return nil, status.Errorf(
			codes.Unimplemented,
			"cannot process move of type %v", t)
	}
}

func (v *Visitor) visitFSM(node *move.Action) error {
	s, err := node.State()
	if err != nil {
		return err
	}

	tick := v.status.Tick()

	switch s {
	case commonstate.Executing:
		e := node.Component()
		c := e.PositionCurve()
		if c == nil {
			return nil
		}

		ticksPerSecond := float64(time.Second / v.status.TickDuration())
		ticksPerTile := id.Tick(ticksPerSecond / e.MoveVelocity())

		p, err := v.generatePath(node)
		if err != nil {
			return err
		}

		// Add to the existing curve, while smoothing out the existing
		// trajectory.
		prevPos := e.Position(tick)

		cv := linearmove.New(e.ID(), tick)
		cv.Add(tick, prevPos)

		tickOffset := id.Tick(0)
		for _, tile := range p {
			curPos := position(tile.Val.GetCoordinate())
			distance := utils.Euclidean(prevPos, curPos)
			tickOffset += ticksPerTile * id.Tick(distance)
			cv.Add(tick+tickOffset, curPos)
			prevPos = curPos
		}

		if err := v.dirty.AddCurve(dirty.Curve{
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
		// We only need to do this for A*-generated paths, since direct
		// paths do not recalculating.
		//
		// TODO(minkezhang): Add test for scheduling here.
		if node.MoveType() == move.Default {
			if err := node.SchedulePartialMove(tick + tickOffset); err != nil {
				// TODO(minkezhang): Handle error by logging and continuing.
				return err
			}
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

	if node, ok := a.(*move.Action); ok {
		return v.visitFSM(node)
	}
	return nil
}
