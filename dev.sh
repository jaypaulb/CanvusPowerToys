#!/bin/bash
# Development script for Canvus PowerToys
# Runs the app with live reload using air

set -e

echo "🚀 Starting Canvus PowerToys in development mode..."
echo ""

# Process WebUI assets first (if needed)
if [ ! -d "webui/public" ] || [ "webui/src" -nt "webui/public" ]; then
    echo "📦 Processing WebUI assets..."
    go run webui/build/process-assets.go
fi

# Set development mode for WebUI
export WEBUI_DEV_MODE=1

# Create tmp directory if it doesn't exist
mkdir -p tmp

echo "🚀 Starting application (auto-reload disabled)..."
echo "   - Changes to .go files will NOT trigger automatic rebuild"
echo "   - Changes to WebUI files in webui/src will require asset processing"
echo "   - Press Ctrl+C to stop"
echo ""

# Run the app directly without live reload
go run ./cmd/powertoys

