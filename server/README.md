# DownFlux Server

We are abstracting away the communications layer of the DownFlux backend from
the actual implementation. This way, we have some flexibility in migrating to a
different server layer in the future if necessary, e.g. REST over HTTP/2.

The implementation of the backend may be found at
[//server/service](/server/service).

Currently, the only API layer implemented is gRPC, at
[//server/grpc](/server/grpc).
