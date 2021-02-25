load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "executor",
    srcs = ["executor.go"],
    importpath = "github.com/downflux/game/engine/server/executor/executor",
    deps = [
        "//api:api_go_proto",
        "//api:data_go_proto",
        "//engine/gamestate:dirty",
        "//engine/gamestate:gamestate",
        "//engine/id:id",
        "//engine/fsm:action",
        "//engine/fsm:schedule",
        "//engine/fsm/api:constants_go_proto",
        "//engine/server/client:list",
        "//engine/visitor:list",
        "//engine/visitor/api:constants_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "executor_test",
    srcs = [":executor_test.go"],
    importpath = "github.com/downflux/game/engine/server/executor/executor_test",
    embed = [":executor"],
    deps = [
        "//engine/entity:list",
        "//engine/fsm:schedule",
        "//engine/fsm/api:constants_go_proto",
        "//engine/fsm/mock:simple",
        "//engine/gamestate:dirty",
        "//engine/gamestate:gamestate",
        "//engine/id:id",
        "//engine/status:status",
        "//engine/visitor:list",
        "//engine/visitor:visitor",
        "//engine/visitor/api:constants_go_proto",
        "//engine/visitor/mock:simple",
    ],
)
