package client

import (
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	scpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

type Client struct {
	id string // read-only

	mux    sync.Mutex
	ch     chan *apipb.StreamDataResponse
	status scpb.ClientStatus
}

func invalidTransitionError(clientStatus, targetStatus scpb.ClientStatus) error {
	return status.Errorf(
		codes.FailedPrecondition,
		"there is no defined transition path from the current client status %v to %v",
		clientStatus,
		targetStatus)
}

func (c *Client) ID() string { return c.id }
func (c *Client) Status() scpb.ClientStatus {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.status
}

// SetStatus changes the client with side-effects to the desired status.
//
// The current expected transition graph is of the form
//   NEW -> DESYNCED -> OK -> NEW
//
// NEW: The client does not have an associated channel, and will not broadcast
// any data.
//
// DESYNCED: The client has an associated channel, but is not synced to the
// current game tick.
//
// OK: The client has an associated channel and is synced.
func (c *Client) SetStatus(s scpb.ClientStatus) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	switch s {
	case scpb.ClientStatus_CLIENT_STATUS_NEW:
		switch c.status {
		case scpb.ClientStatus_CLIENT_STATUS_NEW:
		case scpb.ClientStatus_CLIENT_STATUS_OK:
			close(c.ch)
			c.ch = nil
		default:
			return invalidTransitionError(c.status, s)
		}
	case scpb.ClientStatus_CLIENT_STATUS_DESYNCED:
		switch c.status {
		case scpb.ClientStatus_CLIENT_STATUS_DESYNCED:
		case scpb.ClientStatus_CLIENT_STATUS_NEW:
			c.ch = make(chan *apipb.StreamDataResponse)
		default:
			return invalidTransitionError(c.status, s)
		}
	case scpb.ClientStatus_CLIENT_STATUS_OK:
		switch c.status {
		case scpb.ClientStatus_CLIENT_STATUS_OK:
		case scpb.ClientStatus_CLIENT_STATUS_DESYNCED:
		default:
			return invalidTransitionError(c.status, s)
		}
	default:
		return status.Errorf(codes.InvalidArgument, "invalid target client status %v", s)
	}

	c.status = s

	return nil
}

func (c *Client) Channel() (<-chan *apipb.StreamDataResponse, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.status == scpb.ClientStatus_CLIENT_STATUS_UNKNOWN || c.status == scpb.ClientStatus_CLIENT_STATUS_NEW {
		return nil, status.Errorf(codes.FailedPrecondition, "no client channel exists with client status %v", c.status)
	}

	return c.ch, nil
}

func (c *Client) Send(m *apipb.StreamDataResponse) error {
	ch, err := func () (chan<- *apipb.StreamDataResponse, error) {
		c.mux.Lock()
		defer c.mux.Unlock()

		if c.status == scpb.ClientStatus_CLIENT_STATUS_UNKNOWN || c.status == scpb.ClientStatus_CLIENT_STATUS_NEW {
			return nil, status.Errorf(
				codes.FailedPrecondition,
				"no client channel exists with client status %v",
				c.status)
		}
		return c.ch, nil
	}()
	if err != nil {
		return err
	}

	ch <- m
	return nil
}

func New(cid string) *Client {
	return &Client{
		id:     cid,
		status: scpb.ClientStatus_CLIENT_STATUS_NEW,
	}
}
