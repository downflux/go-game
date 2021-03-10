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
