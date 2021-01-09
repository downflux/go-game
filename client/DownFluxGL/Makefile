all:
	dotnet publish \
		-c Release \
		-r linux-x64 \
		/p:PublishReadyToRun=false \
		/p:TieredCompilation=false \
		--self-contained

run: all
	./bin/Release/netcoreapp3.1/linux-x64/DownFluxGL
