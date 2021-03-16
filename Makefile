gen_types: protocol.proto
	protoc --go_out=types --go-grpc_out=types protocol.proto

server: gen_types $(wildcard server/*.go)
	go build -o server.a ./server

client: gen_types $(wildcard client/*.go)
	go build -o client.a ./client

all: server client
