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
# Note: Fyne requires CGO for Windows cross-compilation
# Requires: sudo apt-get install gcc-mingw-w64-x86-64
# Optional: goversioninfo for .exe version info (go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest)
build-windows: process-assets
	@echo "Building for Windows (requires mingw-w64 for CGO)..."
	@echo "Note: Version increment is optional. Use scripts/increment-version.sh if needed."
	@if ! command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1; then \
		echo "ERROR: x86_64-w64-mingw32-gcc not found. Install with: sudo apt-get install gcc-mingw-w64-x86-64"; \
		exit 1; \
	fi
	@if command -v goversioninfo >/dev/null 2>&1; then \
		echo "Generating version info resource in cmd/powertoys/..."; \
		goversioninfo -64 -o cmd/powertoys/resource.syso versioninfo.json; \
	else \
		echo "WARNING: goversioninfo not found. Install with: go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest"; \
		echo "WARNING: Building without Windows icon and version info."; \
	fi
	@VERSION=$$(grep 'Version.*=' internal/atoms/version/version.go | sed -n 's/.*"\([^"]*\)".*/\1/p'); \
	BUILD_DATE=$$(date -u +"%Y-%m-%dT%H:%M:%SZ"); \
	GIT_COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	OUTPUT_FILE="canvus-powertoys.$$VERSION.exe"; \
	if CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build \
		-trimpath \
		-ldflags="-s -w -H windowsgui -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.Version=$$VERSION -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.BuildDate=$$BUILD_DATE -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.GitCommit=$$GIT_COMMIT" \
		-o $$OUTPUT_FILE ./cmd/powertoys; then \
		echo "Built: $$OUTPUT_FILE"; \
		if [ -f cmd/powertoys/resource.syso ]; then \
			echo "Cleaning up resource.syso..."; \
			rm -f cmd/powertoys/resource.syso; \
		fi \
	else \
		echo "ERROR: Build failed! Version was incremented to $$VERSION but build did not complete."; \
		exit 1; \
	fi

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

# Check WebUI asset sizes (development tool - moved to dross/)
# check-assets:
#	@./webui/build/size-check.sh

# Run all quality checks
check: fmt vet test

