load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    "//server/service:__subpackages__",
])

go_library(
    name = "curve",
    srcs = ["curve.go"],
    importpath = "github.com/downflux/game/curve/curve",
    deps = [
        "//api:constants_go_proto",
    ],
)

