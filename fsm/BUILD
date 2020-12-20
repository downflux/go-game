load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    "//fsm/instance:__pkg__",
    "//server/service:__pkg__",
    "//server/schedule:__pkg__",
    "//server/visitor:__pkg__",
])

go_test(
    name = "schedule_test",
    srcs = ["schedule_test.go"],
    importpath = "github.com/downflux/game/fsm/schedule_test",
    embed = [":schedule"],
    deps = [
	":instance",
        ":move",
        "//fsm/api:constants_go_proto",
        "//server/entity:tank",
        "//server/service:status",
    ],
)

go_library(
    name = "schedule",
    srcs = ["schedule.go"],
    importpath = "github.com/downflux/game/fsm/schedule",
    deps = [
        "//fsm/api:constants_go_proto",
        "//fsm:instance",
        "//fsm:list",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "list_test",
    srcs = ["list_test.go"],
    importpath = "github.com/downflux/game/fsm/list_test",
    embed = [":list"],
    deps = [
        ":fsm",
        ":instance",
        ":move",
        "//api:data_go_proto",
        "//fsm/api:constants_go_proto",
        "//server/entity:tank",
        "//server/service:status",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "github.com/downflux/game/fsm/list",
    deps = [
        ":fsm",
        ":instance",
        "//fsm/api:constants_go_proto",
        "//server:id",
        "//server/visitor/api:constants_go_proto",
        "//server/visitor:visitor",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)

go_library(
    name = "instance",
    srcs = ["instance.go"],
    importpath = "github.com/downflux/game/fsm/instance",
    deps = [
        "//fsm/api:constants_go_proto",
        "//fsm:fsm",
        "//server:id",
        "//server/visitor/api:constants_go_proto",
        "//server/visitor:visitor",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_library(
    name = "move",
    srcs = ["move.go"],
    importpath = "github.com/downflux/game/fsm/move",
    deps = [
        ":fsm",
        ":instance",
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//fsm/api:constants_go_proto",
        "//server:id",
        "//server/entity:entity",
        "//server/service:status",
        "//server/visitor:visitor",
        "@org_golang_google_protobuf//proto:go_default_library",
    ],
)

go_library(
    name = "fsm",
    srcs = ["fsm.go"],
    importpath = "github.com/downflux/game/fsm/fsm",
    deps = [
        "//fsm/api:constants_go_proto",
    ],
)
