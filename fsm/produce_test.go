package produce

import (
	"testing"

	"github.com/downflux/game/engine/fsm/fsm"
	"github.com/downflux/game/engine/fsm/instance"
	"github.com/downflux/game/engine/status/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

var (
	_ instance.Instance = &Instance{}
)

func TestConstructor(t *testing.T) {
	s := status.New(0)

	testConfigs := []struct {
		name string
		i    instance.Instance
		want fsm.State
	}{
		{
			name: "NewPending",
			i:    New(s, s.Tick()+1, gcpb.EntityType_ENTITY_TYPE_TANK, &gdpb.Position{X: 0, Y: 0}),
			want: fsm.State(fcpb.CommonState_COMMON_STATE_PENDING.String()),
		},
		{
			name: "NewExecuting",
			i:    New(s, s.Tick(), gcpb.EntityType_ENTITY_TYPE_TANK, &gdpb.Position{X: 0, Y: 0}),
			want: fsm.State(fcpb.CommonState_COMMON_STATE_EXECUTING.String()),
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := c.i.State(); err != nil || got != c.want {
				t.Errorf("State() == %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestFinish(t *testing.T) {
	s := status.New(0)
	i := New(s, s.Tick(), gcpb.EntityType_ENTITY_TYPE_TANK, &gdpb.Position{X: 0, Y: 0})

	if err := i.Finish(); err != nil {
		t.Fatalf("Finish() = %v, want = nil", err)
	}
}

func TestCancel(t *testing.T) {
	s := status.New(0)
	i := New(s, s.Tick()+1, gcpb.EntityType_ENTITY_TYPE_TANK, &gdpb.Position{X: 0, Y: 0})

	if err := i.Cancel(); err != nil {
		t.Fatalf("Cancel() = %v, want = nil", err)
	}
}
