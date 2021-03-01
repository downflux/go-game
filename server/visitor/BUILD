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
    name = "move_test",
    srcs = ["move_test.go"],
    importpath = "github.com/downflux/game/server/visitor/move_test",
    embed = [":move"],
    deps = [
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/entity:acl",
        "//engine/fsm:action",
        "//engine/gamestate:dirty",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "//map:map",
        "//server/entity:tank",
        "//server/fsm:move",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "move",
    srcs = ["move.go"],
    importpath = "github.com/downflux/game/server/visitor/move",
    deps = [
        "//api:data_go_proto",
        "//engine/curve/common:linearmove",
        "//engine/entity:entity",
	"//engine/fsm/api:constants_go_proto",
        "//engine/gamestate:dirty",
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//map:map",
        "//map:utils", 
        "//pathing/hpf:astar",
        "//pathing/hpf:graph",
        "//engine/id:id",
        "//server/fsm:commonstate",
        "//server/fsm:move",
    ],
)

go_library(
    name = "chase",
    srcs = ["chase.go"],
    importpath = "github.com/downflux/game/server/visitor/chase",
    deps = [
        "//engine/status:status",
        "//engine/visitor:visitor",
        "//engine/fsm:schedule",
        "//server/fsm:chase",
        "//engine/fsm/api:constants_go_proto",
    ],
)

go_test(
    name = "chase_test",
    srcs = ["chase_test.go"],
    importpath = "github.com/downflux/game/server/visitor/chase_test",
    embed = [":chase"],
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
