load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    # This is only true for intermediate testing.
    # TODO(minkezhang): Remove visibility.
    "//engine/fsm:__pkg__",
    "//server/grpc:__pkg__",
    "//server/fsm:__pkg__",

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
        "//engine/curve:curve",
        "//engine/curve:list",
        "//engine/curve/common:linearmove",
        "//engine/entity:entity",
	"//engine/id:id",
        "//server/entity/component:attackable",
        "//server/entity/component:moveable",
        "//server/entity/component:positionable",
        "//server/entity/component:targetable",
    ],
)

go_test(
    name = "tank_test",
    srcs = ["tank_test.go"],
    importpath = "github.com/downflux/game/server/entity/tank_test",
    embed = [":tank"],
    deps = [
        "//engine/entity:entity",
        "//server/entity/component:attackable",
        "//server/entity/component:moveable",
        "//server/entity/component:positionable",
        "//server/entity/component:targetable",
    ],
)
