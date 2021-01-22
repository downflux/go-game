// Package server implements the apipb.DownFluxServer server API.
//
// TODO(minkezhang): Add tests for connection flakiness handling:
//
// 1. if the server cannot send messages in N seconds, the server should close
//    the client stream, and mark the client as dirty.
package server

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/server/grpc/client"
	"github.com/downflux/game/server/grpc/executorutils"
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
	d *gdpb.Coordinate,
	tickDuration time.Duration,
	minPathLength int) (*ServerWrapper, error) {
	sw := &ServerWrapper{}

	gRPCServerImpl, err := NewDownFluxServer(pb, d, tickDuration, minPathLength)
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
	s.eg.Go(s.gRPCServerImpl.utils.Executor().Run)

	for isStarted := false; !isStarted; isStarted = s.gRPCServerImpl.utils.Executor().Status().GetIsStarted() {
		time.Sleep(time.Second)
	}

	return nil
}

func (s *ServerWrapper) Stop() error {
	if err := s.gRPCServerImpl.utils.Executor().Stop(); err != nil {
		return err
	}
	s.gRPCServer.GracefulStop()

	return s.eg.Wait()
}

func NewDownFluxServer(pb *mdpb.TileMap, d *gdpb.Coordinate, tickDuration time.Duration, minPathLength int) (*DownFluxServer, error) {
	utils, err := executorutils.New(pb, d, tickDuration, minPathLength)
	if err != nil {
		return nil, err
	}
	return &DownFluxServer{
		utils: utils,
	}, nil
}

type DownFluxServer struct {
	utils *executorutils.Utils
}

func (s *DownFluxServer) validateClient(cid id.ClientID) error {
	if !s.utils.Executor().ClientExists(cid) {
		return status.Errorf(codes.NotFound, "client %v not found", cid)
	}
	return nil
}

func (s *DownFluxServer) Utils() *executorutils.Utils { return s.utils }

func (s *DownFluxServer) GetStatus(ctx context.Context, req *apipb.GetStatusRequest) (*apipb.GetStatusResponse, error) {
	return &apipb.GetStatusResponse{
		Status: s.utils.Executor().Status(),
	}, nil
}

func (s *DownFluxServer) Attack(ctx context.Context, req *apipb.AttackRequest) (*apipb.AttackResponse, error) {
	if err := s.validateClient(id.ClientID(req.GetClientId())); err != nil {
		return nil, err
	}
	return &apipb.AttackResponse{}, s.utils.Attack(req)
}

func (s *DownFluxServer) Move(ctx context.Context, req *apipb.MoveRequest) (*apipb.MoveResponse, error) {
	if err := s.validateClient(id.ClientID(req.GetClientId())); err != nil {
		return nil, err
	}
	return &apipb.MoveResponse{}, s.utils.Move(req)
}

func (s *DownFluxServer) AddClient(ctx context.Context, req *apipb.AddClientRequest) (*apipb.AddClientResponse, error) {
	log.Println("new Client request")
	cid, err := s.utils.Executor().AddClient()
	if err != nil {
		return nil, err
	}

	resp := &apipb.AddClientResponse{
		ClientId: cid.Value(),
	}
	return resp, nil
}

func (s *DownFluxServer) StreamData(req *apipb.StreamDataRequest, stream apipb.DownFlux_StreamDataServer) error {
	log.Println("new StreamData request")
	cid := id.ClientID(req.GetClientId())

	if err := s.validateClient(cid); err != nil {
		return err
	}

	md := client.New()
	defer func() {
		s.utils.Executor().StopClientStreamError(cid)
		md.Close()
		log.Println("closing StreamData request")
	}()

	if err := s.utils.Executor().StartClientStream(cid); err != nil {
		return err
	}

	ch, err := s.utils.Executor().ClientChannel(cid)
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
