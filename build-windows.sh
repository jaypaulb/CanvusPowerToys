#!/bin/bash
# Build script for Windows (cross-compilation from Linux)

set -e

echo "Building Canvus PowerToys for Windows..."

GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o canvus-powertoys.exe ./cmd/powertoys

echo "Build complete: canvus-powertoys.exe"

