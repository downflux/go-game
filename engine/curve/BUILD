load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "curve",
    srcs = ["curve.go"],
    importpath = "github.com/downflux/game/engine/curve/curve",
    deps = [
        ":data",
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/id:id",
    ],
)

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "github.com/downflux/game/engine/curve/list",
    deps = [
        ":curve",
        "//api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_library(
    name = "data",
    srcs = ["data.go"],
    importpath = "github.com/downflux/game/engine/curve/data",
    deps = [
         "//engine/id:id",
    ],
)

go_test(
    name = "data_test",
    srcs = ["data_test.go"],
    importpath = "github.com/downflux/game/engine/curve/data_test",
    embed = [":data"],
    deps = [
         "//engine/id:id",
    ],
)
