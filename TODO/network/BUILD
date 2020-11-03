load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library", "go_grpc_library")

proto_library(
    name = "demo_proto",
    srcs = ["demo.proto"],
    deps = [],
)

go_grpc_library(
    name = "demo_go_proto",
    importpath = "github.com/downflux/game/TODO/network/demo_go_proto",
    proto = ":demo_proto",
    deps = [],
)

go_binary(
    name = "server",
    importpath = "github.com/downflux/game/TODO/network/server",
    srcs = ["server.go"],
    deps = [
        ":demo_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//keepalive:go_default_library",
    ],
)

go_binary(
    name = "client",
    importpath = "github.com/downflux/game/TODO/network/client",
    srcs = ["client.go"],
    deps = [
        ":demo_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//keepalive:go_default_library",
    ],
)
