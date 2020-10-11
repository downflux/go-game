load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "executor",
    srcs = ["executor.go"],
    importpath = "github.com/downflux/game/server/service/executor",
    visibility = ["//server/grpc:__pkg__"],
    deps = [
        "//server/service/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
