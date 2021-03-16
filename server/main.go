package main

import (
	"context"
	pb "github.com/kawaemon/learn-grpc/types"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	pb.UnimplementedChatServer
}

func (s *Server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	name := request.GetName()

	log.Println("Said hello from " + name)

	return &pb.HelloReply{Message: "Hello " + name}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterChatServer(server, &Server{})

	err = server.Serve(listener)

	if err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}