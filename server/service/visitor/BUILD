load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//server/service:__pkg__"])

go_library(
    name = "dirty",
    srcs = ["dirty.go"],
    importpath = "github.com/downflux/game/server/service/visitor/dirty",
    deps = [
        "//api:constants_go_proto",
    ],
)

go_library(
    name = "visitor",
    srcs = ["visitor.go"],
    importpath = "github.com/downflux/game/server/service/visitor/visitor",
    deps = [
        "//curve:curve",
    ],
)

go_library(
    name = "move",
    srcs = ["move.go"],
    importpath = "github.com/downflux/game/server/service/visitor/move",
    deps = [
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//curve:curve",
	"//curve:linearmove",
        "//entity:entity",
        "//map:map",
        "//map:utils",
        "//pathing/hpf:astar",
        "//pathing/hpf:graph",
    ],
)
