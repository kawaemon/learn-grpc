package main

import (
    "time"
	"context"
	pb "github.com/kawaemon/learn-grpc/types"
	"google.golang.org/grpc"
	"log"
)

func main() {
	connection, err := grpc.Dial("localhost:8000", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Panicf("failed to dial: %v", err)
	}

	defer connection.Close()

	client := pb.NewChatClient(connection)

	name := "hoge"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	response, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})

	if err != nil {
		log.Panicf("failed to send request: %v", err)
	}

	log.Printf("Response: %s", response.GetMessage())
}
