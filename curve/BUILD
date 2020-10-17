load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    "//entity:__pkg__",
    "//server/service:__subpackages__",
])

go_library(
    name = "curve",
    srcs = ["curve.go"],
    importpath = "github.com/downflux/game/curve/curve",
    deps = [
        "//api:constants_go_proto",
        "//api:data_go_proto",
    ],
)

go_library(
    name = "linearmove",
    srcs = ["linearmove.go"],
    importpath = "github.com/downflux/game/curve/linearmove",
    deps = [
        ":curve",
        "//api:data_go_proto",
        "//api:constants_go_proto",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "linearmove_test",
    srcs = ["linearmove_test.go"],
    importpath = "github.com/downflux/game/curve/linearmove_test",
    embed = [":linearmove"],
    deps = [
        "//api:data_go_proto",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)
