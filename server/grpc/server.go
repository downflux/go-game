package server

import (
	"context"
	"log"
	"net"
	"sync"
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
	ch, err := s.validateClient(req.GetClientId())
	if err != nil {
		return err
	}

	done := make(chan struct{})
	defer func() {
		log.Println("Exiting streaming loop, marking client as dirty.")
		close(done)
	}()

	var l sync.Mutex
	var q []*apipb.StreamDataResponse
	var chanClosed bool

	go func() {
		defer func() {
			l.Lock()
			chanClosed = true
			l.Unlock()
		}()
		for {
			log.Println("LISTENING CH")
			select {
				case <-done:
					log.Println("GOT DONE")
					return
				case m, ok := <-ch:
					log.Println("RECV MESSAGE: ", m, ok)
					if !ok {
						return
					}
					l.Lock()
					log.Println("GOROUTINE ACQ L")
					q = append(q, m)
					l.Unlock()
					log.Println("GOROUTINE RELEASE L")
			}
		}
	}()

	for {
		l.Lock()
		tq := q
		q = nil
		ok := !chanClosed
		l.Unlock()

		if tq == nil && !ok {
			return nil
		}

		for _, m := range tq {
			log.Println("sending: ", m)
			// Send does not block on flakey network connection. See gRPC
			// docs. On server keepalive failure, StreamData will return
			// with connection error. On client close, StreamData will
			// return with connection error.
			if err := stream.Send(m); err != nil {
				log.Println("----------------- sent")
				return err
			}
			log.Println("----------------- sent")
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
