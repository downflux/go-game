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

1. Add `//client/DownFlux/` as a project.

1. Set .NET 4.x for gRPC compatibility

   Edit > Project Settings > Player > Other Settings > Api Compatibility Layer

## Update C# gRPC Version

1. Install [NuGetForUnity](https://github.com/GlitchEnzo/NuGetForUnity)
   v2.0.1.

1. Download from [packages.grpc.io](https://packages.grpc.io/). Latest verified
   version is
   [2.26.0-dev](https://packages.grpc.io/archive/2019/12/a02d6b9be81cbadb60eed88b3b44498ba27bcba9-edd81ac6-e3d1-461a-a263-2b06ae913c3f/index.xml)

   See [#22251](https://github.com/grpc/grpc/issues/22251) for explanation and
   architecture updates.

1. Unzip `grpc_unity_package` inside the `//client/DownFlux/Assets` directory.

1. Delete the `ios`, `android`, and `osx` directories from `Grpc.Core/runtimes`
   (as we aren't supporting these builds in Unity) and they contain very large
   binaries.

1. Untar relevant `grpc_protoc` plugin dir from same verified gRPC version, and
   copy `grpc_csharp_plugin` to `/usr/local/bin/`.

1. Install a verified working `protoc` package (latest verified with gRPC
   v2.26.0-dev is protoc
   [v3.8](https://github.com/protocolbuffers/protobuf/releases/tag/v3.8.0)).
   See [StackOverflow](https://askubuntu.com/questions/1072683/).

## Install Visual Studio Code

1. Install [VSCode](https://code.visualstudio.com/docs/setup/linux), taking
   care to go the `apt` update route.

1. Set as default Unity editor

   Edit > Preferences > External Tools > External Script Editor

1. Install the
   [C# extension](https://marketplace.visualstudio.com/items?itemName=ms-dotnettools.csharp)
   extension for VSCode, via Extensions in the side panel.

1. Install the
   [.NET SDK](https://docs.microsoft.com/en-us/dotnet/core/install/linux).
   Currently verified with Ubuntu 20.04.

## Protobuf Generation

1. Generate protobufs from root GitHub repo directory

   ```bash
   protoc -I=${PWD} \
     --grpc_out=${PWD}/client/DownFlux/Assets/Protos/Api \
     --csharp_out=${PWD}/client/DownFlux/Assets/Protos/Api \
     --plugin=protoc-gen-grpc=/usr/local/bin/grpc_csharp_plugin \
     $(find ${PWD}/api -iname "*.proto")
   ```

   We will need to do this per proto directory.
