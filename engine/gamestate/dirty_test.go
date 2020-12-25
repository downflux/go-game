package dirty

import (
	"fmt"
	"testing"

	"github.com/downflux/game/engine/id/id"
	"golang.org/x/sync/errgroup"
)

func TestPop(t *testing.T) {
	l := New()
	want := Curve{
		EntityID: "some-entity",
	}

	if err := l.AddCurve(want); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got := l.Pop().Curves()[0]; got != want {
		t.Fatalf("Curves() = %v, want = %v", got, want)
	}

	if got := l.Pop().Curves(); got != nil {
		t.Errorf("Curves() = %v, want = %v", got, nil)
	}
}

func TestAdd(t *testing.T) {
	const nClients = 1000

	l := New()

	var eg errgroup.Group
	for i := 0; i < nClients; i++ {
		i := i
		eg.Go(func() error { return l.AddCurve(Curve{EntityID: id.EntityID(fmt.Sprintf("entity-%d", i))}) })
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("Wait() = %v, want = nil", err)
	}

	if got := len(l.Pop().Curves()); got != nClients {
		t.Errorf("len() = %v, want = %v", got, nClients)
	}
}
