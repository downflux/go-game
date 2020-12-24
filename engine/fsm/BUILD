load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_test(
    name = "schedule_test",
    srcs = ["schedule_test.go"],
    importpath = "github.com/downflux/game/engine/fsm/schedule_test",
    embed = [":schedule"],
    deps = [
	":instance",
        "//engine/fsm/api:constants_go_proto",
        "//engine/status:status",
        "//fsm:move",
        "//server/entity:tank",
    ],
)

go_library(
    name = "schedule",
    srcs = ["schedule.go"],
    importpath = "github.com/downflux/game/engine/fsm/schedule",
    deps = [
        ":instance",
        ":list",
        "//engine/fsm/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "list_test",
    srcs = ["list_test.go"],
    importpath = "github.com/downflux/game/engine/fsm/list_test",
    embed = [":list"],
    deps = [
        ":fsm",
        ":instance",
        "//api:data_go_proto",
        "//engine/fsm/api:constants_go_proto",
        "//engine/status:status",
        "//fsm:move",
        "//server/entity:tank",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "github.com/downflux/game/engine/fsm/list",
    deps = [
        ":fsm",
        ":instance",
        "//engine/fsm/api:constants_go_proto",
        "//engine/visitor:visitor",
        "//engine/visitor/api:constants_go_proto",
        "//server:id",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)

go_library(
    name = "instance",
    srcs = ["instance.go"],
    importpath = "github.com/downflux/game/engine/fsm/instance",
    deps = [
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/visitor:visitor",
        "//engine/visitor/api:constants_go_proto",
        "//server:id",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_library(
    name = "fsm",
    srcs = ["fsm.go"],
    importpath = "github.com/downflux/game/engine/fsm/fsm",
    deps = [
        "//engine/fsm/api:constants_go_proto",
    ],
)
