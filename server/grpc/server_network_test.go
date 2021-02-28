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
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"

	tpc "github.com/Shopify/toxiproxy/client"
	apipb "github.com/downflux/game/api/api_go_proto"
)

const (
	toxiHost = "localhost"
	// TODO(minkezhang): Change to a random test-selected port instead.
	// This will be blocked on writing a custom Listen(s *ApiServer)
	// method.
	toxiPort = "50000"

	minPathLength = 8
)

var (
	testGlobal = networkImpairmentProxy{}

	tickDuration = 100 * time.Millisecond
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
	rand.Seed(time.Now().UnixNano())
	p := fmt.Sprintf("%d", randomPort())
	testGlobal.toxi = toxiproxy.NewServer()
	go testGlobal.toxi.Listen(toxiHost, p)

	time.Sleep(time.Second) // Wait for toxiproxy server to be up.
	testGlobal.toxiClient = tpc.NewClient(fmt.Sprintf("%s:%s", toxiHost, p))
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

func TestClientCloseStream(t *testing.T) {
	// TODO(minkezhang): Write this test for client detecting a network
	// flake.
}

func TestServerCloseStream(t *testing.T) {
	serverOptionConfig := option.ServerOptionConfig{
		MinimumClientInterval:   10 * time.Second,
		ServerHeartbeatInterval: time.Second,
		ServerHeartbeatTimeout:  time.Second,
	}

	testConfigs := []struct {
		name  string
		toxic tpc.Toxic
		want  codes.Code
	}{
		{
			name: "TestClientHighLatency",
			toxic: tpc.Toxic{
				Name:     "downstream_latency",
				Type:     "latency",
				Stream:   "downstream",
				Toxicity: 1.0,
				Attributes: tpc.Attributes{
					"latency": (4 * serverOptionConfig.ServerHeartbeatTimeout) / time.Millisecond,
				},
			},
			want: codes.Unavailable,
		},
		{
			name: "TestClientTimeout",
			toxic: tpc.Toxic{
				Name:       "downstream_timeout",
				Type:       "timeout",
				Stream:     "downstream",
				Toxicity:   1.0,
				Attributes: tpc.Attributes{"timeout": 0},
			},
			want: codes.Unavailable,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			listenerAddr := testGlobal.newAddress()
			serverAddr := testGlobal.newAddress()

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
				nil,
				tickDuration,
				minPathLength)
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
				ClientId: clientResp.GetClientId().GetClientId(),
			})
			if err != nil {
				t.Fatalf("StreamData() = _, %v, want = nil", err)
			}

			p.AddToxic(c.toxic.Name, c.toxic.Type, c.toxic.Stream, c.toxic.Toxicity, c.toxic.Attributes)

			var eg errgroup.Group
			eg.Go(func() error {
				var m *apipb.StreamDataResponse
				var err error
				for m, err = nil, nil; err == nil; m, err = stream.Recv() {
					fmt.Println("received message: ", m)
				}
				return err
			})

			s, ok := status.FromError(eg.Wait())
			if !ok {
				t.Fatalf("FromError() = _, %v, want = true", ok)
			}

			if s.Code() != c.want {
				t.Errorf("Code() = %v, want = %v", s.Code(), c.want)
			}
		})
	}
}
