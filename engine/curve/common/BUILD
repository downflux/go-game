load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    "//entity:__pkg__",
    "//server/entity:__pkg__",
    "//server/server:__pkg__",
    "//server/visitor:__pkg__",
])

go_library(
    name = "linearmove",
    srcs = ["linearmove.go"],
    importpath = "github.com/downflux/game/engine/curve/common/linearmove",
    deps = [
        "//api:data_go_proto",
        "//api:constants_go_proto",
        "//engine/curve:curve",
        "//server:id",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "linearmove_test",
    srcs = ["linearmove_test.go"],
    importpath = "github.com/downflux/game/engine/curve/common/linearmove_test",
    embed = [":linearmove"],
    deps = [
        "//api:data_go_proto",
        "//server:id",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)
