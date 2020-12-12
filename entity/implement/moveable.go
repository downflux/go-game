package moveable

import (
	"github.com/downflux/game/server/id"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

type Moveable struct {
	moveClient id.ClientID
	moveTarget *gdpb.Position
}

func (m *Moveable) ScheduleMove(cid id.ClientID, dest *gdpb.Position) error {
	m.moveClient = cid
	m.moveTarget = dest
	return nil
}

func (m *Moveable) MoveClient() id.ClientID    { return m.moveClient }
func (m *Moveable) MoveTarget() *gdpb.Position { return m.moveTarget }
