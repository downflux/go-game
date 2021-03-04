load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//server:__subpackages__"])

go_library(
    name = "produce",
    srcs = ["produce.go"],
    importpath = "github.com/downflux/game/server/fsm/produce",
    deps = [
        ":commonstate",
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:visitor",
    ],
)

go_test(
    name = "produce_test",
    srcs = ["produce_test.go"],
    importpath = "github.com/downflux/game/server/fsm/produce_test",
    embed = [":produce"],
    deps = [
        ":commonstate",
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/status:status",
    ],
)

go_library(
    name = "commonstate",
    srcs = ["commonstate.go"],
    importpath = "github.com/downflux/game/server/fsm/commonstate",
    deps = [
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
    ],
)

go_library(
    name = "attack",
    srcs = ["attack.go"],
    importpath = "github.com/downflux/game/server/fsm/attack",
    deps = [
        ":commonstate",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//map:utils",
        "//server/entity/component:attackable",
        "//server/entity/component:targetable",
        "//server/fsm/move:chase",
    ],
)

go_test(
    name = "attack_test",
    srcs = ["attack_test.go"],
    importpath = "github.com/downflux/game/server/fsm/attack_test",
    embed = [":attack"],
    deps = [
        ":commonstate",
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/id:id",
        "//engine/status:status",
        "//server/entity:tank",
        "//server/fsm/move:chase",
    ],
)
