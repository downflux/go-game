load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_binary(
    name = "main",
    srcs = ["main.go"],
    importpath = "github.com/downflux/game/server/grpc/main",
    data = [
        "//data/map:map_data",
    ],
    deps = [
        ":server",
        "//api:api_go_proto",
        "//api:data_go_proto",
        "//entity:entity",
	"//map/api:data_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",

    ],
)

go_library(
    name = "server",
    srcs = ["server.go"],
    importpath = "github.com/downflux/game/server/grpc/server",
    deps = [
        "//api:api_go_proto",
        "//api:data_go_proto",
        "//map/api:data_go_proto",
        "//server/service:executor",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)

go_library(
    name = "option",
    srcs = ["option.go"],
    importpath = "github.com/downflux/game/server/grpc/option",
    deps = [
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//keepalive:go_default_library",
    ],
)

go_test(
    name = "server_network_test",
    srcs = ["server_network_test.go"],
    importpath = "github.com/downflux/game/server/grpc/server_network_test",
    embed = [":server"],
    deps = [
        ":option",
        "//api:api_go_proto",
        "@com_github_shopify_toxiproxy//:go_default_library",
        "@com_github_shopify_toxiproxy//client:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//connectivity:go_default_library",
    ],
)

go_test(
    name = "server_test",
    srcs = ["server_test.go"],
    importpath = "github.com/downflux/game/server/grpc/server_test",
    embed = [":server"],
    deps = [
        "//api:api_go_proto",
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//entity:entity",
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "//server:id",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//test/bufconn:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)
