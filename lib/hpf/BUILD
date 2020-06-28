
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
	"//lib/proto:constants_go_proto",
	"//lib/proto:structs_go_proto",
	"@org_golang_google_grpc//status:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_test(
    name = "tile_test",
    srcs = ["tile_test.go"],
    importpath = "github.com/cripplet/rts-pathing/lib/hpf/tile-test",
    embed = [":tile"],
    deps = [
	"//lib/proto:constants_go_proto",
	"//lib/proto:structs_go_proto",
    ],
)
