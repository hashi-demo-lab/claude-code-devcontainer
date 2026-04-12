BINARY := ansible-mcp-server
VERSION := 0.1.0

.PHONY: build test lint install clean

build:
	go build -ldflags="-X main.version=$(VERSION)" -o $(BINARY) ./...

test:
	go test ./... -v -race -count=1

lint:
	golangci-lint run
	gosec ./...

validate: build test lint

install:
	go install -ldflags="-X main.version=$(VERSION)" ./...

clean:
	rm -f $(BINARY)
