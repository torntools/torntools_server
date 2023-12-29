package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	proto "github.com/panaka13/torntools_server/gen/torntools_proto"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	hostFlag    = flag.String("host", "localhost:3333", "host of server")
	userFlag    = flag.Int("user", 0, "torn user id")
	tornApiFlag = flag.String("tornapi", "", "torn API key")
	apiFlag     = flag.String("api", "", "torntools function")
	fromFlag    = flag.String("from", "1", "date")
	toFlag      = flag.String("to", "2", "date")
)

func parseDate(layout string, message string) (time.Time, error) {
	return time.Parse(layout, message)
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*hostFlag, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := proto.NewTornToolClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var request proto.ViewBookieResultResquest
	request.User = int32(*userFlag)
	request.Api = *tornApiFlag

	fromValue, err := time.Parse("2006-01-02", *fromFlag)
	if err != nil {
		fmt.Printf("fail to parse from date: %s\n", err.Error())
		return
	}
	request.From = timestamppb.New(fromValue)

	toValue, err := time.Parse("2006-01-02", *toFlag)
	if err != nil {
		fmt.Printf("fail to parse from date: %s\n", err.Error())
		return
	}
	request.To = timestamppb.New(toValue)

	response, err := client.ViewBookieResult(ctx, &request)
	if err != nil {
		fmt.Printf("fail to call server: %s\n", err.Error())
		return
	}
	if response == nil {
		fmt.Printf("nil response\n")
		return
	}
	for _, bookie := range response.Results {
		fmt.Println(bookie)
	}
}
