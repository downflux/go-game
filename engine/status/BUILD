load("@io_bazel_rules_go//go:def.bzl", "go_library")

package(default_visibility=["//visibility:public"])

go_library(
    name = "status",
    srcs = ["status.go"],
    importpath = "github.com/downflux/game/engine/status/status",
    deps = [
        "//api:data_go_proto",
        "//engine/id:id",
        "//server/service/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//types/known/durationpb:go_default_library",
        "@org_golang_google_protobuf//types/known/timestamppb:go_default_library",
    ],
)
