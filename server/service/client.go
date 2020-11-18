// Package client represents the server-specific view of a player and its
// associated metadata. This package is used to keep track of notably the
// connection status for StreamData calls and to encapsulate some reconnect
// logic.
package client

import (
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

const (
	clientBufSize = 5
)

// Client is the server-specific representation of a player in a specific game.
type Client struct {
	// id is the UUID of the connecting client. This is immutable.
	id string

	// mux guards ch and status, and must be acquired before any R/W
	// operations occur.
	mux    sync.Mutex

	// ch is an open connection for streaming data -- this is hooked up to
	// the gRPC server, which attempts to read from this channel as fast as
	// possible. This channel should not be blocked on writes.
	ch     chan *apipb.StreamDataResponse

	// status surfaces the server view of a specific client connection
	// state.
	status sscpb.ClientStatus
}

// invalidTransitionError constructs an appropriate gRPC error for when
// external logic attempts a transition between ClientStatus states with no
// edge.
func invalidTransitionError(clientStatus, targetStatus sscpb.ClientStatus) error {
	return status.Errorf(
		codes.FailedPrecondition,
		"there is no defined transition path from the current client status %v to %v",
		clientStatus,
		targetStatus)
}

// ID returns the UUID of the client.
func (c *Client) ID() string { return c.id }

// Status returns the current client connection state.
func (c *Client) Status() sscpb.ClientStatus {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.status
}

// SetStatus changes the client with side-effects to the desired status.
//
// This call is atomic.
func (c *Client) SetStatus(s sscpb.ClientStatus) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.setStatusUnsafe(s)
}

// setStatusUnsafe changes the client with side-effects to the desired status.
//
// The current expected transition graph is of the form
//   NEW -> DESYNCED -> OK -> NEW -> ... TEARDOWN
//                   -> NEW
//
// NEW: The client does not have an associated channel, and will not broadcast
// any data.
//
// DESYNCED: The client has an associated channel, but is not synced to the
// current game tick.
//
// OK: The client has an associated channel and is synced.
//
// TEARDOWN: The client has closed the channel and is not accepting any further
// status transitions.
func (c *Client) setStatusUnsafe(s sscpb.ClientStatus) error {
	switch s {
	case sscpb.ClientStatus_CLIENT_STATUS_NEW:
		switch c.status {
		case sscpb.ClientStatus_CLIENT_STATUS_OK:
			close(c.ch)
			c.ch = nil
		case sscpb.ClientStatus_CLIENT_STATUS_DESYNCED:
			close(c.ch)
			c.ch = nil
		default:
			return invalidTransitionError(c.status, s)
		}
	case sscpb.ClientStatus_CLIENT_STATUS_DESYNCED:
		switch c.status {
		case sscpb.ClientStatus_CLIENT_STATUS_NEW:
			c.ch = make(chan *apipb.StreamDataResponse, clientBufSize)
		default:
			return invalidTransitionError(c.status, s)
		}
	case sscpb.ClientStatus_CLIENT_STATUS_OK:
		switch c.status {
		case sscpb.ClientStatus_CLIENT_STATUS_OK:
		case sscpb.ClientStatus_CLIENT_STATUS_DESYNCED:
		default:
			return invalidTransitionError(c.status, s)
		}
	case sscpb.ClientStatus_CLIENT_STATUS_TEARDOWN:
		switch c.status {
		case sscpb.ClientStatus_CLIENT_STATUS_NEW:
		case sscpb.ClientStatus_CLIENT_STATUS_OK:
			close(c.ch)
			c.ch = nil
		case sscpb.ClientStatus_CLIENT_STATUS_DESYNCED:
			close(c.ch)
			c.ch = nil
		default:
			return invalidTransitionError(c.status, s)
		}
	default:
		return status.Errorf(codes.InvalidArgument, "invalid target client status %v", s)
	}

	c.status = s

	return nil
}

// Channel surfaces a read-only channel of game states. This data is generally
// consumed by the gRPC server and forwarded to the corresponding client.
func (c *Client) Channel() (<-chan *apipb.StreamDataResponse, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if (c.status == sscpb.ClientStatus_CLIENT_STATUS_UNKNOWN) || (c.status == sscpb.ClientStatus_CLIENT_STATUS_NEW) {
		return nil, status.Errorf(
			codes.FailedPrecondition,
			"client channel is not defined for clients with status %v",
			c.status)
	}

	return c.ch, nil
}

// Send will write the associated game state to the internal channel.
func (c *Client) Send(m *apipb.StreamDataResponse) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if (c.status == sscpb.ClientStatus_CLIENT_STATUS_UNKNOWN) || (c.status == sscpb.ClientStatus_CLIENT_STATUS_NEW) {
		return status.Errorf(
			codes.FailedPrecondition,
			"no client channel exists with client status %v",
			c.status)
	}

	c.ch <- m
	return c.setStatusUnsafe(sscpb.ClientStatus_CLIENT_STATUS_OK)
}

// New constructs a new Client instance.
func New(cid string) *Client {
	return &Client{
		id:     cid,
		status: sscpb.ClientStatus_CLIENT_STATUS_NEW,
	}
}
