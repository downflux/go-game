load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "github.com/downflux/game/server/service/list",
    deps = [
        ":client",
        "//api:api_go_proto",
        "//server:id",
        "//server/service/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)

go_library(
    name = "client",
    srcs = ["client.go"],
    importpath = "github.com/downflux/game/server/service/client",
    deps = [
        "//api:api_go_proto",
        "//server:id",
        "//server/service/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "client_test",
    srcs = ["client_test.go"],
    importpath = "github.com/downflux/game/server/service/client_test",
    embed = [":client"],
    deps = [
        "//api:api_go_proto",
        "//api:data_go_proto",
        "//server:id",
        "//server/service/api:constants_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)

go_library(
    name = "executor",
    srcs = ["executor.go"],
    importpath = "github.com/downflux/game/server/service/executor",
    visibility = ["//server/grpc:__pkg__"],
    deps = [
        ":list",
        "//api:api_go_proto",
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/entity:list",
        "//engine/status:status",
        "//engine/visitor:list",
        "//engine/visitor:visitor",
        "//engine/visitor/api:constants_go_proto",
        "//engine/fsm:schedule",
        "//engine/fsm/api:constants_go_proto",
        "//map:map",
        "//map/api:data_go_proto",
        "//pathing/hpf:graph",
        "//server:id",
        "//server/fsm:move",
        "//server/fsm:produce",
        "//server/visitor:dirty",
        "//server/visitor:move",
        "//server/visitor:produce",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "executor_test",
    srcs = [":executor_test.go"],
    importpath = "github.com/downflux/game/server/service/executor_test",
    embed = [":executor"],
    deps = [
        "//api:api_go_proto",
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "//server:id",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)
