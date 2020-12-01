// Package client contains the gRPC-specific view of a connected client.
package client

import (
	"sync"

	apipb "github.com/downflux/game/api/api_go_proto"
)

// Connection encapsulates client-specific connection metadata.
type Connection struct {
	// done indicates to any Goroutine that the underlying client has
	// closed its physical connection and should exit.
	done chan struct{}

	// mux guards the responses and status properties.
	mux       sync.Mutex

	// responses is a cache of upstream Executor responses.
	responses []*apipb.StreamDataResponse

	// status indicates if the client has closed gracefully.
	status    bool
}

// New creates a new instance of the Connection object.
func New() *Connection {
	return &Connection{
		done: make(chan struct{}),
	}
}

// Done returns a signal channel. The channel is closed when the physical
// connection to the remote client is broken.
func (c *Connection) Done() <-chan struct{} {
	return c.done
}

// Close closes the internal signal channel.
func (c *Connection) Close() {
	close(c.done)
}

// SetChannelClosed sets the Connection status to the indicated value. A true
// value indicates that no further responses are expected to be added to the
// response queue.
//
// TODO(minkezhang): Decide if we can merge with Close.
func (c *Connection) SetChannelClosed(s bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.status = s
}

// Responses returns the internal cache of Executor responses. The internal
// cache is reset.
//
// TODO(minkezhang): Rename to PopResponses.
func (c *Connection) Responses() ([]*apipb.StreamDataResponse, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	resp := c.responses
	c.responses = nil
	ok := !c.status

	return resp, ok
}

// AddMessage appends a new Executor game state delta to the internal cache.
func (c *Connection) AddMessage(m *apipb.StreamDataResponse) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.responses = append(c.responses, m)
}
