BINARY_NAME=server

build:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME} ./cmd/main/main.go

run:
	go run ./cmd/main/main.go

clean:
	go clean
	rm ./bin/${BINARY_NAME}
