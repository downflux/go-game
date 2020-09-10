
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test", "go_binary")

go_library(
    name = "astar",
    srcs = ["astar.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/astar",
    deps = [
        ":tile",
        "//lib/proto:constants_go_proto",
	"//lib/proto:structs_go_proto",
	"@com_github_fzipp_astar//:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "astar_test",
    srcs = ["astar_test.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/astar_test",
    embed = [":astar"],
    deps = [
	":tile",
	"//lib/proto:constants_go_proto",
	"//lib/proto:structs_go_proto",
	"@com_github_google_go_cmp//cmp:go_default_library",
	"@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "hpf",
    srcs = ["hpf.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/hpf",
    deps = [
	"//lib/proto:constants_go_proto",
    ],
)

go_library(
    name = "tile",
    srcs = ["tile.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/tile",
    deps = [
	":utils",
	"//lib/proto:constants_go_proto",
	"//lib/proto:structs_go_proto",
	"@com_github_fzipp_astar//:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "tile_test",
    srcs = ["tile_test.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/tile_test",
    embed = [":tile"],
    deps = [
	":utils",
	"//lib/proto:constants_go_proto",
	"//lib/proto:structs_go_proto",
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
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/cluster",
    deps = [
        ":tile",
	":utils",
        "//lib/proto:constants_go_proto",
	"//lib/proto:structs_go_proto",
	"@org_golang_google_grpc//status:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name ="cluster_test",
    srcs = ["cluster_test.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/cluster_test",
    embed = [":cluster"],
    deps = [
	":utils",
        "//lib/proto:constants_go_proto",
        "//lib/proto:structs_go_proto",
	"@com_github_google_go_cmp//cmp:go_default_library",
	"@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "entrance",
    srcs = ["entrance.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/entrance",
    deps = [
	":cluster",
        ":tile",
	":utils",
        "//lib/proto:constants_go_proto",
	"//lib/proto:structs_go_proto",
	"@com_github_golang_protobuf//proto:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "entrance_test",
    srcs = ["entrance_test.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/entrance_test",
    embed = [":entrance"],
    deps = [
	":cluster",
	":tile",
        "//lib/proto:constants_go_proto",
	"//lib/proto:structs_go_proto",
	"@org_golang_google_protobuf//testing/protocmp:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
	"@com_github_google_go_cmp//cmp:go_default_library",
    ],
)

go_library(
    name = "abstractedgemap",
    srcs = ["abstractedgemap.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/abstractedgemap",
    deps = [
        ":utils",
	"//lib/proto:structs_go_proto",
	"@org_golang_google_grpc//codes:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "abstractedgemap_test",
    srcs = ["abstractedgemap_test.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/abstractedgemap_test",
    embed = [":abstractedgemap"],
    deps = [
	":utils",
	"//lib/proto:structs_go_proto",
	"@com_github_google_go_cmp//cmp:go_default_library",
	"@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
	"@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "abstractnodemap",
    srcs = ["abstractnodemap.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/abstractnodemap",
    deps = [
        ":cluster",
        ":utils",
	"//lib/proto:structs_go_proto",
	"@org_golang_google_grpc//codes:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "abstractnodemap_test",
    srcs = ["abstractnodemap_test.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/abstractnodemap_test",
    embed = [":abstractnodemap"],
    deps = [
        ":cluster",
        ":utils",
	"//lib/proto:structs_go_proto",
	"@com_github_google_go_cmp//cmp:go_default_library",
	"@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
	"@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "abstractgraph",
    srcs = ["abstractgraph.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/abstractgraph",
    deps = [
        ":abstractedgemap",
	":abstractnodemap",
	":astar",
	":cluster",
	":entrance",
        ":utils",
	":tile",
        "//lib/proto:constants_go_proto",
	"//lib/proto:structs_go_proto",
	"@com_github_golang_protobuf//proto:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "abstractgraph_test",
    srcs = ["abstractgraph_test.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/abstractgraph_test",
    embed = [":abstractgraph"],
    deps = [
        ":abstractedgemap",
        ":abstractnodemap",
        ":cluster",
	":tile",
        ":utils",
	"//lib/proto:structs_go_proto",
	"@com_github_google_go_cmp//cmp:go_default_library",
	"@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
	"@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "utils",
    srcs = ["utils.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/utils",
    deps = [
	"//lib/proto:structs_go_proto",
    ],
)
