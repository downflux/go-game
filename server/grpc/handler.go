package handler

import (
	"context"
	"log"

	"google.golang.org/grpc/stats"
)

type DownFluxHandler struct {
}

func (h *DownFluxHandler) TagRPC(ctx context.Context, s *stats.RPCTagInfo) context.Context { return ctx }
func (h *DownFluxHandler) HandleRPC(context.Context, stats.RPCStats) {}
func (h *DownFluxHandler) TagConn(ctx context.Context, s *stats.ConnTagInfo) context.Context { return ctx }
func (h *DownFluxHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	switch s.(type) {
		case *stats.ConnEnd:
			log.Println("handler detected connection break")
	}
}
