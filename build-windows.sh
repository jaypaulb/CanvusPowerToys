#!/bin/bash
# Build script for Windows (cross-compilation from Linux)

set -e

echo "Building Canvus PowerToys for Windows..."

# Get version info
VERSION=$(grep 'Version.*=' internal/atoms/version/version.go | sed -n 's/.*"\([^"]*\)".*/\1/p')
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo "Version: $VERSION"
echo "Build Date: $BUILD_DATE"
echo "Git Commit: $GIT_COMMIT"

# Generate version info resource if goversioninfo is available
if command -v goversioninfo >/dev/null 2>&1; then
    echo "Generating version info resource..."
    goversioninfo -64 versioninfo.json
fi

# Use full version number in filename
OUTPUT_FILE="canvus-powertoys.$VERSION.exe"

# Build with version info
GOOS=windows GOARCH=amd64 go build \
    -ldflags="-s -w -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.Version=$VERSION -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.BuildDate=$BUILD_DATE -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.GitCommit=$GIT_COMMIT" \
    -o "$OUTPUT_FILE" ./cmd/powertoys

echo "Build complete: $OUTPUT_FILE"

