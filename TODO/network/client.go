package main

import (
	"context"
	"io"
	"log"
	"time"

        "google.golang.org/grpc"
        "google.golang.org/grpc/keepalive"

	dpb "github.com/downflux/game/TODO/network/demo_go_proto"
)

const (
	addr = "localhost:4445"
	clientKeepAliveTime = 10 * time.Second
	clientKeepAliveTimeout = 1 * time.Second
)

var (
	dialOpts = grpc.WithKeepaliveParams(
                        keepalive.ClientParameters{
                                Time: clientKeepAliveTime,
                                Timeout: clientKeepAliveTimeout,
                                PermitWithoutStream: true,
                        },
                )
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), dialOpts) // keepalive here
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()  // ?

	c := dpb.NewNetworkDemoClient(conn)
	log.Println("establishing server stream")
	stream, err := c.StreamData(context.Background(), &dpb.StreamDataRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		log.Println("calling server stream")
		r, err := stream.Recv()
		log.Println("read from server stream")
		if err != nil {
			if err == io.EOF {
				log.Println("received EOF, exiting normally")
			} else {
				log.Print(err)
			}
			break
		}

		log.Print(r)
	}
}
