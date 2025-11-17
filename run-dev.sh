#!/bin/bash
# Quick development run script
# Processes assets and runs the app once (no live reload)
# Use ./dev.sh for live reload with air

set -e

echo "🚀 Starting Canvus PowerToys (quick dev mode)..."
echo ""

# Process WebUI assets first
if [ ! -d "webui/public" ] || [ "webui/src" -nt "webui/public" ]; then
    echo "📦 Processing WebUI assets..."
    go run webui/build/process-assets.go
fi

# Set development mode for WebUI
export WEBUI_DEV_MODE=1

echo "▶️  Running application..."
echo "   - WebUI changes will be visible after browser refresh"
echo "   - Go code changes require restart"
echo ""

# Run the app
go run ./cmd/powertoys

