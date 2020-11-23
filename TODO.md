# Refactoring

* [x] change MoveRequest.Destination to Coordinate
* [x] replace tickID with tick for now
* [ ] add server documentation
* [ ] add client documentation
* [ ] add API documentation
* [ ] add client tests
* [x] add server tests
* [x] add AddClientResponse.Tick
* [x] add oneof curve subclass
* [x] make entity interface more formal
* [x] broadcast entities along with curves in StreamCurves
* [x] refactor server / executor

# Feature
* [x] make channel send from server nonblocking
* [x] make channel send from server with timeout
* [x] support Reconnect() and share game state at specific tick
* [ ] parallelize move commands
* [x] add Run() which iterates through Tick()
* [x] consider what to do with partial moves -- refactor out commandQueue with cancellation, etc.
* [ ] **add Produce() Visitor**
