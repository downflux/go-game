package server

import (
	"context"

	"github.com/downflux/game/server/service/commands/move"
	"github.com/downflux/game/server/service/executor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

func NewDownFluxService() *DownFluxService {
	return &DownFluxService{
		ex: executor.New(),
	}
}

type DownFluxService struct {
	ex *executor.Executor
}

func (s *DownFluxService) Move(ctx context.Context, req *apipb.MoveRequest) (*apipb.MoveResponse, error) {
	return nil, executor.AddCommand(s.ex, move.Import(req))
}
