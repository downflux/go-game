package client

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/id/id"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"

	apipb "github.com/downflux/game/api/api_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	ccpb "github.com/downflux/game/engine/server/client/api/constants_go_proto"
)

func TestNew(t *testing.T) {
	const cid = "client-id"
	state := fsm.State(ccpb.ClientState_CLIENT_STATE_NEW.String())

	c := New(id.ClientID(cid))

	if c.ID() != cid {
		t.Fatalf("ID() = %v, want = %v", c.ID(), cid)
	}

	if got, err := c.State(); err != nil || got != state {
		t.Errorf("State() = %v, %v, want = %v, nil", got, err, state)
	}
}

func TestGetChannelInvald(t *testing.T) {
	c := New("client-id")

	if _, err := c.Channel(); err == nil {
		t.Fatal("Channel() = nil, want a non-nil error")
	}

	if err := c.Send(nil); err == nil {
		t.Error("Send() = nil, want a non-nil-error")
	}
}

func TestSend(t *testing.T) {
	const nClients = 1000
	message := &apipb.StreamDataResponse{
		State: &gdpb.GameState{
			Entities: []*gdpb.Entity{
				{EntityId: "eid"},
			},
		},
	}

	var clients []*Client
	var channels []<-chan *apipb.StreamDataResponse
	for i := 0; i < nClients; i++ {
		c := New(id.ClientID(fmt.Sprintf("client-%d", i)))
		if err := c.SetState(ccpb.ClientState_CLIENT_STATE_DESYNCED); err != nil {
			t.Fatalf("SetState() = %v, want = nil", err)
		}
		if err := c.SetState(ccpb.ClientState_CLIENT_STATE_OK); err != nil {
			t.Fatalf("SetState() = %v, want = nil", err)
		}
		clients = append(clients, c)

		ch, err := c.Channel()
		if err != nil {
			t.Fatalf("Channel() = %v, want = nil", err)
		}
		channels = append(channels, ch)
	}

	var eg errgroup.Group
	for i := 0; i < nClients; i++ {
		i := i
		eg.Go(func() error { return clients[i].Send(message) })
	}

	for i := 0; i < nClients; i++ {
		ch := channels[i]
		eg.Go(func() error {
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			m := <-ch
			if diff := cmp.Diff(m, message, protocmp.Transform()); diff != "" {
				return status.Errorf(codes.Internal, "<-ch mismatch (-want +got):\n%v", diff)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("Wait() = %v, want = nil", err)
	}
}
