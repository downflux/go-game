load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    # This is only true for intermediate testing.
    # TODO(minkezhang): Remove visibility.
    "//server/grpc:__pkg__",

    "//server/service:__pkg__",
    "//server/service/visitor:__pkg__",
])

go_library(
    name = "tank",
    srcs = ["tank.go"],
    importpath = "github.com/downflux/game/server/service/visitor/entity/tank",
    deps = [
	"//api:constants_go_proto",
	"//api:data_go_proto",
        "//curve:curve",
        "//curve:linearmove",
        "//server/service/visitor:visitor",
        "//server/service/visitor/entity:entity",
    ],
)

go_library(
    name = "entity",
    srcs = ["entity.go"],
    importpath = "github.com/downflux/game/server/service/visitor/entity/entity",
    deps = [
        "//api:constants_go_proto",
        "//curve:curve"
    ],
)

go_library(
    name = "entitylist",
    srcs = ["entitylist.go"],
    importpath = "github.com/downflux/game/server/service/visitor/entity/entitylist",
    deps = [
        "//api:constants_go_proto",
        "//server/service/visitor:visitor",
        "//server/service/visitor/entity:entity",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)
