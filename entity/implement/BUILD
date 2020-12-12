load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    "//server/entity:__pkg__",
])

go_library(
    name = "implement",
    srcs = ["implement.go"],
    importpath = "github.com/downflux/game/entity/implement/implement",
    deps = [
        "//entity/implement/api:constants_go_proto",
    ],
)

go_library(
    name = "moveable",
    srcs = ["moveable.go"],
    importpath = "github.com/downflux/game/entity/implement/moveable",
    deps = [
        "//server:id",
        "//api:data_go_proto",
    ],
)
