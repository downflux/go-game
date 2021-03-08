package projectile

import (
	"testing"
	"time"

	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/server/entity/projectile"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move/move"
	"github.com/google/go-cmp/cmp"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	projectileaction "github.com/downflux/game/server/fsm/attack/projectile"
)

var (
	_ visitor.Visitor = &Visitor{}
)

func newTank(t *testing.T, eid id.EntityID, tick id.Tick, pos *gdpb.Position, cid id.ClientID) *tank.Entity {
	e, err := tank.New(eid, tick, pos, cid)
	if err != nil {
		t.Fatalf("New() = %v, want = nil", err)
	}
	return e
}

func TestAttack(t *testing.T) {
	cid := id.ClientID("client-id")
	eid0 := id.EntityID("entity-id-1")
	eid1 := id.EntityID("entity-id-2")
	peid := id.EntityID("projectile-of-eid1")
	p0 := &gdpb.Position{X: 0, Y: 0}
	p1 := &gdpb.Position{X: 1, Y: 0}
	t0 := id.Tick(0)

	s := status.New(time.Millisecond)
	d := dirty.New()

	projectileVisitor := New(s, d)

	tank0 := newTank(t, eid0, t0, p0, cid)
	tank1 := newTank(t, eid1, t0, p1, cid)

	hp := tank1.TargetHealth(s.Tick())

	// TODO(minkezhang): Link to tank.
	shell, err := projectile.New(peid, t0, p1, cid)
	if err != nil {
		t.Fatalf("New() = %v, want = nil", err)
	}

	moveFSM := move.New(shell, s, p1)
	projectileFSM := projectileaction.New(tank0, tank1, moveFSM)

	// Verify preconditions.
	// TODO(minkezhang): Move to projectile FSM test instead.
	if got, err := moveFSM.State(); err != nil || got != commonstate.Finished {
		t.Fatalf("State() = %v, %v, want = %v, nil", got, err, commonstate.Finished)
	}
	if got, err := projectileFSM.State(); err != nil || got != commonstate.Executing {
		t.Fatalf("State() = %v, %v, want = %v, nil", got, err, commonstate.Executing)
	}

	if err := projectileVisitor.Visit(projectileFSM); err != nil {
		t.Fatalf("Visit() = %v, want = nil", err)
	}

	// Verify target was damaged.
	if got := tank1.TargetHealth(s.Tick()); got >= hp {
		t.Fatalf("GetHealth() = %v, want < %v", got, hp)
	}

	// Verify dirty curve.
	dirtyCurves := []dirty.Curve{
		{EntityID: eid1, Property: gcpb.EntityProperty_ENTITY_PROPERTY_HEALTH},
	}
	got := projectileVisitor.dirty.Pop().Curves()
	if diff := cmp.Diff(dirtyCurves, got); diff != "" {
		t.Fatalf("Curves() mismatch (-want +got):\n%v", diff)
	}

	if got, err := projectileFSM.State(); err != nil || got != commonstate.Finished {
		t.Fatalf("State() = %v, %v, want = %v, nil", got, err, commonstate.Finished)
	}
}
