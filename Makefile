# Define phony targets to avoid conflicts with existing files/directories
.PHONY: all build test run clean fmt vet help proto

# Default target - builds and runs the application
all: build run

# Build the application - specify output directory for multiple packages
build:
	@go build -o bin/main ./...

# Run tests - customize the pattern for test files
test:
	@go test ./... -coverprofile=coverage.out

# Run the application
run: build
	@./bin/main

# Clean up build artifacts
clean:
	@rm -rf bin coverage.out

# Code formatting
fmt:
	@go fmt ./...

# Static code analysis
vet:
	@go vet ./...

# Help target - displays available commands
help:
	@grep -E '^[^ ]+:.*?#?\{2\}' $(MAKEFILE_LIST) | sort

proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto
