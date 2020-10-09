package edge

import (
	"testing"

	rtsspb "github.com/downflux/game/pathing/proto/structs_go_proto"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/downflux/game/pathing/hpf/utils"
	"google.golang.org/protobuf/testing/protocmp"
)

func coordinateLess(c1, c2 *rtsspb.Coordinate) bool {
	return c1.GetX() < c2.GetX() || c1.GetX() == c2.GetX() && c1.GetY() < c2.GetY()
}

func abstractEdgeLess(e1, e2 *rtsspb.AbstractEdge) bool {
	return coordinateLess(e1.GetSource(), e2.GetSource()) || cmp.Equal(e1.GetSource(), e2.GetSource(), protocmp.Transform()) && coordinateLess(e1.GetDestination(), e2.GetDestination())
}

func TestMapAddError(t *testing.T) {
	addAlreadyExistMap := &Map{}
	if err := addAlreadyExistMap.Add(&rtsspb.AbstractEdge{
		Source:      &rtsspb.Coordinate{X: 0, Y: 0},
		Destination: &rtsspb.Coordinate{X: 1, Y: 1},
	}); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	testConfigs := []struct {
		name string
		m    *Map
		n    *rtsspb.AbstractEdge
	}{
		{
			name: "AddAlreadyExist",
			m:    addAlreadyExistMap,
			n: &rtsspb.AbstractEdge{
				Source:      &rtsspb.Coordinate{X: 0, Y: 0},
				Destination: &rtsspb.Coordinate{X: 1, Y: 1},
				Weight:      1,
			},
		},
		{
			name: "AddSelfLoop",
			m:    &Map{},
			n: &rtsspb.AbstractEdge{
				Source:      &rtsspb.Coordinate{X: 0, Y: 0},
				Destination: &rtsspb.Coordinate{X: 0, Y: 0},
				Weight:      1,
			},
		},
	}
	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if err := c.m.Add(c.n); err == nil {
				t.Error("Add() = nil, want a non-nil error")
			}
		})
	}
}

func TestMapAdd(t *testing.T) {
	want := &rtsspb.AbstractEdge{
		Source:      &rtsspb.Coordinate{X: 0, Y: 0},
		Destination: &rtsspb.Coordinate{X: 1, Y: 1},
		Weight:      1,
	}

	em := &Map{}
	if err := em.Add(want); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got, err := em.Get(
		utils.MC(want.GetSource()),
		utils.MC(want.GetDestination())); err != nil || !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("Get() = %v, %v, want = %v, nil", got, err, want)
	}
}

func TestMapGet(t *testing.T) {
	want := &rtsspb.AbstractEdge{
		Source:      &rtsspb.Coordinate{X: 0, Y: 0},
		Destination: &rtsspb.Coordinate{X: 1, Y: 1},
		Weight:      1,
	}

	em := &Map{}
	if err := em.Add(want); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	got1, err := em.Get(
		utils.MC(want.GetSource()),
		utils.MC(want.GetDestination()))
	if err != nil {
		t.Fatalf("Get() = _, %v, want = _, nil", err)
	}
	got2, err := em.Get(
		utils.MC(want.GetDestination()),
		utils.MC(want.GetSource()))
	if err != nil {
		t.Fatalf("Get() = _, %v, want = _, nil", err)
	}

	if diff := cmp.Diff(got1, got2, protocmp.Transform()); diff != "" {
		t.Errorf("Get() mismatch (-want +got):\n%s", diff)
	}
}

func TestMapPop(t *testing.T) {
	want := &rtsspb.AbstractEdge{
		Source:      &rtsspb.Coordinate{X: 0, Y: 0},
		Destination: &rtsspb.Coordinate{X: 1, Y: 1},
		Weight:      1,
	}

	em := &Map{}
	if err := em.Add(want); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got, err := em.Pop(
		utils.MC(want.GetSource()),
		utils.MC(want.GetDestination())); err != nil || !cmp.Equal(want, got, protocmp.Transform()) {
		t.Errorf("Pop() = %v, %v, want = %v, nil", got, err, want)
	}
}

func TestMapGetBySource(t *testing.T) {
	s := &rtsspb.Coordinate{X: 0, Y: 0}
	want := []*rtsspb.AbstractEdge{
		{
			Source:      s,
			Destination: &rtsspb.Coordinate{X: 1, Y: 1},
			Weight:      1,
		},
		{
			Source:      &rtsspb.Coordinate{X: 2, Y: 2},
			Destination: s,
			Weight:      1,
		},
	}

	em := &Map{}
	for _, n := range want {
		if err := em.Add(n); err != nil {
			t.Fatalf("Add() = %v, want = nil", err)
		}
	}
	if err := em.Add(&rtsspb.AbstractEdge{
		Source:      &rtsspb.Coordinate{X: 1, Y: 1},
		Destination: &rtsspb.Coordinate{X: 2, Y: 2},
		Weight:      1,
	}); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	got, err := em.GetBySource(utils.MC(s))
	if err != nil {
		t.Fatalf("GetBySource() = _, %v, want = _, nil", err)
	}
	if diff := cmp.Diff(want, got, protocmp.Transform(), cmpopts.SortSlices(abstractEdgeLess)); diff != "" {
		t.Errorf("GetBySource() mismatch (-want +got):\n%s", diff)
	}
}
