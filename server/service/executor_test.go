package executor

import (
	"sync"
	"testing"

	"github.com/downflux/game/server/service/visitor/entity/tank"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/testing/protocmp"

	apipb "github.com/downflux/game/api/api_go_proto"
	gcpb "github.com/downflux/game/api/constants_go_proto"
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

// TODO(minkezhang): Test sending Move request on a stale tick -- what should
// actually occur in the response?

func TestDoTick(t *testing.T) {
	const (
		eid      = "entity-id"
		t1       = float64(0)
		nClients = 1000
	)
	dest := &gdpb.Position{X: 3, Y: 0}
	src := &gdpb.Position{X: 0, Y: 0}

	want := &apipb.StreamDataResponse{
		Tick:     0,
		Entities: []*gdpb.Entity{{EntityId: eid, Type: gcpb.EntityType_ENTITY_TYPE_TANK}},
		Curves: []*gdpb.Curve{{
			EntityId: eid,
			Type:     gcpb.CurveType_CURVE_TYPE_LINEAR_MOVE,
			Category: gcpb.CurveCategory_CURVE_CATEGORY_MOVE,
			Data: []*gdpb.CurveDatum{
				{Datum: &gdpb.CurveDatum_PositionDatum{&gdpb.Position{X: 0, Y: 0}}},
				{Datum: &gdpb.CurveDatum_PositionDatum{&gdpb.Position{X: 1, Y: 0}}},
				{Datum: &gdpb.CurveDatum_PositionDatum{&gdpb.Position{X: 2, Y: 0}}},
				{Datum: &gdpb.CurveDatum_PositionDatum{&gdpb.Position{X: 3, Y: 0}}},
			},
		}},
	}

	e, err := New(simpleLinearMapProto, &gdpb.Coordinate{X: 2, Y: 1})
	if err != nil {
		t.Fatalf("New() = _, %v, want = nil", err)
	}

	if err := e.AddEntity(tank.New(eid, t1, src)); err != nil {
		t.Fatalf("AddEntity() = %v, want = nil", err)
	}

	var cids []string
	for i := 0; i < nClients; i++ {
		cid, err := e.AddClient()
		if err != nil {
			t.Fatalf("AddClient() = _, %v, want = nil", err)
		}
		cids = append(cids, cid)
	}

	if err := e.AddMoveCommands(&apipb.MoveRequest{
		Tick:        t1,
		ClientId:    cids[0],
		EntityIds:   []string{eid},
		Destination: dest,
		MoveType:    gcpb.MoveType_MOVE_TYPE_FORWARD,
	}); err != nil {
		t.Fatalf("AddMoveCommands() = %v, want = nil", err)
	}

	// Connect to server and signal intent to start listening for messages.
	for i := 0; i < nClients; i++ {
		if err := e.StartClientStream(cids[i]); err != nil {
			t.Fatalf("StartClientStream() = %v, want = nil", err)
		}
	}

	var eg errgroup.Group

	var streamResponsesMux sync.Mutex
	var streamResponses []*apipb.StreamDataResponse
	for i := 0; i < nClients; i++ {
		ch, err := e.ClientChannel(cids[i])
		if err != nil {
			t.Fatalf("ClientChannel() = _, %v, want = _, nil", err)
		}
		// Assuming all clients will receive messages in a timely
		// manner. Start listening for messages before the tick starts
		// to guarantee we will recieve a message during
		// broadcastCurves.
		eg.Go(func() error {
			m := <-ch

			streamResponsesMux.Lock()
			defer streamResponsesMux.Unlock()

			streamResponses = append(streamResponses, m)
			return nil
		})
	}

	eg.Go(e.doTick)

	if err := eg.Wait(); err != nil {
		t.Fatalf("Wait() = %v, want = nil", err)
	}

	for _, streamResponse := range streamResponses {
		if diff := cmp.Diff(
			want,
			streamResponse,
			protocmp.Transform(),
			protocmp.IgnoreFields(&apipb.StreamDataResponse{}, "tick"),
			protocmp.IgnoreFields(&gdpb.CurveDatum{}, "tick"),
		); diff != "" {
			t.Errorf("<-e.ClientChannel() mismatch (-want +got):\n%v", diff)
		}
	}
}

func TestAddEntity(t *testing.T) {
	e, err := New(simpleLinearMapProto, &gdpb.Coordinate{X: 2, Y: 1})
	if err != nil {
		t.Fatalf("New() = _, %v, want = nil", err)
	}

	if err := e.AddEntity(tank.New("simple", 100, &gdpb.Position{X: 0, Y: 0})); err != nil {
		t.Fatalf("AddEntity() = %v, want = nil", err)
	}

	if err := e.AddEntity(tank.New("simple", 0, nil)); err == nil {
		t.Error("AddEntity() = nil, want a non-nil error")
	}
}
