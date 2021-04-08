load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "group",
    srcs = ["group.go"],
    importpath = "github.com/downflux/game/engine/fsm/group/group",
    deps = [
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
    ],
)

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "github.com/downflux/game/engine/fsm/group/list",
    deps = [
        ":group",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
