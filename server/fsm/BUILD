load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    "//engine:__subpackages__",  # test-only
    "//server:__subpackages__",
])

go_test(
    name = "move_test",
    srcs = ["move_test.go"],
    importpath = "github.com/downflux/game/server/fsm/move_test",
    embed = [":move"],
    deps = [
        ":commonstate",
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/id:id",
        "//engine/status:status",
        "//server/entity:tank",
    ],
)

go_library(
    name = "move",
    srcs = ["move.go"],
    importpath = "github.com/downflux/game/server/fsm/move",
    deps = [
        ":commonstate",
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//server/entity/component:moveable",
        "@org_golang_google_protobuf//proto:go_default_library",
    ],
)

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
    name = "chase",
    srcs = ["chase.go"],
    importpath = "github.com/downflux/game/server/fsm/chase",
    deps = [
        "//engine/entity:entity",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "//engine/visitor:visitor",
        "//server/entity/component:moveable",
        "//server/fsm:commonstate",
        "//server/fsm:move",
    ],
)

go_test(
    name = "chase_test",
    srcs = ["chase_test.go"],
    importpath = "github.com/downflux/game/server/fsm/chase_test",
    embed = [":chase"],
    deps = [
        "//engine/fsm:action",
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
        ":chase",
    ],
)

go_test(
    name = "attack_test",
    srcs = ["attack_test.go"],
    importpath = "github.com/downflux/game/server/fsm/attack_test",
    embed = [":attack"],
    deps = [
        "//engine/fsm:action",
    ],
)
