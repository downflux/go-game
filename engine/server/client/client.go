package client

import (
	"sync"

	"github.com/downflux/game/engine/fsm/action"
	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
	ccpb "github.com/downflux/game/engine/server/client/api/constants_go_proto"
)

const (
	fsmType = fcpb.FSMType_FSM_TYPE_CLIENT

	// TODO(minkezhang): Change to a buffered value (e.g. 5) and verify
	// tests do not break.
	clientBufSize = 0
)

var (
	unknown = fsm.State(ccpb.ClientState_CLIENT_STATE_UNKNOWN.String())

	// newState indicates the client does not have an associated channel,
	// and will not broadcast any data.
	newState = fsm.State(ccpb.ClientState_CLIENT_STATE_NEW.String())

	// desynced indicates the client has an associated channel, but is not
	// synced to the current game tick.
	desynced = fsm.State(ccpb.ClientState_CLIENT_STATE_DESYNCED.String())

	// ok indicates the client has an associated channel and is synced.
	ok = fsm.State(ccpb.ClientState_CLIENT_STATE_OK.String())

	// teardown indicates the client has closed the channel and is not
	// accepting any further status transitions.
	teardown = fsm.State(ccpb.ClientState_CLIENT_STATE_TEARDOWN.String())

	transitions = []fsm.Transition{
		{From: newState, To: desynced},
		{From: desynced, To: newState},
		{From: desynced, To: ok},
		{From: ok, To: ok},
		{From: ok, To: newState},
		{From: newState, To: teardown},
		{From: desynced, To: teardown},
		{From: ok, To: teardown},
	}

	FSM = fsm.New(transitions, fsmType)
)

type Client struct {
	*action.Base

	// id is the UUID of the connecting client.
	id id.ClientID // Read-only.

	// mux guards the Base and ch properties.
	mux sync.Mutex

	// ch is an open connection for streaming data -- this is hooked up to
	// the gRPC server, which attempts to read from this channel as fast as
	// possible. This channel should not be blocked on writes.
	ch chan *apipb.StreamDataResponse
}

// New constructs a new Client instance.
func New(cid id.ClientID) *Client {
	return &Client{
		Base: action.New(FSM, newState),
		id:   cid,
	}
}

func (c *Client) ID() id.ClientID { return c.id }

// State returns the current client connection state.
func (c *Client) State() (fsm.State, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.Base.State()
}

// SetState changes the client with side-effects to the desired status.
//
// This call is atomic.
func (c *Client) SetState(s ccpb.ClientState) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.setStateUnsafe(s)
}

// setStateUnsafe changes the client with side-effects to the desired status.
func (c *Client) setStateUnsafe(s ccpb.ClientState) error {
	f, err := c.Base.State()
	if err != nil {
		return err
	}

	t := fsm.State(s.String())

	switch t {
	case fsm.State(ccpb.ClientState_CLIENT_STATE_NEW.String()):
		if f == fsm.State(ccpb.ClientState_CLIENT_STATE_OK.String()) || f == fsm.State(ccpb.ClientState_CLIENT_STATE_DESYNCED.String()) {
			close(c.ch)
			c.ch = nil
		}
	case fsm.State(ccpb.ClientState_CLIENT_STATE_DESYNCED.String()):
		if f == fsm.State(ccpb.ClientState_CLIENT_STATE_NEW.String()) {
			c.ch = make(chan *apipb.StreamDataResponse, clientBufSize)
		}
	case fsm.State(ccpb.ClientState_CLIENT_STATE_TEARDOWN.String()):
		if f == fsm.State(ccpb.ClientState_CLIENT_STATE_NEW.String()) || f == fsm.State(ccpb.ClientState_CLIENT_STATE_DESYNCED.String()) {
			close(c.ch)
			c.ch = nil
		}
	}

	return c.Base.To(f, t, false)
}

// Channel surfaces a read-only channel of game states. This data is generally
// consumed by the gRPC server and forwarded to the corresponding client.
func (c *Client) Channel() (<-chan *apipb.StreamDataResponse, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	s, err := c.Base.State()
	if err != nil {
		return nil, err
	}

	if s == fsm.State(ccpb.ClientState_CLIENT_STATE_UNKNOWN.String()) || s == fsm.State(ccpb.ClientState_CLIENT_STATE_NEW.String()) {
		return nil, status.Errorf(codes.FailedPrecondition, "client channel is not defined for clients with status %v", s)
	}

	return c.ch, nil
}

// Send will write the associated game state to the internal channel.
func (c *Client) Send(m *apipb.StreamDataResponse) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	s, err := c.Base.State()
	if err != nil {
		return err
	}

	if s == fsm.State(ccpb.ClientState_CLIENT_STATE_UNKNOWN.String()) || s == fsm.State(ccpb.ClientState_CLIENT_STATE_NEW.String()) {
		return status.Errorf(codes.FailedPrecondition, "no client channel exists with client status %v", s)
	}

	// Only send data if there is interesting data to send.
	if m.GetEntities() != nil || m.GetCurves() != nil {
		c.ch <- m
		return c.setStateUnsafe(ccpb.ClientState_CLIENT_STATE_OK)
	}
	return nil
}
