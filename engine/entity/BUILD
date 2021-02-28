load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "acl",
    srcs = ["acl.go"],
    importpath = "github.com/downflux/game/engine/entity/acl",
    deps = [
        "//api:data_go_proto",
        "//engine/id:id",
    ],
)

go_test(
    name = "acl_test",
    srcs = ["acl_test.go"],
    importpath = "github.com/downflux/game/engine/entity/acl_test",
    embed = [":acl"],
    deps = [
        "//api:data_go_proto",
        "//engine/id:id",
    ],
)

go_library(
    name = "entity",
    srcs = ["entity.go"],
    importpath = "github.com/downflux/game/engine/entity/entity",
    deps = [
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "//engine/curve:list",
        "//engine/id:id",
    ],
)

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "github.com/downflux/game/engine/entity/list",
    deps = [
        ":entity",
	"//engine/id:id",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
