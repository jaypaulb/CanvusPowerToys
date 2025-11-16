#!/bin/bash
# Development server script for WebUI
# This script sets WEBUI_DEV_MODE and runs the application

# First, process assets to ensure webui/public exists
echo "Processing WebUI assets..."
go run webui/build/process-assets.go

# Set development mode and run
echo "Starting development server..."
export WEBUI_DEV_MODE=1
go run ./cmd/powertoys

