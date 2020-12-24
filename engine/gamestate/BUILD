load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_test(
    name = "dirty_test",
    srcs = ["dirty_test.go"],
    embed = [":dirty"],
    importpath = "github.com/downflux/game/engine/gamestate/dirty",
    deps = [
        "//engine/id:id",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)

go_library(
    name = "dirty",
    srcs = ["dirty.go"],
    importpath = "github.com/downflux/game/engine/gamestate/dirty",
    deps = [
        "//api:constants_go_proto",
        "//engine/id:id",
    ],
)
