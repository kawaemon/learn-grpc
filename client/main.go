package main

import (
	"bufio"
	"context"
	pb "github.com/kawaemon/learn-grpc/types"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

func main() {
	connection, err := grpc.Dial("localhost:8000", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Panicf("failed to dial: %v", err)
	}

	defer connection.Close()

	client := pb.NewChatClient(connection)

	ctx := context.Background()

	stream, err := client.Chat(ctx)

	if err != nil {
		log.Panicf("failed to create stream: %v", err)
	}

	log.Println("now chat is connected")

	stopChan := make(chan struct{})
	stdinChan := make(chan string)
	recvChan := make(chan string)

	go (func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			stdinChan <- text
		}
	})()

	go (func() {
		for {
			msg, e := stream.Recv()

			if e == io.EOF {
				stopChan <- struct{}{}
				return
			}

			if e != nil {
				log.Fatalf("failed to receive message: %v", e)
				return
			}

			recvChan <- msg.GetMessage()
		}
	})()

	running := true
	for running {
		select {
		case <-stopChan:
			running = false

		case msg := <-recvChan:
			log.Printf("server said: %s\n", msg)

		case input := <-stdinChan:
			err := stream.Send(&pb.ChatMessage{Message: input})

			if err != nil {
				log.Fatalf("failed to send to server")
			}
		}
	}
}
