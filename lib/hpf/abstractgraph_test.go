package abstractgraph

import (
	"testing"

	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/utils"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestAbstractNodeMapAdd(t *testing.T) {
	want := &rtsspb.AbstractNode{
		TileCoordinate: &rtsspb.Coordinate{
			X: 1,
			Y: 1,
		},
	}

	nm := AbstractNodeMap{}
	if err := nm.Add(want); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got, err := nm.Get(utils.MC(want.GetTileCoordinate())); err != nil || !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("Get() = %v, %v, want = %v, nil", got, err, want)
	}
}

func TestAbstractNodeMapRemove(t *testing.T) {
	c := &rtsspb.Coordinate{X: 1, Y: 1}
	nm := AbstractNodeMap{utils.MC(c): &rtsspb.AbstractNode{TileCoordinate: c}}
	nm.Remove(utils.MC(c))

	if got, err := nm.Get(utils.MC(c)); err != nil || got != nil {
		t.Errorf("Get() = %v, %v, want = nil, nil", got, err)
	}
}

func TestAbstractEdgeMapAdd(t *testing.T) {
	want := &rtsspb.AbstractEdge{
		Source: &rtsspb.Coordinate{
			X: 0,
			Y: 1,
		},
		Destination: &rtsspb.Coordinate{
			X: 1,
			Y: 0,
		},
	}

	em := AbstractEdgeMap{}
	if err := em.Add(want); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got, err := em.Get(utils.MC(want.GetSource()), utils.MC(want.GetDestination())); err != nil || !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("Get() = %v, %v, want = %v, nil", got, err, want)
	}
}

func TestAbstractEdgeMapAddError(t *testing.T) {
	s := &rtsspb.Coordinate{X: 0, Y: 1}
	d := &rtsspb.Coordinate{X: 1, Y: 0}

	em := AbstractEdgeMap{}
	if err := em.Add(&rtsspb.AbstractEdge{Source: s, Destination: d}); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if err := em.Add(&rtsspb.AbstractEdge{Source: s, Destination: d}); err == nil {
		t.Errorf("Add() = nil, want a non-nil error")
	}

}

func TestAbstractEdgeMapGetCommutative(t *testing.T) {
	s := &rtsspb.Coordinate{X: 0, Y: 1}
	d := &rtsspb.Coordinate{X: 1, Y: 0}

	em := AbstractEdgeMap{}
	if err := em.Add(&rtsspb.AbstractEdge{Source: s, Destination: d}); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	got1, err := em.Get(utils.MC(s), utils.MC(d))
	if err != nil {
		t.Fatalf("Get() = _, %v, want = _, nil", err)
	}
	got2, err := em.Get(utils.MC(d), utils.MC(s))
	if err != nil {
		t.Fatalf("Get() = _, %v, want = _, nil", err)
	}

	if diff := cmp.Diff(got1, got2, protocmp.Transform()); diff != "" {
		t.Errorf("Get() mismatch (-want +got):\n%s", diff)
	}
}

func TestAbstractEdgeMapRemove(t *testing.T) {
	s := &rtsspb.Coordinate{X: 0, Y: 1}
	d := &rtsspb.Coordinate{X: 1, Y: 0}

	em := AbstractEdgeMap{}
	if err := em.Add(&rtsspb.AbstractEdge{Source: s, Destination: d}); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if err := em.Remove(utils.MC(s), utils.MC(d)); err != nil {
		t.Fatalf("Remove() = %v, want = nil", err)
	}

	if got, err := em.Get(utils.MC(s), utils.MC(d)); err != nil || got != nil {
		t.Errorf("Get() = %v, %v, want = nil, nil", got, err)
	}
}
