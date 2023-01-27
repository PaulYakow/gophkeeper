test:
	go test ./...

#Пример использования: make proto name=user
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/$(name).proto
.PHONY: proto

lint:
	gofumpt -d .
	golangci-lint run

grpc_test:
	grpcui -proto ./proto/user.proto -plaintext localhost:9090