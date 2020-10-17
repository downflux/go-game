package server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	// "time"

	"github.com/downflux/game/entity/entity"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/executor"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	apipb "github.com/downflux/game/api/api_go_proto"
	// gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	mcpb "github.com/downflux/game/map/api/constants_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
)

const bufSize = 1024 * 1024

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
		ctx: context.Background(),
	}, nil
}

func TestSendMoveCommand(t *testing.T) {
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
	// eg.Go(func() error { return s.gRPCServer.Serve(s.listener) })
	go s.gRPCServer.Serve(s.listener)

	client := apipb.NewDownFluxClient(conn)
	addClientResp, err := client.AddClient(s.ctx, &apipb.AddClientRequest{})
	if err != nil {
		t.Fatalf("AddPlayer() = _, %v, want = nil", err)
	}

	cid := addClientResp.GetClientId()

	// TODO(minkezhang): This is a hack -- clients should get the entities via broadcast.
	e := entity.NewSimpleEntity(id.RandomString(32), 0, &gdpb.Position{X: 0, Y: 0})
	executor.AddEntity(s.gRPCServerImpl.ex, e)

	stream, err := client.StreamCurves(s.ctx, &apipb.StreamCurvesRequest{
		ClientId: cid,
	})
	if err != nil {
		t.Fatalf("StreamCurves() = _, %v, want = nil", err)
	}

	var streamResp []*apipb.StreamCurvesResponse
	var streamRespMux sync.Mutex

	eg.Go(func() error {
		fmt.Println("listening")
		m, err := stream.Recv()
		fmt.Println("recv", m)
		if err != nil {
			return err
		}
		streamRespMux.Lock()
		defer streamRespMux.Unlock()

		streamResp = append(streamResp, m)
		return nil
	})
	fmt.Println("HI")

	fmt.Println("TICKING")
	executor.Tick(s.gRPCServerImpl.ex)
	fmt.Println("CLOSING STREAMS")
	executor.CloseStreams(s.gRPCServerImpl.ex)

	fmt.Println("eg.WAIT")
	s.gRPCServer.GracefulStop()
	if err := eg.Wait(); err != nil {
		t.Fatalf("StreamCurvesResponse() = %v, want = nil", err)
	}

	fmt.Println(streamResp)

	/*
	nextTickID := ""
	for nextTickID == "" {
		mux.Lock()
		nextTickID = tickID
		mux.Unlock()
		time.Sleep(time.Second)
	}
	 */
	/*
	moveResp, err := client.Move(s.ctx, &apipb.MoveRequest{
		ClientId: cid,
		EntityIds: []string{e.ID()},
		// TODO(minkezhang): Fill out.
		TickId: "",
		Destination: &gdpb.Position{X: 3, Y: 0},
		MoveType: gcpb.MoveType_MOVE_TYPE_FORWARD,
	})
	 */
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
