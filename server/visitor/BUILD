load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//server:__subpackages__"])

go_library(
    name = "produce",
    srcs = ["produce.go"],
    importpath = "github.com/downflux/game/server/visitor/produce",
    deps = [
        "//api:constants_go_proto",
        "//engine/entity:entity",
        "//engine/entity:list",
	"//engine/fsm/api:constants_go_proto",
        "//engine/gamestate:dirty",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//server/entity:tank",
        "//server/fsm:commonstate",
        "//server/fsm:produce",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ]
)

go_test(
    name = "produce_test",
    srcs = ["produce_test.go"],
    importpath = "github.com/downflux/game/server/visitor/produce_test",
    embed = [":produce"],
    deps = [
        "//engine/visitor:visitor",
    ],
)

go_test(
    name = "attack_test",
    srcs = ["attack_test.go"],
    importpath = "github.com/downflux/game/server/visitor/attack_test",
    embed = [":attack"],
    deps = [
        "//engine/visitor:visitor",
    ],
)

go_library(
    name = "attack",
    srcs = ["attack.go"],
    importpath = "github.com/downflux/game/server/visitor/attack",
    deps = [
        "//engine/gamestate:dirty",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//server/fsm:attack",
        "//server/fsm:commonstate",
        "//engine/fsm/api:constants_go_proto",
    ],
)
