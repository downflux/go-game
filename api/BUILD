load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library", "go_grpc_library")
load("@rules_proto_grpc//csharp:defs.bzl", "csharp_grpc_compile", "csharp_proto_compile")

package(default_visibility=["//visibility:public"])

csharp_grpc_compile(
    name = "api_csharp_proto",
    deps = [":api_proto"],
)

proto_library(
    name = "api_proto",
    srcs = ["api.proto"],
    deps = [
        ":constants_proto",
        ":data_proto",
    ],
)

go_grpc_library(
    name = "api_go_proto",
    importpath = "github.com/downflux/game/api/api_go_proto",
    proto = ":api_proto",
    deps = [
        ":constants_go_proto",
        ":data_go_proto",
    ],
)

proto_library(
    name = "data_proto",
    srcs = ["data.proto"],
    deps = [
        ":constants_proto",
	"@com_google_protobuf//:duration_proto",
	"@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "data_go_proto",
    importpath = "github.com/downflux/game/api/data_go_proto",
    proto = ":data_proto",
    deps = [
        ":constants_go_proto",
        "@io_bazel_rules_go//proto/wkt:duration_go_proto",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
    ],
)

proto_library(
    name = "constants_proto",
    srcs = ["constants.proto"],
)

go_proto_library(
    name = "constants_go_proto",
    importpath = "github.com/downflux/game/api/constants_go_proto",
    proto = ":constants_proto",
)
