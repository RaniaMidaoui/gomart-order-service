proto:
	protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. pkg/pb/*.proto

server:
	go run cmd/main.go