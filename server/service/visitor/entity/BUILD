load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=[
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
        "//server/service/visitor/entity:entity",
    ],
)

go_library(
    name = "entity",
    srcs = ["entity.go"],
    importpath = "github.com/downflux/game/server/service/visitor/entity/entity",
    deps = [
        "//api:constants_go_proto",
        "//curve:curve",
        "//server/service/visitor:visitor",
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
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)
