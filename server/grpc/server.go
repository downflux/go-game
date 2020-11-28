package server

import (
	"context"
	"net"
	"time"

	"github.com/downflux/game/server/grpc/client"
	"github.com/downflux/game/server/id"
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
	eg             errgroup.Group
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

	s.eg.Go(func() error { return s.gRPCServer.Serve(lis) })
	s.eg.Go(s.gRPCServerImpl.ex.Run)

	for isStarted := false; !isStarted; isStarted = s.gRPCServerImpl.ex.Status().GetIsStarted() {
		time.Sleep(time.Second)
	}

	return nil
}

func (s *ServerWrapper) Stop() error {
	if err := s.gRPCServerImpl.ex.Stop(); err != nil {
		return err
	}
	s.gRPCServer.GracefulStop()

	return s.eg.Wait()
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

func (s *DownFluxServer) validateClient(cid id.ClientID) error {
	if !s.ex.ClientExists(cid) {
		return status.Errorf(codes.NotFound, "client %v not found", cid)
	}
	return nil
}

// Executor returns the internal executor.Executor instance. This is a debug
// function.
//
// TODO(minkezhang): Delete this function.
func (s *DownFluxServer) Executor() *executor.Executor { return s.ex }

func (s *DownFluxServer) GetStatus(ctx context.Context, req *apipb.GetStatusRequest) (*apipb.GetStatusResponse, error) {
	return &apipb.GetStatusResponse{
		Status: s.ex.Status(),
	}, nil
}

func (s *DownFluxServer) Move(ctx context.Context, req *apipb.MoveRequest) (*apipb.MoveResponse, error) {
	if err := s.validateClient(id.ClientID(req.GetClientId())); err != nil {
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
		ClientId: cid.Value(),
	}
	return resp, nil
}

func (s *DownFluxServer) StreamData(req *apipb.StreamDataRequest, stream apipb.DownFlux_StreamDataServer) error {
	cid := id.ClientID(req.GetClientId())

	if err := s.validateClient(cid); err != nil {
		return err
	}

	md := client.New()
	defer func() {
		s.ex.StopClientStreamError(cid)
		md.Close()
	}()

	if err := s.ex.StartClientStream(cid); err != nil {
		return err
	}

	ch, err := s.ex.ClientChannel(cid)
	if err != nil {
		return err
	}

	go func(md *client.Connection) {
		defer md.SetChannelClosed(true)
		for {
			select {
			case <-md.Done():
				return
			case m, ok := <-ch:
				if !ok {
					return
				}
				md.AddMessage(m)
			}
		}
	}(md)

	for {
		resp, ok := md.Responses()
		if resp == nil && !ok {
			return nil
		}

		for _, m := range resp {
			// Send does not block on flakey network connection. See gRPC
			// docs. On server keepalive failure, StreamData will return
			// with connection error. On client close, StreamData will
			// return with connection error.
			if err := stream.Send(m); err != nil {
				return err
			}
		}
	}
	return nil
}

// Failure modes to consider -- NOT latency
// Server cannot send for N seconds -- close stream, mark client as dirty
// IGNORE Server sends in (N - 1) seconds and needs to resync client
// 	TODO(minkezhang): Skip messages in queue < current tick.
//	This is actually okay -- we'll need to resync because messages are always deltas.
// 	unless we were to merge all skipped messages.
// Server sends normally (happy path)
