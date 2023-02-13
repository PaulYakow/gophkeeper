# Запуск тестов (без покрытия)
test:
	go test ./...

# Генерация *.pb.go файлов из соответствующего proto-файла
# Пример использования: make proto name=user
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/$(name).proto
.PHONY: proto

# Запуск линтеров (статическая проверка кода)
lint:
	gofumpt -d .
	golangci-lint run

# Запуск утилиты grpcui (должна быть установлена) для проверки gRPC-сервера с помощью веб-интерфейса
# Пример использования: make grpc_test name=user
grpc_test:
	grpcui -proto ./proto/$(name).proto -plaintext localhost:9090

# Компиляция клиента для разных платформ
CLIENT_BINARY_NAME=gophkeeper-tui
client_build:
	GOARCH=amd64 GOOS=linux go build -o ${CLIENT_BINARY_NAME}-linux ./cmd/client/main.go
	GOARCH=amd64 GOOS=windows go build -o ${CLIENT_BINARY_NAME}-windows ./cmd/client/main.go
	GOARCH=amd64 GOOS=darwin go build -o ${CLIENT_BINARY_NAME}-darwin ./cmd/client/main.go

# Компиляция сервера для разных платформ
SERVER_BINARY_NAME=gophkeeper-srv
server_build:
	GOARCH=amd64 GOOS=linux go build -o ${SERVER_BINARY_NAME}-linux ./cmd/server/main.go
	GOARCH=amd64 GOOS=windows go build -o ${SERVER_BINARY_NAME}-windows ./cmd/server/main.go
	GOARCH=amd64 GOOS=darwin go build -o ${SERVER_BINARY_NAME}-darwin ./cmd/server/main.go
