package client

import (
	"sync"

	apipb "github.com/downflux/game/api/api_go_proto"
)

type Client struct {
	mux      sync.Mutex
	id       string // read-only
	ch       chan *apipb.StreamDataResponse
	isSynced bool
}

func (c *Client) ID() string {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.id
}
func (c *Client) Channel() chan *apipb.StreamDataResponse {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.ch
}
func (c *Client) NewChannel() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.ch = make(chan *apipb.StreamDataResponse)
}
func (c *Client) CloseChannel() {
	c.mux.Lock()
	defer c.mux.Unlock()

	close(c.ch)
	c.ch = nil
}
func (c *Client) IsSynced() bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.isSynced
}
func (c *Client) SetIsSynced(s bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.isSynced = s
	if !s {
		close(c.ch)
		c.ch = nil
	}
}

func New(cid string) *Client {
	c := &Client{
		id:       cid,
		isSynced: false,
	}
	c.NewChannel()
	return c
}
