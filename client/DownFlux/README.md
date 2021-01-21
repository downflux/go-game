# DownFlux Unity Client

*N.B.*: These instructions are specific for Ubuntu Linux; instructions may be
different for Windows installs.

## Install Unity

1. See [official documentation](https://unity3d.com/get-unity/download).

1. Set Unity Editors folder for Unity Hub.

1. The client is tested on `Unity 2019.4.19f1 (LTS)`; install the relevant
   version and ensure Linux and Windows build support is installed.

1. Due to a
   [bug with Unity Hub](https://forum.unity.com/threads/no-editor-installed.893605),
   we need to create a throwaway project first.

1. Add `DownFlux/` as a project.

1. Set .NET 4.x for gRPC compatibility

   Edit > Project Settings > Player > Other Settings > Api Compatibility Layer

## Update C# gRPC Version

1. Install [NuGetForUnity](https://github.com/GlitchEnzo/NuGetForUnity).
   Currently verified with 2.0.1.

   Assets > Import Package

1. Ensure the correct gRPC packages are installed

   | Package               | Version |
   | --------------------- | ------- |
   | Google.Protobuf       | 3.14.0  |
   | Grpc                  | 2.34.1  |
   | Grpc.Core             | 2.34.1  |
   | Grpc.Core.Api         | 2.34.1  |
   | Grpc.Net.Client       | 2.34.0  |
   | Grpc.Tools            | 2.34.1  |
   | Google.Protobuf.Tools | 3.14.0  |

1. Copy the gRPC binary from `Packages/Grpc.Tools` directory into `/usr/local/bin`,
   set executable bit, and `root:root` as owner for the binaries. See
   [#694](https://github.com/golang/protobuf/issues/694) for potential
   pitfalls with linking native `.proto` files.

## [Optional] Install Visual Studio Code

1. Install [VSCode](https://code.visualstudio.com/docs/setup/linux), taking
   care to go the `apt` update route.

1. Set as default Unity editor

   Edit > Preferences > External Tools > External Script Editor

1. Install the
   [C# extension](https://marketplace.visualstudio.com/items?itemName=ms-dotnettools.csharp)
   extension for VSCode, via Extensions in the side panel.

1. Per extension notes, ensure we set

   ```json
   "omnisharp.useGlobalMono": "never"
   ```

   in `DownFlux/.vscode/settings.json`

1. Install the
   [.NET 5.0 SDK](https://docs.microsoft.com/en-us/dotnet/core/install/linux).
   Currently verified with Ubuntu 20.04.

1. Install `mono-complete` via
   [official source](https://www.mono-project.com/download/stable/#download-lin).

1. Edit `~/.bashrc`:

   ```bash
   export FrameworkPathOverride=/lib/mono/4.7.1-api
   ```

   This seems necessary due to an archaic unfixed bug
   [#335](https://github.com/dotnet/sdk/issues/335).

1. `FrameworkPathOverride` is not passed in via the Unity GUI; we have to
   start VSCode separately in the terminal

   ```bash
   code
   ```

## Protobuf Generation

Generate protobufs from root GitHub repo directory. See
[official documentation](https://developers.google.com/protocol-buffers/docs/reference/csharp-generated#compiler_options)
for `protoc` flag explanation.

```bash
rm -rf ${PWD}/client/DownFlux/Assets/Protos/*
protoc \
  -I=${PWD} \
  -I=${PWD}/client/DownFlux/Packages/Google.Protobuf.Tools.3.14.0/tools/ \
  --grpc_out=${PWD}/client/DownFlux/Assets/Protos \
  --csharp_out=${PWD}/client/DownFlux/Assets/Protos \
  --csharp_opt=file_extension=.g.cs,base_namespace=DF \
  --plugin=protoc-gen-grpc=/usr/local/bin/grpc_csharp_plugin \
  $(find ${PWD}/ -iname "*.proto" -print -o -path ${PWD}/client -prune)
```
