load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//server:__subpackages__"])

go_library(
    name = "projectile",
    srcs = ["projectile.go"],
    importpath = "github.com/downflux/game/server/fsm/attack/projectile",
    deps = [
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "//engine/visitor:visitor",
        "//server/entity/component:attackable",
        "//server/entity/component:targetable",
        "//server/fsm:commonstate",
        "//server/fsm/move:move",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "projectile_test",
    srcs = ["projectile_test.go"],
    importpath = "github.com/downflux/game/server/fsm/attack/projectile_test",
    embed = [":projectile"],
    deps = [
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/id:id",
        "//engine/status:status",
        "//server/entity:projectile",
	"//server/entity:tank",
        "//server/fsm:commonstate",
        "//server/fsm/move:move",
    ],
)

go_library(
    name = "attack",
    srcs = ["attack.go"],
    importpath = "github.com/downflux/game/server/fsm/attack/attack",
    deps = [
        ":projectile",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//map:utils",
        "//server/entity/component:attackable",
        "//server/entity/component:targetable",
        "//server/fsm:commonstate",
        "//server/fsm/move:chase",
    ],
)

go_test(
    name = "attack_test",
    srcs = ["attack_test.go"],
    importpath = "github.com/downflux/game/server/fsm/attack/attack_test",
    embed = [":attack"],
    deps = [
        ":projectile",
        "//api:data_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/id:id",
        "//engine/status:status",
        "//server/entity:projectile",
        "//server/entity:tank",
        "//server/fsm:commonstate",
        "//server/fsm/move:chase",
        "//server/fsm/move:move",
    ],
)
