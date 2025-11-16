# WebUI Development Mode

## Overview

Development mode allows you to test WebUI changes without recompiling the entire application. Files are served directly from `webui/public` instead of embedded assets.

## Usage

### Option 1: Using the dev server script

```bash
./webui/dev-server.sh
```

This script will:
1. Process assets (minify and copy to `webui/public`)
2. Set `WEBUI_DEV_MODE=1`
3. Run the application

### Option 2: Manual setup

1. Process assets first:
```bash
go run webui/build/process-assets.go
```

2. Set environment variable and run:
```bash
export WEBUI_DEV_MODE=1
go run ./cmd/powertoys
```

Or in one line:
```bash
WEBUI_DEV_MODE=1 go run ./cmd/powertoys
```

## How It Works

- **Development Mode** (`WEBUI_DEV_MODE=1`): Files are served directly from `webui/public` directory. Changes to HTML/CSS/JS files are immediately visible after refreshing the browser.

- **Production Mode** (default): Files are embedded in the binary using `//go:embed`. You must rebuild the application to see changes.

## Workflow

1. Make changes to files in `webui/src/`
2. Run `go run webui/build/process-assets.go` to process assets (or use the dev script)
3. Start the server with `WEBUI_DEV_MODE=1`
4. Refresh your browser to see changes
5. No need to recompile Go code unless you change backend logic

## Notes

- The dev server automatically finds `webui/public` relative to the project root
- If `webui/public` is not found, it falls back to embedded assets
- Always process assets before starting dev mode to ensure `webui/public` exists
- Changes to Go code still require restarting the server

