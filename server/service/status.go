// Package status contains logic reguarding the current Executor state.
package status

import (
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

// Status represents the internal Executor state.
type Status struct {
	// tickDuration is the target maximum interval between successive
	// ticks, typically ~10Hz. This is immutable.
	tickDuration  time.Duration

	// isStoppedImpl represents the boolean value of if the Executor
	// should stop running the core game loop logic. This boolean is set
	// to true at teardown.
	//
	// TODO(minkezhang): Remove this logic and combile with isStartedImpl.
	isStoppedImpl int32

	// isStartedImpl represents the boolean value of if the Executor
	// has started running the core game loop logic. This boolean is set
	// to true at the beginning of Executor.Run().
	isStartedImpl int32

	// tickImpl represents the internal server tick counter. This is
	// advanced once per game loop.
	tickImpl      int64

	// startTimeMux guards the startTimeImpl variable.
	startTimeMux  sync.Mutex

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
	}
}

// PB exports the Status instance into an associated protobuf.
func (s *Status) PB() *gdpb.ServerStatus {
	return &gdpb.ServerStatus{
		Tick:         s.Tick(),
		IsStarted:    s.IsStarted(),
		TickDuration: durationpb.New(s.tickDuration),
		StartTime:    timestamppb.New(s.StartTime()),
	}
}

func (s *Status) Tick() float64   { return float64(atomic.LoadInt64(&(s.tickImpl))) }
func (s *Status) IncrementTick()  { atomic.AddInt64(&(s.tickImpl), 1) }
func (s *Status) IsStarted() bool { return atomic.LoadInt32(&(s.isStartedImpl)) != 0 }
func (s *Status) SetIsStarted()   { atomic.StoreInt32(&(s.isStartedImpl), 1) }
func (s *Status) IsStopped() bool { return atomic.LoadInt32(&(s.isStoppedImpl)) != 0 }
func (s *Status) SetIsStopped()   { atomic.StoreInt32(&(s.isStoppedImpl), 1) }

func (s *Status) StartTime() time.Time {
	s.startTimeMux.Lock()
	defer s.startTimeMux.Unlock()

	return s.startTimeImpl
}

func (s *Status) SetStartTime() {
	s.startTimeMux.Lock()
	defer s.startTimeMux.Unlock()

	s.startTimeImpl = time.Now()
}
