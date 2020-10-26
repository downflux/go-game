package server

import (
	"context"
	"io"
	"log"
	"net"
	"sync"
	"testing"

	"github.com/downflux/game/entity/entity"
	"github.com/downflux/game/server/id"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"

	apipb "github.com/downflux/game/api/api_go_proto"
	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	mcpb "github.com/downflux/game/map/api/constants_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

const (
	bufSize = 1024 * 1024
	idLen   = 8
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

type sut struct {
	gRPCServer     *grpc.Server
	gRPCServerImpl *DownFluxServer
	listener       *bufconn.Listener
	ctx            context.Context
}

func newConn(s *sut) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		s.ctx,
		"bufnet",
		grpc.WithContextDialer(
			func(ctx context.Context, c string) (net.Conn, error) {
				return s.listener.Dial()
			},
		),
		grpc.WithInsecure(),
	)
}

func newSUT() (*sut, error) {
	gRPCServer := grpc.NewServer()
	gRPCServerImpl, err := NewDownFluxServer(simpleLinearMapProto, &gdpb.Coordinate{X: 2, Y: 1})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not create SUT: %v", err)
	}
	apipb.RegisterDownFluxServer(gRPCServer, gRPCServerImpl)
	listener := bufconn.Listen(bufSize)

	return &sut{
		gRPCServer:     gRPCServer,
		gRPCServerImpl: gRPCServerImpl,
		listener:       listener,
		ctx:            context.Background(),
	}, nil
}

func TestSendMoveCommand(t *testing.T) {
	const expectedStreamMessageLength = 1
	dest := &gdpb.Position{X: 3, Y: 0}
	want := &apipb.StreamCurvesResponse{
		Curves: []*gdpb.Curve{
			&gdpb.Curve{
				Type: gcpb.CurveType_CURVE_TYPE_LINEAR_MOVE,
				Data: []*gdpb.CurveDatum{
					&gdpb.CurveDatum{
						Datum: &gdpb.CurveDatum_PositionDatum{
							&gdpb.Position{X: 0, Y: 0},
						},
					},
					&gdpb.CurveDatum{
						Datum: &gdpb.CurveDatum_PositionDatum{
							&gdpb.Position{X: 1, Y: 0},
						},
					},
					&gdpb.CurveDatum{
						Datum: &gdpb.CurveDatum_PositionDatum{
							&gdpb.Position{X: 2, Y: 0},
						},
					},
					&gdpb.CurveDatum{
						Datum: &gdpb.CurveDatum_PositionDatum{
							&gdpb.Position{X: 3, Y: 0},
						},
					},
				},
			},
		},
	}
	e := entity.NewSimpleEntity(id.RandomString(idLen), 0, &gdpb.Position{X: 0, Y: 0})

	s, err := newSUT()
	if err != nil {
		t.Fatalf("newSut() = _, %v, want = nil", err)
	}
	conn, err := newConn(s)
	if err != nil {
		t.Fatalf("newConn() = _, %v, want = nil", err)
	}
	defer conn.Close()

	// TODO(minkezhang): This is a hack -- clients should get the entities
	// via broadcast.
	s.gRPCServerImpl.ex.AddEntity(e)

	var eg errgroup.Group
	eg.Go(func() error { return s.gRPCServer.Serve(s.listener) })
	eg.Go(func() error { return s.gRPCServerImpl.Executor().Run() })

	client := apipb.NewDownFluxClient(conn)
	resp, err := client.AddClient(s.ctx, &apipb.AddClientRequest{})
	if err != nil {
		t.Fatalf("AddPlayer() = _, %v, want = nil", err)
	}

	cid := resp.GetClientId()

	stream, err := client.StreamCurves(s.ctx, &apipb.StreamCurvesRequest{
		ClientId: cid,
	})
	if err != nil {
		t.Fatalf("StreamCurves() = _, %v, want = nil", err)
	}

	var streamResp []*apipb.StreamCurvesResponse
	var streamRespMux sync.Mutex

	eg.Go(func() error {
		for {
			m, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}

			streamRespMux.Lock()
			streamResp = append(streamResp, m)
			streamRespMux.Unlock()
		}
		return nil
	})

	var serverReady bool
	var tick float64
	for !serverReady {
		s, err := client.GetStatus(s.ctx, &apipb.GetStatusRequest{})
		if err != nil {
			t.Fatalf("GetStatus() = _, %v, want = nil", err)
		}
		serverReady = s.GetStatus().GetIsStarted()
		tick = s.GetStatus().GetTick()
	}

	if _, err := client.Move(s.ctx, &apipb.MoveRequest{
		ClientId:    cid,
		EntityIds:   []string{e.ID()},
		Tick:        tick,
		Destination: dest,
		MoveType:    gcpb.MoveType_MOVE_TYPE_FORWARD,
	}); err != nil {
		t.Fatalf("Move() = _, %v, want = nil", err)
	}

	var nMessages int
	for nMessages < expectedStreamMessageLength {
		streamRespMux.Lock()
		nMessages = len(streamResp)
		streamRespMux.Unlock()
	}

	s.gRPCServerImpl.Executor().Stop()
	s.gRPCServer.GracefulStop()

	if err := eg.Wait(); err != nil {
		t.Fatalf("StreamCurvesResponse() = %v, want = nil", err)
	}

	streamRespMux.Lock()
	defer streamRespMux.Unlock()

	if diff := cmp.Diff(
		want,
		streamResp[0],
		protocmp.Transform(),
		protocmp.IgnoreFields(&apipb.StreamCurvesResponse{}, "tick"),
		protocmp.IgnoreFields(&gdpb.Curve{}, "entity_id", "curve_id"),
		protocmp.IgnoreFields(&gdpb.CurveDatum{}, "tick"),
	); diff != "" {
		t.Errorf("StreamCurvesResponse() mismatch (-want +got):\n%v", diff)
	}
}

func TestAddClient(t *testing.T) {
	s, err := newSUT()
	if err != nil {
		t.Fatalf("newSut() = _, %v, want = nil", err)
	}
	conn, err := newConn(s)
	if err != nil {
		t.Fatalf("newConn() = _, %v, want = nil", err)
	}
	defer conn.Close()
	var eg errgroup.Group
	eg.Go(func() error { return s.gRPCServer.Serve(s.listener) })

	client := apipb.NewDownFluxClient(conn)
	resp, err := client.AddClient(s.ctx, &apipb.AddClientRequest{})
	if err != nil {
		t.Fatalf("AddPlayer() = _, %v, want = nil", err)
	}

	if resp.GetClientId() == "" {
		t.Fatalf("GetClientId() = %v, want a non-empty value", err)
	}

	s.gRPCServer.GracefulStop()
	if err := eg.Wait(); err != nil {
		t.Errorf("Wait() = %v, want = nil", err)
	}
}
