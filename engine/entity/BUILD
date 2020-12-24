load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "entity",
    srcs = ["entity.go"],
    importpath = "github.com/downflux/game/engine/entity/entity",
    deps = [
        "//api:constants_go_proto",
        "//curve:curve",
        "//server:id",
    ],
)

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "github.com/downflux/game/engine/entity/list",
    deps = [
        ":entity",
	"//server:id",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
