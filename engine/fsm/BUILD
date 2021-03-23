load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_test(
    name = "schedule_test",
    srcs = ["schedule_test.go"],
    importpath = "github.com/downflux/game/engine/fsm/schedule_test",
    embed = [":schedule"],
    deps = [
	":action",
        "//engine/fsm/api:constants_go_proto",
        "//engine/fsm/mock:simple",
        "//engine/id:id",
    ],
)

go_library(
    name = "schedule",
    srcs = ["schedule.go"],
    importpath = "github.com/downflux/game/engine/fsm/schedule",
    deps = [
        ":action",
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
        ":action",
        ":fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/fsm/mock:dependent",
        "//engine/fsm/mock:simple",
        "//engine/id:id",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
    ],
)

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "github.com/downflux/game/engine/fsm/list",
    deps = [
        ":action",
        ":fsm",
        "//engine/id:id",
        "//engine/fsm/api:constants_go_proto",
        "//engine/visitor:visitor",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)

go_library(
    name = "action",
    srcs = ["action.go"],
    importpath = "github.com/downflux/game/engine/fsm/action",
    deps = [
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/visitor:visitor",
        "//engine/id:id",
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
