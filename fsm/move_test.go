package move

import (
	"testing"

	"github.com/downflux/game/fsm/fsm"
	"github.com/downflux/game/fsm/instance"
	"github.com/downflux/game/server/entity/tank"
	"github.com/downflux/game/server/service/status"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

var (
	_ instance.Instance = &Instance{}
)

func TestState(t *testing.T) {
	const eid = "entity-id"
	const t0 = 0
	p0 := &gdpb.Position{X: 0, Y: 0}

	executingNewEntity := tank.New(eid, t0, p0)
	executingNewStatus := status.New(0)
	executingNewI1 := New(executingNewEntity, executingNewStatus, &gdpb.Position{X: 1, Y: 1})

	scheduleEntity := tank.New(eid, t0, p0)
	scheduleStatus := status.New(0)
	scheduleI1 := New(scheduleEntity, scheduleStatus, &gdpb.Position{X: 1, Y: 1})
	scheduleI1.Schedule(100)

	cancelEntity := tank.New(eid, t0, p0)
	cancelStatus := status.New(0)
	cancelI1 := New(cancelEntity, cancelStatus, &gdpb.Position{X: 1, Y: 1})
	cancelI1.Cancel()

	finishedEntity := tank.New(eid, t0, p0)
	finishedStatus := status.New(0)
	finishedI1 := New(finishedEntity, finishedStatus, p0)

	pendingCanceledEntity := tank.New(eid, t0, p0)
	pendingCanceledStatus := status.New(0)
	pendingCanceledI1 := New(pendingCanceledEntity, pendingCanceledStatus, &gdpb.Position{X: 1, Y: 1})
	pendingCanceledI1.Schedule(100)
	pendingCanceledI1.Cancel()

	testConfigs := []struct {
		name string
		i    *Instance
		want fsm.State
	}{
		{name: "NewExecutingTest", i: executingNewI1, want: executing},
		{name: "ScheduleTest", i: scheduleI1, want: pending},
		{name: "CanceledTest", i: cancelI1, want: canceled},
		{name: "FinishedTest", i: finishedI1, want: finished},
		{name: "PendingCanceledTest", i: pendingCanceledI1, want: canceled},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			got, err := c.i.State()
			if err != nil || got != c.want {
				t.Fatalf("State() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestPrecedence(t *testing.T) {
	const eid = "entity-id"
	const t0 = 0
	p0 := &gdpb.Position{X: 0, Y: 0}

	sameTickEntity := tank.New(eid, t0, p0)
	sameTickStatus := status.New(0)
	sameTickI1 := New(sameTickEntity, sameTickStatus, &gdpb.Position{X: 1, Y: 1})
	sameTickI2 := New(sameTickEntity, sameTickStatus, &gdpb.Position{X: 2, Y: 2})

	diffTickSamePosEntity := tank.New(eid, t0, p0)
	diffTickSamePosStatus := status.New(0)
	diffTickSamePosI1 := New(diffTickSamePosEntity, diffTickSamePosStatus, &gdpb.Position{X: 1, Y: 1})
	diffTickSamePosStatus.IncrementTick()
	diffTickSamePosI2 := New(diffTickSamePosEntity, diffTickSamePosStatus, &gdpb.Position{X: 1, Y: 1})

	precedenceEntity := tank.New(eid, t0, p0)
	precedenceStatus := status.New(0)
	precedenceI1 := New(precedenceEntity, precedenceStatus, &gdpb.Position{X: 1, Y: 1})
	precedenceStatus.IncrementTick()
	precedenceI2 := New(precedenceEntity, precedenceStatus, &gdpb.Position{X: 2, Y: 2})

	// We are testing if i1 < i2.
	testConfigs := []struct {
		name string
		i1   *Instance
		i2   *Instance
		want bool
	}{
		{name: "SameTickNoPrecedenceTest", i1: sameTickI1, i2: sameTickI2, want: false},
		{name: "DiffTickSamePosNoPrecedenceTest", i1: diffTickSamePosI1, i2: diffTickSamePosI2, want: false},
		{name: "DiffTickSamePosNoPrecedenceCommutativeTest", i1: diffTickSamePosI2, i2: diffTickSamePosI1, want: false},
		{name: "PrecedenceTest", i1: precedenceI1, i2: precedenceI2, want: false},
		{name: "PrecedenceCommutativeTest", i1: precedenceI2, i2: precedenceI1, want: true},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := c.i1.Precedence(c.i2); got != c.want {
				t.Errorf("Precedence() = %v, want = %v", got, c.want)
			}
		})
	}
}
