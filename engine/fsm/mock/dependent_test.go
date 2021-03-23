package dependent

import (
	"testing"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/fsm/mock/simple"
	"github.com/downflux/game/engine/id/id"
)

var (
	_ action.Action = &Action{}
)

func TestState(t *testing.T) {
	c := New(id.ActionID("child-id"), 0, nil)
	p := New(id.ActionID("parent-id"), 0, c)

	for _, n := range []action.Action{c, p} {
		want := fsm.State(simple.Pending)
		if got, err := n.State(); err != nil || got != want {
			t.Fatalf("State() = %v, %v, want = %v, nil", got, err, want)
		}
	}

	for _, n := range []action.Action{c, p} {
		if err := n.Cancel(); err != nil {
			t.Fatalf("Cancel() = %v, want = nil", err)
		}
	}

	for _, n := range []action.Action{c, p} {
		want := fsm.State(simple.Canceled)
		if got, err := n.State(); err != nil || got != want {
			t.Fatalf("State() = %v, %v, want = %v, nil", got, err, want)
		}
	}
}
