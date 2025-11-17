# Canvus PowerToys

A comprehensive desktop application for configuring and managing Multitaction Canvus installations. Canvus PowerToys provides intuitive GUI interfaces for configuring screen layouts, Canvus settings, CSS customizations, custom menus, and remote management via integrated WebUI.

## Features

### 🖥️ Screen.xml Creator
- Visual 10x5 grid editor for multi-display configurations
- GPU output assignment with number drawing method
- Touch area assignment via click/drag or manual entry
- Automatic resolution discovery with manual override
- Visual feedback with color coding and layout area indicators
- Integration with mt-canvus.ini for video-output areas

### ⚙️ Canvus Config Editor
- Auto-detection of mt-canvus.ini in standard locations
- Searchable and filterable interface for all INI options
- Grouped options by section or functionality
- Hover tooltips with embedded schema documentation
- Save to user config (`%APPDATA%\MultiTaction\canvus\mt-canvus.ini`)
- Save to system config (`%ProgramData%\MultiTaction\canvus\mt-canvus.ini`)
- Smart backup creation with rotation

### 🎨 CSS Options Manager
- Enable/disable rotation (temporary, reverts on canvas close)
- Enable/disable video looping (with memory warning)
- Kiosk mode with validation:
  - Validates default-canvas is set
  - Validates auto-pin=0
  - Hides UI layers appropriately
- Kiosk plus mode (kiosk + finger menu)
- Auto-generates CSS plugin files (.canvusplugin JSON + styles.css)
- Updates plugin-folders in mt-canvus.ini
- Handles plugin API version compatibility

### 📋 Custom Menu Designer
- Visual menu structure editor (hierarchical/tree view)
- Form-based item creation (not raw YAML editing)
- Icon picker with included icon set (16 icons in 4 categories)
- Content file browser for PDFs, videos, images
- Coordinate system configuration (viewport vs canvas)
- Positioning tools (location, size, origin, offset)
- Import existing menu.yml files
- Generates YAML format matching CUSTOM_MENU_MANUAL.md

### 🌐 WebUI Integration
- Integrated HTTP server for remote management
- On-demand activation (only when enabled)
- LAN accessible (bind to 0.0.0.0, default port 8080)
- Canvas tracking via ClientID/Workspace subscription
- Real-time canvas updates via Server-Sent Events (SSE)
- Secure token storage (encrypted)
- Mobile-responsive interface with dark mode support

#### WebUI Pages
- **Main Page**: Navigation hub with canvas header and connection status
- **Pages Management**: Create and manage canvas pages/zones
- **Macros**: Move, copy, grouping, and pinning operations
- **Remote Content Upload**: File upload interface for admins
- **RCU**: Remote content upload interface

### 🔧 Application Core
- Tabbed GUI interface (one tab per feature)
- System tray integration (close button hides to tray)
- Auto-detect configuration files in standard locations
- Support Windows and Ubuntu file paths
- Form-based validation (prevents invalid configurations)
- Real-time validation feedback
- Smart backup creation with rotation
- Comprehensive logging in dev mode
- Simple error dialogs in production mode

## System Requirements

### Windows (Primary Platform)
- Windows 10 or later
- .NET Framework not required (standalone executable)
- Administrator privileges for system config writes

### Linux/Ubuntu (Future Support)
- Ubuntu 20.04 or later
- GTK3 libraries (for Fyne framework)

## Installation

### Windows
1. Download the latest `canvus-powertoys.X.X.X.exe` from releases
2. Place the executable in your desired location (no installer required)
3. Run the executable - it will auto-detect configuration files

### Linux (Future)
1. Download the latest `canvus-powertoys-linux` from releases
2. Make executable: `chmod +x canvus-powertoys-linux`
3. Run: `./canvus-powertoys-linux`

## Quick Start

1. **Launch the application** - The app will auto-detect `mt-canvus.ini` and `screen.xml` in standard locations
2. **Choose a feature tab**:
   - **Screen.xml Creator**: Configure multi-display setups
   - **Canvus Config Editor**: Modify Canvus settings
   - **CSS Options Manager**: Enable kiosk mode and CSS customizations
   - **Custom Menu Designer**: Create custom menus
   - **WebUI Settings**: Enable remote management
3. **Make your changes** - All changes are validated in real-time
4. **Save** - Backups are created automatically before saving

## Project Structure

```
CanvusPowerToys/
├── cmd/powertoys/          # Main application entry point
├── internal/
│   ├── atoms/              # Basic building blocks (<100 lines)
│   │   ├── backup/         # Backup manager
│   │   ├── config/         # INI/XML/YAML handlers
│   │   ├── logger/         # Logging infrastructure
│   │   ├── paths/          # File path utilities
│   │   ├── validation/     # Validation logic
│   │   └── webui/          # WebUI client components
│   ├── molecules/          # Composed components (<200 lines)
│   │   ├── configeditor/   # Config Editor UI
│   │   ├── cssoptions/     # CSS Options Manager UI
│   │   ├── custommenu/     # Custom Menu Designer UI
│   │   ├── screenxml/      # Screen.xml Creator UI
│   │   ├── tray/           # System tray integration
│   │   └── webui/           # WebUI server components
│   └── organisms/          # Complex modules (<500 lines)
│       ├── app/            # Main window and application
│       └── services/       # File service
├── webui/                  # WebUI frontend assets
│   ├── public/             # Minified production assets
│   └── src/                # Source assets (CSS, JS, HTML)
├── assets/                 # Embedded application assets
├── tests/                  # Test suite (mirrors source structure)
├── build.sh                # Linux build script
├── build-windows.sh        # Windows build script
├── Makefile                # Build automation
└── README.md               # This file
```

## Documentation

- **[Build Instructions](BUILD.md)** - How to build from source
- **[Deployment Guide](DEPLOYMENT.md)** - Deployment and configuration guide
- **[Product Requirements](dev/PRD.md)** - Complete feature specifications
- **[Technical Stack](dev/TECH_STACK.md)** - Technology choices and architecture

## Development

### Prerequisites
- Go 1.24.5 or later
- Fyne framework dependencies
- For Windows cross-compilation: `gcc-mingw-w64-x86-64` (Linux)
- Optional: `goversioninfo` for Windows version info

### Building
See [BUILD.md](BUILD.md) for detailed build instructions.

Quick build:
```bash
# Linux
make build-linux

# Windows (from Linux)
make build-windows
```

### Testing
```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Run tests with verbose output
make test-verbose
```

### Code Quality
```bash
# Format code
make fmt

# Run go vet
make vet

# Run linter (if golangci-lint installed)
make lint

# Run all quality checks
make check
```

## Architecture

### Atomic Design Structure
The codebase follows atomic design principles:
- **Atoms**: Single-purpose, indivisible components (<100 lines)
- **Molecules**: Composed components (<200 lines)
- **Organisms**: Complex modules (<500 lines)
- **Templates**: Reusable patterns (<200 lines)
- **Pages**: Complete features (<500 lines)

### Test-Driven Development
- Tests written before implementation (Red-Green-Refactor)
- Test structure mirrors source structure
- Minimum 80% test coverage (100% for critical paths)

## Configuration Files

### Standard Locations

**Windows:**
- User config: `%APPDATA%\MultiTaction\canvus\mt-canvus.ini`
- System config: `%ProgramData%\MultiTaction\canvus\mt-canvus.ini`
- Logs: `%LOCALAPPDATA%\MultiTaction\Canvus\logs\`

**Linux/Ubuntu:**
- User config: `~/.config/MultiTaction/canvus/mt-canvus.ini`
- System config: `/etc/MultiTaction/canvus/mt-canvus.ini`
- Logs: `~/.local/share/MultiTaction/Canvus/logs/`

## WebUI Configuration

1. Open the **WebUI Settings** tab
2. Enable WebUI server
3. Configure Canvus Server URL (e.g., `https://canvus.example.com`)
4. Enter Private-Token (stored securely, only last 6 digits displayed)
5. Access WebUI from LAN devices at `http://<your-ip>:8080`

## Backup System

- Automatic backups created before all file saves
- Backup rotation (keeps last N backups)
- Backup naming: `filename.YYYYMMDD-HHMMSS.backup`
- Smart backup (only creates backup if file changed)

## Logging

- Development mode: Comprehensive logging to console
- Production mode: Simple error dialogs, logs to file
- Log location: `%LOCALAPPDATA%\MultiTaction\Canvus\logs\` (Windows)
- Log rotation: Automatic (keeps last N log files)

## Contributing

1. Follow atomic design principles
2. Write tests before implementation (TDD)
3. Maintain >80% test coverage
4. Follow Go code style guidelines
5. Update documentation as needed

## License

[Add license information here]

## Support

For issues, feature requests, or questions:
- Check [TASKS.md](dev/TASKS.md) for current development status
- Review [PRD.md](dev/PRD.md) for feature specifications
- Check logs in `%LOCALAPPDATA%\MultiTaction\Canvus\logs\` for troubleshooting

## Version

Current version: See `internal/atoms/version/version.go`

---

**Note**: This application is designed for AV technicians and meeting facilitators managing Multitaction Canvus installations in experience centres.
