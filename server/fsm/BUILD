load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    "//engine:__subpackages__",
    "//server/service:__pkg__",
    "//server/schedule:__pkg__",
    "//server/visitor:__pkg__",
])

go_test(
    name = "move_test",
    srcs = ["move_test.go"],
    importpath = "github.com/downflux/game/server/fsm/move_test",
    embed = [":move"],
    deps = [
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/status:status",
        "//server/entity:tank",
    ],
)

go_library(
    name = "move",
    srcs = ["move.go"],
    importpath = "github.com/downflux/game/server/fsm/move",
    deps = [
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/entity:entity",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "@org_golang_google_protobuf//proto:go_default_library",
    ],
)

go_library(
    name = "produce",
    srcs = ["produce.go"],
    importpath = "github.com/downflux/game/server/fsm/produce",
    deps = [
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
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/status:status",
    ],
)
