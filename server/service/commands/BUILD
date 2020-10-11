load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//server/grpc:__pkg__"])

go_library(
    name = "move",
    srcs = ["move.go"],
    importpath = "github.com/downflux/game/server/service/commands/move",
    deps = [
        "//api:api_go_proto",
        "//server/service/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],    
)
