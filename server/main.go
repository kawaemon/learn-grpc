package main

import (
	"bufio"
	"context"
	pb "github.com/kawaemon/learn-grpc/types"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
)

var (
	stdinChan = make(chan string)
)

type Server struct {
	pb.UnimplementedChatServer
}

func (s *Server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	name := request.GetName()

	log.Println("Said hello from " + name)

	return &pb.HelloReply{Message: "Hello " + name}, nil
}


func (s *Server) Chat(stream pb.Chat_ChatServer) error {
	log.Println("Client connected.")

	stopChan := make(chan struct{})
	recvChan := make(chan string)

	go (func() {
		for {
			msg, e := stream.Recv()

			if e == io.EOF {
				stopChan <- struct{}{}
				return
			}

			if e != nil {
				log.Printf("failed to receive message: %v", e)
				stopChan <- struct{}{}
				return
			}

			recvChan <- msg.GetMessage()
		}
	})()

	running := true
	for running {
		select {
		case <- stopChan:
			running = false

		case msg := <- recvChan:
			log.Printf("server said: %s\n", msg)

		case input := <- stdinChan:
			err := stream.Send(&pb.ChatMessage{Message: input})

			if err != nil {
				log.Fatalf("failed to send to server")
			}
		}
	}

	log.Println("Client disconnected")

	return nil
}

func main() {
	go (func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			stdinChan <- text
		}
	})()

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
