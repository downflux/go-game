load("@io_bazel_rules_go//go:def.bzl", "go_library")

package(default_visibility=["//visibility:public"])

go_library(
    name = "lifecycle",
    srcs = ["lifecycle.go"],
    importpath = "github.com/downflux/game/engine/entity/component/lifecycle",
    deps = [
        "//engine/id:id",
    ],
)

go_library(
    name = "curve",
    srcs = ["curve.go"],
    importpath = "github.com/downflux/game/engine/entity/component/curve",
    deps = [
        "//engine/curve:list",
    ],
)
