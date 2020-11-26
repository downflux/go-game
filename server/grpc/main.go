package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/downflux/game/server/grpc/server"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

	apipb "github.com/downflux/game/api/api_go_proto"
	gcpb "github.com/downflux/game/api/constants_go_proto"
	gdpb "github.com/downflux/game/api/data_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
)

var (
	port    = flag.Int("port", 4444, "gRPC server listener port")
	mapFile = flag.String("map_file", "data/map/demo.textproto", "game map textproto file")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	addr := fmt.Sprintf("localhost:%d", *port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("could not open host addr %s: %v", addr, err)
	}

	d, err := ioutil.ReadFile(*mapFile)
	if err != nil {
		log.Fatalf("could not open map file %s: %v", *mapFile, err)
	}

	mapPB := &mdpb.TileMap{}

	if err := proto.UnmarshalText(string(d), mapPB); err != nil {
		log.Fatalf("could not parse map file: %v", err)
	}

	downFluxServer, err := server.NewDownFluxServer(mapPB, &gdpb.Coordinate{X: 5, Y: 5})
	if err != nil {
		log.Fatal("could not construct DownFlux server instance: %v", err)
	}

	log.Printf("serving on %s", addr)

	s := grpc.NewServer()
	apipb.RegisterDownFluxServer(s, downFluxServer)

	downFluxServer.Executor().AddEntity(gcpb.EntityType_ENTITY_TYPE_TANK, &gdpb.Position{X: 0, Y: 0})

	go s.Serve(lis)
	downFluxServer.Executor().Run()
}
