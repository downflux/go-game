load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//server:__subpackages__"])

go_library(
    name = "tank",
    srcs = ["tank.go"],
    importpath = "github.com/downflux/game/server/entity/tank",
    deps = [
	"//api:constants_go_proto",
	"//api:data_go_proto",
        "//engine/curve:curve",
        "//engine/curve:list",
        "//engine/curve/common:delta",
        "//engine/curve/common:linearmove",
        "//engine/curve/common:step",
        "//engine/curve/common:timer",
        "//engine/entity:entity",
        "//engine/entity/component:curve",
        "//engine/entity/component:lifecycle",
	"//engine/id:id",
        "//server/entity/component:attackable",
        "//server/entity/component:moveable",
        "//server/entity/component:positionable",
        "//server/entity/component:targetable",
    ],
)

go_test(
    name = "tank_test",
    srcs = ["tank_test.go"],
    importpath = "github.com/downflux/game/server/entity/tank_test",
    embed = [":tank"],
    deps = [
        "//engine/entity:entity",
        "//server/entity/component:attackable",
        "//server/entity/component:moveable",
        "//server/entity/component:positionable",
        "//server/entity/component:targetable",
    ],
)

go_library(
    name = "projectile",
    srcs = ["projectile.go"],
    importpath = "github.com/downflux/game/server/entity/projectile",
    deps = [
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/curve:curve",
        "//engine/curve:list",
        "//engine/curve/common:linearmove",
        "//engine/curve/common:step",
        "//engine/entity:entity",
        "//engine/entity/component:curve",
        "//engine/entity/component:lifecycle",
        "//engine/id:id",
        "//server/entity/component:moveable",
        "//server/entity/component:positionable",
    ],
)

go_test(
    name = "projectile_test",
    srcs = ["projectile_test.go"],
    importpath = "github.com/downflux/game/server/entity/projectile_test",
    embed = [":projectile"],
    deps = [
        "//engine/entity:entity",
        "//server/entity/component:moveable",
        "//server/entity/component:positionable",
    ],
)
