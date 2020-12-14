load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
    # This is only true for intermediate testing.
    # TODO(minkezhang): Remove visibility.
    "//server/grpc:__pkg__",

    "//fsm:__pkg__",
    "//server/service:__pkg__",
    "//server/visitor:__pkg__",
])

go_library(
    name = "tank",
    srcs = ["tank.go"],
    importpath = "github.com/downflux/game/server/entity/tank",
    deps = [
	"//api:constants_go_proto",
	"//api:data_go_proto",
        "//curve:curve",
        "//curve:linearmove",
        "//entity/implement:implement",
        "//entity/implement:moveable",
        "//entity/implement/api:constants_go_proto",
	"//server:id",
        "//server/visitor:visitor",
        "//server/entity:entity",
    ],
)

go_library(
    name = "entity",
    srcs = ["entity.go"],
    importpath = "github.com/downflux/game/server/entity/entity",
    deps = [
        "//api:constants_go_proto",
        "//curve:curve",
        "//entity/implement:implement",
        "//entity/implement/api:constants_go_proto",
        "//server:id",
        "//server/visitor:visitor",
        "//server/visitor/api:constants_go_proto",
    ],
)

go_library(
    name = "entitylist",
    srcs = ["entitylist.go"],
    importpath = "github.com/downflux/game/server/entity/entitylist",
    deps = [
        "//api:constants_go_proto",
        "//entity/implement:implement",
        "//entity/implement/api:constants_go_proto",
	"//server:id",
        "//server/visitor:visitor",
        "//server/entity:entity",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)
