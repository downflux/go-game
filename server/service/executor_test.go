package executor

import (
	"testing"

	"github.com/downflux/game/entity/entity"
	"github.com/downflux/game/server/service/commands/move"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"

	gdpb "github.com/downflux/game/api/data_go_proto"
	mcpb "github.com/downflux/game/map/api/constants_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
)

var (
	/**
	 * Y = 0 - - - -
	 *   X = 0
	 */
	simpleLinearMapProto = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 4, Y: 1},
		Tiles: []*mdpb.Tile{
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 3, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
	}
)

func TestNewExecutor(t *testing.T) {
	_, err := New(simpleLinearMapProto, &gdpb.Coordinate{X: 2, Y: 1})
	if err != nil {
		t.Errorf("New() = _, %v, want = nil", err)
	}
}

func TestAddEntity(t *testing.T) {
	e, err := New(simpleLinearMapProto, &gdpb.Coordinate{X: 2, Y: 1})
	if err != nil {
		t.Fatalf("New() = _, %v, want = nil", err)
	}

	if err := AddEntity(e, entity.NewSimpleEntity("simple", 100, &gdpb.Position{X: 0, Y: 0})); err != nil {
		t.Fatalf("AddEntity() = %v, want = nil", err)
	}

	if err := AddEntity(e, entity.NewSimpleEntity("simple", 0, nil)); err == nil {
		t.Error("AddEntity() = nil, want a non-nil error")
	}
}

func TestBuildMoveCommands(t *testing.T) {
	testConfigs := []struct {
		name string
		cid  string
		eid  string
		t1   float64
		t2   float64
		p1   *gdpb.Position
		p2   *gdpb.Position
		want []*move.Command
	}{
		{
			name: "SimpleSingleton",
			cid:  "random-client",
			eid:  "some-entity",
			t1:   0,
			t2:   1,
			p1:   &gdpb.Position{X: 0, Y: 0},
			p2:   &gdpb.Position{X: 1, Y: 0},
			want: []*move.Command{
				move.New(nil, nil, "random-client", 1, &gdpb.Position{X: 0, Y: 0}, &gdpb.Position{X: 1, Y: 0}),
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			e, err := New(simpleLinearMapProto, &gdpb.Coordinate{X: 2, Y: 1})
			if err != nil {
				t.Fatalf("New() = _, %v, want = nil", err)
			}

			if err := AddEntity(e, entity.NewSimpleEntity(c.eid, c.t1, c.p1)); err != nil {
				t.Fatalf("AddEntity() = %v, want = nil", err)
			}

			got := buildMoveCommands(e, c.cid, c.t2, c.p2, []string{c.eid})
			if diff := cmp.Diff(
				got,
				c.want,
				cmp.AllowUnexported(move.Command{}),
				cmpopts.IgnoreFields(move.Command{}, "tileMap", "abstractGraph"),
				protocmp.Transform(),
			); diff != "" {
				t.Errorf("buildEntities() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
