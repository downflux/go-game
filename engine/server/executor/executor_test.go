package executor

import (
	"sync"
	"testing"
	"time"

	"github.com/downflux/game/engine/gamestate/dirty"
	"github.com/downflux/game/engine/gamestate/gamestate"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/visitor/visitor"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/entity/component/moveable"
	"github.com/downflux/game/server/visitor/move"
	"github.com/downflux/game/server/visitor/produce"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/testing/protocmp"

	apipb "github.com/downflux/game/api/api_go_proto"
	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	entitylist "github.com/downflux/game/engine/entity/list"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
	serverstatus "github.com/downflux/game/engine/status/status"
	vcpb "github.com/downflux/game/engine/visitor/api/constants_go_proto"
	visitorlist "github.com/downflux/game/engine/visitor/list"
	mcpb "github.com/downflux/game/map/api/constants_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	moveaction "github.com/downflux/game/server/fsm/move"
	produceaction "github.com/downflux/game/server/fsm/produce"
)

const (
	minPathLength = 8
)

var (
	tickDuration = 100 * time.Millisecond

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

func newTestExecutor(pb *mdpb.TileMap, d *gdpb.Coordinate) (*Executor, error) {
	tm, err := tile.ImportMap(pb)
	if err != nil {
		return nil, err
	}
	g, err := graph.BuildGraph(tm, d)
	if err != nil {
		return nil, err
	}

	dirties := dirty.New()

	state := gamestate.New(serverstatus.New(tickDuration), entitylist.New())

	visitors, err := visitorlist.New([]visitor.Visitor{
		produce.New(state.Status(), state.Entities(), dirties),
		move.New(tm, g, state.Status(), dirties, minPathLength),
	})
	if err != nil {
		return nil, err
	}

	return New(visitors, state, dirties, map[vcpb.VisitorType]fcpb.FSMType{
		vcpb.VisitorType_VISITOR_TYPE_MOVE:    fcpb.FSMType_FSM_TYPE_MOVE,
		vcpb.VisitorType_VISITOR_TYPE_PRODUCE: fcpb.FSMType_FSM_TYPE_PRODUCE,
	})
}

func TestNewExecutor(t *testing.T) {
	_, err := newTestExecutor(simpleLinearMapProto, &gdpb.Coordinate{X: 2, Y: 1})
	if err != nil {
		t.Errorf("newTestExecutor() = _, %v, want = nil", err)
	}
}

func TestAddEntity(t *testing.T) {
	const (
		t0       = float64(0)
		nClients = 1000
	)
	src := &gdpb.Position{X: 0, Y: 0}

	want := &apipb.StreamDataResponse{
		Tick: t0 + 1,
		State: &gdpb.GameState{
			Entities: []*gdpb.Entity{{Type: gcpb.EntityType_ENTITY_TYPE_TANK}},
			Curves: []*gdpb.Curve{{
				Type:     gcpb.CurveType_CURVE_TYPE_LINEAR_MOVE,
				Property: gcpb.EntityProperty_ENTITY_PROPERTY_POSITION,
				Data: []*gdpb.CurveDatum{
					{Datum: &gdpb.CurveDatum_PositionDatum{&gdpb.Position{X: 0, Y: 0}}},
				},
				Tick: t0 + 1,
			},
			}},
	}

	e, err := newTestExecutor(simpleLinearMapProto, &gdpb.Coordinate{X: 2, Y: 1})
	if err != nil {
		t.Fatalf("newTestExecutor() = _, %v, want = nil", err)
	}
	if err := e.Schedule(produceaction.New(
		e.gamestate.Status(),
		e.gamestate.Status().Tick(),
		gcpb.EntityType_ENTITY_TYPE_TANK,
		src,
	)); err != nil {
		t.Fatalf("Schedule() = %v, want = nil", err)
	}

	var cids []id.ClientID
	for i := 0; i < nClients; i++ {
		cid, err := e.AddClient()
		if err != nil {
			t.Fatalf("AddClient() = _, %v, want = nil", err)
		}
		cids = append(cids, cid)
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
			protocmp.IgnoreFields(&gdpb.CurveDatum{}, "tick"),
			protocmp.IgnoreFields(&gdpb.Curve{}, "entity_id"),
			protocmp.IgnoreFields(&gdpb.Entity{}, "entity_id"),
		); diff != "" {
			t.Errorf("<-e.ClientChannel() mismatch (-want +got):\n%v", diff)
		}
	}
}

func TestDoMove(t *testing.T) {
	const (
		t0       = float64(0)
		nClients = 1000
	)
	dest := &gdpb.Position{X: 3, Y: 0}
	src := &gdpb.Position{X: 0, Y: 0}

	e, err := newTestExecutor(simpleLinearMapProto, &gdpb.Coordinate{X: 2, Y: 1})
	if err != nil {
		t.Fatalf("newTestExecutor() = _, %v, want = nil", err)
	}

	if err := e.Schedule(
		produceaction.New(
			e.gamestate.Status(),
			e.gamestate.Status().Tick(),
			gcpb.EntityType_ENTITY_TYPE_TANK,
			src,
		)); err != nil {
		t.Fatalf("Schedule() = %v, want = nil", err)
	}

	var eg errgroup.Group
	var streamResponsesMux sync.Mutex
	var streamResponses []*apipb.StreamDataResponse
	var cids []id.ClientID
	chs := map[id.ClientID]<-chan *apipb.StreamDataResponse{}
	for i := 0; i < nClients; i++ {
		// Add client -- emulate AddClient gRPC call.
		cid, err := e.AddClient()
		if err != nil {
			t.Fatalf("AddClient() = _, %v, want = nil", err)
		}
		cids = append(cids, cid)

		// As in the StreamData gRPC endpoint, first mark the client as
		// reachable.
		if err := e.StartClientStream(cids[i]); err != nil {
			t.Fatalf("StartClientStream() = %v, want = nil", err)
		}

		ch, err := e.ClientChannel(cid)
		if err != nil {
			t.Fatalf("ClientChannel() = _, %v, want = _, nil", err)
		}

		chs[cid] = ch
	}

	// Listen for the initial tick execution, which will add the scheduled
	// entity.
	for _, ch := range chs {
		ch := ch
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

	if gotLen := len(streamResponses); gotLen != nClients {
		t.Fatalf("len() = %v, want = %v", gotLen, nClients)
	}

	eid := streamResponses[0].GetState().GetEntities()[0].GetEntityId()
	streamResponses = nil

	want := &apipb.StreamDataResponse{
		Tick: t0 + 2,
		State: &gdpb.GameState{
			Curves: []*gdpb.Curve{{
				EntityId: eid,
				Type:     gcpb.CurveType_CURVE_TYPE_LINEAR_MOVE,
				Property: gcpb.EntityProperty_ENTITY_PROPERTY_POSITION,
				Data: []*gdpb.CurveDatum{
					// First element is the current position of
					// the entity. This is necessary for the client
					// to do a smooth interpolation.
					{Datum: &gdpb.CurveDatum_PositionDatum{&gdpb.Position{X: 0, Y: 0}}},

					// Following elements relate to the actual tile
					// coordinates for the path.
					{Datum: &gdpb.CurveDatum_PositionDatum{&gdpb.Position{X: 0, Y: 0}}},
					{Datum: &gdpb.CurveDatum_PositionDatum{&gdpb.Position{X: 1, Y: 0}}},
					{Datum: &gdpb.CurveDatum_PositionDatum{&gdpb.Position{X: 2, Y: 0}}},
					{Datum: &gdpb.CurveDatum_PositionDatum{&gdpb.Position{X: 3, Y: 0}}},
				},
				Tick: t0 + 2,
			},
			}},
	}

	if err := e.Schedule(moveaction.New(
		e.gamestate.Entities().Get(id.EntityID(eid)).(moveable.Component),
		e.gamestate.Status(),
		dest,
	)); err != nil {
		t.Fatalf("Schedule() = %v, want = nil", err)
	}

	// Listen for the move command broadcast.
	for _, ch := range chs {
		ch := ch
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
			protocmp.IgnoreFields(&gdpb.CurveDatum{}, "tick"),
		); diff != "" {
			t.Errorf("<-e.ClientChannel() mismatch (-want +got):\n%v", diff)
		}
	}
}
