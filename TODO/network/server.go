package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	/*
	   "google.golang.org/grpc/codes"
	   "google.golang.org/grpc/keepalive"
	   "google.golang.org/grpc/status"
	*/

	dpb "github.com/downflux/game/TODO/network/demo_go_proto"
)

const (
	addr                   = "localhost:4444"
	serverKeepAliveTime    = time.Second
	serverKeepAliveTimeout = 30 * time.Second
)

var (
	serverOpts = []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             1 * time.Second, // serverKeepAliveTime,
				PermitWithoutStream: true,
			},
		),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    serverKeepAliveTime,
				Timeout: serverKeepAliveTimeout,
			},
		),
	}
)

type DemoServer struct{}

func (s *DemoServer) Single(ctx context.Context, req *dpb.SingleRequest) (*dpb.SingleResponse, error) {
	return &dpb.SingleResponse{}, nil
}

func (s *DemoServer) StreamData(req *dpb.StreamDataRequest, stream dpb.NetworkDemo_StreamDataServer) error {
	log.Println("recieved client StreamData request")
	for {
		log.Println("sending client data")
		if err := stream.Send(&dpb.StreamDataResponse{}); err != nil {
			log.Print(err)
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	demoServer := &DemoServer{}
	s := grpc.NewServer(serverOpts...)
	dpb.RegisterNetworkDemoServer(s, demoServer)
	s.Serve(lis)
}
