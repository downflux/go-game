package clientlist

import (
	"sync"

	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/client"
	"golang.org/x/sync/errgroup"
        "google.golang.org/grpc/codes"
        "google.golang.org/grpc/status"

	apipb "github.com/downflux/game/api/api_go_proto"
	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

var (
	notFound = status.Errorf(codes.NotFound, "specified client not found in client list")
)

type List struct {
	idLen int

	mux sync.RWMutex
	clients map[string]*client.Client
}

func New(idLen int) *List {
	return &List{
		idLen: idLen,
		clients: map[string]*client.Client{},
	}
}

func (l *List) In(cid string) bool {
        l.mux.RLock()
	defer l.mux.RUnlock()

	return l.inUnsafe(cid)
}

func (l *List) Broadcast(partialGenerator, fullGenerator func() *apipb.StreamDataResponse) error {
        l.mux.RLock()
	defer l.mux.RUnlock()

	partial := partialGenerator()
	var full *apipb.StreamDataResponse

        desyncedClients := l.filterUnsafe(sscpb.ClientStatus_CLIENT_STATUS_DESYNCED)
        if desyncedClients == nil && partial.GetCurves() == nil && partial.GetEntities() == nil {
		return nil
        }
	if desyncedClients != nil {
		full = fullGenerator()
	}

        var eg errgroup.Group
        for _, c := range l.clients {
                c := c
		switch c.Status() {
		case sscpb.ClientStatus_CLIENT_STATUS_OK:
			eg.Go(func() error { return c.Send(partial) })
		case sscpb.ClientStatus_CLIENT_STATUS_DESYNCED:
			eg.Go(func() error { return c.Send(full) })
		}
        }
        return eg.Wait()
}

func (l *List) Channel(cid string) (<-chan *apipb.StreamDataResponse, error) {
	l.mux.RLock()
	defer l.mux.RUnlock()

	if !l.inUnsafe(cid) {
		return nil, notFound
	}

	return l.clients[cid].Channel()
}

func (l *List) Add() (string, error) {
	// TODO(minkezhang): Add maxClients check.
        l.mux.Lock()
        defer l.mux.Unlock()

        cid := id.RandomString(l.idLen)
        for _, found := l.clients[cid]; found; cid = id.RandomString(l.idLen) {
        }
        l.clients[cid] = client.New(cid)

        return cid, nil

}

func (l *List) Start(cid string) error {
	l.mux.RLock()
	defer l.mux.RUnlock()

	if !l.inUnsafe(cid) {
		return notFound
	}

        return l.clients[cid].SetStatus(sscpb.ClientStatus_CLIENT_STATUS_DESYNCED)
}

func (l *List) Stop(cid string, success bool) error {
	l.mux.RLock()
	defer l.mux.RUnlock()

	return l.stopUnsafe(cid, success)
}

func (l *List) StopAll() error {
	l.mux.RLock()
	defer l.mux.RUnlock()

	for cid, _ := range l.clients {
		if err := l.stopUnsafe(cid, true); err != nil {
			return err
		}
	}

	return nil
}

func (l *List) stopUnsafe(cid string, success bool) error {
	if !l.inUnsafe(cid) {
		return notFound
	}

	if success {
		return l.clients[cid].SetStatus(sscpb.ClientStatus_CLIENT_STATUS_TEARDOWN)
	}
        return l.clients[cid].SetStatus(sscpb.ClientStatus_CLIENT_STATUS_NEW)
}

func (l *List) inUnsafe(cid string) bool {
	_, found := l.clients[cid]
        return found
}

func (l *List) filterUnsafe(status sscpb.ClientStatus) map[string]bool {
	cids := map[string]bool{}

	for _, c := range l.clients {
		if c.Status() == status {
			cids[c.ID()] = true
		}
	}

	return cids
}
