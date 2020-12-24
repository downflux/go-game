load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library", "go_grpc_library")

package(default_visibility=["//visibility:public"])

proto_library(
    name = "constants_proto",
    srcs = ["constants.proto"],
)

go_proto_library(
    name = "constants_go_proto",
    importpath = "github.com/downflux/game/engine/server/client/api/constants_go_proto",
    proto = ":constants_proto",
)
