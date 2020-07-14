
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test", "go_binary")

go_library(
    name = "hpf",
    srcs = ["hpf.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/hpf",
    deps = [
	"//lib/proto:constants_go_proto",
        "@com_github_beefsack_go_astar//:go_default_library",
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
	"@com_github_golang_protobuf//proto:go_default_library",
	"@com_github_google_go_cmp//cmp:go_default_library",
	"@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
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
	"@com_github_golang_protobuf//proto:go_default_library",
	"@com_github_google_go_cmp//cmp:go_default_library",
    ],
)

go_library(
    name = "entrance",
    srcs = ["entrance.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/entrance",
    deps = [
	":cluster",
        ":tile",
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
    name = "abstractgraph",
    srcs = ["abstractgraph.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/abstractgraph",
    deps = [
        ":utils",
	"//lib/proto:structs_go_proto",
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
