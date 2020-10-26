package executor

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/entity/entity"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/commands/move"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	apipb "github.com/downflux/game/api/api_go_proto"
	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
	sscpb "github.com/downflux/game/server/service/api/constants_go_proto"
)

const (
	idLen = 8
	tickDuration = 100 * time.Millisecond
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

// TODO(minkezhang): Move to command/ directory.
type Command interface {
	Type() sscpb.CommandType
	ClientID() string
	Tick() float64

	// TODO(minkezhang): Refactor Curve interface to not be dependent on curve.Curve.
	Execute() (curve.Curve, error)
}

func New(pb *mdpb.TileMap, d *gdpb.Coordinate) (*Executor, error) {
	tm, err := tile.ImportMap(pb)
	if err != nil {
		return nil, err
	}
	g, err := graph.BuildGraph(tm, d)
	if err != nil {
		return nil, err
	}

	return &Executor{
		tileMap:       tm,
		abstractGraph: g,
		entities:      map[string]entity.Entity{},
		commandQueue:  nil,
		clientChannel: map[string]chan *apipb.StreamCurvesResponse{},
	}, nil
}

type Executor struct {
	tileMap       *tile.Map
	abstractGraph *graph.Graph

	isStoppedImpl int32
	isStartedImpl int32
	tickImpl int64

	// Add-only.
	dataMux  sync.RWMutex
	entities map[string]entity.Entity

	// Add and delete. Reset per tick.
	commandQueueMux sync.Mutex
	commandQueue    []Command

	// Add and delete. Reset per tick.
	curveQueueMux sync.RWMutex
	curveQueue    []curve.Curve

	clientChannelMux sync.RWMutex
	clientChannel    map[string]chan *apipb.StreamCurvesResponse
}

func (e *Executor) tick() float64 { return float64(atomic.LoadInt64(&(e.tickImpl))) }
func (e *Executor) incrementTick() { atomic.AddInt64(&(e.tickImpl), 1) }
func (e *Executor) isStarted() bool { return atomic.LoadInt32(&(e.isStartedImpl)) != 0 }
func (e *Executor) setIsStarted() { atomic.StoreInt32(&(e.isStartedImpl), 1) }
func (e *Executor) isStopped() bool { return atomic.LoadInt32(&(e.isStoppedImpl)) != 0 }
func (e *Executor) setIsStopped() { atomic.StoreInt32(&(e.isStoppedImpl), 1) }

func (e *Executor) Status() *gdpb.ServerStatus {
	return &gdpb.ServerStatus{
		Tick: e.tick(),
		IsStarted: e.isStarted(),
		TickDuration: durationpb.New(tickDuration),
	}
}

func (e *Executor) AddClient() (string, error) {
	// TODO(minkezhang): Add Client struct.
	// TODO(minkezhang): Add maxClients check.
	e.clientChannelMux.Lock()
	defer e.clientChannelMux.Unlock()

	cid := id.RandomString(idLen)
	for _, found := e.clientChannel[cid]; found; cid = id.RandomString(idLen) {
	}
	e.clientChannel[cid] = make(chan *apipb.StreamCurvesResponse)

	return cid, nil
}

func (e *Executor) ClientChannel(cid string) <-chan *apipb.StreamCurvesResponse {
	e.clientChannelMux.RLock()
	defer e.clientChannelMux.RUnlock()

	return e.clientChannel[cid]
}

func (e *Executor) popCommandQueue() ([]Command, error) {
	e.commandQueueMux.Lock()
	defer e.commandQueueMux.Unlock()

	commands := e.commandQueue
	e.commandQueue = nil
	return commands, nil
}

func (e *Executor) processCommand(cmd Command) error {
	if cmd.Type() == sscpb.CommandType_COMMAND_TYPE_MOVE {
		c, err := cmd.Execute()
		if err != nil {
			return err
		}

		if err := func() error {
			e.dataMux.RLock()
			defer e.dataMux.RUnlock()

			if err := e.entities[c.EntityID()].Curve(gcpb.CurveCategory_CURVE_CATEGORY_MOVE).ReplaceTail(c); err != nil {
				return err
			}

			log.Printf("server appending adding curve to queue: %v", c)

			e.curveQueueMux.Lock()
			e.curveQueue = append(e.curveQueue, c)
			e.curveQueueMux.Unlock()

			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) popCurveQueue() []curve.Curve {
	e.curveQueueMux.Lock()
	defer e.curveQueueMux.Unlock()

	curves := e.curveQueue
	e.curveQueue = nil
	return curves
}

func (e *Executor) broadcastCurves() error {
	curves := e.popCurveQueue()

	// TODO(minkezhang): Implement a server status endpoint.
	// server_test is relying on this as a proxy for IsAlive().
	if curves == nil {
		return nil
	}

	resp := &apipb.StreamCurvesResponse{
		Tick: e.tick(),
	}

	for _, c := range curves {
		d, err := c.ExportDelta()
		if err != nil {
			return err
		}
		resp.Curves = append(resp.GetCurves(), d)
	}

	// TODO(minkezhang): Make concurrent, and timeout.
	// Once timeout, client needs to resync.
	e.clientChannelMux.RLock()
	defer e.clientChannelMux.RUnlock()
	for _, ch := range e.clientChannel {
		log.Println("server sending message ", resp)
		ch <- resp
		log.Println("server sent message")
	}
	return nil
}

// TODO(minkezhang): Make private as part of loop.
func CloseStreams(e *Executor) error {
	e.clientChannelMux.Lock()
	defer e.clientChannelMux.Unlock()
	for _, ch := range e.clientChannel {
		close(ch)
	}
	return nil
}

func (e *Executor) Stop() { e.setIsStopped() }

func (e *Executor) Run() error {
	e.setIsStarted()
	for !e.isStopped() {
		t := time.Now()
		e.incrementTick()

		if err := e.t(); err != nil {
			return err
		}

		if d := time.Now().Sub(t); d < tickDuration {
			time.Sleep(tickDuration - d)
		}
	}
	return nil
}

// TODO(minkezhang): Test this and make private.
func (e *Executor) t() error {
	commands, err := e.popCommandQueue()
	if err != nil {
		return err
	}

	log.Printf("[%d] processing commands", e.tick())
	for _, cmd := range commands {
		// TODO(minkezhang): Add actual error handling here -- only
		// Only return early if error is very bad.
		if err := e.processCommand(cmd); err != nil {
			log.Printf("[%d] error while processing command %v: %v", cmd, err)
			return err
		}
	}

	// TODO(minkezhang): Broadcast new entities.

	log.Printf("[%d] broadcasting curves", e.tick())
	if err := e.broadcastCurves(); err != nil {
		return err
	}

	return nil
}

// TODO(minkezhang): Make this method private -- this is currently public for
// debugging purposes.
func (e *Executor) AddEntity(en entity.Entity) error {
	e.dataMux.Lock()
	defer e.dataMux.Unlock()

	if _, found := e.entities[en.ID()]; found {
		return status.Errorf(codes.AlreadyExists, "given entity ID %v already exists in the entity list", en.ID())
	}

	// TODO(minkezhang): Broadcast new entities by adding to new entity
	// queue.

	e.entities[en.ID()] = en
	return nil
}

func (e *Executor) addCommands(cs []Command) error {
	e.commandQueueMux.Lock()
	defer e.commandQueueMux.Unlock()

	e.commandQueue = append(e.commandQueue, cs...)

	// TODO(minkezhang): Add client validation as per design doc.
	return nil
}

// buildMoveCommands
//
// Is expected to be called concurrently.
//
// TODO(minkezhang): Decide how / when / if we want to deal with click
// spamming (same eids, multiple move commands per tick).
func (e *Executor) buildMoveCommands(cid string, t float64, dest *gdpb.Position, eids []string) []*move.Command {
	e.dataMux.RLock()
	defer e.dataMux.RUnlock()

	var res []*move.Command
	for _, eid := range eids {
		en, found := e.entities[eid]
		if found {
			p := en.Curve(gcpb.CurveCategory_CURVE_CATEGORY_MOVE).Get(t)
			res = append(
				res,
				move.New(
					e.tileMap,
					e.abstractGraph,
					cid,
					eid,
					t,
					p.(*gdpb.Position),
					dest))
		} else {
			log.Printf("entity ID %s not found in server entity lookup, could not build Move command", eid)
		}
	}
	return res
}

// AddMoveCommands
//
// Is expected to be called concurrently.
func (e *Executor) AddMoveCommands(req *apipb.MoveRequest) error {
	// TODO(minkezhang): If tick outside window, return error.
	var cs []Command

	for _, c := range e.buildMoveCommands(req.GetClientId(), e.tick(), req.GetDestination(), req.GetEntityIds()) {
		cs = append(cs, c)
	}

	return e.addCommands(cs)
}
