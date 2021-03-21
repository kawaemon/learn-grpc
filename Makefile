GO_TYPE_FILES=types/protocol.pb.go types/protocol_grpc.pb.go
TS_TYPE_FILES=web/src/grpc/protocol_pb.js web/src/grpc/protocol_pb.d.ts \
			  web/src/grpc/protocol_pb_service.js web/src/grpc/protocol_pb_service.d.ts
TS_TYPE_OUT_DIR=web/src/grpc/

$(GO_TYPE_FILES): ./protocol.proto
	protoc --go_out=types --go-grpc_out=types protocol.proto

$(TS_TYPE_FILES): ./protocol.proto
	protoc --plugin="protoc-gen-ts=./web/node_modules/.bin/protoc-gen-ts" \
	    --js_out="import_style=commonjs,binary:${TS_TYPE_OUT_DIR}" \
	    --ts_out="service=grpc-web:${TS_TYPE_OUT_DIR}" \
	   	./protocol.proto

proxy:
	grpcwebproxy \
		--allow_all_origins \
		--backend_addr localhost:4000 \
		--use_websockets \
		--server_tls_cert_file ./certs/cert.pem \
		--server_tls_key_file ./certs/privkey.pem \
		--server_http_tls_port  3000

server.a: $(GO_TYPE_FILES) server/main.go
	go build -o server.a ./server

client.a: $(GO_TYPE_FILES) client/main.go
	go build -o client.a ./client

webdev: $(TS_TYPE_FILES)
	cd web && pnpm start

all: server.a client.a
