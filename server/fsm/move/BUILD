load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//server:__subpackages__"])

go_test(
    name = "move_test",
    srcs = ["move_test.go"],
    importpath = "github.com/downflux/game/server/fsm/move/move_test",
    embed = [":move"],
    deps = [
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/id:id",
        "//engine/status:status",
        "//server/entity:tank",
        "//server/fsm:commonstate",
    ],
)

go_library(
    name = "move",
    srcs = ["move.go"],
    importpath = "github.com/downflux/game/server/fsm/move/move",
    deps = [
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//server/entity/component:moveable",
        "//server/fsm:commonstate",
        "@org_golang_google_protobuf//proto:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_library(
    name = "chase",
    srcs = ["chase.go"],
    importpath = "github.com/downflux/game/server/fsm/move/chase",
    deps = [
        ":move",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//map:utils",
        "//server/entity/component:moveable",
        "//server/entity/component:targetable",
        "//server/fsm:commonstate",
    ],
)

go_test(
    name = "chase_test",
    srcs = ["chase_test.go"],
    importpath = "github.com/downflux/game/server/fsm/move/chase_test",
    embed = [":chase"],
    deps = [
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/id:id",
        "//engine/status:status",
        "//server/entity:tank",
        "//server/fsm:commonstate",
    ],
)
