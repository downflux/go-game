package client

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/downflux/game/server/id"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"

	apipb "github.com/downflux/game/api/api_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

func TestNew(t *testing.T) {
	const cid = "client-id"
	const status = sscpb.ClientStatus_CLIENT_STATUS_NEW

	c := New(id.ClientID(cid))

	if c.ID() != cid {
		t.Fatalf("ID() = %v, want = %v", c.ID(), cid)
	}

	if c.Status() != status {
		t.Errorf("Status() = %v, want = %v", c.Status, status)
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
		Entities: []*gdpb.Entity{
			{EntityId: "eid"},
		},
	}

	var clients []*Client
	var channels []<-chan *apipb.StreamDataResponse
	for i := 0; i < nClients; i++ {
		c := New(id.ClientID(fmt.Sprintf("client-%d", i)))
		if err := c.SetStatus(sscpb.ClientStatus_CLIENT_STATUS_DESYNCED); err != nil {
			t.Fatalf("SetStatus() = %v, want = nil", err)
		}
		if err := c.SetStatus(sscpb.ClientStatus_CLIENT_STATUS_OK); err != nil {
			t.Fatalf("SetStatus() = %v, want = nil", err)
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
		eg.Go(func() error {
			err := clients[i].Send(message)
			fmt.Println("DEBUG: ERR %v", err)
			return err
		})
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
