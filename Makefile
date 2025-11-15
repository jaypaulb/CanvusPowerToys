.PHONY: build build-linux build-windows test test-cover lint fmt vet clean process-assets

# Process WebUI assets (minify CSS, JS, HTML)
process-assets:
	@echo "Processing WebUI assets..."
	@go run webui/build/process-assets.go

# Build for current platform
build: process-assets
	go build -ldflags="-s -w" -o canvus-powertoys ./cmd/powertoys

# Build for Linux
build-linux: process-assets
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o canvus-powertoys-linux ./cmd/powertoys

# Build for Windows
build-windows: process-assets
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

# Check WebUI asset sizes
check-assets:
	@./webui/build/size-check.sh

# Run all quality checks
check: fmt vet test check-assets

