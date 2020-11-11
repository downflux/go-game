package server

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Shopify/toxiproxy"
	"github.com/downflux/game/server/grpc/handler"
	"github.com/downflux/game/server/grpc/option"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/connectivity"

	tpc "github.com/Shopify/toxiproxy/client"
	apipb "github.com/downflux/game/api/api_go_proto"
)

const (
	toxiHost = "localhost"
	// TODO(minkezhang): Change to a random test-selected port instead.
	// This will be blocked on writing a custom Listen(s *ApiServer)
	// method.
	toxiPort = "50000"
)

var (
	testGlobal = networkImpairmentProxy{}
)

type networkImpairmentProxy struct {
	toxi       *toxiproxy.ApiServer
	toxiClient *tpc.Client
	usedPorts  map[int32]bool
}

func randomPort() int32 {
	// Recommended ephemeral port range is 49152 - 65535.
	// Pick a random unused one in this range.
	//
	// TODO(minkezhang): Use actually unused ephemeral ports, vs. picking
	// randomly as we are doing here.
	return 49152 + rand.Int31n(65535-49152)
}

func (p networkImpairmentProxy) newAddress() string {
	if p.usedPorts == nil {
		p.usedPorts = map[int32]bool{}
	}
	var port int32
	for port = randomPort(); p.usedPorts[port]; port = randomPort() {
	}
	p.usedPorts[port] = true
	return fmt.Sprintf("%s:%d", toxiHost, port)
}

func setup() {
	testGlobal.toxi = toxiproxy.NewServer()
	go testGlobal.toxi.Listen(toxiHost, toxiPort)

	time.Sleep(time.Second) // Wait for toxiproxy server to be up.
	testGlobal.toxiClient = tpc.NewClient(fmt.Sprintf("%s:%s", toxiHost, toxiPort))
}

func teardown() {
	testGlobal.toxi.Collection.Clear()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)

}

func newGRPCClient(hostAddr string) (*grpc.ClientConn, apipb.DownFluxClient, error) {
	conn, err := grpc.Dial(hostAddr, grpc.WithInsecure()) // TODO(minkezhang): Add DialOpts here.
	if err != nil {
		return nil, nil, err
	}

	return conn, apipb.NewDownFluxClient(conn), nil
}

func TestClientCloseStream(t *testing.T)     {}
func TestServerDetectedTimeout(t *testing.T) {}
func TestServerDetectedLatency(t *testing.T) {
	listenerAddr := testGlobal.newAddress()
	serverAddr := testGlobal.newAddress()
	serverOptionConfig := option.ServerOptionConfig{
		MinimumClientInterval:   10 * time.Second,
		ServerHeartbeatInterval: time.Second,
		ServerHeartbeatTimeout:  time.Second,
	}

	p, err := testGlobal.toxiClient.CreateProxy("downflux", listenerAddr, serverAddr)
	if err != nil {
		t.Fatalf("CreateProxy() = _, %v, want = nil", err)
	}
	defer p.Delete()

	// Create gRPC server.
	sw, err := NewServerWrapper(
		append(
			option.ServerOptions(serverOptionConfig),
			grpc.StatsHandler(&handler.DownFluxHandler{})),
		nil,
		nil)
	if err != nil {
		t.Fatalf("NewServerWrapper() = _, %v, want = nil", err)
	}
	sw.Start(serverAddr)
	defer sw.Stop()

	// Create gRPC client.
	conn, client, err := newGRPCClient(listenerAddr)
	if err != nil {
		t.Fatalf("newGRPCClient() = _, _, %v, want = nil", err)
	}
	defer conn.Close()

	clientResp, err := client.AddClient(context.Background(), &apipb.AddClientRequest{})
	if err != nil {
		t.Fatalf("AddClient() = _, %v, want = nil", err)
	}

	stream, err := client.StreamData(context.Background(), &apipb.StreamDataRequest{
		ClientId: clientResp.GetClientId(),
	})
	if err != nil {
		t.Fatalf("StreamData() = _, %v, want = nil", err)
	}

	p.AddToxic("latency_downstream", "latency", "downstream", 1.0, tpc.Attributes{
		"latency": (10 * serverOptionConfig.ServerHeartbeatTimeout) / time.Millisecond,
	})

	go func() {
		var m *apipb.StreamDataResponse
		var err error
		for m, err = nil, nil; err == nil; m, err = stream.Recv() {
			fmt.Println("received message: ", m)
		}
		fmt.Println("final error:", err)
		// conn.Close()
	}()

	// Register for a conn.WaitForStateChange -- at this point, inspect server and ensure
	// 1. it's still ticking
	// 2. it's disconnected client stream / etc. after N retries
	//   a. (non-transient latency / termination)
	//   b. connection error detection
	// 3. it's ready for reconnect (mark client as dirty)
	for {
		s := conn.GetState()
		fmt.Println(s)
		/*
			if s != connectivity.Ready {
				p.RemoveToxic("latency_downstream")
			}
		*/
		time.Sleep(time.Second)
	}

	t.Error(err)
}
