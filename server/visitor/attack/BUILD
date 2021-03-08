load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//server:__subpackages__"])

go_library(
    name = "projectile",
    srcs = [":projectile.go"],
    importpath = "github.com/downflux/game/server/visitor/attack/projectile",
    deps = [
        "//engine/fsm/api:constants_go_proto",
        "//engine/gamestate:dirty",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//server/fsm/attack:projectile",
        "//server/fsm:commonstate",
    ],
)

go_test(
    name = "projectile_test",
    srcs = ["projectile_test.go"],
    importpath = "github.com/downflux/game/server/visitor/attack/projectile_test",
    embed = [":projectile"],
    deps = [
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/gamestate:dirty",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//server/entity:projectile",
        "//server/entity:tank",
        "//server/fsm:commonstate",
        "//server/fsm/attack:projectile",
        "//server/fsm/move:move",
	"@com_github_google_go_cmp//cmp:go_default_library",
    ],
)

go_test(
    name = "attack_test",
    srcs = ["attack_test.go"],
    importpath = "github.com/downflux/game/server/visitor/attack/attack_test",
    embed = [":attack"],
    deps = [
        "//engine/visitor:visitor",
    ],
)

go_library(
    name = "attack",
    srcs = ["attack.go"],
    importpath = "github.com/downflux/game/server/visitor/attack/attack",
    deps = [
        "//engine/fsm/api:constants_go_proto",
        "//engine/gamestate:dirty",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//server/fsm:commonstate",
        "//server/fsm/attack:attack",
    ],
)
