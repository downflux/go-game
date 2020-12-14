package move

import (
	"sync"

	"github.com/downflux/game/fsm/fsm"
	"github.com/downflux/game/fsm/instance"
	"github.com/downflux/game/server/entity/entity"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/status"
	"google.golang.org/protobuf/proto"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

const (
	pending   = "PENDING"
	executing = "EXECUTING"
	canceled  = "CANCELED"
	finished  = "FINISHED"
)

var (
	transitions = []fsm.Transition{
		{From: pending, To: executing},
		{From: pending, To: canceled},
		{From: pending, To: finished},
		{From: executing, To: pending},
	}

	FSM = fsm.New(transitions)
)

type Instance struct {
	*instance.Base

	// mux guards e, dfStatus, scheduledTick, and destination properties
	mux      sync.Mutex
	e        entity.Entity
	dfStatus *status.Status

	// TODO(minkezhang): Move eid, scheduledTick, and destination into
	// separate external cache.
	scheduledTick id.Tick
	destination   *gdpb.Position
}

func New(
	eid id.EntityID,
	dfStatus *status.Status,
	destination *gdpb.Position) *Instance {
	return &Instance{
		Base: instance.New(FSM, pending),
	}
}

/**
 * // TODO(minkezhang): Implement.
 * func (n *Instance) Accept(v visitor.Visitor) error { return v.Visit(n) }
 */

func (n *Instance) State() (fsm.State, error) {
	n.mux.Lock()
	defer n.mux.Unlock()

	tick := n.dfStatus.Tick()

	s, err := n.Base.State()
	if err != nil {
		return "", err
	}

	switch s {
	case pending:
		c := n.e.Curve(gcpb.EntityProperty_ENTITY_PROPERTY_POSITION)
		if proto.Equal(n.destination, c.Get(tick).(*gdpb.Position)) {
			return finished, nil
		}
		if n.scheduledTick <= tick {
			return executing, nil
		}
		if n.destination == nil {
			return canceled, nil
		}
		return pending, nil
	default:
		return s, nil
	}
}
