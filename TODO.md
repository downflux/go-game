# Bugs

* [ ] check why Unity game exiting does not trigger server stream exit
* [ ] check why server stream uses 1 CPU per stream

# Refactoring

* [x] change MoveRequest.Destination to Coordinate
* [x] replace tickID with tick for now
* [x] add AddClientResponse.Tick
* [x] add oneof curve subclass
* [x] make entity interface more formal
* [x] broadcast entities along with curves in StreamCurves
* [x] refactor server / executor
* [x] make Executor.visitors an ordered map object for easy referencing in code
* [x] merge Executor.isStarted and isStopped
* [x] add ClientID, EntityID, Tick data types
* [x] add server documentation
* [x] make `maps/astar.Path` take in `Position`; last move is to the Position offset
* [x] ensure client still works with refactors
* [x] add tests for move.Visitor
* [x] rename CurveCategory to EntityProperty
* [x] migrate to FSM Visitor model
* [x] migrate engine components
* [x] migrate instance -> action
* [x] add Entity Implements layer
* [ ] add generator for IDs instead of manually testing for random string
* [ ] add server / engine documentation
* [ ] add API documentation
* [ ] add tests for produce.Visitor
* [ ] add tests for chase.Visitor
* [ ] add tests for Status
* [ ] add tests for Entity
* [ ] add client documentation
* [ ] add client tests

# Feature
* [x] make channel send from server nonblocking
* [x] make channel send from server with timeout
* [x] support Reconnect() and share game state at specific tick
* [x] add Run() which iterates through Tick()
* [x] consider what to do with partial moves -- refactor out commandQueue with cancellation, etc.
* [x] add Produce() Visitor
* [x] add entity click-drag-select action
* [x] add Attack command
* [ ] **add AttackTarget curve**
* [ ] add FSM design document
* [ ] add game state export / import
* [ ] log to file, log non-fatal errors instead of erroring out (e.g. `Run`)
* [ ] replace pathfinding with flow fields
* [ ] parallelize move commands
* [ ] make client load map from server (new API)
* [ ] make map layered (terrain vs. rendering data; pathing data for size > `1x1` tiles)
