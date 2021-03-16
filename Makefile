gen_types: protocol.proto
	protoc --go_out=types --go-grpc_out=types protocol.proto

server: gen_types $(wildcard server/*.go)
	go build -o server.a ./server