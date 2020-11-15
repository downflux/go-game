# Add client disconnect detection, salient points

* Okay for `broadcastCurves` to close channel -- only `broadcastCurves` and
  server shutdown should close, only `AddClient` should create channel
* `broadcastCurve may be send-non-blocking, just make sure to read in a
  blocking manner. Data may be sent out of order (due to concurrency) but is
  okay because data is add-only and curves have merge-idempotence.
* `broadcastCurves` will wait for TIMEOUT secs before closing
* `AddClient(old_id)` -- wait for TIMEOUT secs before checking for `IsSync` (?)
  if `!IsSync`, create channel, otherwise null op

## Heartbeats

* gRPC heartbeats
* custom bidi heartbeat implementation
* context.WithDeadline (?)
* [gRPC keepalive](https://pkg.go.dev/google.golang.org/grpc/keepalive)
