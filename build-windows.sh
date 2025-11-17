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
# The resource.syso file must be in the same directory as main.go for Go to include it
if command -v goversioninfo >/dev/null 2>&1; then
    echo "Generating version info resource in cmd/powertoys/..."
    goversioninfo -64 -o cmd/powertoys/resource.syso versioninfo.json
else
    echo "WARNING: goversioninfo not found. Install with: go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest"
    echo "WARNING: Building without Windows icon and version info."
fi

# Use full version number in filename
OUTPUT_FILE="canvus-powertoys.$VERSION.exe"

# Build with version info
# -H windowsgui: Sets Windows subsystem to GUI (hides console window)
#   To unhide console for debugging: Remove -H windowsgui from -ldflags
# -trimpath: Remove file system paths from binary for smaller size and reproducibility
GOOS=windows GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -H windowsgui -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.Version=$VERSION -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.BuildDate=$BUILD_DATE -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.GitCommit=$GIT_COMMIT" \
    -o "$OUTPUT_FILE" ./cmd/powertoys

# Clean up resource.syso file after build
if [ -f cmd/powertoys/resource.syso ]; then
    echo "Cleaning up resource.syso..."
    rm -f cmd/powertoys/resource.syso
fi

echo "Build complete: $OUTPUT_FILE"

