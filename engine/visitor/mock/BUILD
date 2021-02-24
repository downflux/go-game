load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "simple",
    srcs = ["simple.go"],
    importpath = "github.com/downflux/game/engine/visitor/mock/simple",
    deps = [
        "//engine/visitor:visitor",
        "//engine/visitor/api:constants_go_proto",
    ],
)

go_test(
    name = "simple_test",
    srcs = [":simple_test.go"],
    importpath = "github.com/downflux/game/engine/visitor/mock/simple_test",
    embed = [":simple"],
    deps = [
        "//engine/visitor:visitor",
    ],
)
