package main

import (
	"context"
	"fmt"
	"time"

	proto "github.com/panaka13/torntools_server/gen/torntools_proto"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:3333", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := proto.NewTornToolClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel();

	var request proto.ViewBookieResultResquest

	response, err := client.ViewBookieResult(ctx, &request)
	if err != nil {
		fmt.Printf("fail to call server: %s\n", err.Error())
		return
	}
	if response == nil {
		fmt.Printf("nil response\n")
		return
	}
	fmt.Printf("result size: %d\n", len(response.GetResults()))
}