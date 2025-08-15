.PHONY: init
init:
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/server ./cmd/server

.PHONY: swag
swag:
	swag init  -g cmd/server/main.go -o ./docs