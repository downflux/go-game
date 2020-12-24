// Package list encapsulates logic for managing multiple clients.
package list

import (
	"sync"

	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/downflux/game/engine/server/client/client"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	ccpb "github.com/downflux/game/engine/server/client/api/constants_go_proto"
)

var (
	notFound = status.Errorf(codes.NotFound, "specified client not found in client list")
)

// List contains an iterable of Client instances.
type List struct {
	// idLen is a configurable value for the length of generated Client
	// UUIDs.
	idLen int

	// mux guards the clients iterable.
	mux sync.RWMutex

	// clients is an internal iterable of Client instances, hashed by the
	// Client UUID.
	clients map[id.ClientID]*client.Client
}

// New returns a new List instance.
func New(idLen int) *List {
	return &List{
		idLen:   idLen,
		clients: map[id.ClientID]*client.Client{},
	}
}

// In atomically checks for if the given Client with corresponding UUID exists
// in the List.
func (l *List) In(cid id.ClientID) bool {
	l.mux.RLock()
	defer l.mux.RUnlock()

	return l.inUnsafe(cid)
}

// Broadcast atomically sends data to all available Client instances.
//
// the partialGenerator and fullGenerator functions are invoked when there is
// at least one Client which needs this data; we're passing in the functions
// instead of the actual messages, as generating the messages themselves may
// be expensive.
//
// The full game state created by the fullGenerator function is necessary when
// a Client is in state DESYNCED.
func (l *List) Broadcast(partialGenerator, fullGenerator func() *apipb.StreamDataResponse) error {
	l.mux.RLock()
	defer l.mux.RUnlock()

	partial := partialGenerator()
	var full *apipb.StreamDataResponse

	desyncedClients := l.filterUnsafe(ccpb.ClientState_CLIENT_STATE_DESYNCED)
	if desyncedClients == nil && partial.GetCurves() == nil && partial.GetEntities() == nil {
		return nil
	}
	if desyncedClients != nil {
		full = fullGenerator()
	}

	var eg errgroup.Group
	for _, c := range l.clients {
		c := c
		s, err := c.State()
		if err != nil {
			return err
		}
		switch s {
		case fsm.State(ccpb.ClientState_CLIENT_STATE_OK.String()):
			eg.Go(func() error { return c.Send(partial) })
		case fsm.State(ccpb.ClientState_CLIENT_STATE_DESYNCED.String()):
			eg.Go(func() error { return c.Send(full) })
		}
	}
	return eg.Wait()
}

// Channel returns a read-only channel of game states. This is generally passed
// to the gRPC server to be forwarded to the client.
func (l *List) Channel(cid id.ClientID) (<-chan *apipb.StreamDataResponse, error) {
	l.mux.RLock()
	defer l.mux.RUnlock()

	if !l.inUnsafe(cid) {
		return nil, notFound
	}

	return l.clients[cid].Channel()
}

// Add creates a new Client instance and inserts it into the List.
func (l *List) Add() (id.ClientID, error) {
	// TODO(minkezhang): Add maxClients check.
	l.mux.Lock()
	defer l.mux.Unlock()

	cid := id.ClientID(id.RandomString(l.idLen))
	for _, found := l.clients[cid]; found; cid = id.ClientID(id.RandomString(l.idLen)) {
	}
	l.clients[cid] = client.New(cid)

	return cid, nil

}

// Start will indicate to the associated Client instance that a channel
// instance should be created, and allows Client.Send() calls to occur.
func (l *List) Start(cid id.ClientID) error {
	l.mux.RLock()
	defer l.mux.RUnlock()

	if !l.inUnsafe(cid) {
		return notFound
	}

	return l.clients[cid].SetState(ccpb.ClientState_CLIENT_STATE_DESYNCED)
}

// Stop will indicate to the associated Client that the game state channel
// should be torn down, either because the game ended (success == true) or a
// network disconnect has occurred (success == false).
func (l *List) Stop(cid id.ClientID, success bool) error {
	l.mux.RLock()
	defer l.mux.RUnlock()

	return l.stopUnsafe(cid, success)
}

// StopAll will iteratively teardown all Client channels. This is typically
// done when the game ends.
func (l *List) StopAll() error {
	l.mux.RLock()
	defer l.mux.RUnlock()

	for cid := range l.clients {
		if err := l.stopUnsafe(cid, true); err != nil {
			return err
		}
	}

	return nil
}

// stopUnsafe implements the Client channel disconnection logic.
func (l *List) stopUnsafe(cid id.ClientID, success bool) error {
	if !l.inUnsafe(cid) {
		return notFound
	}

	if success {
		return l.clients[cid].SetState(ccpb.ClientState_CLIENT_STATE_TEARDOWN)
	}
	return l.clients[cid].SetState(ccpb.ClientState_CLIENT_STATE_NEW)
}

// inUnsafe implements the Client membership test logic.
func (l *List) inUnsafe(cid id.ClientID) bool {
	_, found := l.clients[cid]
	return found
}

// filterUnsafe retuns a list of Client instances which are currently in the
// specified ClientStatus.
func (l *List) filterUnsafe(filterState ccpb.ClientState) map[id.ClientID]bool {
	cids := map[id.ClientID]bool{}

	for _, c := range l.clients {
		s, err := c.State()
		if err == nil && s == fsm.State(filterState.String()) {
			cids[c.ID()] = true
		}
	}

	return cids
}
