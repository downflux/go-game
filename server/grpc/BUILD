load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "server",
    srcs = ["server.go"],
    importpath = "github.com/downflux/game/server/grpc/server",
    deps = [
        "//api:api_go_proto",
        "//server/service/commands:move",
        "//server/service:executor",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)
