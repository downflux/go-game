load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test", "go_binary")

go_library(
    name = "astar",
    srcs = ["astar.go"],
    importpath = "github.com/downflux/game/pathing/hpf/astar",
    deps = [
        ":cluster",
        ":graph",
        ":graphastar",
        "//map:astar",
        "//map:map",
        "//map:utils",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "astar_test",
    srcs = ["astar_test.go"],
    importpath = "github.com/downflux/game/pathing/hpf/astar_test",
    embed = [":astar"],
    deps = [
        "//api:data_go_proto",
        "//map:map",
        "//map:utils",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "graphastar",
    srcs = ["graphastar.go"],
    importpath = "github.com/downflux/game/pathing/hpf/graphastar",
    deps = [
        ":graph",
        "//map:astar",
        "//map:map",
        "//map:utils",
        "//pathing/api:constants_go_proto",
        "//pathing/api:data_go_proto",
        "@com_github_fzipp_astar//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "graphastar_test",
    srcs = ["graphastar_test.go"],
    importpath = "github.com/downflux/game/pathing/hpf/graphastar_test",
    embed = [":graphastar"],
    deps = [
        "//api:data_go_proto",
        "//map:astar",
        "//map:map",
        "//map:utils",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "//pathing/api:data_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "cluster",
    srcs = ["cluster.go"],
    importpath = "github.com/downflux/game/pathing/hpf/cluster",
    deps = [
        "//api:data_go_proto",
        "//map:map",
        "//map:utils",
        "//pathing/api:constants_go_proto",
        "//pathing/api:data_go_proto",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name ="cluster_test",
    srcs = ["cluster_test.go"],
    importpath = "github.com/downflux/game/pathing/hpf/cluster_test",
    embed = [":cluster"],
    deps = [
        "//map:utils",
        "//pathing/api:constants_go_proto",
        "//pathing/api:data_go_proto",
    ],
)

go_library(
    name = "entrance",
    srcs = ["entrance.go"],
    importpath = "github.com/downflux/game/pathing/hpf/entrance",
    deps = [
        ":cluster",
        "//api:data_go_proto",
        "//map:map",
        "//map:utils",
        "//pathing/api:constants_go_proto",
        "//pathing/api:data_go_proto",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "entrance_test",
    srcs = ["entrance_test.go"],
    importpath = "github.com/downflux/game/pathing/hpf/entrance_test",
    embed = [":entrance"],
    deps = [
        ":cluster",
        "//map:map",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "//pathing/api:constants_go_proto",
        "//pathing/api:data_go_proto",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_google_go_cmp//cmp:go_default_library",
    ],
)

go_library(
    name = "edge",
    srcs = ["edge.go"],
    importpath = "github.com/downflux/game/pathing/hpf/edge",
    deps = [
        "//map:utils",
        "//pathing/api:data_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "edge_test",
    srcs = ["edge_test.go"],
    importpath = "github.com/downflux/game/pathing/hpf/edge_test",
    embed = [":edge"],
    deps = [
        "//api:data_go_proto",
        "//map:utils",
        "//pathing/api:data_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "node",
    srcs = ["node.go"],
    importpath = "github.com/downflux/game/pathing/hpf/node",
    deps = [
        ":cluster",
        "//map:utils",
        "//pathing/api:data_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "node_test",
    srcs = ["node_test.go"],
    importpath = "github.com/downflux/game/pathing/hpf/node_test",
    embed = [":node"],
    deps = [
        ":cluster",
        "//api:data_go_proto",
        "//map:utils",
        "//pathing/api:data_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "graph",
    srcs = ["graph.go"],
    importpath = "github.com/downflux/game/pathing/hpf/graph",
    deps = [
        ":cluster",
        ":edge",
        ":entrance",
        ":node",
        "//map:astar",
        "//map:map",
        "//map:utils",
        "//api:data_go_proto",
        "//pathing/api:constants_go_proto",
        "//pathing/api:data_go_proto",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "graph_test",
    srcs = ["graph_test.go"],
    importpath = "github.com/downflux/game/pathing/hpf/graph_test",
    embed = [":graph"],
    deps = [
        ":cluster",
        ":edge",
        ":entrance",
        ":node",
        "//map:map",
        "//map:utils",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "//pathing/api:constants_go_proto",
        "//pathing/api:data_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)
