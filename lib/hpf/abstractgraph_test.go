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
	nm.Add(want)

	if got, err := nm.Get(utils.MC(want.GetTileCoordinate())); err != nil || !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("Get() = %v, %v, want = %v, nil", got, err, want)
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
	em.Add(want)

	if got, err := em.Get(utils.MC(want.GetSource()), utils.MC(want.GetDestination())); err != nil || !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("Get() = %v, %v, want = %v, nil", got, err, want)
	}
}
