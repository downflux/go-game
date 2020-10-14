load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//server/service:__pkg__"])

go_library(
    name = "move",
    srcs = ["move.go"],
    importpath = "github.com/downflux/game/server/service/commands/move",
    deps = [
        "//api:data_go_proto",
	"//curve:curve",
	"//map:map",
	"//pathing/hpf:astar",
	"//pathing/hpf:graph",
        "//server/service/api:constants_go_proto",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],    
)
