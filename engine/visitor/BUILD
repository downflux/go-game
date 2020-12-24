load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "github.com/downflux/game/engine/visitor/list",
    deps = [
        ":visitor",
	"//engine/visitor/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_library(
    name = "visitor",
    srcs = ["visitor.go"],
    importpath = "github.com/downflux/game/engine/visitor/visitor",
    deps = [
	"//engine/visitor/api:constants_go_proto",
    ],
)
