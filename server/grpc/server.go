package server

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/downflux/game/server/service/executor"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

type ServerWrapper struct {
	gRPCServer     *grpc.Server
	gRPCServerImpl *DownFluxServer
	wg             errgroup.Group
}

func NewServerWrapper(
	serverOptions []grpc.ServerOption,
	pb *mdpb.TileMap,
	d *gdpb.Coordinate) (*ServerWrapper, error) {
	sw := &ServerWrapper{}

	gRPCServerImpl, err := NewDownFluxServer(pb, d)
	if err != nil {
		return nil, err
	}

	sw.gRPCServerImpl = gRPCServerImpl
	sw.gRPCServer = grpc.NewServer(serverOptions...)
	apipb.RegisterDownFluxServer(sw.gRPCServer, sw.gRPCServerImpl)

	return sw, nil
}

func (s *ServerWrapper) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.wg.Go(func() error { return s.gRPCServer.Serve(lis) })
	s.wg.Go(s.gRPCServerImpl.ex.Run)

	for isStarted := false; !isStarted; isStarted = s.gRPCServerImpl.ex.Status().GetIsStarted() {
		time.Sleep(time.Second)
	}

	return nil
}

func (s *ServerWrapper) Stop() error {
	s.gRPCServerImpl.ex.Stop()
	s.gRPCServer.GracefulStop()

	return s.wg.Wait()
}

func NewDownFluxServer(pb *mdpb.TileMap, d *gdpb.Coordinate) (*DownFluxServer, error) {
	ex, err := executor.New(pb, d)
	if err != nil {
		return nil, err
	}
	return &DownFluxServer{
		ex: ex,
	}, nil
}

type DownFluxServer struct {
	ex *executor.Executor
}

// Executor returns the internal executor.Executor instance. This is a debug
// function.
//
// TODO(minkezhang): Delete this function.
func (s *DownFluxServer) Executor() *executor.Executor { return s.ex }

func (s *DownFluxServer) validateClient(cid string) (<-chan *apipb.StreamDataResponse, error) {
	ch := s.ex.ClientChannel(cid)
	if ch == nil {
		return nil, status.Errorf(codes.NotFound, "client %v not found", cid)
	}
	return ch, nil
}

func (s *DownFluxServer) GetStatus(ctx context.Context, req *apipb.GetStatusRequest) (*apipb.GetStatusResponse, error) {
	return &apipb.GetStatusResponse{
		Status: s.ex.Status(),
	}, nil
}

func (s *DownFluxServer) Move(ctx context.Context, req *apipb.MoveRequest) (*apipb.MoveResponse, error) {
	if _, err := s.validateClient(req.GetClientId()); err != nil {
		return nil, err
	}
	return &apipb.MoveResponse{}, s.ex.AddMoveCommands(req)
}

func (s *DownFluxServer) AddClient(ctx context.Context, req *apipb.AddClientRequest) (*apipb.AddClientResponse, error) {
	cid, err := s.ex.AddClient()
	if err != nil {
		return nil, err
	}

	resp := &apipb.AddClientResponse{
		ClientId: cid,
	}
	return resp, nil
}

func (s *DownFluxServer) StreamData(req *apipb.StreamDataRequest, stream apipb.DownFlux_StreamDataServer) error {
	// TODO(minkezhang): Make this a loop -- we guarantee Executor will
	// always be able to send, blocking.
	//
	// We detect timeouts and disconnects in the server via gRPC keepalive
	// operations and StatsHandler (https://stackoverflow.com/q/62654489).
	defer log.Println("Exiting streaming loop, marking client as dirty.")
	_, err := s.validateClient(req.GetClientId())
	if err != nil {
		return err
	}

	for {
		// TODO(minkezhang): Add timeout on send.
		if err := stream.Send(&apipb.StreamDataResponse{}); err != nil {
			log.Println("StreamData: sending message resulted in error", err)
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	/*
	for r := range ch {
		if err := stream.Send(r); err != nil {
			log.Println("StreamData: sending message resulted in error", err)
			return err
		}
	} */
	return nil
}

// Failure modes to consider -- NOT latency
// Server cannot send for N seconds -- close stream, mark client as dirty
// IGNORE Server sends in (N - 1) seconds and needs to resync client
// 	TODO(minkezhang): Skip messages in queue < current tick.
//	This is actually okay -- we'll need to resync because messages are always deltas.
// 	unless we were to merge all skipped messages.
// Server sends normally (happy path)
