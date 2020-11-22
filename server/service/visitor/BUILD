load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
	"//server/service:__pkg__",
	"//server/service/visitor/entity:__pkg__",
])

go_test(
    name = "dirty_test",
    srcs = ["dirty_test.go"],
    embed = [":dirty"],
    importpath = "github.com/downflux/game/server/service/visitor/dirty_test",
    deps = [
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)

go_library(
    name = "dirty",
    srcs = ["dirty.go"],
    importpath = "github.com/downflux/game/server/service/visitor/dirty",
    deps = [
        "//api:constants_go_proto",
    ],
)

go_library(
    name = "visitor",
    srcs = ["visitor.go"],
    importpath = "github.com/downflux/game/server/service/visitor/visitor",
    deps = [
        "//api:constants_go_proto",
    ],
)

go_test(
    name = "move_test",
    srcs = ["move_test.go"],
    importpath = "github.com/downflux/game/server/service/visitor/move_test",
    embed = [":move"],
    deps = [
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "//map:map",
        "//server/service:status",
        "//server/service/visitor:dirty",
        "//server/service/visitor/entity:tank",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)

go_library(
    name = "move",
    srcs = ["move.go"],
    importpath = "github.com/downflux/game/server/service/visitor/move",
    deps = [
        ":dirty",
        "//api:constants_go_proto",
        "//api:data_go_proto",
	"//curve:linearmove",
        "//entity:entity",
        "//map:map",
        "//map:utils",
        "//pathing/hpf:astar",
        "//pathing/hpf:graph",
	"//server/service:status",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
