package client

import (
	"sync"

	apipb "github.com/downflux/game/api/api_go_proto"
)

type Connection struct {
	done chan struct{}

	mux       sync.Mutex
	responses []*apipb.StreamDataResponse
	status    bool
}

func New() *Connection {
	return &Connection{
		done: make(chan struct{}),
	}
}

func (c *Connection) Done() <-chan struct{} {
	return c.done
}

func (c *Connection) Close() {
	close(c.done)
}

func (c *Connection) SetChannelClosed(s bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.status = s
}

func (c *Connection) Responses() ([]*apipb.StreamDataResponse, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	resp := c.responses
	c.responses = nil
	ok := !c.status

	return resp, ok
}

func (c *Connection) AddMessage(m *apipb.StreamDataResponse) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.responses = append(c.responses, m)
}
