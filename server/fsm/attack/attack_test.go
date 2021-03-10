package attack

import (
	"testing"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/server/entity/projectile"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move/chase"
	"github.com/downflux/game/server/fsm/move/move"

	gdpb "github.com/downflux/game/api/data_go_proto"
	projectileaction "github.com/downflux/game/server/fsm/attack/projectile"
)

var (
	_ action.Action = &Action{}
)

func newTank(
	t *testing.T,
	eid id.EntityID,
	tick id.Tick,
	p *gdpb.Position,
	proj *projectile.Entity) *tank.Entity {
	cid := id.ClientID("client-id")
	tankEntity, err := tank.New(eid, tick, p, cid, proj)
	if err != nil {
		t.Fatalf("New() = %v, want = nil", err)
	}
	return tankEntity
}

func newAction(source *tank.Entity, dest *tank.Entity) *Action {
	s := status.New(0)
	chaseAction := chase.New(s, source, dest)
	return New(s, source, dest, chaseAction)
}

func TestState(t *testing.T) {
	pendingStatus := status.New(0)
	pendingSourceShell, err := projectile.New(
		"sourceShell", 0, &gdpb.Position{X: 0, Y: 0}, id.ClientID("client-id"))
	if err != nil {
		t.Fatalf("New() = %v, want = %v", err)
	}
	pendingSource := newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}, pendingSourceShell)
	pendingTarget := newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 1}, nil)
	pendingMoveAction := move.New(
		pendingSource,
		pendingStatus,
		pendingTarget.Position(pendingStatus.Tick()),
		move.Direct,
	)
	pendingProjectileAction := projectileaction.New(pendingSource, pendingTarget, pendingMoveAction)
	pendingChaseAction := chase.New(pendingStatus, pendingSource, pendingTarget)
	pendingAttackAction := New(pendingStatus, pendingSource, pendingTarget, pendingChaseAction)
	pendingAttackAction.SetProjectileMove(pendingProjectileAction)

	chaseCanceledAction := newAction(
		newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}, nil),
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 1}, nil),
	)
	if err := chaseCanceledAction.chase.Cancel(); err != nil {
		t.Fatalf("Cancel() = %v, want = nil", err)
	}

	targetDeadAction := newAction(
		newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}, nil),
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 1}, nil),
	)
	if err := targetDeadAction.target.TargetHealthCurve().Add(0, -1*targetDeadAction.target.TargetHealth(0)); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	targetNotReady := newAction(
		newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}, nil),
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 1}, nil),
	)
	if err := targetNotReady.source.AttackTimerCurve().Add(0, true); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	targetOutOfRangeSource := newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}, nil)
	targetOutOfRange := newAction(
		targetOutOfRangeSource,
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: targetOutOfRangeSource.AttackRange() + 1}, nil),
	)

	attackCanceled := newAction(
		newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}, nil),
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 1}, nil),
	)
	if err := attackCanceled.Cancel(); err != nil {
		t.Fatalf("Cancel() = %v, want = nil", err)
	}
	if s, err := attackCanceled.chase.State(); err != nil || s != commonstate.Canceled {
		t.Fatalf("State() = %v, %v, want = %v, nil", s, err, commonstate.Canceled)
	}

	testConfigs := []struct {
		name string
		a    *Action
		want fsm.State
	}{
		{name: "TestChaseCanceled", a: chaseCanceledAction, want: commonstate.Canceled},
		{name: "TestTargetDead", a: targetDeadAction, want: commonstate.Finished},
		{name: "TestPendingNotReady", a: targetNotReady, want: commonstate.Pending},
		{name: "TestPendingOutOfRange", a: targetOutOfRange, want: commonstate.Pending},
		{
			name: "TestExecuting",
			a: newAction(
				newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}, nil),
				newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 0}, nil),
			),
			want: commonstate.Executing,
		},
		{name: "TestCancel", a: attackCanceled, want: commonstate.Canceled},
		{name: "TestPendingMove", a: pendingAttackAction, want: commonstate.Pending},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := c.a.State(); err != nil || got != c.want {
				t.Errorf("State() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}
