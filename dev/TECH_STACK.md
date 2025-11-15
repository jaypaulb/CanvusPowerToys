# Technology Stack & Architecture
## Canvus PowerToys

**Version:** 1.0
**Date:** 2024

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
│  │   - Shows only enabled pages       │ │
│  └───────────────────────────────────┘ │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │   Core Services                   │ │
│  │   - File I/O (INI, XML, YAML)      │ │
│  │   - Validation Engine              │ │
│  │   - Backup Manager                 │ │
│  │   - System Tray                    │ │
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

**Version:** Go 1.21+ (or latest stable)

### GUI Framework: Fyne

**Rationale:**
- Cross-platform GUI framework for Go
- Modern, native-looking interfaces
- Good documentation and community
- Supports system tray integration
- Mature enough for production use

**Version:** Fyne v2.4+ (or latest stable)

**Key Features Used:**
- Tabbed container for main UI
- Form widgets for configuration
- System tray integration
- File dialogs
- Custom widgets for grid display (Screen.xml Creator)

### HTTP Server: net/http (stdlib)

**Rationale:**
- Built into Go standard library
- No external dependencies
- Sufficient for WebUI needs
- Good performance
- Easy to integrate with Fyne

**Usage:**
- Integrated WebUI server
- LAN accessible (bind to 0.0.0.0)
- On-demand activation
- Serves static files and API endpoints

---

## Third-Party Dependencies

### Required Libraries

#### 1. INI File Parsing
**Library:** `gopkg.in/ini.v1`
**Version:** v1.67.0+
**Purpose:** Parse and write mt-canvus.ini files
**Rationale:** Mature library, handles Windows INI format well, actively maintained

#### 2. YAML Processing
**Library:** `gopkg.in/yaml.v3`
**Version:** v3.0.1+
**Purpose:** Parse and generate custom menu YAML files
**Rationale:** Standard Go YAML library, reliable, well-tested

#### 3. XML Generation
**Library:** `encoding/xml` (stdlib)
**Purpose:** Generate screen.xml files
**Rationale:** Built into Go, no dependency needed

#### 4. System Tray
**Library:** `github.com/getlantern/systray`
**Version:** Latest
**Purpose:** System tray integration (minimize to tray)
**Rationale:** Cross-platform, works well with Fyne, actively maintained

#### 5. GPU/Display Detection (Windows)
**Library:** TBD - May need `github.com/lxn/walk` or WMI calls
**Purpose:** Detect GPU outputs and resolutions for Screen.xml Creator
**Rationale:** Windows-specific functionality, requires native API access
**Status:** To be determined during implementation

### Optional/Development Libraries

#### Logging
**Library:** `github.com/sirupsen/logrus` or standard `log` package
**Purpose:** Structured logging with levels (DEBUG, INFO, WARN, ERROR)
**Rationale:** Better than standard log for development, can use stdlib for production

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
- ⚠️ No automatic updates
- ⚠️ No system integration (start menu, etc.)

**Alternatives Considered:**
- Windows installer (more professional, but complex)
- Both options (maintenance overhead)

---

## Project Structure

```
canvus-powertoys/
├── cmd/
│   └── powertoys/
│       └── main.go              # Application entry point
├── internal/
│   ├── atoms/                   # Basic building blocks (<100 lines)
│   │   ├── config/
│   │   │   ├── ini_parser.go
│   │   │   ├── xml_generator.go
│   │   │   └── yaml_handler.go
│   │   ├── validation/
│   │   │   └── validator.go
│   │   └── backup/
│   │       └── manager.go
│   ├── molecules/               # Composed components (<200 lines)
│   │   ├── screenxml/
│   │   │   ├── grid_widget.go
│   │   │   └── creator.go
│   │   ├── configeditor/
│   │   │   └── editor.go
│   │   ├── cssoptions/
│   │   │   └── manager.go
│   │   ├── custommenu/
│   │   │   └── designer.go
│   │   └── webui/
│   │       └── server.go
│   ├── organisms/               # Complex modules (<500 lines)
│   │   ├── app/
│   │   │   ├── main_window.go
│   │   │   └── tabs.go
│   │   └── services/
│   │       ├── file_service.go
│   │       └── validation_service.go
│   └── templates/              # Reusable patterns (<200 lines)
│       └── form_template.go
├── webui/                       # WebUI static files and routes
│   ├── static/
│   │   ├── css/
│   │   ├── js/
│   │   └── html/
│   └── routes/
│       └── handlers.go
├── assets/                      # Application assets
│   ├── icons/
│   └── templates/
├── tests/                       # Test files (mirror structure)
│   ├── atoms/
│   ├── molecules/
│   └── organisms/
├── docs/                        # Documentation
│   ├── PRD.md
│   ├── TECH_STACK.md
│   ├── TASKS.md
│   └── QnA.md
├── go.mod
├── go.sum
├── README.md
└── .gitignore
```

---

## Build & Distribution

### Development Build
```bash
go build -o canvus-powertoys ./cmd/powertoys
```

### Windows Release Build
```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o canvus-powertoys.exe ./cmd/powertoys
```

### Linux Release Build
```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o canvus-powertoys ./cmd/powertoys
```

### Build Flags
- `-s`: Omit symbol table and debug information
- `-w`: Omit DWARF symbol table
- Reduces binary size significantly

---

## Development Environment

### Required Tools
- **Go:** 1.21 or later
- **Git:** For version control
- **Code Editor:** VS Code, GoLand, or preferred editor with Go support

### Recommended Tools
- **golangci-lint:** Linting and code quality
- **goimports:** Automatic import formatting
- **Delve:** Debugger for Go

### Development Workflow
1. Develop on Linux
2. Build Windows executable: `GOOS=windows GOARCH=amd64 go build ...`
3. Deploy to Windows for testing
4. Iterate based on feedback

---

## Platform-Specific Considerations

### Windows
- File paths: `%APPDATA%`, `%ProgramData%`, `%LOCALAPPDATA%`
- GPU detection: WMI or DXGI APIs
- System tray: Native Windows API via systray library
- Logs: `%LOCALAPPDATA%\MultiTaction\Canvus\logs\`

### Ubuntu (Future)
- File paths: `~/MultiTaction/canvus/`, `/etc/MultiTaction/canvus/`
- GPU detection: xrandr or similar
- System tray: Native Linux API via systray library
- Logs: `~/.local/share/MultiTaction/Canvus/logs/` or similar

---

## Security Considerations

### Token Storage
- Store Canvus Server Private-Token securely
- Encrypt token in app configuration
- Display only last 6 digits when saved
- Never log full tokens

### File Operations
- Validate all file paths
- Check file permissions before operations
- Create backups before modifications
- Handle errors gracefully

### WebUI Security
- No password/access control (LAN only, trusted network)
- Token used only for Canvus Server API authentication
- Validate all inputs from WebUI
- Sanitize file paths

---

## Performance Considerations

### Startup Time
- Lazy load tabs (load on first access)
- Defer heavy operations
- Cache configuration file reads

### Memory Usage
- Release resources when tabs are closed
- Limit backup retention
- Efficient file I/O operations

### UI Responsiveness
- Run file operations in goroutines
- Update UI on main thread only
- Show progress indicators for long operations

---

## Future Considerations

### WebUI Migration
- Separate project step to analyze `canvus-webui`
- Define scope of features to port
- Determine integration approach
- Estimate effort and timeline

### Ubuntu Support
- Test on Ubuntu
- Verify file path handling
- Test GPU detection (xrandr)
- Validate system tray functionality

### Additional Features
- Auto-update mechanism (optional)
- Plugin system for extensions
- Advanced screen.xml configurations
- Multi-language support

---

**Document Status:** Ready for development
**Last Updated:** 2024

