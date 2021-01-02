load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    "//engine:__subpackages__",  # Test-only.

    "//server:__subpackages__",
])

go_library(
    name = "moveable",
    srcs = ["moveable.go"],
    importpath = "github.com/downflux/game/server/entity/component/moveable",
    deps = [
        ":positionable",
        "//engine/curve:curve",
        "//engine/id:id",
    ],
)

go_library(
    name = "attackable",
    srcs = ["attackable.go"],
    importpath = "github.com/downflux/game/server/entity/component/attackable",
    deps = [
        "//api:data_go_proto",
        "//engine/curve:curve",
        "//engine/id:id",
    ],
)

go_library(
    name = "targetable",
    srcs = ["targetable.go"],
    importpath = "github.com/downflux/game/server/entity/component/targetable",
)

go_library(
    name = "positionable",
    srcs = ["positionable.go"],
    importpath = "github.com/downflux/game/server/entity/component/positionable",
    deps = [
        "//api:data_go_proto",
        "//engine/curve:curve",
        "//engine/id:id",
    ],
)
