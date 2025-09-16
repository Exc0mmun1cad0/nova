.PHONY: cover test lint build run

cover:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

test:
	go clean -testcache
	go test -v ./...

lint:
	golangci-lint run

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -extldflags '-static'" -o nova

run: build
	./nova