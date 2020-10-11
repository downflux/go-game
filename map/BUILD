load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test", "go_binary")

go_library(
    name = "map",
    srcs = ["map.go"],
    importpath = "github.com/downflux/game/map/map",
    deps = [
        ":utils",
        "//api:data_go_proto",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "@com_github_fzipp_astar//:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "map_test",
    srcs = ["map_test.go"],
    importpath = "github.com/downflux/game/map/map_test",
    embed = [":map"],
    deps = [
        ":utils",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "utils",
    srcs = ["utils.go"],
    importpath = "github.com/downflux/game/map/utils",
    deps = [
        "//api:data_go_proto",
    ],
)

go_library(
    name = "astar",
    srcs = ["astar.go"],
    importpath = "github.com/downflux/game/map//astar",
    deps = [
        ":tile",
        ":utils",
        "//api:data_go_proto",
        "//map/api:constants_go_proto",
        "@com_github_fzipp_astar//:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "astar_test",
    srcs = ["astar_test.go"],
    importpath = "github.com/downflux/game/map//astar_test",
    embed = [":astar"],
    deps = [
        ":tile",
        ":utils",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)
