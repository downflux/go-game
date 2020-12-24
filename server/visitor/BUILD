load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
	"//fsm:__pkg__",
	"//server/entity:__pkg__",
	"//server/service:__pkg__",
])

go_library(
    name = "produce",
    srcs = ["produce.go"],
    importpath = "github.com/downflux/game/server/visitor/produce",
    deps = [
        ":dirty",
        "//api:constants_go_proto",
        "//engine/entity:entity",
        "//engine/entity:list",
        "//engine/fsm:fsm",
        "//engine/fsm:instance",
        "//engine/fsm/api:constants_go_proto",
        "//engine/status:status",
        "//engine/visitor:visitor",
	"//engine/visitor/api:constants_go_proto",
        "//server:id",
        "//server/entity:tank",
        "//server/fsm:produce",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ]
)

go_test(
    name = "dirty_test",
    srcs = ["dirty_test.go"],
    embed = [":dirty"],
    importpath = "github.com/downflux/game/server/visitor/dirty_test",
    deps = [
        "//server:id",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)

go_library(
    name = "dirty",
    srcs = ["dirty.go"],
    importpath = "github.com/downflux/game/server/visitor/dirty",
    deps = [
        "//api:constants_go_proto",
        "//server:id",
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
        "//engine/fsm:instance",
        "//engine/status:status",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "//map:map",
        "//server/entity:tank",
        "//server/fsm:move",
        "//server/visitor:dirty",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "move",
    srcs = ["move.go"],
    importpath = "github.com/downflux/game/server/visitor/move",
    deps = [
        ":dirty",
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/curve/common:linearmove",
        "//engine/entity:entity",
        "//engine/status:status",
        "//engine/visitor:visitor",
	"//engine/visitor/api:constants_go_proto",
        "//engine/fsm:fsm",
        "//engine/fsm:instance",
        "//engine/fsm/api:constants_go_proto",
        "//map:map",
        "//map:utils", 
        "//pathing/hpf:astar",
        "//pathing/hpf:graph",
        "//server:id",
        "//server/fsm:move",
    ],
)
