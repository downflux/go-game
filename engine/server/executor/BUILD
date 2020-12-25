load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "executor",
    srcs = ["executor.go"],
    importpath = "github.com/downflux/game/engine/server/executor/executor",
    deps = [
        "//api:api_go_proto",
        "//api:data_go_proto",
        "//engine/gamestate:dirty",
        "//engine/gamestate:gamestate",
        "//engine/id:id",
        "//engine/fsm:action",
        "//engine/fsm:schedule",
        "//engine/fsm/api:constants_go_proto",
        "//engine/server/client:list",
        "//engine/visitor:list",
        "//engine/visitor/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "executor_test",
    srcs = [":executor_test.go"],
    importpath = "github.com/downflux/game/engine/server/executor/executor_test",
    embed = [":executor"],
    deps = [
        "//api:api_go_proto",
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/entity:list",
        "//engine/fsm/api:constants_go_proto",
        "//engine/gamestate:dirty",
        "//engine/gamestate:gamestate",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:list",
        "//engine/visitor:visitor",
        "//engine/visitor/api:constants_go_proto",
        "//map:map",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "//pathing/hpf:graph",
        "//server/fsm:move",
        "//server/fsm:produce",
        "//server/visitor:move",
        "//server/visitor:produce",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)
