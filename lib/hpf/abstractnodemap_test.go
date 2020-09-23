package abstractnodemap

import (
	"testing"

	rtsspb "github.com/minkezhang/rts-pathing/lib/proto/structs_go_proto"

	"github.com/minkezhang/rts-pathing/lib/hpf/cluster"
	"github.com/minkezhang/rts-pathing/lib/hpf/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	trivialClusterMap, _ = cluster.ImportMap(
		&rtsspb.ClusterMap{
			TileDimension:    &rtsspb.Coordinate{X: 1, Y: 1},
			TileMapDimension: &rtsspb.Coordinate{X: 1, Y: 1},
		})
)

func coordinateLess(c1, c2 *rtsspb.Coordinate) bool {
	return c1.GetX() < c2.GetX() || (c1.GetX() == c2.GetX() && c1.GetY() < c2.GetY())
}

func abstractNodeLess(n1, n2 *rtsspb.AbstractNode) bool {
	return coordinateLess(n1.GetTileCoordinate(), n2.GetTileCoordinate())
}

func TestMapAdd(t *testing.T) {
	want := &rtsspb.AbstractNode{
		TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
	}

	m := Map{ClusterMap: trivialClusterMap}

	if err := m.Add(want); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got, err := m.Get(utils.MC(want.GetTileCoordinate())); err != nil || !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("Get() = %v, %v, want = %v, nil", got, err, want)
	}
}

func TestMapGet(t *testing.T) {
	m := Map{ClusterMap: trivialClusterMap}
	if got, err := m.Get(utils.MapCoordinate{X: 0, Y: 0}); err != nil || got != nil {
		t.Errorf("Get() = %v, %v, want = nil, nil", got, err)
	}
}

func TestMapGetByCluster(t *testing.T) {
	cm, err := cluster.ImportMap(&rtsspb.ClusterMap{
		TileDimension:    &rtsspb.Coordinate{X: 2, Y: 2},
		TileMapDimension: &rtsspb.Coordinate{X: 4, Y: 4},
	})
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}

	m := Map{ClusterMap: cm}
	want := []*rtsspb.AbstractNode{
		{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0}},
		{TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 1}},
	}

	for _, n := range want {
		if err := m.Add(n); err != nil {
			t.Fatalf("Add() = %v, want = nil", err)
		}
	}

	if err := m.Add(
		&rtsspb.AbstractNode{TileCoordinate: &rtsspb.Coordinate{X: 3, Y: 3}}); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	got, err := m.GetByCluster(utils.MC(&rtsspb.Coordinate{X: 0, Y: 0}))
	if err != nil {
		t.Fatalf("GetByCluster() = _, %v, want = _, nil", err)
	}

	if diff := cmp.Diff(want, got, protocmp.Transform(), cmpopts.SortSlices(abstractNodeLess)); diff != "" {
		t.Errorf("GetByCluster() mismatch (-want +got):\n%s", diff)
	}
}

func TestMapAddError(t *testing.T) {
	clusterCoordinate := &rtsspb.Coordinate{X: 0, Y: 0}
	tileCoordinate := &rtsspb.Coordinate{X: 0, Y: 0}
	testConfigs := []struct {
		name string
		m    Map
		n    *rtsspb.AbstractNode
	}{
		{name: "LevelMismatch", m: Map{ClusterMap: trivialClusterMap}, n: &rtsspb.AbstractNode{Level: 1, TileCoordinate: tileCoordinate}},
		{
			name: "AlreadyExist",
			m: Map{
				ClusterMap: trivialClusterMap,
				nodes: map[utils.MapCoordinate]map[utils.MapCoordinate]*rtsspb.AbstractNode{
					utils.MC(clusterCoordinate): {
						utils.MC(tileCoordinate): {
							TileCoordinate: tileCoordinate,
						},
					},
				},
			}, n: &rtsspb.AbstractNode{
				Level: 1, TileCoordinate: tileCoordinate},
		},
	}
	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if err := c.m.Add(c.n); err == nil {
				t.Fatal("Add() = nil, want a non-nil error")
			}
		})
	}
}

func TestMapRemove(t *testing.T) {
	want := &rtsspb.AbstractNode{
		TileCoordinate: &rtsspb.Coordinate{X: 0, Y: 0},
	}

	m := Map{ClusterMap: trivialClusterMap}

	if err := m.Add(want); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got, err := m.Pop(utils.MC(want.GetTileCoordinate())); err != nil || !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("Pop() = %v, %v, want = %v, nil", got, err, want)
	}
}
