load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    "//server/service:__pkg__",
])

go_library(
    name = "entity",
    srcs = ["entity.go"],
    importpath = "github.com/downflux/game/entity/entity",
    deps = [
        "//api:constants_go_proto",
    ],
)

