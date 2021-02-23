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
