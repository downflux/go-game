# DownFlux Plugins

## gRPC DLLs

gRPC requires some additionally packaged object files, e.g.
`libgrpc_csharp_ext.so` in order for Unity to run. These shared libraries may be
found in the `Grpc.Core` NuGet package.

Remember to rename the runtime to the expected filename.

See [#25223](https://github.com/grpc/grpc/issues/25223) for more information.
