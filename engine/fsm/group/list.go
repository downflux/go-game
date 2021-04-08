package list

import (
	"github.com/downflux/game/engine/fsm/group/group"
	"github.com/downflux/game/engine/id/id"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

type List struct {
	groups map[id.GroupID]group.Group // Read-only.
}

func New(groups []group.Group) (*List, error) {
	l := &List{}
	for _, g := range groups {
		for _, o := range groups {
			if g.ID() != o.ID() && (g.Collide(o) || o.Collide(g)) {
				return nil, status.Errorf(
					codes.FailedPrecondition,
					"cannot create FSM group list with conflicting groups %v and %v",
					g.ID(),
					o.ID(),
				)
			}
		}
		l.groups[g.ID()] = g
	}
	return l, nil
}

func (l List) Group(gid id.GroupID) group.Group { return l.groups[gid] }
func (l List) ID(t fcpb.FSMType) id.GroupID {
	for _, g := range l.groups {
		if g.Contains(t) {
			return g.ID()
		}
	}
	return ""
}
