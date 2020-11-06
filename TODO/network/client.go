package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/toxiproxy"
	toxiproxyclient "github.com/Shopify/toxiproxy/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	dpb "github.com/downflux/game/TODO/network/demo_go_proto"
)

const (
	addr                   = "localhost:4444"
	clientKeepAliveTime    = 10 * time.Second
	clientKeepAliveTimeout = 1 * time.Second
)

var (
	dialOpts = grpc.WithKeepaliveParams(
		keepalive.ClientParameters{
			Time:                clientKeepAliveTime,
			Timeout:             clientKeepAliveTimeout,
			PermitWithoutStream: true,
		},
	)
)

func main() {
	tox := toxiproxy.NewServer()
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM)
	go func() {
		<-signals
		os.Exit(0)
	}()

	go func() {
		tox.Listen("localhost", "50001")
		select {}
	}()

	proxycli := toxiproxyclient.NewClient("localhost:50001")
	p, err := proxycli.CreateProxy("downflux", "localhost:50000", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure(), dialOpts) // keepalive here
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close() // ?

	c := dpb.NewNetworkDemoClient(conn)
	stream, err := c.StreamData(context.Background(), &dpb.StreamDataRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		r, err := stream.Recv()
		p.AddToxic("latency_downstream", "latency", "downstream", 1.0, toxiproxyclient.Attributes{
			"latency": 10000,
		})
		defer p.Delete()
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
