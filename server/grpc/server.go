package server

import (
	"context"

	"github.com/downflux/game/server/service/executor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

func NewDownFluxServer() (*DownFluxServer, error) {
	ex, err := executor.New(nil, nil)
	if err != nil {
		return nil, err
	}
	return &DownFluxServer{
		ex: ex,
	}, notImplemented
}

type DownFluxServer struct {
	ex *executor.Executor
}

func (s *DownFluxServer) Move(ctx context.Context, req *apipb.MoveRequest) (*apipb.MoveResponse, error) {
	return nil, executor.AddCommand(s.ex, executor.NewMoveCommand(s.ex, req))
}

func (s *DownFluxServer) StreamCurves(ctx context.Context, req *apipb.StreamCurvesRequest, stream apipb.DownFlux_StreamCurvesServer) error {
	return notImplemented
}
