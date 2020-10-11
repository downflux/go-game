package node

import (
	"testing"

	gdpb "github.com/downflux/game/api/data_go_proto"
	pdpb "github.com/downflux/game/pathing/api/data_go_proto"

	"github.com/downflux/game/pathing/hpf/cluster"
	"github.com/downflux/game/pathing/hpf/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	trivialClusterMap, _ = cluster.ImportMap(
		&pdpb.ClusterMap{
			TileDimension:    &gdpb.Coordinate{X: 1, Y: 1},
			TileMapDimension: &gdpb.Coordinate{X: 1, Y: 1},
		})
)

func coordinateLess(c1, c2 *gdpb.Coordinate) bool {
	return c1.GetX() < c2.GetX() || (c1.GetX() == c2.GetX() && c1.GetY() < c2.GetY())
}

func abstractNodeLess(n1, n2 *pdpb.AbstractNode) bool {
	return coordinateLess(n1.GetTileCoordinate(), n2.GetTileCoordinate())
}

func TestMapAdd(t *testing.T) {
	want := &pdpb.AbstractNode{
		TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0},
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
	cm, err := cluster.ImportMap(&pdpb.ClusterMap{
		TileDimension:    &gdpb.Coordinate{X: 2, Y: 2},
		TileMapDimension: &gdpb.Coordinate{X: 4, Y: 4},
	})
	if err != nil {
		t.Fatalf("ImportMap() = _, %v, want = _, nil", err)
	}

	m := Map{ClusterMap: cm}
	want := []*pdpb.AbstractNode{
		{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0}},
		{TileCoordinate: &gdpb.Coordinate{X: 0, Y: 1}},
	}

	for _, n := range want {
		if err := m.Add(n); err != nil {
			t.Fatalf("Add() = %v, want = nil", err)
		}
	}

	if err := m.Add(
		&pdpb.AbstractNode{TileCoordinate: &gdpb.Coordinate{X: 3, Y: 3}}); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	got, err := m.GetByCluster(utils.MC(&gdpb.Coordinate{X: 0, Y: 0}))
	if err != nil {
		t.Fatalf("GetByCluster() = _, %v, want = _, nil", err)
	}

	if diff := cmp.Diff(want, got, protocmp.Transform(), cmpopts.SortSlices(abstractNodeLess)); diff != "" {
		t.Errorf("GetByCluster() mismatch (-want +got):\n%s", diff)
	}
}

func TestMapAddError(t *testing.T) {
	clusterCoordinate := &gdpb.Coordinate{X: 0, Y: 0}
	tileCoordinate := &gdpb.Coordinate{X: 0, Y: 0}
	testConfigs := []struct {
		name string
		m    Map
		n    *pdpb.AbstractNode
	}{
		{
			name: "AlreadyExist",
			m: Map{
				ClusterMap: trivialClusterMap,
				nodes: map[utils.MapCoordinate]map[utils.MapCoordinate]*pdpb.AbstractNode{
					utils.MC(clusterCoordinate): {
						utils.MC(tileCoordinate): {
							TileCoordinate: tileCoordinate,
						},
					},
				},
			}, n: &pdpb.AbstractNode{TileCoordinate: tileCoordinate},
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
	want := &pdpb.AbstractNode{
		TileCoordinate: &gdpb.Coordinate{X: 0, Y: 0},
	}

	m := Map{ClusterMap: trivialClusterMap}

	if err := m.Add(want); err != nil {
		t.Fatalf("Add() = %v, want = nil", err)
	}

	if got, err := m.Pop(utils.MC(want.GetTileCoordinate())); err != nil || !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("Pop() = %v, %v, want = %v, nil", got, err, want)
	}
}
