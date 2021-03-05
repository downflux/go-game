package direct

import (
	"time"

	"github.com/downflux/game/engine/curve/common/linearmove"
	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move/move"

	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

const (
	// fsmType is the registered FSMType of the move visitor.
	fsmType = fcpb.FSMType_FSM_TYPE_DIRECT_MOVE
)

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

	// status is a shared object with the game engine and indicates
	// current tick, etc.
	status status.ReadOnlyStatus

	// dirty is a shared object between the game engine and the Visitor.
	dirty *dirty.List
}

// New constructs a new move Visitor instance.
func New(dfStatus status.ReadOnlyStatus, dcs *dirty.List) *Visitor {
	return &Visitor{
		Base:   *visitor.New(fsmType),
		status: dfStatus,
		dirty:  dcs,
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

		// TODO(minkezhang): Check for collisions, e.g. walls.
		p := []*gdpb.Position{e.Position(tick), node.Destination()}

		cv := linearmove.New(e.ID(), tick)
		cv.Add(tick, e.Position(tick))
		cv.Add(tick+ticksPerTile*id.Tick(utils.Euclidean(p[0], p[1])), p[1])

		if err := v.dirty.AddCurve(dirty.Curve{
			EntityID: e.ID(),
			Property: c.Property(),
		}); err != nil {
			return err
		}
		if err := c.Merge(cv); err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}

// Visit mutates the specified entity's position curve.
func (v *Visitor) Visit(a visitor.Agent) error {
	if node, ok := a.(*move.Action); ok {
		return v.visitFSM(node)
	}
	return nil
}
