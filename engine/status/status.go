// Package status contains logic reguarding the current Executor state.
package status

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/downflux/game/engine/id/id"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	gdpb "github.com/downflux/game/api/data_go_proto"
	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

// Status represents the internal Executor state.
type Status struct {
	// tickDuration is the target maximum interval between successive
	// ticks, typically ~10Hz. This is immutable.
	tickDuration time.Duration

	// statusEnumMux guards the statusEnum property.
	statusEnumMux sync.Mutex

	// statusEnum represents the current executor internal run-state.
	statusEnum sscpb.ServerStatus

	// tickImpl represents the internal server tick counter. This is
	// advanced once per game loop.
	tickImpl int64

	// startTimeMux guards the startTimeImpl variable.
	startTimeMux sync.Mutex

	// startTimeImpl represents the time at which Executor.Run() was
	// called. This is useful client-side along with the tick duration to
	// estimate current server tick.
	//
	// This should only be set once.
	//
	// TODO(minkezhang): Actually ensure this is only set once.
	startTimeImpl time.Time
}

// New returns a new Status instance.
func New(tickDuration time.Duration) *Status {
	return &Status{
		tickDuration: tickDuration,
		statusEnum:   sscpb.ServerStatus_SERVER_STATUS_NOT_STARTED,
	}
}

// PB exports the Status instance into an associated protobuf.
func (s *Status) PB() *gdpb.ServerStatus {
	return &gdpb.ServerStatus{
		Tick:         s.Tick().Value(),
		IsStarted:    s.IsStarted(),
		TickDuration: durationpb.New(s.tickDuration),
		StartTime:    timestamppb.New(s.StartTime()),
	}
}

func (s *Status) TickDuration() time.Duration { return s.tickDuration }

// Tick returns the current game tick.
func (s *Status) Tick() id.Tick { return id.Tick(atomic.LoadInt64(&(s.tickImpl))) }

// IncrementTick adds one to the current game tick -- this is called at the
// beginning of each tick loop.
func (s *Status) IncrementTick() { atomic.AddInt64(&(s.tickImpl), 1) }

// IsStarted returns if the Executor is currently executing ticks.
func (s *Status) IsStarted() bool {
	s.statusEnumMux.Lock()
	defer s.statusEnumMux.Unlock()
	return s.statusEnum == sscpb.ServerStatus_SERVER_STATUS_RUNNING
}

// SetIsStarted sets the internal Executor status as running the core loop.
func (s *Status) SetIsStarted() error {
	s.statusEnumMux.Lock()
	defer s.statusEnumMux.Unlock()

	target := sscpb.ServerStatus_SERVER_STATUS_RUNNING
	if s.statusEnum != sscpb.ServerStatus_SERVER_STATUS_NOT_STARTED {
		return status.Errorf(codes.Aborted, "cannot set server status to %v from %v")
	}
	s.statusEnum = target
	return nil
}

// IsStopped returns if the Executor has finished running the game loop.
func (s *Status) IsStopped() bool {
	s.statusEnumMux.Lock()
	defer s.statusEnumMux.Unlock()
	return s.statusEnum == sscpb.ServerStatus_SERVER_STATUS_STOPPED
}

// SetIsStopped sets the internal Executor status as having finished the core
// loop. This may be called because the game ended, or if it crashed.
func (s *Status) SetIsStopped() error {
	s.statusEnumMux.Lock()
	defer s.statusEnumMux.Unlock()

	s.statusEnum = sscpb.ServerStatus_SERVER_STATUS_STOPPED
	return nil
}

// StartTime returns the wall-clock time at which the server started executing
// the core loop.
func (s *Status) StartTime() time.Time {
	s.startTimeMux.Lock()
	defer s.startTimeMux.Unlock()

	return s.startTimeImpl
}

// SetStartTime sets the time at which the server initially started running the
// core loop.
func (s *Status) SetStartTime() {
	s.startTimeMux.Lock()
	defer s.startTimeMux.Unlock()

	s.startTimeImpl = time.Now()
}
