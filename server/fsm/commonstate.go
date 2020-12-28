package commonstate

import (
	"github.com/downflux/game/engine/fsm/fsm"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

var (
	Unknown   = fsm.State(fcpb.CommonState_COMMON_STATE_UNKNOWN.String())
	Pending   = fsm.State(fcpb.CommonState_COMMON_STATE_PENDING.String())
	Executing = fsm.State(fcpb.CommonState_COMMON_STATE_EXECUTING.String())
	Canceled  = fsm.State(fcpb.CommonState_COMMON_STATE_CANCELED.String())
	Finished  = fsm.State(fcpb.CommonState_COMMON_STATE_FINISHED.String())
)
