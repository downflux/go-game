package status

import (
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	gdpb "github.com/downflux/game/api/data_go_proto"
)

type Status struct {
	tickDuration  time.Duration
	isStoppedImpl int32
	isStartedImpl int32
	tickImpl      int64

	startTimeMux  sync.Mutex
	startTimeImpl time.Time
}

func New(tickDuration time.Duration) *Status {
	return &Status{
		tickDuration: tickDuration,
	}
}

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
