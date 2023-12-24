package main

import (
	"context"
	"fmt"
	"log"
	"net"

	proto "github.com/panaka13/torntools_server/gen/torntools_proto"
	grpc "google.golang.org/grpc"
)

type tornToolServer struct {
	proto.UnimplementedTornToolServer
}

func (s tornToolServer) ViewBookieResult(context.Context, *proto.ViewBookieResultResquest) (*proto.ViewBookieResultResponse, error) {
	fmt.Println("ViewBookieResult API")
	var response proto.ViewBookieResultResponse
	var bookieResult  proto.BookieResult
	response.Results = append(response.Results, bookieResult)
	return &response, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 3333))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	var server tornToolServer

	proto.RegisterTornToolServer(grpcServer, server)
	grpcServer.Serve(lis)

}
