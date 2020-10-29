package executor

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/downflux/game/entity/entity"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/downflux/game/server/id"
	"github.com/downflux/game/server/service/command/command"
	"github.com/downflux/game/server/service/command/move"
	"golang.org/x/sync/errgroup"
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
	idLen        = 8
	tickDuration = 100 * time.Millisecond
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

type dirtyCurve struct {
	eid      string
	category gcpb.CurveCategory
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
	tickImpl      int64

	// Add-only. Acquire first.
	dataMux  sync.RWMutex
	entities map[string]entity.Entity

	// Add and delete. Reset per tick.
	commandQueueMux sync.Mutex
	commandQueue    []command.Command

	// Add and delete. Reset per tick. Acquire second.
	entityQueueMux sync.RWMutex
	entityQueue    []string

	// Add and delete. Reset per tick. Acquire last.
	curveQueueMux sync.RWMutex
	curveQueue    []dirtyCurve

	clientChannelMux sync.RWMutex
	clientChannel    map[string]chan *apipb.StreamCurvesResponse
}

func (e *Executor) tick() float64   { return float64(atomic.LoadInt64(&(e.tickImpl))) }
func (e *Executor) incrementTick()  { atomic.AddInt64(&(e.tickImpl), 1) }
func (e *Executor) isStarted() bool { return atomic.LoadInt32(&(e.isStartedImpl)) != 0 }
func (e *Executor) setIsStarted()   { atomic.StoreInt32(&(e.isStartedImpl), 1) }
func (e *Executor) isStopped() bool { return atomic.LoadInt32(&(e.isStoppedImpl)) != 0 }
func (e *Executor) setIsStopped()   { atomic.StoreInt32(&(e.isStoppedImpl), 1) }

func (e *Executor) Status() *gdpb.ServerStatus {
	return &gdpb.ServerStatus{
		Tick:         e.tick(),
		IsStarted:    e.isStarted(),
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

func (e *Executor) popCommandQueue() ([]command.Command, error) {
	e.commandQueueMux.Lock()
	defer e.commandQueueMux.Unlock()

	commands := e.commandQueue
	e.commandQueue = nil
	return commands, nil
}

func (e *Executor) processCommand(cmd command.Command) error {
	if cmd.Type() == sscpb.CommandType_COMMAND_TYPE_MOVE {
		c, err := cmd.Execute(move.Args{
			Tick:   e.tick(),
			Source: e.entities[cmd.(*move.Command).EntityID()].Curve(gcpb.CurveCategory_CURVE_CATEGORY_MOVE).Get(e.tick()).(*gdpb.Position)})
		if err != nil {
			return err
		}

		if err := func() error {
			e.dataMux.RLock()
			defer e.dataMux.RUnlock()

			if err := e.entities[c.EntityID()].Curve(c.Category()).ReplaceTail(c); err != nil {
				return err
			}

			e.curveQueueMux.Lock()
			e.curveQueue = append(e.curveQueue, dirtyCurve{
				eid:      c.EntityID(),
				category: c.Category(),
			})
			e.curveQueueMux.Unlock()

			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) popTickQueue() ([]*gdpb.Curve, []*gdpb.Entity) {
	e.dataMux.Lock()
	e.entityQueueMux.Lock()
	e.curveQueueMux.Lock()
	defer e.curveQueueMux.Unlock()
	defer e.entityQueueMux.Unlock()
	defer e.dataMux.Unlock()

	var processedCurves = map[string]map[gcpb.CurveCategory]bool{}

	var curves []*gdpb.Curve
	var entities []*gdpb.Entity

	// TODO(minkezhang): Make concurrent.
	for _, eid := range e.entityQueue {
		entities = append(entities, &gdpb.Entity{
			EntityId: eid,
			Type:     e.entities[eid].Type(),
		})
	}
	for _, dc := range e.curveQueue {
		if _, found := processedCurves[dc.eid]; !found {
			processedCurves[dc.eid] = map[gcpb.CurveCategory]bool{}
		}
		if _, found := processedCurves[dc.eid][dc.category]; !found {
			processedCurves[dc.eid][dc.category] = true
			// TODO(minkezhang): Consider and implement what
			// happens on late / re-sync. Simply populate e.tick()
			// with last known good client tick, and mark all
			// curves and entities as dirty.
			curves = append(curves, e.entities[dc.eid].Curve(dc.category).ExportTail(e.tick()))
		}
	}

	e.curveQueue = nil
	e.entityQueue = nil

	return curves, entities
}

func (e *Executor) broadcastCurves() error {
	curves, entities := e.popTickQueue()

	if curves == nil && entities == nil {
		return nil
	}

	// TODO(minkezhang): Decide if it's okay that the reported tick may not
	// coincide with the ticks of the curve and entities.
	resp := &apipb.StreamCurvesResponse{
		Tick:     e.tick(),
		Curves:   curves,
		Entities: entities,
	}

	e.clientChannelMux.RLock()
	defer e.clientChannelMux.RUnlock()

	var eg errgroup.Group
	for _, ch := range e.clientChannel {
		ch := ch
		eg.Go(func() error {
			ch <- resp
			// TODO(minkezhang): Add timeout support.
			// Will need to implement resync logic once timeout is
			// added.
			return nil
		})
	}
	return eg.Wait()
}

func (e *Executor) closeStreams() error {
	e.clientChannelMux.Lock()
	defer e.clientChannelMux.Unlock()

	for _, ch := range e.clientChannel {
		close(ch)
	}
	return nil
}

func (e *Executor) Stop() {
	e.setIsStopped()
	e.closeStreams()
}

func (e *Executor) Run() error {
	e.setIsStarted()
	for !e.isStopped() {
		t := time.Now()
		e.incrementTick()

		if err := e.doTick(); err != nil {
			return err
		}

		// TODO(minkezhang): Add metrics collection here for tick
		// distribution.
		if d := time.Now().Sub(t); d < tickDuration {
			time.Sleep(tickDuration - d)
		}
	}
	return nil
}

func (e *Executor) doTick() error {
	commands, err := e.popCommandQueue()
	if err != nil {
		return err
	}

	log.Printf("[%.f] processing commands", e.tick())
	for _, cmd := range commands {
		// TODO(minkezhang): Add actual error handling here -- only
		// Only return early if error is very bad.
		if err := e.processCommand(cmd); err != nil {
			log.Printf("[%.f] error while processing command %v: %v", cmd, err)
			return err
		}
	}

	log.Printf("[%.f] broadcasting curves", e.tick())
	if err := e.broadcastCurves(); err != nil {
		return err
	}

	return nil
}

// TODO(minkezhang): Make this method private -- this is currently public for
// debugging purposes.
// TODO(minkezhang): Make this generate a Command instead.
func (e *Executor) AddEntity(en entity.Entity) error {
	e.dataMux.Lock()
	e.entityQueueMux.Lock()
	e.curveQueueMux.Lock()
	defer e.curveQueueMux.Unlock()
	defer e.entityQueueMux.Unlock()
	defer e.dataMux.Unlock()

	if _, found := e.entities[en.ID()]; found {
		return status.Errorf(codes.AlreadyExists, "given entity ID %v already exists in the entity list", en.ID())
	}

	e.entities[en.ID()] = en

	e.entityQueue = append(e.entityQueue, en.ID())
	for _, cat := range en.CurveCategories() {
		e.curveQueue = append(e.curveQueue, dirtyCurve{
			eid:      en.ID(),
			category: cat,
		})
	}

	return nil
}

func (e *Executor) addCommands(cs []command.Command) error {
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
func (e *Executor) buildMoveCommands(cid string, dest *gdpb.Position, eids []string) []*move.Command {
	e.dataMux.RLock()
	defer e.dataMux.RUnlock()

	var res []*move.Command
	for _, eid := range eids {
		_, found := e.entities[eid]
		if found {
			res = append(
				res,
				move.New(
					e.tileMap,
					e.abstractGraph,
					cid,
					eid,
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
	var cs []command.Command

	for _, c := range e.buildMoveCommands(req.GetClientId(), req.GetDestination(), req.GetEntityIds()) {
		cs = append(cs, c)
	}

	return e.addCommands(cs)
}
