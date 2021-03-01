package chase

import (
	"testing"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/fsm/commonstate"

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
	return New(status.New(0), source, dest)
}

func TestState(t *testing.T) {
	actionWithMoveInRange := newAction(
		newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}),
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: chaseRadius}),
	)
	if err := actionWithMoveInRange.SetMove(GenerateMove(actionWithMoveInRange)); err != nil {
		t.Fatalf("SetMove() = %v, want = nil", err)
	}

	actionWithMoveOutOfRange := newAction(
		newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}),
		newTank(t, "target", 0, &gdpb.Position{X: 0, Y: chaseRadius + 1}),
	)
	if err := actionWithMoveOutOfRange.SetMove(GenerateMove(actionWithMoveOutOfRange)); err != nil {
		t.Fatalf("SetMove() = %v, want = nil", err)
	}

	s0 := &gdpb.Position{X: 0, Y: 0}
	d0 := &gdpb.Position{X: 0, Y: chaseRadius + 1}
	d1 := &gdpb.Position{X: 0, Y: (chaseRadius + 1) + (chaseRadius + 1)}
	actionWithFinishedOutOfRange := newAction(
		newTank(t, "source", 0, s0),
		newTank(t, "target", 0, d0),
	)
	if err := actionWithFinishedOutOfRange.SetMove(GenerateMove(actionWithFinishedOutOfRange)); err != nil {
		t.Fatalf("SetMove() = %v, want = nil", err)
	}
	actionWithFinishedOutOfRange.Source().PositionCurve().Add(actionWithFinishedOutOfRange.Status().Tick(), d0)
	actionWithFinishedOutOfRange.Destination().PositionCurve().Add(actionWithFinishedOutOfRange.Status().Tick(), d1)
	if got, err := actionWithFinishedOutOfRange.move.State(); err != nil || got != commonstate.Finished {
		t.Fatalf("State() = %v, %v, want = %v, nil", err, got, commonstate.Finished)
	}

	actionWithPropagatedCancel := newAction(newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}), newTank(t, "target", 0, &gdpb.Position{X: 0, Y: 1}))
	if err := actionWithPropagatedCancel.SetMove(GenerateMove(actionWithPropagatedCancel)); err != nil {
		t.Fatalf("SetMove() = %v, want = nil", err)
	}
	if err := actionWithPropagatedCancel.move.Cancel(); err != nil {
		t.Fatalf("Cancel() = %v, want = nil", err)
	}

	testConfigs := []struct {
		name string
		a    *Action
		want fsm.State
	}{
		{
			name: "NewWithinRange",
			a: newAction(
				newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}),
				newTank(t, "target", 0, &gdpb.Position{X: 0, Y: chaseRadius})),
			want: commonstate.Pending,
		},
		{
			name: "NewOutOfRange",
			a: newAction(
				newTank(t, "source", 0, &gdpb.Position{X: 0, Y: 0}),
				newTank(t, "target", 0, &gdpb.Position{X: 0, Y: chaseRadius + 1})),
			want: OutOfRange,
		},
		{name: "MoveInRange", a: actionWithMoveInRange, want: commonstate.Pending},
		{name: "MoveOutOfRange", a: actionWithMoveOutOfRange, want: commonstate.Pending},
		{name: "FinishedOutOfRange", a: actionWithFinishedOutOfRange, want: OutOfRange},
		{name: "PropagatedMoveCancel", a: actionWithPropagatedCancel, want: commonstate.Canceled},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := c.a.State(); err != nil || got != c.want {
				t.Errorf("State() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}

}
