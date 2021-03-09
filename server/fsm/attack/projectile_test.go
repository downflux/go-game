package projectile

import (
	"testing"
	"time"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/server/entity/projectile"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move/move"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

var (
	_ action.Action = &Action{}
)

func newProjectile(
	t *testing.T,
	eid id.EntityID,
	tick id.Tick,
	pos *gdpb.Position,
	cid id.ClientID) *projectile.Entity {
	e, err := projectile.New(eid, tick, pos, cid)
	if err != nil {
		t.Fatalf("New() = %v, want = nil", err)
	}
	return e
}

func newTank(
	t *testing.T,
	eid id.EntityID,
	tick id.Tick,
	pos *gdpb.Position,
	cid id.ClientID,
	proj *projectile.Entity) *tank.Entity {
	e, err := tank.New(eid, tick, pos, cid, proj)
	if err != nil {
		t.Fatalf("New() = %v, want = nil", err)
	}
	return e
}

func TestState(t *testing.T) {
	cid := id.ClientID("client-id")
	s := status.New(time.Millisecond)

	p0 := &gdpb.Position{X: 0, Y: 0}
	p1 := &gdpb.Position{X: 0, Y: 1}

	shell := newProjectile(t, id.EntityID("source-shell"), 0, p0, cid)
	source := newTank(t, id.EntityID("source-entity"), 0, p0, cid, shell)
	target := newTank(t, id.EntityID("target-entity"), 0, p1, cid, nil)

	canceledMove := move.New(shell, s, target.Position(s.Tick()))
	if err := canceledMove.Cancel(); err != nil {
		t.Fatalf("Cancel() = %v, want = nil", err)
	}

	testConfigs := []struct {
		name string
		i    *Action
		want fsm.State
	}{
		{
			name: "TestMoving",
			i:    New(source, target, move.New(shell, s, target.Position(s.Tick()))),
			want: commonstate.Pending,
		},
		{
			name: "TestReachedDestination",
			i:    New(source, target, move.New(shell, s, shell.Position(s.Tick()))),
			want: commonstate.Executing,
		},
		{
			name: "TestCanceledMove",
			i:    New(source, target, canceledMove),
			want: commonstate.Canceled,
		},
	}

	for _, c := range testConfigs {
		if got, err := c.i.State(); err != nil || got != c.want {
			t.Errorf("State() = %v, %v, want = %v, nil", got, err, c.want)
		}
	}
}
