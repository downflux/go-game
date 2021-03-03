load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//engine:__subpackages__"])

go_library(
    name = "simple",
    srcs = [":simple.go"],
    importpath = "github.com/downflux/game/engine/entity/mock/simple",
    deps = [
        "//api:constants_go_proto",
        "//engine/entity:entity",
        "//engine/entity/component:curve",
        "//engine/entity/component:lifecycle",
        "//engine/id:id",
    ]
)

go_test(
    name = "simple_test",
    srcs = [":simple_test.go"],
    importpath = "github.com/downflux/game/engine/entity/mock/simple_test",
    embed = [":simple"],
    deps = [
        "//engine/entity:entity",
    ],
)
