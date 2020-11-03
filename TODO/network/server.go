package main

import (
	"context"
	"log"
	"net"

        "google.golang.org/grpc"
	/*
        "google.golang.org/grpc/codes"
        "google.golang.org/grpc/keepalive"
        "google.golang.org/grpc/status"
	 */

	dpb "github.com/downflux/game/TODO/network/demo_go_proto"
)

const (
	addr = "localhost:4444"
)

type DemoServer struct {}

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
		return nil
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	demoServer := &DemoServer{}
	s := grpc.NewServer()  // ...
	dpb.RegisterNetworkDemoServer(s, demoServer)
	s.Serve(lis)
}
