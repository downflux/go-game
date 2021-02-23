load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//engine:__subpackages__"])

go_library(
    name = "simple",
    srcs = [":simple.go"],
    importpath = "github.com/downflux/game/engine/fsm/mock/simple",
    deps = [
        "//engine/fsm:action",
        "//engine/fsm/api:constants_go_proto",
        "//engine/fsm:fsm",
        "//engine/id:id",
        "//engine/visitor:visitor",
    ]
)

go_test(
    name = "simple_test",
    srcs = ["simple_test.go"],
    importpath = "github.com/downflux/game/engine/fsm/mock/simple_test",
    embed = [":simple"],
    deps = [
        "//engine/fsm:action",
        "//engine/id:id",
    ],
)
