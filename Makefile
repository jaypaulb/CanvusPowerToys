.PHONY: build build-linux build-windows test test-cover lint fmt vet clean

# Build for current platform
build:
	go build -ldflags="-s -w" -o canvus-powertoys ./cmd/powertoys

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o canvus-powertoys-linux ./cmd/powertoys

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o canvus-powertoys.exe ./cmd/powertoys

# Run tests
test:
	go test ./...

# Run tests with coverage
test-cover:
	go test -cover ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run linter (if golangci-lint is installed)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

# Clean build artifacts
clean:
	rm -f canvus-powertoys canvus-powertoys.exe canvus-powertoys-linux

# Run all quality checks
check: fmt vet test

