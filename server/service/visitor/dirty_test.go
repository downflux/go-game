package dirty

import (
	"fmt"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestPop(t *testing.T) {
	l := New()
	want := Curve{
		EntityID: "some-entity",
	}

	if err := l.Add(want); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got := l.Pop(); got[0] != want {
		t.Fatalf("Pop() = %v, want = %v", got, want)
	}

	if got := l.Pop(); got != nil {
		t.Errorf("Pop() = %v, want = %v", got, nil)
	}
}

func TestAdd(t *testing.T) {
	const nClients = 1000

	l := New()

	var eg errgroup.Group
	for i := 0; i < nClients; i++ {
		i := i
		eg.Go(func() error { return l.Add(Curve{EntityID: fmt.Sprintf("entity-%d", i)}) })
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("Wait() = %v, want = nil", err)
	}

	if got := len(l.Pop()); got != nClients {
		t.Errorf("len() = %v, want = %v", got, nClients)
	}
}
