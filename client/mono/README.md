# Mono-based DownFlux Client

See [Setting up your development environment for Ubuntu 20.04](https://docs.monogame.net/articles/getting_started/1_setting_up_your_development_environment_ubuntu.html).

To build, see [Package games for distribution](https://docs.monogame.net/articles/packaging_games.html).

```bash
dotnet publish -c Release -r linux-x64 /p:PublishReadyToRun=false /p:TieredCompilation=false --self-contained
./bin/Release/netcoreapp3.1/linux-x64/DownFluxGL
```
