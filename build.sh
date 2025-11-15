#!/bin/bash
# Build script for Linux

set -e

echo "Building Canvus PowerToys for Linux..."

go build -ldflags="-s -w" -o canvus-powertoys ./cmd/powertoys

echo "Build complete: canvus-powertoys"

