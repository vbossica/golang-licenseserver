CLIENT_BINARY_NAME=license-client
SERVER_BINARY_NAME=license-server

build:
	go mod tidy
	go build -o ${CLIENT_BINARY_NAME} cmd/license_client/main.go
	go build -o ${SERVER_BINARY_NAME} cmd/license_server/main.go

test:
	go test ./client ./server

clean:
	go clean
	rm ${CLIENT_BINARY_NAME}
	rm ${SERVER_BINARY_NAME}