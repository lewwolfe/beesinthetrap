.PHONY: all build test lint clean run windows linux darwin

# Define binary names
BINARY_NAME = beesinthetrap
WINDOWS_BINARY = $(BINARY_NAME)-windows-amd64.exe
LINUX_BINARY = $(BINARY_NAME)-linux-amd64
DARWIN_BINARY = $(BINARY_NAME)-darwin-amd64

# Default: Build and test
all: lint test build

# Build the application for current platform
build:
	go build -o $(BINARY_NAME) ./cmd/beesinthetrap/main.go

# Build for Windows
windows:
	GOOS=windows GOARCH=amd64 go build -o $(WINDOWS_BINARY) ./cmd/beesinthetrap/main.go

# Build for Linux
linux:
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_BINARY) ./cmd/beesinthetrap/main.go

# Build for macOS
darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(DARWIN_BINARY) ./cmd/beesinthetrap/main.go

# Build for all platforms
build-all: windows linux darwin

# Run tests
test:
	go test ./...

# Lint the code
lint:
	golangci-lint run ./...

# Clean generated files
clean:
	rm -f $(BINARY_NAME) $(WINDOWS_BINARY) $(LINUX_BINARY) $(DARWIN_BINARY)

# Run the application
run: build
	./$(BINARY_NAME)