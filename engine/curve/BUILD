load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "curve",
    srcs = ["curve.go"],
    importpath = "github.com/downflux/game/engine/curve/curve",
    deps = [
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/id:id",
    ],
)
