# Technology Stack & Architecture
## Canvus PowerToys

**Version:** 1.0.66
**Last Updated:** 2025-01-27
**Status:** Active Development

---

## Architecture Overview

Canvus PowerToys is a desktop application built with Go and the Fyne GUI framework. It provides a unified interface for managing Multitaction Canvus configurations through a tabbed GUI. The application includes an integrated, on-demand WebUI server for remote management capabilities.

### High-Level Architecture

```
┌─────────────────────────────────────────┐
│         Canvus PowerToys App            │
│  (Single Go Process - Portable .exe)    │
├─────────────────────────────────────────┤
│                                         │
│  ┌───────────────────────────────────┐ │
│  │      Fyne GUI Framework            │ │
│  │  ┌─────────────────────────────┐  │ │
│  │  │  Tab 1: Screen.xml Creator  │  │ │
│  │  ├─────────────────────────────┤  │ │
│  │  │  Tab 2: Config Editor        │  │ │
│  │  ├─────────────────────────────┤  │ │
│  │  │  Tab 3: CSS Options Manager  │  │ │
│  │  ├─────────────────────────────┤  │ │
│  │  │  Tab 4: Custom Menu Designer │  │ │
│  │  ├─────────────────────────────┤  │ │
│  │  │  Tab 5: WebUI Settings       │  │ │
│  │  └─────────────────────────────┘  │ │
│  └───────────────────────────────────┘ │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │   HTTP Server (On-Demand)          │ │
│  │   - Starts when WebUI enabled      │ │
│  │   - LAN accessible                 │ │
│  │   - Embedded static assets         │ │
│  │   - SSE for real-time updates      │ │
│  │   - REST API endpoints             │ │
│  └───────────────────────────────────┘ │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │   Core Services                   │ │
│  │   - File I/O (INI, XML, YAML)      │ │
│  │   - Validation Engine              │ │
│  │   - Backup Manager                 │ │
│  │   - System Tray                    │ │
│  │   - Canvas Tracking                │ │
│  │   - Workspace Subscription         │ │
│  └───────────────────────────────────┘ │
└─────────────────────────────────────────┘
         │                    │
         ▼                    ▼
┌─────────────────┐  ┌──────────────────┐
│  File System    │  │  MT Canvus       │
│  - mt-canvus.ini│  │  Server API      │
│  - screen.xml   │  │  (via WebUI)     │
│  - menu.yml     │  └──────────────────┘
│  - CSS plugins  │
└─────────────────┘
```

---

## Technology Choices

### Core Language: Go

**Rationale:**
- Single binary deployment (no dependencies)
- Excellent cross-platform support (Windows primary, Ubuntu future)
- Strong standard library
- Good performance
- Static compilation simplifies distribution
- Excellent concurrency support (goroutines, channels)
- Built-in HTTP server and JSON/XML/YAML support

**Version:** Go 1.24.5

**Key Standard Library Packages Used:**
- `net/http` - HTTP server and client
- `encoding/xml` - XML generation for screen.xml
- `encoding/json` - JSON processing for API communication
- `context` - Context management for cancellation and timeouts
- `embed` - Asset embedding (icons, WebUI static files)
- `os`, `path/filepath` - File system operations
- `fmt`, `log` - Logging and formatting

### GUI Framework: Fyne

**Rationale:**
- Cross-platform GUI framework for Go
- Modern, native-looking interfaces
- Good documentation and community
- Supports system tray integration
- Mature enough for production use
- Material Design-inspired components
- Custom theme support

**Version:** Fyne v2.7.1

**Key Features Used:**
- Tabbed container for main UI (`container.NewAppTabs`)
- Form widgets for configuration (`widget.NewForm`, `widget.NewEntry`)
- System tray integration (via `github.com/getlantern/systray`)
- File dialogs (`dialog.NewFileOpen`, `dialog.NewFileSave`)
- Custom widgets for grid display (Screen.xml Creator)
- Custom theme implementation (`theme.NewMTTheme()`)
- Embedded assets (icons, images)
- Window management and layout containers

**Theme Customization:**
- Custom MT theme with dark blue background
- MT pink text for entries
- Consistent styling across all tabs

### HTTP Server: net/http (stdlib)

**Rationale:**
- Built into Go standard library
- No external dependencies
- Sufficient for WebUI needs
- Good performance
- Easy to integrate with Fyne
- Supports Server-Sent Events (SSE) for real-time updates
- Built-in context support for graceful shutdown

**Usage:**
- Integrated WebUI server
- LAN accessible (bind to 0.0.0.0)
- On-demand activation
- Serves embedded static files
- REST API endpoints for canvas management
- SSE endpoints for real-time workspace updates
- File upload handling
- Graceful shutdown with context timeout

**Server Configuration:**
- ReadTimeout: 15 seconds
- WriteTimeout: 15 seconds
- IdleTimeout: 60 seconds
- Graceful shutdown: 5 seconds

---

## Third-Party Dependencies

### Required Libraries

#### 1. INI File Parsing
**Library:** `gopkg.in/ini.v1`
**Version:** v1.67.0
**Purpose:** Parse and write mt-canvus.ini files
**Rationale:** Mature library, handles Windows INI format well, actively maintained, supports section-based configuration

#### 2. YAML Processing
**Library:** `gopkg.in/yaml.v3`
**Version:** v3.0.1
**Purpose:** Parse and generate custom menu YAML files
**Rationale:** Standard Go YAML library, reliable, well-tested, supports full YAML specification

#### 3. XML Generation
**Library:** `encoding/xml` (stdlib)
**Purpose:** Generate screen.xml files
**Rationale:** Built into Go, no dependency needed, supports XML marshaling/unmarshaling

#### 4. System Tray
**Library:** `github.com/getlantern/systray`
**Version:** v1.2.2
**Purpose:** System tray integration (minimize to tray only)
**Rationale:** Cross-platform, works well with Fyne, actively maintained, supports Windows and Linux

#### 5. Asset Minification
**Library:** `github.com/tdewolff/minify/v2`
**Version:** v2.24.7
**Purpose:** Minify CSS, JavaScript, and HTML assets for WebUI
**Rationale:** Reduces asset size, improves load times, supports CSS/JS/HTML minification, used in build process

**Usage:**
- Minifies WebUI assets during build (`webui/build/process-assets.go`)
- Processes files from `webui/src/` to `webui/public/`
- Calculates size reduction statistics
- Integrated into Makefile build process

### Indirect Dependencies (via Fyne and other libraries)

#### Fyne Dependencies
- `fyne.io/systray` - System tray support for Fyne
- `github.com/go-gl/gl` - OpenGL bindings
- `github.com/go-gl/glfw/v3.3/glfw` - GLFW window management
- `github.com/go-text/render` - Text rendering
- `github.com/go-text/typesetting` - Text typesetting
- `golang.org/x/image` - Image processing
- `golang.org/x/sys` - System-specific functionality
- `golang.org/x/text` - Text processing and internationalization

#### System Integration
- `github.com/godbus/dbus/v5` - D-Bus integration (Linux)
- `github.com/rymdport/portal` - XDG Desktop Portal (Linux)

#### Testing
- `github.com/stretchr/testify` - Testing framework (for test files)

---

## Development Tools & Build System

### Build System: Makefile

**Purpose:** Standardized build commands and asset processing

**Key Targets:**
- `build` - Build for current platform with asset processing
- `build-linux` - Build Linux executable (amd64)
- `build-windows` - Build Windows executable with version info
- `test` - Run all tests
- `test-cover` - Run tests with coverage
- `lint` - Run golangci-lint
- `fmt` - Format code with go fmt
- `vet` - Run go vet
- `clean` - Remove build artifacts
- `process-assets` - Minify WebUI assets

**Build Flags:**
- `-ldflags="-s -w"` - Strip symbols and debug info (reduces binary size)
- `-H windowsgui` - Windows GUI application (no console window)
- Version injection via `-X` flags for build metadata

### Development Tools

#### 1. Air (Live Reload)
**Tool:** `github.com/cosmtrek/air`
**Purpose:** Live reload during development
**Configuration:** `.air.toml`
**Features:**
- Automatic rebuild on file changes
- Excludes test files and build artifacts
- Watches Go files, templates, and HTML
- Configurable build commands and delays

#### 2. golangci-lint
**Tool:** `github.com/golangci/golangci-lint`
**Purpose:** Comprehensive Go linting and code quality
**Configuration:** `.golangci.yml`
**Enabled Linters:**
- `errcheck` - Check for unchecked errors
- `goconst` - Find repeated strings
- `gocritic` - Advanced code analysis
- `gofmt` - Code formatting
- `goimports` - Import formatting
- `golint` - Go linting
- `gosec` - Security analysis
- `gosimple` - Simplify code
- `govet` - Go vet checks
- `ineffassign` - Detect ineffectual assignments
- `misspell` - Spell checking
- `staticcheck` - Static analysis
- `typecheck` - Type checking
- `unused` - Detect unused code
- `revive` - Fast, configurable linter

**Settings:**
- Timeout: 5 minutes
- Cyclomatic complexity threshold: 15
- Spell check locale: US

#### 3. goversioninfo
**Tool:** `github.com/josephspurrier/goversioninfo/cmd/goversioninfo`
**Purpose:** Generate Windows version info resources
**Configuration:** `versioninfo.json`
**Features:**
- Embeds version information in Windows executables
- Sets file description, company name, copyright
- Configures application icon
- Optional (build works without it, but without version info)

### Asset Processing

**Build Step:** `webui/build/process-assets.go`
**Purpose:** Minify and prepare WebUI assets for embedding

**Process:**
1. Reads source files from `webui/src/`
2. Minifies CSS, JavaScript, and HTML using `tdewolff/minify`
3. Writes minified files to `webui/public/`
4. Calculates and reports size reduction
5. Maintains directory structure

**Asset Structure:**
```
webui/
├── src/              # Source files (development)
│   ├── atoms/        # Atomic components (CSS, JS)
│   ├── molecules/    # Composed components
│   ├── organisms/    # Complex modules
│   ├── pages/        # Page templates and logic
│   └── templates/    # Reusable templates
└── public/           # Minified assets (production)
    └── [mirrors src structure]
```

**Embedding:**
- Assets embedded at compile time using `embed` package
- No external file dependencies for WebUI
- Single binary includes all static assets

---

## Project Structure

### Atomic Design Architecture

The project follows atomic design principles with strict file size limits:

```
canvus-powertoys/
├── cmd/
│   └── powertoys/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── atoms/                    # Basic building blocks (<100 lines)
│   │   ├── assets/               # Asset embedding
│   │   ├── backup/
│   │   │   └── manager.go        # Backup management
│   │   ├── config/
│   │   │   ├── ini_parser.go     # INI file parsing
│   │   │   ├── xml_generator.go # XML generation
│   │   │   └── yaml_handler.go  # YAML processing
│   │   ├── errors/
│   │   │   └── errors.go         # Error definitions
│   │   ├── logger/
│   │   │   └── logger.go         # Logging utilities
│   │   ├── paths/
│   │   │   └── paths.go          # Path management
│   │   ├── shortcut/
│   │   │   ├── shortcut_stub.go  # Stub for non-Windows
│   │   │   └── shortcut_windows.go # Windows shortcuts
│   │   ├── theme/
│   │   │   └── mt_theme.go       # Custom MT theme
│   │   ├── validation/
│   │   │   └── validator.go      # Validation utilities
│   │   ├── version/
│   │   │   └── version.go        # Version information
│   │   └── webui/
│   │       ├── api_client.go     # Canvus API client
│   │       ├── canvas_tracker.go # Canvas tracking
│   │       ├── client_resolver.go # Client ID resolution
│   │       ├── device_name.go    # Device name utilities
│   │       ├── widget.go         # Widget definitions
│   │       ├── workspace_subscriber.go # Workspace SSE
│   │       └── zone_bounding_box.go # Zone calculations
│   │
│   ├── molecules/                # Composed components (<200 lines)
│   │   ├── configeditor/
│   │   │   ├── compound_entry_group.go # Compound entries
│   │   │   ├── config_option.go  # Config option widget
│   │   │   ├── config_schema.go  # Configuration schema
│   │   │   ├── editor.go         # Main editor component
│   │   │   ├── form_control.go   # Form control widgets
│   │   │   ├── ini_file_parser.go # INI file parsing
│   │   │   ├── schema_embedded.go # Embedded schema
│   │   │   └── section_group.go  # Section grouping
│   │   ├── cssoptions/
│   │   │   └── manager.go        # CSS options management
│   │   ├── custommenu/
│   │   │   ├── designer.go       # Menu designer UI
│   │   │   └── icons.go          # Icon management
│   │   ├── screenxml/
│   │   │   ├── cell_editor.go    # Cell editing
│   │   │   ├── cell_widget.go    # Cell widget
│   │   │   ├── creator.go        # Main creator component
│   │   │   ├── fast_index.go     # Fast index assignment
│   │   │   ├── gpu_assignment.go # GPU output assignment
│   │   │   ├── grid_container.go # Grid container
│   │   │   ├── grid_widget.go    # Grid widget
│   │   │   ├── ini_integration.go # INI integration
│   │   │   ├── resolution.go     # Resolution detection
│   │   │   ├── touch_area.go     # Touch area management
│   │   │   └── xml_generator.go   # XML generation
│   │   ├── tray/
│   │   │   ├── tray_stub.go      # Stub for non-Windows
│   │   │   └── tray.go           # System tray implementation
│   │   └── webui/
│   │       ├── admin_handler.go  # Admin endpoints
│   │       ├── api_routes.go     # API route registration
│   │       ├── api_routes_integration_test.go # Integration tests
│   │       ├── canvas_service.go # Canvas service
│   │       ├── canvas_service_test.go # Service tests
│   │       ├── embed_assets.go   # Asset embedding
│   │       ├── macros_handler.go # Macros endpoints
│   │       ├── macros_handlers_common.go # Macros utilities
│   │       ├── macros_operations.go # Macro operations
│   │       ├── manager.go        # WebUI manager
│   │       ├── pages_handler.go  # Pages endpoints
│   │       ├── pages_zones.go    # Zones management
│   │       ├── pages_zones_test.go # Zones tests
│   │       ├── rcu_handler.go    # RCU endpoints
│   │       ├── sse_handler.go    # SSE handler
│   │       ├── static_handler.go # Static file serving
│   │       ├── upload_handler.go # File upload
│   │       └── workspace_subscription_test.go # Subscription tests
│   │
│   ├── organisms/               # Complex modules (<500 lines)
│   │   ├── app/
│   │   │   └── main_window.go    # Main application window
│   │   ├── services/
│   │   │   └── file_service.go   # File service
│   │   └── webui/
│   │       └── server.go         # WebUI HTTP server
│   │
│   └── templates/              # Reusable patterns (<200 lines)
│       └── [future templates]
│
├── webui/                       # WebUI static files and build
│   ├── build/
│   │   └── process-assets.go    # Asset minification tool
│   ├── embed.go                 # Asset embedding
│   ├── public/                  # Minified assets (generated)
│   │   ├── atoms/               # Atomic components
│   │   ├── molecules/           # Composed components
│   │   ├── organisms/           # Complex modules
│   │   ├── pages/               # Page templates
│   │   └── templates/           # Reusable templates
│   ├── src/                     # Source assets (development)
│   │   └── [mirrors public structure]
│   └── README-DEV.md            # WebUI development docs
│
├── assets/                      # Application assets
│   ├── embed.go                 # Asset embedding
│   ├── icons/                   # Application icons
│   │   ├── actions/             # Action icons
│   │   ├── categories/          # Category icons
│   │   ├── documents/           # Document icons
│   │   ├── navigation/          # Navigation icons
│   │   ├── CanvusPowerToysIcon.ico
│   │   ├── CanvusPowerToysIcon.png
│   │   └── README.md
│   └── templates/               # Template files
│
├── tests/                       # Test files (mirror structure)
│   ├── atoms/
│   ├── molecules/
│   └── organisms/
│
├── dev/                         # Development documentation
│   ├── PRD.md                   # Product Requirements
│   ├── TASKS.md                 # Task tracking
│   └── TECH_STACK.md            # This file
│
├── go.mod                       # Go module definition
├── go.sum                       # Go dependency checksums
├── Makefile                     # Build system
├── .air.toml                    # Air configuration
├── .golangci.yml                # Linter configuration
├── versioninfo.json             # Windows version info
├── build.sh                     # Build script
├── build-windows.sh             # Windows build script
├── README.md                    # Project README
└── .gitignore                   # Git ignore rules
```

---

## WebUI Architecture

### Overview

The WebUI is a fully integrated web interface served by the Go application. It provides remote management capabilities for Canvus systems via a web browser.

### Architecture Components

#### 1. Static Asset Serving
- **Handler:** `webuimolecules.StaticHandler`
- **Purpose:** Serve embedded static files (HTML, CSS, JS)
- **Location:** Files embedded from `webui/public/`
- **Features:**
  - Embedded at compile time (no external files)
  - Minified assets for optimal performance
  - Atomic design structure (atoms, molecules, organisms, pages, templates)

#### 2. API Routes
- **Handler:** `webuimolecules.APIRoutes`
- **Purpose:** REST API endpoints for canvas management
- **Endpoints:**
  - `/api/canvas/*` - Canvas operations
  - `/api/pages/*` - Page management
  - `/api/zones/*` - Zone management
  - `/api/macros/*` - Macro operations
  - `/api/rcu/*` - RCU (Remote Control Unit) operations
  - `/api/admin/*` - Administrative operations
  - `/api/upload` - File upload

#### 3. Server-Sent Events (SSE)
- **Handler:** `webuimolecules.SSEHandler`
- **Purpose:** Real-time workspace updates
- **Features:**
  - Subscribes to Canvus Server workspace events
  - Pushes updates to connected clients
  - Graceful connection handling
  - Context-aware cancellation

#### 4. Canvas Service
- **Service:** `webuimolecules.CanvasService`
- **Purpose:** Canvas tracking and management
- **Features:**
  - Resolves client_id from device name
  - Tracks canvas state
  - Manages workspace subscriptions
  - Handles canvas lifecycle

#### 5. API Client
- **Client:** `webuiatoms.APIClient`
- **Purpose:** Communication with Canvus Server API
- **Features:**
  - HTTP client with authentication
  - Token-based authentication
  - Error handling and retries
  - Request/response logging

### WebUI Pages

#### Implemented Pages:
1. **Main/Dashboard** (`pages/main.html`)
   - Overview and navigation
   - Canvas status
   - Quick actions

2. **Pages Management** (`pages/pages.html`)
   - Page listing and management
   - Zone configuration
   - Page operations

3. **Macros** (`pages/macros.html`)
   - Macro creation and execution
   - Macro library
   - Macro scheduling

4. **RCU** (`pages/rcu.html`)
   - Remote Control Unit interface
   - RCU operations
   - Device management

5. **Remote Upload** (`pages/remote-upload.html`)
   - File upload interface
   - Content management
   - Upload history

### Asset Structure (Atomic Design)

```
webui/public/
├── atoms/              # Basic building blocks
│   ├── css/           # Atomic CSS components
│   │   ├── badge.css
│   │   ├── button.css
│   │   ├── card.css
│   │   ├── input.css
│   │   └── link.css
│   └── js/            # Atomic JavaScript
│       └── error-handler.js
│
├── molecules/         # Composed components
│   ├── css/           # Molecule CSS
│   │   ├── canvas-header.css
│   │   ├── form-group.css
│   │   ├── navbar.css
│   │   └── page-card.css
│   └── js/            # Molecule JavaScript
│       └── workspace-client.js
│
├── organisms/         # Complex modules
│   ├── css/           # Organism CSS
│   └── js/            # Organism JavaScript
│
├── pages/             # Complete pages
│   ├── css/           # Page-specific CSS
│   │   ├── macros.css
│   │   ├── rcu.css
│   │   └── remote-upload.css
│   ├── html/          # Page HTML
│   │   ├── main.html
│   │   ├── pages.html
│   │   ├── macros.html
│   │   ├── rcu.html
│   │   └── remote-upload.html
│   └── js/            # Page JavaScript
│       ├── common.js
│       ├── pages.js
│       ├── macros.js
│       ├── rcu.js
│       └── remote-upload.js
│
├── templates/         # Reusable templates
│   ├── css/           # Template CSS
│   │   ├── modal-template.css
│   │   └── page-template.css
│   └── html/          # Template HTML
│       ├── modal-template.html
│       └── page-template.html
│
└── css/               # Global styles
    ├── dark-theme.css
    ├── design-system.css
    └── responsive.css
```

---

## Build & Distribution

### Development Build
```bash
# Process assets and build
make build

# Or manually:
go run webui/build/process-assets.go
go build -ldflags="-s -w" -o canvus-powertoys ./cmd/powertoys
```

### Windows Release Build
```bash
# Full build with version info
make build-windows

# Requirements:
# - mingw-w64: sudo apt-get install gcc-mingw-w64-x86-64
# - goversioninfo (optional): go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest

# Build process:
# 1. Process WebUI assets
# 2. Generate version info resource (if goversioninfo available)
# 3. Cross-compile with CGO enabled
# 4. Clean up resource files
```

**Build Flags:**
- `CGO_ENABLED=1` - Required for Fyne on Windows
- `GOOS=windows GOARCH=amd64` - Windows 64-bit target
- `CC=x86_64-w64-mingw32-gcc` - MinGW cross-compiler
- `-ldflags="-s -w -H windowsgui -X ..."` - Strip symbols, hide console, inject version
- `-H windowsgui` - Windows GUI application (no console)

### Linux Release Build
```bash
make build-linux

# Or manually:
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o canvus-powertoys-linux ./cmd/powertoys
```

### Build Flags Explained
- `-s`: Omit symbol table and debug information
- `-w`: Omit DWARF symbol table
- `-H windowsgui`: Windows GUI application (no console window)
- `-X`: Inject build-time variables (version, build date, git commit)

**Size Reduction:**
- Symbol stripping reduces binary size by ~20-30%
- Minified WebUI assets reduce total size significantly
- Single binary includes all assets (no external dependencies)

---

## Development Environment

### Required Tools
- **Go:** 1.24.5 or later
- **Git:** For version control
- **Make:** For build system (optional, can use scripts directly)
- **Code Editor:** VS Code, GoLand, or preferred editor with Go support

### Recommended Tools
- **golangci-lint:** Comprehensive linting (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)
- **Air:** Live reload during development (`go install github.com/cosmtrek/air@latest`)
- **goversioninfo:** Windows version info (optional) (`go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest`)
- **Delve:** Debugger for Go (`go install github.com/go-delve/delve/cmd/dlv@latest`)

### Development Workflow
1. **Start Development Server:**
   ```bash
   air  # Live reload on file changes
   ```

2. **Run Tests:**
   ```bash
   make test          # Run all tests
   make test-cover    # Run with coverage
   make test-verbose  # Verbose output
   ```

3. **Code Quality:**
   ```bash
   make fmt   # Format code
   make vet   # Run go vet
   make lint  # Run linter
   make check # Run all checks
   ```

4. **Build for Testing:**
   ```bash
   make build         # Current platform
   make build-linux   # Linux executable
   make build-windows # Windows executable
   ```

### Development on Linux, Deploy on Windows

**Cross-Compilation Setup:**
```bash
# Install MinGW for Windows cross-compilation
sudo apt-get install gcc-mingw-w64-x86-64

# Verify installation
x86_64-w64-mingw32-gcc --version
```

**Build Process:**
1. Develop and test on Linux
2. Build Windows executable: `make build-windows`
3. Transfer to Windows for testing
4. Iterate based on feedback

**Testing Considerations:**
- File path differences (Windows vs Linux)
- System tray behavior
- GPU detection (Windows-specific)
- File permissions and locations

---

## Platform-Specific Considerations

### Windows

**File Paths:**
- User config: `%APPDATA%\MultiTaction\canvus\mt-canvus.ini`
- System config: `%ProgramData%\MultiTaction\canvus\mt-canvus.ini`
- Logs: `%LOCALAPPDATA%\MultiTaction\Canvus\logs\`
- Upload directory: User-configurable

**GPU Detection:**
- Windows-specific implementation needed
- May use WMI or DXGI APIs
- Fallback to manual entry if detection fails

**System Tray:**
- Native Windows API via `github.com/getlantern/systray`
- Minimize to tray functionality
- Context menu support

**Shortcuts:**
- Windows shortcut creation (`internal/atoms/shortcut/shortcut_windows.go`)
- Desktop and Start Menu shortcuts

**Version Info:**
- Embedded via `goversioninfo` and `resource.syso`
- Includes icon, version, company, copyright
- Visible in Windows file properties

### Ubuntu (Future)

**File Paths:**
- User config: `~/MultiTaction/canvus/mt-canvus.ini`
- System config: `/etc/MultiTaction/canvus/mt-canvus.ini`
- Logs: `~/.local/share/MultiTaction/Canvus/logs/` or similar
- XDG-compliant paths

**GPU Detection:**
- Use `xrandr` or similar Linux tools
- May require additional dependencies

**System Tray:**
- Native Linux API via `github.com/getlantern/systray`
- D-Bus integration for system tray
- XDG Desktop Portal support

**Shortcuts:**
- Desktop entry files (`.desktop`)
- Application menu integration

---

## Security Considerations

### Token Storage
- Store Canvus Server Private-Token securely
- Encrypt token in app configuration (future enhancement)
- Display only last 6 digits when saved
- Never log full tokens
- Token used only for Canvus Server API authentication

### File Operations
- Validate all file paths
- Check file permissions before operations
- Create backups before modifications
- Handle errors gracefully
- Sanitize user-provided paths
- Prevent directory traversal attacks

### WebUI Security
- **No password/access control** - LAN only, trusted network assumption
- Token used only for Canvus Server API authentication
- Validate all inputs from WebUI
- Sanitize file paths in uploads
- CORS considerations (if needed in future)
- Rate limiting (future enhancement)

### Input Validation
- All user inputs validated
- Type checking for configuration values
- Range validation for numeric inputs
- Format validation for file paths
- Prevent injection attacks (XML, YAML, INI)

### Error Handling
- Don't expose sensitive information in errors
- Log errors securely (no secrets in logs)
- Return generic error messages to users
- Handle errors gracefully without crashing

---

## Performance Considerations

### Startup Time
- Lazy load tabs (load on first access)
- Defer heavy operations
- Cache configuration file reads
- Minimize asset loading time

### Memory Usage
- Release resources when tabs are closed
- Limit backup retention (configurable)
- Efficient file I/O operations
- Embedded assets loaded on demand

### UI Responsiveness
- Run file operations in goroutines
- Update UI on main thread only (Fyne requirement)
- Show progress indicators for long operations
- Non-blocking operations where possible

### WebUI Performance
- Minified assets reduce load time
- Embedded assets eliminate file I/O
- SSE connections managed efficiently
- Connection pooling for API requests
- Graceful degradation on errors

### Build Performance
- Asset minification as separate build step
- Parallel test execution
- Incremental builds (Go compiler)
- Caching of build artifacts

---

## Testing Strategy

### Test Structure
- Tests mirror source structure (atomic design)
- Unit tests for atoms
- Integration tests for molecules
- System tests for organisms
- E2E tests for complete features

### Test Coverage
- Target: >80% coverage
- Critical paths: 100% coverage
- Public APIs: 100% coverage
- Edge cases: Must be tested (zero values, empty states, nulls)
- Error conditions: All error paths tested

### Test Tools
- Standard `testing` package
- `github.com/stretchr/testify` for assertions
- Mock external dependencies
- Test fixtures for complex scenarios

### Running Tests
```bash
make test          # Run all tests
make test-cover    # Run with coverage report
make test-verbose  # Verbose output
go test ./...      # Run all tests manually
go test -v ./internal/molecules/webui/...  # Test specific package
```

---

## Version Management

### Version Information
- **Current Version:** 1.0.66
- **Location:** `internal/atoms/version/version.go`
- **Build Metadata:** Injected at build time (BuildDate, GitCommit)

### Version Components
- **Version:** Semantic versioning (Major.Minor.Patch)
- **BuildDate:** UTC timestamp of build
- **GitCommit:** Short git commit hash
- **AppName:** "Canvus PowerToys"
- **AppID:** "com.canvus.powertoys"

### Windows Version Info
- **Configuration:** `versioninfo.json`
- **Tool:** `goversioninfo` (optional)
- **Includes:** File version, product version, company, copyright, icon

### Version Display
- Shown in application (if implemented)
- Embedded in Windows executable properties
- Available via `version.GetFullVersion()`

---

## Architecture Decision Records (ADR)

### ADR-001: Go + Fyne for Desktop GUI

**Status:** Accepted
**Context:** Need cross-platform desktop application with modern GUI
**Decision:** Use Go with Fyne framework
**Consequences:**
- ✅ Single binary deployment
- ✅ Cross-platform support
- ✅ Good performance
- ✅ Custom theme support
- ⚠️ Less mature ecosystem than Qt/PyQt
- ⚠️ Limited widget library compared to web frameworks

**Alternatives Considered:**
- Python + PyQt/PySide (larger runtime, slower startup)
- C# + WPF (Windows-only, complex cross-platform)
- Electron (large memory footprint, overkill)
- C++ + Qt (more complex, longer development time)

---

### ADR-002: Integrated On-Demand WebUI Server

**Status:** Accepted
**Context:** Need WebUI for remote management, but not always needed
**Decision:** Integrate HTTP server in same Go process, start on-demand
**Consequences:**
- ✅ Single process, simpler deployment
- ✅ Resource efficient (only runs when needed)
- ✅ Easy to enable/disable
- ✅ Embedded assets (no external files)
- ✅ Real-time updates via SSE
- ⚠️ WebUI migration complexity (deferred to separate project)

**Alternatives Considered:**
- Separate process (more complex, resource overhead)
- Always running (wasteful if not used)
- External service (deployment complexity)

---

### ADR-003: Form-Based Validation (No Raw Text Editing)

**Status:** Accepted
**Context:** Prevent configuration errors, improve user experience
**Decision:** Form-based UI prevents invalid configurations, no raw text editing
**Consequences:**
- ✅ Zero configuration errors from invalid input
- ✅ Better user experience
- ✅ Real-time validation feedback
- ✅ Type-safe configuration
- ⚠️ More complex UI development
- ⚠️ Less flexibility for advanced users

**Alternatives Considered:**
- Text editor with validation (still allows errors)
- Hybrid approach (complexity)

---

### ADR-004: Smart Backup with Rotation

**Status:** Accepted
**Context:** Need to protect user configurations, but avoid backup bloat
**Decision:** Only backup if file changed, maintain rotation of last N backups
**Consequences:**
- ✅ Efficient storage usage
- ✅ Protection against data loss
- ✅ Automatic cleanup
- ⚠️ Slightly more complex backup logic

**Alternatives Considered:**
- Always backup (wasteful, many identical backups)
- No backups (risky)
- User-configurable (complexity)

---

### ADR-005: Single Portable Executable

**Status:** Accepted
**Context:** Easy distribution, no installation complexity
**Decision:** Single .exe file, no installer
**Consequences:**
- ✅ Simple distribution
- ✅ No installation required
- ✅ Can run from anywhere
- ✅ Embedded assets (no external files)
- ⚠️ No automatic updates
- ⚠️ No system integration (start menu, etc.)

**Alternatives Considered:**
- Windows installer (more professional, but complex)
- Both options (maintenance overhead)

---

### ADR-006: Atomic Design Architecture

**Status:** Accepted
**Context:** Need maintainable, scalable codebase with clear structure
**Decision:** Follow atomic design principles with strict file size limits
**Consequences:**
- ✅ Small, focused files (easy to understand)
- ✅ Clear hierarchy (atoms → molecules → organisms)
- ✅ Reusability (atoms and molecules reusable)
- ✅ Testability (each level testable independently)
- ✅ Maintainability (changes localized)
- ⚠️ More files to manage
- ⚠️ Requires discipline to maintain structure

**File Size Limits:**
- Atoms: <100 lines (ideally <50)
- Molecules: <200 lines (ideally <150)
- Organisms: <500 lines (ideally <300)
- Templates: <200 lines
- Pages: <500 lines

---

### ADR-007: Asset Minification and Embedding

**Status:** Accepted
**Context:** Need efficient WebUI assets, single binary deployment
**Decision:** Minify assets at build time, embed in binary
**Consequences:**
- ✅ Reduced asset size (CSS, JS, HTML minified)
- ✅ Single binary (no external files)
- ✅ Faster load times
- ✅ Simplified deployment
- ⚠️ Build step required (asset processing)
- ⚠️ Larger binary size (but acceptable)

**Process:**
- Source files in `webui/src/`
- Minified to `webui/public/` during build
- Embedded at compile time via `embed` package

---

## Future Considerations

### WebUI Migration
- Separate project step to analyze `canvus-webui`
- Define scope of features to port
- Determine integration approach
- Estimate effort and timeline
- Current status: Core WebUI features implemented, full migration deferred

### Ubuntu Support
- Test on Ubuntu
- Verify file path handling
- Test GPU detection (xrandr)
- Validate system tray functionality
- Test desktop integration

### Additional Features
- Auto-update mechanism (optional)
- Plugin system for extensions
- Advanced screen.xml configurations
- Multi-language support (i18n)
- Enhanced error reporting
- Telemetry (opt-in)

### Performance Optimizations
- Lazy loading of WebUI pages
- Asset compression (gzip)
- Connection pooling improvements
- Caching strategies
- Background processing for heavy operations

---

## Dependencies Summary

### Direct Dependencies
```
fyne.io/fyne/v2 v2.7.1
github.com/getlantern/systray v1.2.2
github.com/tdewolff/minify/v2 v2.24.7
gopkg.in/ini.v1 v1.67.0
gopkg.in/yaml.v3 v3.0.1
```

### Key Indirect Dependencies
- Fyne ecosystem (systray, gl, glfw, text rendering)
- System integration (dbus, portal)
- Testing (testify)

### Standard Library (Key Packages)
- `net/http` - HTTP server and client
- `encoding/xml` - XML processing
- `encoding/json` - JSON processing
- `context` - Context management
- `embed` - Asset embedding
- `os`, `path/filepath` - File operations
- `fmt`, `log` - Logging

---

**Document Status:** Active - Updated with current implementation
**Last Updated:** 2025-01-27
**Maintainer:** Development Team
