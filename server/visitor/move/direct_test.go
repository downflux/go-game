package direct

import (
	"testing"
	"time"

	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/status/status"
	"github.com/downflux/game/server/entity/projectile"
	"github.com/downflux/game/server/fsm/move/move"
	"github.com/google/go-cmp/cmp"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
)

func TestMoveError(t *testing.T) {
	eid := id.EntityID("entity-id")
	cid := id.ClientID("client-id")
	t0 := id.Tick(0)
	p0 := &gdpb.Position{X: 0, Y: 0}
	p1 := &gdpb.Position{X: 0, Y: 11}
	dimension := &gdpb.Coordinate{X: 10, Y: 10}

	s := status.New(time.Millisecond)

	v := New(s, dirty.New(), dimension)
	p, err := projectile.New(eid, t0, p0, cid)
	if err != nil {
		t.Fatalf("New() = %v, want = %v", err)
	}
	i := move.New(p, s, p1)

	if err := v.Visit(i); err == nil {
		t.Error("Visit() = nil, want a non-nil error")
	}
}

func TestMove(t *testing.T) {
	eid := id.EntityID("entity-id")
	cid := id.ClientID("client-id")
	t0 := id.Tick(0)
	p0 := &gdpb.Position{X: 0, Y: 0}
	p1 := &gdpb.Position{X: 0, Y: 9}
	dimension := &gdpb.Coordinate{X: 10, Y: 10}

	s := status.New(time.Millisecond)

	v := New(s, dirty.New(), dimension)
	p, err := projectile.New(eid, t0, p0, cid)
	if err != nil {
		t.Fatalf("New() = %v, want = %v", err)
	}
	i := move.New(p, s, p1)

	if err := v.Visit(i); err != nil {
		t.Fatalf("Visit() = %v, want = nil", err)
	}

	want := []dirty.Curve{
		{EntityID: eid, Property: gcpb.EntityProperty_ENTITY_PROPERTY_POSITION},
	}

	got := v.dirty.Pop().Curves()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Pop() mismatch (-want +got):\n%v", diff)
	}
}
