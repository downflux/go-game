load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "client",
    srcs = ["client.go"],
    importpath = "github.com/downflux/game/server/grpc/client",
    deps = [
        "//api:api_go_proto",
    ],
)

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
        "//api:constants_go_proto",
        "//api:data_go_proto",
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
        ":client",
        ":executorutils",
        "//api:api_go_proto",
        "//api:data_go_proto",
        "//engine/id:id",
        "//map/api:data_go_proto",
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
    timeout = "short",
    embed = [":server"],
    deps = [
        ":option",
        ":handler",
        "//api:api_go_proto",
        "@com_github_shopify_toxiproxy//:go_default_library",
        "@com_github_shopify_toxiproxy//client:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//connectivity:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
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
        "//map/api:constants_go_proto",
        "//map/api:data_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//test/bufconn:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "handler",
    srcs = ["handler.go"],
    importpath = "github.com/downflux/game/server/grpc/handler",
    deps = [
	"@org_golang_google_grpc//stats:go_default_library",
    ],
)

go_library(
    name = "executorutils",
    srcs = ["executorutils.go"],
    importpath = "github.com/downflux/game/server/grpc/executorutils",
    deps = [
        "//api:api_go_proto",
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/gamestate:dirty",
        "//engine/gamestate:gamestate",
        "//engine/entity:list",
        "//engine/fsm:schedule",
        "//engine/fsm/api:constants_go_proto",
        "//engine/id:id",
        "//engine/server/executor:executor",
        "//engine/status:status",
        "//engine/visitor:list",
        "//engine/visitor:visitor",
        "//engine/visitor/api:constants_go_proto",
        "//map:map",
        "//map/api:data_go_proto",
        "//pathing/hpf:graph",
        "//server/entity/component:moveable",
        "//server/fsm:move",
        "//server/fsm:produce",
        "//server/visitor:chase",
        "//server/visitor:move",
        "//server/visitor:produce",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
