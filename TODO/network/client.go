package main

import (
	"context"
	"io"
	"log"

        "google.golang.org/grpc"

	dpb "github.com/downflux/game/TODO/network/demo_go_proto"
)

const (
	addr = "localhost:4445"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) // keepalive here
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()  // ?

	c := dpb.NewNetworkDemoClient(conn)
	log.Println("establishing server stream")
	stream, err := c.StreamData(context.Background(), &dpb.StreamDataRequest{})
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
