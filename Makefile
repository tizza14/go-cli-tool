.PHONY: build test test-unit test-integration test-coverage clean install run help

# Binary name
BINARY_NAME=go-cli-tool
VERSION=1.0.0
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-X github.com/yourusername/go-cli-tool/cmd.Version=$(VERSION) \
                  -X github.com/yourusername/go-cli-tool/cmd.BuildDate=$(BUILD_DATE) \
                  -X github.com/yourusername/go-cli-tool/cmd.GitCommit=$(GIT_COMMIT)"

## help: Display this help message
help:
	@echo "Available targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

## build-all: Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe .
	@echo "Cross-platform build complete"

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) .
	@echo "Installation complete"

## run: Build and run the application
run: build
	./$(BINARY_NAME)

## test: Run all tests
test:
	@echo "Running all tests..."
	go test ./... -v -race

## test-unit: Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test ./... -v -short

## test-integration: Run integration tests
test-integration: build
	@echo "Running integration tests..."
	go test ./test -tags=integration -v

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## test-bench: Run benchmark tests
test-bench:
	@echo "Running benchmarks..."
	go test ./... -bench=. -benchmem

## lint: Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/"; \
	fi

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## mod: Tidy and verify modules
mod:
	@echo "Tidying modules..."
	go mod tidy
	go mod verify

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html
	go clean -testcache
	@echo "Clean complete"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

## check: Run all checks (fmt, vet, test)
check: fmt vet test
	@echo "All checks passed!"
