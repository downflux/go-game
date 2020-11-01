package server

import (
	"context"
	"time"

	"github.com/downflux/game/server/service/executor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
)

const (
	serverKeepAliveTime = 5 * time.Second
	serverKeepAliveTimeout = 5 * time.Second

	// clientKeepAliveTime is the time between heartbeat pings the client
	// will wait between resending the ping. The minimum interval accepted
	// by the grpc/keepalive package is 10s.
	clientKeepAliveTime = 10 * time.Second
	clientKeepAliveTimeout = 5 * time.Second
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")

	// DefaultServerOptions returns the default options the server will employ
	// for connecting to the client. Notably, these options will allow the server
	// to receive keepalive messages from the client periodically to facilitate
	// detecting network problems early.
	//
	// Example
	//
	// s := grpc.NewServer(DefaultServerOptions...)
	DefaultServerOptions = []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime: serverKeepAliveTime,
				PermitWithoutStream: false,
			},
		),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time: serverKeepAliveTime,
				Timeout: serverKeepAliveTimeout,
			},
		),
	}

	// DefaultClientOptions returns the recommended default client options when
	// connecting to the server. This will mainly be used for client disconnect
	// detection.
	//
	// Example
	//
	// c, err := grpc.Dial("localhost:4444", DefaultClientOptions...)
	DefaultClientOptions = []grpc.DialOption{
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time: clientKeepAliveTime,
				Timeout: clientKeepAliveTimeout,
				PermitWithoutStream: false,
			},
		),
	}

)

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
	ch, err := s.validateClient(req.GetClientId())
	if err != nil {
		return err
	}

	for r := range ch {
		if err := stream.Send(r); err != nil {
			return err
		}
	}
	return nil
}
