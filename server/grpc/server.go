package server

import (
	"context"
	"fmt"

	"github.com/downflux/game/server/service/executor"
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

func (s *DownFluxServer) validateClient(cid string) (<-chan *apipb.StreamCurvesResponse, error) {
	ch := s.ex.ClientChannel(cid)
	if ch == nil {
		return nil, status.Errorf(codes.NotFound, "client %v not found", cid)
	}
	return ch, nil
}

func (s *DownFluxServer) Move(ctx context.Context, req *apipb.MoveRequest) (*apipb.MoveResponse, error) {
	if _, err := s.validateClient(req.GetClientId()); err != nil {
		return nil, err
	}
	return nil, executor.AddMoveCommands(s.ex, req)
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

func (s *DownFluxServer) StreamCurves(req *apipb.StreamCurvesRequest, stream apipb.DownFlux_StreamCurvesServer) error {
	ch, err := s.validateClient(req.GetClientId())
	if err != nil {
		return err
	}

	for r := range ch {
		fmt.Println("sending: %v", r)
		if err := stream.Send(r); err != nil {
			return err
		}
	}
	return nil
}
