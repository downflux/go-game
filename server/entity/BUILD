load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    # This is only true for intermediate testing.
    # TODO(minkezhang): Remove visibility.
    "//server/grpc:__pkg__",

    "//fsm:__pkg__",
    "//server/service:__pkg__",
    "//server/visitor:__pkg__",
])

go_library(
    name = "tank",
    srcs = ["tank.go"],
    importpath = "github.com/downflux/game/server/entity/tank",
    deps = [
	"//api:constants_go_proto",
	"//api:data_go_proto",
        "//curve:linearmove",
        "//engine/curve:curve",
        "//engine/entity:entity",
	"//server:id",
    ],
)
