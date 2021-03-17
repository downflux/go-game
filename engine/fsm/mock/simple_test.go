package simple

import (
	"testing"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/id/id"
)

var (
	_ action.Action = &Action{}
)

func TestNew(t *testing.T) {
	aid := id.ActionID("action-id")
	priority := 1

	a := New(aid, priority)

	if got, err := a.State(); err != nil || got != Pending {
		t.Fatalf("State() = %v, %v, want = %v, nil", got, err, Pending)
	}

	if got := a.ID(); got != aid {
		t.Fatalf("ID() = %v, want = %v", got, aid)
	}
	if got := a.priority; got != priority {
		t.Errorf("a.priority = %v, want = %v", got, priority)
	}
}

func TestPrecedence(t *testing.T) {
	cancelAction := New(id.ActionID("action-id"), 0)
	cancelAction.Cancel()

	low := New(id.ActionID("action-id"), 1)
	high := New(id.ActionID("action-id"), 19)

	highDiffID := New(id.ActionID("diff-action-id"), 19)

	testConfigs := []struct {
		name string
		a1   *Action
		a2   *Action
		want bool
	}{
		{name: "SimpleCase", a1: low, a2: high, want: false},
		{name: "SimpleCaseReverse", a1: high, a2: low, want: true},
		{name: "CancelLowPriority", a1: cancelAction, a2: low, want: true},
		{name: "DifferentID", a1: highDiffID, a2: low, want: false},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := c.a1.Precedence(c.a2); got != c.want {
				t.Errorf("Precedence() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestCancel(t *testing.T) {
	n := New(id.ActionID("entity-id"), 0)

	err := n.Cancel()
	if err != nil {
		t.Fatalf("Cancel() = %v, want = nil", err)
	}
}
