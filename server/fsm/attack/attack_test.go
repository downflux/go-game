package attack

import (
	"testing"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/fsm/commonstate"
	"github.com/downflux/game/server/fsm/move/chase"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

var (
	_ action.Action = &Action{}
)

func newTank(t *testing.T, eid id.EntityID, tick id.Tick, p *gdpb.Position) *tank.Entity {
	cid := id.ClientID("client-id")
	tankEntity, err := tank.New(eid, tick, p, cid)
	if err != nil {
		t.Fatalf("New() = %v, want = nil", err)
	}
	return tankEntity
}

func newAction(source *tank.Entity, dest *tank.Entity) *Action {
	s := status.New(0)
	chaseAction := chase.New(status.New(0), source, dest)
	return New(s, source, dest, chaseAction)
}

func TestState(t *testing.T) {
	chaseCanceledAction := newAction(
		newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}),
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 1}),
	)
	if err := chaseCanceledAction.chase.Cancel(); err != nil {
		t.Fatalf("Cancel() = %v, want = nil", err)
	}

	targetDeadAction := newAction(
		newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}),
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 1}),
	)
	if err := targetDeadAction.target.TargetHealthCurve().Add(0, -1*targetDeadAction.target.TargetHealth(0)); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	targetNotReady := newAction(
		newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}),
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 1}),
	)
	if err := targetNotReady.source.AttackTimerCurve().Add(0, true); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	targetOutOfRangeSource := newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0})
	targetOutOfRange := newAction(
		targetOutOfRangeSource,
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: targetOutOfRangeSource.AttackRange() + 1}),
	)

	attackCanceled := newAction(
		newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}),
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 1}),
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
				newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}),
				newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 0}),
			),
			want: commonstate.Executing,
		},
		{name: "TestCancel", a: attackCanceled, want: commonstate.Canceled},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := c.a.State(); err != nil || got != c.want {
				t.Errorf("State() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}
