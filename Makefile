.PHONY: build test test-verbose test-cover lint clean

# Build the binary
build:
	go build -o gomigrate .

# Run all tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage report
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run linter (requires: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -f gomigrate coverage.out

# Install the CLI locally
install:
	go install .