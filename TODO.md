# Refactoring

* [x] change MoveRequest.Destination to Coordinate
* [x] replace tickID with tick for now
* [x] add server tests
* [x] add AddClientResponse.Tick
* [x] add oneof curve subclass
* [x] make entity interface more formal
* [x] broadcast entities along with curves in StreamCurves
* [x] refactor server / executor
* [x] make Executor.visitors an ordered map object for easy referencing in code
* [ ] **add server documentation**
* [ ] add ClientID, EntityID, CurveID data types
* [ ] add client documentation
* [ ] add API documentation
* [ ] add client tests
* [ ] rename CurveCategory to EntityProperty
* [ ] add tests for Visitors
* [ ] log to file, log non-fatal errors instead of erroring out (e.g. `Run`)

# Feature
* [x] make channel send from server nonblocking
* [x] make channel send from server with timeout
* [x] support Reconnect() and share game state at specific tick
* [x] add Run() which iterates through Tick()
* [x] consider what to do with partial moves -- refactor out commandQueue with cancellation, etc.
* [x] add Produce() Visitor
* [ ] parallelize move commands
