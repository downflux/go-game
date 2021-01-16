workspace(name = "downflux")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "ac03931e56c3b229c145f1a8b2a2ad3e8d8f1af57e43ef28a26123362a1e3c7e",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.24.4/rules_go-v0.24.4.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.24.4/rules_go-v0.24.4.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "b85f48fa105c4403326e9525ad2b2cc437babaa6e15a3fc0b1dbab0ab064bc7c",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.2/bazel-gazelle-v0.22.2.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.2/bazel-gazelle-v0.22.2.tar.gz",
    ],
)

# Match version with the rules_proto_grpc lib.
http_archive(
    name = "com_google_protobuf",
    strip_prefix = "protobuf-3.13.0",
    sha256 = "1c744a6a1f2c901e68c5521bc275e22bdc66256eeb605c2781923365b7087e5f",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v3.13.0.zip"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")
protobuf_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()
go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
gazelle_dependencies()

# Golang external dependencies.
go_repository(
    name = "org_golang_google_protobuf",
    importpath = "google.golang.org/protobuf",
    # TODO(minkezhang): Monitor this repo for when we can move away from ptypes.
    sum = "h1:Ejskq+SyPohKW+1uil0JJMtmHCgJPJ/qWTxr8qp+R4c=",
    version = "v1.25.0",
)

go_repository(
    name = "com_github_golang_protobuf",
    importpath = "github.com/golang/protobuf",
    commit = "d23c5127dc24889085f8ccea5c9d560a57a879d8",
    version = "v1.3.3",
)

go_repository(
    name = "org_golang_google_grpc",
    build_file_proto_mode = "disable",
    importpath = "google.golang.org/grpc",
    commit = "754ee590a4f386d0910d887f3b8776354042260b",  # v1.29.1
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    sum = "h1:oWX7TPOiFAMXLq8o0ikBYfCJVlRHBcsciT5bXOrH628=",
    version = "v0.0.0-20190311183353-d8887717615a",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    sum = "h1:g61tztE5qeGQ89tm6NTjjM9VPIm088od1l6aSorWRWg=",
    version = "v0.3.0",
)

go_repository(
    name = "org_golang_x_sync",
    importpath = "golang.org/x/sync",
    sum = "h1:Bx6FllMpG4NWDOfhMBz1VR2QYNp/SAOHPIAsaVmxfPo=",
    version = "v0.0.0-20201008141435-b3e1573b7520",
)

go_repository(
    name = "com_github_google_go_cmp",
    importpath = "github.com/google/go-cmp",
    commit = "77ae86f624cb174e21763cffcbbf070eb06cb016",  # v0.5.0
)

go_repository(
    name = "com_github_fzipp_astar",
    importpath = "github.com/fzipp/astar",
    commit = "330d9f048c8e9a904c80f32aa3a95ce6d24afbe0",
)

# ToxiProxy dependencies.
# See https://github.com/kubernetes/test-infra/blob/master/repos.bzl.
go_repository(
    name = "com_github_shopify_toxiproxy",
    importpath = "github.com/Shopify/toxiproxy",
    sum = "h1:TKdv8HiTLgE5wdJuEML90aBgNWsokNbMijUGhmcoBJc=",
    version = "v2.1.4+incompatible",
)
go_repository(
    name = "com_github_sirupsen_logrus",
    importpath = "github.com/sirupsen/logrus",
    sum = "h1:juTguoYk5qI21pwyTXY3B3Y5cOTH3ZUyZCg1v/mihuo=",
    version = "v1.2.0",
)
go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    sum = "h1:Pgr17XVTNXAk3q/r4CpKzC5xBM/qW1uVLV+IhRZpIIk=",
    version = "v1.6.2",
)
go_repository(
    name = "in_gopkg_tomb_v1",
    importpath = "gopkg.in/tomb.v1",
    sum = "h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=",
    version = "v1.0.0-20141024135613-dd632973f1e7",
)
go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    sum = "h1:bselrhR0Or1vomJZC8ZIjWtbDmn9OYFLX5Ik9alpJpE=",
    version = "v0.0.0-20190404164418-38d8ce5564a5",
)
go_repository(
    name = "com_github_gorilla_context",
    importpath = "github.com/gorilla/context",
    sum = "h1:AWwleXJkX/nhcU9bZSnZoi3h/qGYqQAGhq6zZe/aQW8=",
    version = "v1.1.1",
)
go_repository(
    name = "org_golang_x_sys",
    importpath = "golang.org/x/sys",
    sum = "h1:nFYrTHrdrAOpShe27kaFHjsqYSEQ0KWqdWLu3xuZJts=",
    version = "v0.0.0-20190403152447-81d4e9dc473e",
)

# CSharp Proto lib

http_archive(
    name = "rules_proto_grpc",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/2.0.0.tar.gz"],
    sha256 = "d771584bbff98698e7cb3cb31c132ee206a972569f4dc8b65acbdd934d156b33",
    strip_prefix = "rules_proto_grpc-2.0.0",
)

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_toolchains")

rules_proto_grpc_toolchains()

load("@rules_proto_grpc//csharp:repositories.bzl", rules_proto_grpc_csharp_repos="csharp_repos")

rules_proto_grpc_csharp_repos()

load("@io_bazel_rules_dotnet//dotnet:deps.bzl", "dotnet_repositories")

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

dotnet_repositories()

load(
    "@io_bazel_rules_dotnet//dotnet:defs.bzl",
    "core_register_sdk",
    "dotnet_register_toolchains",
    "dotnet_repositories_nugets",
)

dotnet_register_toolchains()
dotnet_repositories_nugets()

core_register_sdk()

load("@rules_proto_grpc//csharp/nuget:nuget.bzl", "nuget_rules_proto_grpc_packages")

nuget_rules_proto_grpc_packages()
