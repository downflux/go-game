load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
	"//fsm:__pkg__",
	"//server/entity:__pkg__",
	"//server/service:__pkg__",
])

go_library(
    name = "visitorlist",
    srcs = ["visitorlist.go"],
    importpath = "github.com/downflux/game/server/visitor/visitorlist",
    deps = [
        ":visitor",
	"//server/visitor/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_library(
    name = "produce",
    srcs = ["produce.go"],
    importpath = "github.com/downflux/game/server/visitor/produce",
    deps = [
        ":dirty",
        ":visitor",
        "//api:constants_go_proto",
        "//fsm:fsm",
        "//fsm:instance",
        "//fsm:produce",
        "//fsm/api:constants_go_proto",
        "//server:id",
        "//server/entity:entity",
        "//server/entity:entitylist",
        "//server/entity:tank",
        "//server/service:status",
	"//server/visitor/api:constants_go_proto",
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

go_library(
    name = "visitor",
    srcs = ["visitor.go"],
    importpath = "github.com/downflux/game/server/visitor/visitor",
    deps = [
	"//server/visitor/api:constants_go_proto",
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
        "//fsm:instance",
        "//fsm:move",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "//map:map",
        "//server/entity:tank",
        "//server/service:status",
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
        ":visitor",
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//curve:linearmove",
        "//fsm/api:constants_go_proto",
        "//fsm:fsm",
        "//fsm:instance",
        "//fsm:move",
        "//map:map",
        "//map:utils", 
        "//pathing/hpf:astar",
        "//pathing/hpf:graph",
        "//server:id",
        "//server/service:status",
	"//server/visitor/api:constants_go_proto",
    ],
)
