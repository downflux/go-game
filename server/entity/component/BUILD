load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//server:__subpackages__",])

go_library(
    name = "moveable",
    srcs = ["moveable.go"],
    importpath = "github.com/downflux/game/server/entity/component/moveable",
    deps = [
        "//api:data_go_proto",
        "//engine/curve:curve",
        "//engine/id:id",
    ],
)

go_test(
    name = "moveable_test",
    srcs = ["moveable_test.go"],
    importpath = "github.com/downflux/game/server/entity/component/moveable_test",
    embed = [":moveable"],
)
