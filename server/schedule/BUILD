load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
   name = "schedule",
   srcs = ["schedule.go"],
   importpath = "github.com/downflux/game/server/schedule/schedule",
   deps = [
        "//fsm/api:constants_go_proto",
        "//fsm:instance",
        "//fsm:list",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
   ],
)
