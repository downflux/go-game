
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test", "go_binary")

go_library(
    name = "astar",
    srcs = ["astar.go"],
    importpath = "github.com/downflux/game/pathing/hpf/astar",
    deps = [
        ":cluster",
        ":graph",
        ":graphastar",
        ":tile",
        ":tileastar",
        ":utils",
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
        ":tile",
        ":utils",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
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
        ":tile",
        ":tileastar",
        ":utils",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
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
        ":graph",
        ":tile",
        ":utils",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "tileastar",
    srcs = ["tileastar.go"],
    importpath = "github.com/downflux/game/pathing/hpf/tileastar",
    deps = [
        ":tile",
        ":utils",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
        "@com_github_fzipp_astar//:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "tileastar_test",
    srcs = ["tileastar_test.go"],
    importpath = "github.com/downflux/game/pathing/hpf/tileastar_test",
    embed = [":tileastar"],
    deps = [
        ":tile",
        ":utils",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "tile",
    srcs = ["tile.go"],
    importpath = "github.com/downflux/game/pathing/hpf/tile",
    deps = [
        ":utils",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
        "@com_github_fzipp_astar//:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "tile_test",
    srcs = ["tile_test.go"],
    importpath = "github.com/downflux/game/pathing/hpf/tile_test",
    embed = [":tile"],
    deps = [
        ":utils",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "cluster",
    srcs = ["cluster.go"],
    importpath = "github.com/downflux/game/pathing/hpf/cluster",
    deps = [
        ":tile",
        ":utils",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
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
        ":utils",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
    ],
)

go_library(
    name = "entrance",
    srcs = ["entrance.go"],
    importpath = "github.com/downflux/game/pathing/hpf/entrance",
    deps = [
        ":cluster",
        ":tile",
        ":utils",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
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
        ":tile",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
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
        ":utils",
        "//pathing/proto:structs_go_proto",
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
        ":utils",
        "//pathing/proto:structs_go_proto",
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
        ":utils",
        "//pathing/proto:structs_go_proto",
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
        ":utils",
        "//pathing/proto:structs_go_proto",
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
        ":edge",
        ":node",
        ":tileastar",
        ":cluster",
        ":entrance",
        ":utils",
        ":tile",
        "//pathing/proto:constants_go_proto",
        "//pathing/proto:structs_go_proto",
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
        ":edge",
        ":node",
        ":cluster",
        ":tile",
        ":utils",
        "//pathing/proto:structs_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "utils",
    srcs = ["utils.go"],
    importpath = "github.com/downflux/game/pathing/hpf/utils",
    deps = [
        "//pathing/proto:structs_go_proto",
    ],
)
