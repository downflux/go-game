package acl

import (
	"testing"

	"github.com/downflux/game/engine/id/id"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

func TestValidate(t *testing.T) {
	cid := id.ClientID("client-id")

	testConfigs := []struct {
		name  string
		acl   ACL
		input *gdpb.ClientID
		want  bool
	}{
		{name: "TestPublic", acl: *New(cid, PublicWritable), want: true},
		{
			name:  "TestPublicWithInput",
			acl:   *New(cid, PublicWritable),
			input: &gdpb.ClientID{ClientId: "another-client-id"},
			want:  true,
		},
		{
			name:  "TestClientOK",
			acl:   *New(cid, ClientWritable),
			input: &gdpb.ClientID{ClientId: "client-id"},
			want:  true,
		},
		{
			name:  "TestClientInvalid",
			acl:   *New(cid, ClientWritable),
			input: &gdpb.ClientID{ClientId: "another-client-id"},
			want:  false,
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got := c.acl.Validate(c.input); got != c.want {
				t.Errorf("Validate() = %v, want = %v", got, c.want)
			}
		})
	}
}
