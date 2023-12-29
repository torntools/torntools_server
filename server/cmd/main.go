package main

import (
	"fmt"
	"log"
	"net"

	proto "github.com/panaka13/torntools_server/gen/torntools_proto"
	"github.com/panaka13/torntools_server/server/controller"
	grpc "google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 3333))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	var server controller.TornToolServer

	proto.RegisterTornToolServer(grpcServer, server)
	grpcServer.Serve(lis)

}
