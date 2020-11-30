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
* [ ] rename CurveCategory to EntityProperty
* [ ] **add server documentation**
* [ ] add API documentation
* [ ] ensure client still works with refactors
* [ ] add tests for Visitor
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
* [ ] parallelize move commands
* [ ] log to file, log non-fatal errors instead of erroring out (e.g. `Run`)
* [ ] add entity click-drag-select action
