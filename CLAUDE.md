# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Quick Commands

### Build
```bash
make build           # Build for current platform
make build-linux     # Build for Linux
make build-windows   # Build for Windows (requires mingw-w64)
```

### Development
```bash
make test            # Run all tests
make test-cover      # Run tests with coverage report
make test-verbose    # Run tests with verbose output
make fmt             # Format code (gofmt)
make vet             # Run go vet
make lint            # Run golangci-lint (if installed)
make check           # Run fmt, vet, and test
```

### Single Test Execution
```bash
go test ./internal/molecules/webui -v -run TestCanvasService
go test ./tests/atoms/paths -v
```

### WebUI Asset Processing
```bash
go run webui/build/process-assets.go  # Minifies CSS, JS, HTML before build
```

### Clean
```bash
make clean           # Remove build artifacts
```

## Project Architecture

### Atomic Design Structure
The codebase follows a strict layered architecture (inspired by atomic design):

- **Atoms** (`internal/atoms/`) - Single-responsibility building blocks (<100 lines)
  - `config/` - INI/XML/YAML parsing and generation
  - `paths/` - File path utilities for Windows/Linux
  - `backup/` - Backup manager for file safety
  - `logger/` - Centralized logging infrastructure
  - `validation/` - Validation logic and error handling
  - `webui/` - WebUI API client and canvas tracking
  - `version/` - Version info (also used in build-time flags)
  - `theme/` - Fyne theme customization
  - `shortcut/` - Platform-specific shortcuts (Windows/stub)

- **Molecules** (`internal/molecules/`) - Composed UI/business components (<200 lines)
  - `configeditor/` - Config Editor UI with INI file parsing and schema
  - `screenxml/` - Screen.xml Creator UI with 10x5 grid editor
  - `cssoptions/` - CSS Options Manager for kiosk/rotation modes
  - `custommenu/` - Custom Menu Designer with YAML generation
  - `tray/` - System tray integration (Windows/stub)
  - `webui/` - WebUI server handlers and SSE management

- **Organisms** (`internal/organisms/`) - Complex feature modules (<500 lines)
  - `app/` - Main window orchestration, tab management
  - `services/` - FileService for config path detection and management

### Key Data Flow Patterns

**Configuration File Handling:**
1. FileService detects config paths at startup
2. Config Editor reads INI via ini.v1 parser
3. Schema provides validation rules per section
4. Changes validated before save
5. Backup created before file write

**Screen.xml Creation:**
1. Creator initializes 10x5 grid widget
2. User assigns GPU outputs and touch areas
3. Generator creates XML structure
4. INI video-output areas updated in parallel

**WebUI Server:**
1. Manager initializes HTTP server with routes
2. Canvas service subscribes to workspace updates via API
3. SSE handler broadcasts real-time updates to clients
4. Macros operations (move/copy/pin) send commands back

### Important Design Notes

- **Platform Abstraction:** Use atoms/paths for Windows/Linux compatibility
- **Embedded Assets:** Icons and WebUI assets embedded via `//go:embed`
- **Error Handling:** errors.go provides typed errors for recovery UI
- **Logging:** Always use logger package (not log/fmt direct)
- **Testing:** Mirror source structure in tests/ folder, test atoms thoroughly

## Important Conventions

### File Organization
- Source: `internal/{atoms,molecules,organisms}/feature/`
- Tests: `tests/atoms/feature/` (mirrors atoms structure only)
- WebUI: `webui/src/` for source, `webui/public/` for minified assets
- Assets: `assets/` for icons and resources

### Version Management
- Current version in `internal/atoms/version/version.go`
- Build process sets version via `-X` ldflags from git/date
- Increment manually before building (check scripts/increment-version.sh if exists)

### WebUI Asset Building
- CSS/JS/HTML must be minified before binary build
- `go run webui/build/process-assets.go` handles this
- Processed assets embedded in binary at build time
- Source files in `webui/src/`, processed files generated

### Windows Cross-Compilation
- Requires: `gcc-mingw-w64-x86-64`
- Requires CGO: `CGO_ENABLED=1`
- Optional: `goversioninfo` for .exe version info and icons
- Build flag `-H windowsgui` hides console window

### Configuration File Paths
**Windows:**
- User: `%APPDATA%\MultiTaction\canvus\mt-canvus.ini`
- System: `%ProgramData%\MultiTaction\canvus\mt-canvus.ini`
- Logs: `%LOCALAPPDATA%\MultiTaction\Canvus\logs\`

**Linux/Ubuntu:**
- User: `~/.config/MultiTaction/canvus/mt-canvus.ini`
- System: `/etc/MultiTaction/canvus/mt-canvus.ini`
- Logs: `~/.local/share/MultiTaction/Canvus/logs/`

## Dependencies

**Direct Dependencies:**
- `fyne.io/fyne/v2` - GUI framework (with CGO on Windows)
- `github.com/getlantern/systray` - System tray integration
- `gopkg.in/ini.v1` - INI file parsing
- `gopkg.in/yaml.v3` - YAML parsing/generation
- `github.com/tdewolff/minify/v2` - Asset minification

**Platform-Specific Code:**
- `internal/atoms/shortcut/` - Windows shortcut creation
- `internal/molecules/tray/` - Windows system tray
- `internal/atoms/paths/` - Windows/Linux file paths

## Testing Strategy

- Minimum 80% coverage for non-UI code
- Test atoms thoroughly (they're reused)
- WebUI integration tests in molecules/webui/*_test.go
- Use `go test ./...` to run all tests
- Mock API responses for WebUI tests

## Development Workflow

1. **Make changes** in appropriate layer (atom/molecule/organism)
2. **Update tests** (TDD preferred)
3. **Run `make check`** to verify fmt/vet/tests pass
4. **Process assets** if WebUI changes: `go run webui/build/process-assets.go`
5. **Build target** (linux/windows) with appropriate make target
6. **Test built binary** on target platform before release

## Common Tasks

### Adding a New Feature Tab
1. Create molecule in `internal/molecules/featurename/`
2. Implement UI creation function returning `fyne.CanvasObject`
3. Add to MainWindow tabs in `internal/organisms/app/main_window.go`
4. Create supporting atoms if needed
5. Add tests following existing patterns

### Adding WebUI Routes
1. Add handler in `internal/molecules/webui/handler_name.go`
2. Register route in `api_routes.go`
3. Test with integration tests
4. Update WebUI frontend in `webui/src/`
5. Run asset processing before build

### Modifying Config Schema
1. Update `internal/molecules/configeditor/config_schema.go`
2. Update corresponding atoms/config parsers if needed
3. Test with existing and new config files
4. Update documentation if adding new options

## Debugging Tips

- Set `logger.Log()` calls throughout code (uses stdout in dev mode)
- Remove `-s -w` from build flags to keep debug symbols
- Remove `-H windowsgui` from Windows builds to show console
- Check logs in standard locations for production issues
- Use `go test -v` for verbose test output

## Documentation References

- [BUILD.md](BUILD.md) - Build from source, cross-compilation details
- [DEPLOYMENT.md](DEPLOYMENT.md) - Deployment, configuration, troubleshooting
- [README.md](README.md) - Features, architecture overview, quick start
