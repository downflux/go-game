load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "client",
    srcs = ["client.go"],
    importpath = "github.com/downflux/game/engine/server/client/client",
    deps = [
        "//api:api_go_proto",
        "//engine/fsm:action",
        "//engine/fsm:fsm",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "//engine/server/client/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "client_test",
    srcs = ["client_test.go"],
    importpath = "github.com/downflux/game/engine/server/client/client_test",
    embed = [":client"],
    deps = [
        "//engine/fsm:fsm",
        "//engine/server/client/api:constants_go_proto",
        "//api:api_go_proto",
        "//api:data_go_proto",
        "//engine/id:id",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "github.com/downflux/game/engine/server/client/list",
    deps = [
        ":client",
        "//api:api_go_proto",
        "//engine/fsm:fsm",
        "//engine/id:id",
        "//engine/server/client/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)
