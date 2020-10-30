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
	"google.golang.org/protobuf/types/known/timestamppb"

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
		clients: map[string]*Client{},
	}, nil
}

// TODO(minkezhang): Export out into separate module.
type Client struct {
	mux sync.Mutex
	id string  // read-only
	ch chan *apipb.StreamDataResponse
	isSynced bool
}

func (c *Client) ID() string {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.id
}
func (c *Client) Channel() chan *apipb.StreamDataResponse {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.ch
}
func (c *Client) NewChannel() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.ch = make(chan *apipb.StreamDataResponse)
}
func (c *Client) CloseChannel() {
	c.mux.Lock()
	defer c.mux.Unlock()

	close(c.ch)
	c.ch = nil
}
func (c *Client) IsSynced() bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.isSynced
}
func (c *Client) SetIsSynced(s bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.isSynced = s
	if !s {
		close(c.ch)
		c.ch = nil
	}
}
func NewClient(cid string) *Client {
	c := &Client{
		id: cid,
		isSynced: false,
	}
	c.NewChannel()
	return c
}

type Executor struct {
	tileMap       *tile.Map
	abstractGraph *graph.Graph

	isStoppedImpl int32
	isStartedImpl int32
	tickImpl      int64
	startTimeMux  sync.Mutex
	startTimeImpl time.Time

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

	clientsMux sync.RWMutex
	clients    map[string]*Client
}

func (e *Executor) tick() float64   { return float64(atomic.LoadInt64(&(e.tickImpl))) }
func (e *Executor) incrementTick()  { atomic.AddInt64(&(e.tickImpl), 1) }
func (e *Executor) isStarted() bool { return atomic.LoadInt32(&(e.isStartedImpl)) != 0 }
func (e *Executor) setIsStarted()   { atomic.StoreInt32(&(e.isStartedImpl), 1) }
func (e *Executor) isStopped() bool { return atomic.LoadInt32(&(e.isStoppedImpl)) != 0 }
func (e *Executor) setIsStopped()   { atomic.StoreInt32(&(e.isStoppedImpl), 1) }
func (e *Executor) startTime() time.Time {
	e.startTimeMux.Lock()
	defer e.startTimeMux.Unlock()

	return e.startTimeImpl
}
func (e *Executor) setStartTime() {
	e.startTimeMux.Lock()
	defer e.startTimeMux.Unlock()

	e.startTimeImpl = time.Now()
}

func (e *Executor) Status() *gdpb.ServerStatus {
	return &gdpb.ServerStatus{
		Tick:         e.tick(),
		IsStarted:    e.isStarted(),
		TickDuration: durationpb.New(tickDuration),
		StartTime:    timestamppb.New(e.startTime()),
	}
}

func (e *Executor) AddClient() (string, error) {
	log.Printf("DEBUG: Adding client")
	// TODO(minkezhang): Add maxClients check.
	e.clientsMux.Lock()
	defer e.clientsMux.Unlock()

	cid := id.RandomString(idLen)
	for _, found := e.clients[cid]; found; cid = id.RandomString(idLen) {
	}
	e.clients[cid] = NewClient(cid)

	return cid, nil
}

func (e *Executor) ClientChannel(cid string) <-chan *apipb.StreamDataResponse {
	e.clientsMux.RLock()
	defer e.clientsMux.RUnlock()

	return e.clients[cid].Channel()
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
	e.dataMux.RLock()
	e.entityQueueMux.Lock()
	e.curveQueueMux.Lock()
	defer e.curveQueueMux.Unlock()
	defer e.entityQueueMux.Unlock()
	defer e.dataMux.RUnlock()

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
		// Do not broadcast curve twice.
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

func (e *Executor) allCurvesAndEntities() ([]*gdpb.Curve, []*gdpb.Entity) {
	e.dataMux.RLock()
	defer e.dataMux.RUnlock()

        var curves []*gdpb.Curve
        var entities []*gdpb.Entity

	// TODO(minkezhang): Give some leeway here, broadcast a bit in the
	// past.
	beginningTick := e.tick()

	for _, en := range e.entities {
		entities = append(entities, &gdpb.Entity{
                        EntityId: en.ID(),
                        Type:     en.Type(),
                })
		for _, cat := range en.CurveCategories() {
			curves = append(curves, e.entities[en.ID()].Curve(cat).ExportTail(beginningTick))
		}
	}
	return curves, entities
}

func (e *Executor) broadcastCurves() error {
	curves, entities := e.popTickQueue()

	// TODO(minkezhang): Decide if it's okay that the reported tick may not
	// coincide with the ticks of the curve and entities.
	resp := &apipb.StreamDataResponse{
		Tick:     e.tick(),
		Curves:   curves,
		Entities: entities,
	}
	allResp := &apipb.StreamDataResponse{
		Tick: e.tick(),
	}

	e.clientsMux.RLock()
	defer e.clientsMux.RUnlock()

	var needFullState bool
	for _, c := range e.clients {
		needFullState = needFullState || !c.IsSynced()
	}
	if needFullState {
		allCurves, allEntities := e.allCurvesAndEntities()
		allResp.Curves = allCurves
		allResp.Entities = allEntities
	}

	if !needFullState && curves == nil && entities == nil {
		return nil
	}

	log.Printf("DEBUG: sending curves to %d clients", len(e.clients))
	var eg errgroup.Group
	for _, c := range e.clients {
		c := c
		ch := c.Channel()
		eg.Go(func() error {
			log.Printf("DEBUG: Attempting to send message to client %v", c)
			if c.IsSynced() {
				log.Printf("DEBUG: Sending response to a synced client: %v", resp)
				ch <- resp
			} else {
				log.Printf("DEBUG: Sending response to an unsynced client: %v", allResp)
				ch <- allResp
				c.SetIsSynced(true)
			}
			// TODO(minkezhang): Add timeout support.
			// Will need to implement resync logic once timeout is
			// added.
			return nil
		})
	}
	return eg.Wait()
}

func (e *Executor) closeStreams() error {
	e.clientsMux.Lock()
	defer e.clientsMux.Unlock()

	for _, c := range e.clients {
		c.CloseChannel()
	}
	return nil
}

func (e *Executor) Stop() {
	e.setIsStopped()
	e.closeStreams()
}

func (e *Executor) Run() error {
	e.setStartTime()
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

	for _, cmd := range commands {
		// TODO(minkezhang): Add actual error handling here -- only
		// Only return early if error is very bad.
		if err := e.processCommand(cmd); err != nil {
			return err
		}
	}

	log.Printf("[%.f] Broadcasting curves to all clients", e.tick())
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
