package option

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	ServerKeepAliveTime    = 1 * time.Second
	ServerKeepAliveTimeout = 1 * time.Second

	// ClientKeepAliveTime is the time between heartbeat pings the client
	// will wait between resending the ping. The minimum interval accepted
	// by the grpc/keepalive package is 10s.
	ClientKeepAliveTime    = 10 * time.Second
	ClientKeepAliveTimeout = 5 * time.Second
)

var (
	// DefaultClientOptions returns the recommended default client options when
	// connecting to the server. This will mainly be used for client disconnect
	// detection.
	//
	// Example
	//
	// c, err := grpc.Dial("localhost:4444", DefaultClientOptions...)
	DefaultClientOptions = []grpc.DialOption{
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                ClientKeepAliveTime,
				Timeout:             ClientKeepAliveTimeout,
				PermitWithoutStream: false,
			},
		),
	}

	// DefaultServerOptions returns the default options the server will employ
	// for connecting to the client. Notably, these options will allow the server
	// to receive keepalive messages from the client periodically to facilitate
	// detecting network problems early.
	//
	// Example
	//
	// s := grpc.NewServer(DefaultServerOptions...)
	DefaultServerOptions = ServerOptions(
		ServerOptionConfig{
			MinimumClientInterval:   ServerKeepAliveTime,
			ServerHeartbeatInterval: ServerKeepAliveTime,
			ServerHeartbeatTimeout:  ServerKeepAliveTimeout,
		},
	)
)

type ServerOptionConfig struct {
	MinimumClientInterval   time.Duration
	ServerHeartbeatInterval time.Duration
	ServerHeartbeatTimeout  time.Duration
}

func ServerOptions(c ServerOptionConfig) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             c.MinimumClientInterval,
				PermitWithoutStream: false,
			},
		),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    c.ServerHeartbeatInterval,
				Timeout: c.ServerHeartbeatTimeout,
			},
		),
	}
}
