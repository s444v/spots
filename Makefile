.PHONY: all run build vet test fmt

all: fmt vet test build

run:
	set -a; . ./.env; set +a; go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

fmt:
	gofmt -s -w .

vet:
	go vet ./...

test:
	go test -race ./...
	