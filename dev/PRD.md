# Product Requirements Document (PRD)
## Canvus PowerToys

**Version:** 2.0
**Date:** 2024
**Status:** In Development - Core Features Complete, WebUI In Progress

---

## 1. Problem Statement

AV technicians and Meeting Facilitators managing Multitaction Canvus installations in experience centres need a comprehensive tool to configure and manage Canvus systems efficiently. Currently, configuration requires manual editing of INI files, XML files, and YAML files, which is error-prone and time-consuming. There is no unified interface for managing screen configurations, Canvus settings, CSS customizations, custom menus, and remote access.

---

## 2. Goals & Success Metrics

### Primary Goals
- Provide a unified desktop application for Canvus configuration management
- Simplify complex configuration tasks through intuitive GUI interfaces
- Reduce configuration errors through form-based validation
- Enable remote management capabilities via integrated WebUI
- Support both Windows (primary) and Ubuntu (future) platforms

### Success Metrics
- **Usability:** AV techs and facilitators can configure systems without manual file editing
- **Reliability:** Zero configuration errors due to UI validation
- **Efficiency:** Reduce configuration time by 50% compared to manual editing
- **Adoption:** Successful deployment in experience centre environments

### Current Implementation Status
- ✅ **Core Desktop Features:** Complete (Screen.xml Creator, Config Editor, CSS Options Manager, Custom Menu Designer)
- 🚧 **WebUI Integration:** In Progress (Foundation complete, pages implemented, testing pending)
- ⏳ **Production Polish:** Pending (Error handling, logging refinement, Windows testing)

---

## 3. User Personas & Use Cases

### Primary Persona: AV Technician
- **Role:** Technical staff managing Canvus installations
- **Needs:**
  - Configure screen.xml for multi-display setups
  - Manage mt-canvus.ini settings
  - Enable CSS customizations (kiosk modes, etc.)
  - Create custom menus for specific deployments
- **Pain Points:** Manual file editing, error-prone configuration, no visual feedback

### Secondary Persona: Meeting Facilitator
- **Role:** Operational staff using Canvus in experience centres
- **Needs:**
  - Quick access to configuration tools
  - Simple interfaces for common tasks
  - Remote management capabilities via WebUI
- **Pain Points:** Complex technical configuration, need for IT support

### Use Cases

#### UC1: Configure Multi-Display Setup
**Actor:** AV Technician
**Goal:** Create screen.xml for 10x5 display wall
**Steps:**
1. Open Screen.xml Creator tab
2. See 10x5 grid representation with visual feedback
3. Assign GPU outputs to grid cells (using number drawing method or click-to-assign)
4. Define touch area groups via click/drag or manual entry
5. Generate screen.xml file with validation
6. Optionally update mt-canvus.ini with video-output areas

**Status:** ✅ Implemented

#### UC2: Configure Canvus Settings
**Actor:** AV Technician or Meeting Facilitator
**Goal:** Modify mt-canvus.ini settings
**Steps:**
1. Open Canvus Config Editor tab
2. Search/filter for specific setting
3. View tooltip explanation with schema documentation
4. Modify setting value with real-time validation
5. Save to user or system config location
6. Backup created automatically with rotation

**Status:** ✅ Implemented

#### UC3: Enable Kiosk Mode
**Actor:** AV Technician
**Goal:** Enable kiosk mode with proper validation
**Steps:**
1. Open CSS Options Manager tab
2. Enable "Kiosk Mode"
3. App validates default-canvas is set (shows error if not)
4. App validates auto-pin=0 (shows error if not)
5. CSS plugin generated automatically with proper API version
6. mt-canvus.ini updated with plugin-folders

**Status:** ✅ Implemented

#### UC4: Create Custom Menu
**Actor:** AV Technician
**Goal:** Create custom menu for specific deployment
**Steps:**
1. Open Custom Menu Designer tab
2. Create menu structure (hierarchical tree view)
3. Add menu items with icons from included icon set (16 icons in 4 categories)
4. Configure actions (create notes, PDFs, videos, images, browser, open-folder)
5. Set positioning and layout (coordinate system, location, size, origin, offset)
6. Save menu.yml and update mt-canvus.ini

**Status:** ✅ Implemented

#### UC5: Remote Management via WebUI
**Actor:** Meeting Facilitator or AV Technician
**Goal:** Access Canvus management via web interface
**Steps:**
1. Enable WebUI in desktop app (WebUI Settings tab)
2. Configure Canvus Server URL and Private-Token
3. Access WebUI from LAN device (mobile or desktop browser)
4. View currently tracked canvas (automatic canvas tracking)
5. Use enabled pages/features:
   - Pages Management (create/manage canvas pages/zones)
   - Macros (move, copy, grouping, pinning)
   - Remote Content Upload (upload files to canvas)
   - RCU (Remote Content Upload interface)
6. Manage Canvus remotely with real-time canvas updates

**Status:** 🚧 In Progress (Implementation complete, testing pending)

---

## 4. Functional Requirements

### FR1: Screen.xml Creator
- **FR1.1:** Display 10x5 grid representing video graphics array ✅
- **FR1.2:** Support GPU output assignment via number drawing method ✅
- **FR1.3:** Query OS/GPU for output resolution discovery (placeholder with manual override) ✅
- **FR1.4:** Support click/drag for touch area assignment ✅
- **FR1.5:** Support manual touch area index entry ✅
- **FR1.6:** Visual feedback with color coding and pink frames for layout areas ✅
- **FR1.7:** Generate screen.xml based on screen-xml-rules.md ✅
- **FR1.8:** Offer to update mt-canvus.ini with video-output areas ✅

**Status:** ✅ Complete

### FR2: Canvus Config Editor
- **FR2.1:** Auto-detect mt-canvus.ini in standard locations ✅
- **FR2.2:** Support manual file path override ✅
- **FR2.3:** Searchable/filterable interface for all ini options ✅
- **FR2.4:** Group options by section or functionality ✅
- **FR2.5:** Hover tooltips for each option with embedded schema ✅
- **FR2.6:** Save to user config (%APPDATA%\MultiTaction\canvus\mt-canvus.ini) ✅
- **FR2.7:** Save to system config (%ProgramData%\MultiTaction\canvus\mt-canvus.ini) ✅
- **FR2.8:** Create backup before saving (smart backup with rotation) ✅

**Status:** ✅ Complete

### FR3: CSS Options Manager
- **FR3.1:** Enable/disable rotation (temporary, reverts on canvas close) ✅
- **FR3.2:** Enable/disable video looping (with memory warning) ✅
- **FR3.3:** Enable kiosk mode with validation:
  - Validate default-canvas is set ✅
  - Validate auto-pin=0 ✅
  - Hide UI layers appropriately ✅
- **FR3.4:** Enable kiosk plus mode (kiosk + finger menu) ✅
- **FR3.5:** Auto-generate CSS plugin files (.canvusplugin JSON + styles.css) ✅
- **FR3.6:** Update plugin-folders in mt-canvus.ini ✅
- **FR3.7:** Handle plugin API version compatibility ✅

**Status:** ✅ Complete

**Future Enhancements (Requires Plugin Development):**
- Video looping toggle button (requires C++ plugin)
- Hidden long-press button to toggle kiosk mode (requires C++ plugin)
- Dynamic CSS class toggling (requires plugin code)

### FR4: Custom Menu Designer
- **FR4.1:** Visual menu structure editor (hierarchical/tree view) ✅
- **FR4.2:** Form-based item creation (not raw YAML editing) ✅
- **FR4.3:** Icon picker with minimal icon set included (16 icons in 4 categories) ✅
- **FR4.4:** Content file browser for PDFs, videos, images ✅
- **FR4.5:** Coordinate system configuration (viewport vs canvas) ✅
- **FR4.6:** Positioning tools (location, size, origin, offset) ✅
- **FR4.7:** Import existing menu.yml files ✅
- **FR4.8:** Generate YAML format matching CUSTOM_MENU_MANUAL.md ✅
- **FR4.9:** Save menu.yml to same location as mt-canvus.ini ✅
- **FR4.10:** Update custom-menu entry in mt-canvus.ini ✅

**Status:** ✅ Complete

### FR5: WebUI Server
- **FR5.1:** Integrated HTTP server in same Go process ✅
- **FR5.2:** On-demand activation (only when enabled) ✅
- **FR5.3:** LAN accessible (bind to 0.0.0.0, default port 8080) ✅
- **FR5.4:** Show only enabled pages/functions ✅
- **FR5.5:** Store Canvus Server Private-Token securely ✅
- **FR5.6:** Display only last 6 digits of saved token ✅
- **FR5.7:** Communicate with MT Canvus Server (online or LAN) ✅
- **FR5.8:** No password/access control for WebUI itself (LAN only, trusted network) ✅
- **FR5.9:** Canvas tracking via ClientID/Workspace subscription ✅
- **FR5.10:** Real-time canvas updates via Server-Sent Events (SSE) ✅

**Status:** 🚧 In Progress (Implementation complete, testing pending)

#### FR5.11: WebUI Pages
- **FR5.11.1:** Main Page (Navigation Hub) ✅
  - Canvas header showing installation_name and canvas_name
  - Page descriptions and links
  - Connection status indicator
  - Real-time canvas updates via SSE
- **FR5.11.2:** Pages Management ✅
  - Create/manage canvas pages/zones
  - Atomic design components
  - Mobile-responsive interface
- **FR5.11.3:** Macros Management ✅
  - Manage tab: Move and copy operations (delete removed)
  - Grouping tab: Auto-grid, group by color, group by title
  - Pinning tab: Pin all, unpin all
  - Undelete tab removed (not working with media assets)
- **FR5.11.4:** Remote Content Upload ✅
  - File upload interface for admins
  - Progress indicators
  - Mobile file upload support
  - Improved error handling
- **FR5.11.5:** RCU (Remote Content Upload) ✅
  - Remote content upload interface
  - Atomic design components
  - Mobile-responsive

**Status:** ✅ Pages Implemented (Testing Pending)

#### FR5.12: WebUI Frontend Architecture
- **FR5.12.1:** Atomic Design Structure ✅
  - Atoms: button, input, card, badge, link components
  - Molecules: navbar, canvas-header, form-group, page-card, sse-client
  - Templates: page-template, modal-template
- **FR5.12.2:** Design System ✅
  - MultiTaction brand colors (CSS variables)
  - Typography and spacing
  - Dark theme support
  - Responsive breakpoints (mobile-first)
- **FR5.12.3:** Asset Optimization ✅
  - Minified CSS, JavaScript, and HTML
  - Embedded assets in binary
  - Size monitoring and optimization pipeline

**Status:** ✅ Complete

### FR6: Application Core
- **FR6.1:** Tabbed GUI interface (one tab per feature) ✅
- **FR6.2:** System tray integration (close button hides to tray) ✅
- **FR6.3:** Auto-detect configuration files in standard locations ✅
- **FR6.4:** Support Windows and Ubuntu file paths ✅
- **FR6.5:** Form-based validation (prevent invalid configurations) ✅
- **FR6.6:** Real-time validation feedback ✅
- **FR6.7:** Save button disabled until all conflicts resolved ✅
- **FR6.8:** Smart backup creation (only if file changed) ✅
- **FR6.9:** Backup rotation (keep last N backups) ✅
- **FR6.10:** Comprehensive logging in dev mode ✅
- **FR6.11:** Simple error dialogs in production mode ✅
- **FR6.12:** Logs stored in %LOCALAPPDATA%\MultiTaction\Canvus\logs\ ✅

**Status:** ✅ Complete

---

## 5. Non-Functional Requirements

### NFR1: Performance
- Standard desktop application performance ✅
- No specific startup time requirements
- No specific memory footprint constraints
- Responsive UI (no noticeable lag) ✅
- WebUI asset optimization (minified, embedded) ✅

### NFR2: Reliability
- Zero data loss (backups before all saves) ✅
- Graceful error handling ✅
- No crashes from invalid user input (UI prevents invalid states) ✅
- Canvas tracking reconnection logic ✅

### NFR3: Usability
- Intuitive interface for technical users ✅
- Clear error messages ✅
- Helpful tooltips and documentation ✅
- Consistent UI patterns across tabs ✅
- Mobile-responsive WebUI ✅
- Dark mode support (WebUI) ✅

### NFR4: Maintainability
- Single portable executable (no installer) ✅
- Cross-platform codebase (Windows primary, Ubuntu future) ✅
- Well-structured code following atomic design principles ✅
- Comprehensive logging for troubleshooting ✅
- Embedded documentation and schemas ✅

### NFR5: Security
- Secure token storage (encrypted) ✅
- No secrets in logs ✅
- Validate all inputs ✅
- Safe file operations (backups, validation) ✅
- LAN-only WebUI access (no password, trusted network) ✅

---

## 6. Constraints & Assumptions

### Constraints
- **Development Environment:** Linux (developing on Linux, deploying on Windows) ✅
- **Testing:** Iterative approach (build → deploy → test → feedback) ✅
- **WebUI Migration:** Completed - Full refactor and migration from Node.js to Go ✅
- **Icon Set:** Minimal icon set implemented (16 icons in 4 categories) ✅
- **Windows Testing:** Deferred - requires Windows build and testing environment ⏳

### Assumptions
- Users have access to Canvus Server for WebUI token generation ✅
- Users understand basic Canvus concepts (canvases, widgets, etc.)
- Standard Canvus installation paths are accessible ✅
- Users have appropriate file system permissions for configuration files ✅
- LAN network is trusted (no WebUI password required)

---

## 7. Out of Scope

### Phase 1 (Initial Release)
- **Auto-update:** No automatic update checking
- **Multi-language:** English only
- **Advanced Features:** Complex screen.xml configurations beyond 10x5 grid
- **Plugin Marketplace:** No plugin distribution system
- **WebUI Admin Page:** Not needed for PowerToys WebUI (skipped)

### Future Considerations
- Ubuntu version (Windows primary for now)
- Advanced screen.xml grid sizes/configurations
- Plugin version auto-update
- CSS Options Manager WebUI page (investigation pending)
- Advanced WebUI features (after user feedback)

---

## 8. Dependencies

### External Dependencies
- Multitaction Canvus installation
- Canvus Server (for WebUI features)
- Standard Windows/Linux file system access

### Technical Dependencies
- Go runtime (compiled into executable) - Go 1.21+
- Fyne GUI framework - v2.4+
- Third-party libraries:
  - `gopkg.in/ini.v1` - INI file parsing
  - `gopkg.in/yaml.v3` - YAML file processing
  - `github.com/getlantern/systray` - System tray integration
  - `encoding/xml` (stdlib) - XML generation
  - `tdewolff/minify` - Asset minification
  - See TECH_STACK.md for complete list

---

## 9. Risks & Mitigation

### Risk 1: Cross-Platform Development
- **Risk:** Developing on Linux, deploying on Windows may cause issues
- **Status:** ✅ Mitigated - Cross-compilation working, iterative testing cycle established
- **Mitigation:** Iterative testing cycle, early Windows builds, user feedback

### Risk 2: WebUI Migration Complexity
- **Risk:** Migrating WebUI from Node.js to Go may be complex
- **Status:** ✅ Resolved - Full refactor completed with improvements
- **Mitigation:** Complete refactor with atomic design, canvas tracking improvements, asset optimization

### Risk 3: GPU/Display Detection
- **Risk:** Windows GPU detection may be challenging
- **Status:** ⚠️ Partial - Placeholder with manual override
- **Mitigation:** Use established libraries, fallback to manual entry (implemented)

### Risk 4: Configuration File Compatibility
- **Risk:** Changes to Canvus may break configuration file formats
- **Status:** ✅ Mitigated - Validation and backup strategy implemented
- **Mitigation:** Version validation, comprehensive testing, backup strategy

### Risk 5: Canvas Tracking Reliability
- **Risk:** Canvas tracking may be unreliable or complex
- **Status:** ✅ Resolved - ClientID/Workspace subscription implemented
- **Mitigation:** Improved architecture with ClientID resolution and SSE-based workspace subscription

---

## 10. Architecture Highlights

### Desktop Application
- **Framework:** Go + Fyne GUI
- **Structure:** Atomic Design (atoms/molecules/organisms)
- **Deployment:** Single portable executable
- **Platforms:** Windows (primary), Ubuntu (future)

### WebUI Integration
- **Server:** Integrated HTTP server in Go process
- **Architecture:** Atomic Design (atoms/molecules/templates)
- **Canvas Tracking:** ClientID resolution + Workspace SSE subscription
- **Frontend:** HTML/CSS/JavaScript with embedded assets
- **Design:** MultiTaction brand colors, dark mode, mobile-responsive
- **Optimization:** Minified assets, embedded in binary

### Key Architectural Decisions
- **Atomic Design:** All code follows atoms → molecules → organisms hierarchy
- **Form-Based Validation:** No raw text editing, prevents configuration errors
- **Smart Backups:** Only backup if file changed, rotation system
- **On-Demand WebUI:** Only runs when enabled, resource efficient
- **Canvas Tracking:** Automatic via backend, no frontend UUID feedback loop

---

## 11. Implementation Status Summary

### Completed Features ✅
- **Phase 1:** Foundation & Setup (100%)
- **Phase 2:** Screen.xml Creator (100%)
- **Phase 3:** Canvus Config Editor (100%)
- **Phase 4:** CSS Options Manager (100%)
- **Phase 5:** Custom Menu Designer (100%)
- **Phase 6:** WebUI Foundation & Pages (95% - testing pending)

### In Progress 🚧
- **Phase 6:** WebUI Testing & Refinement
- **Phase 7:** Production Polish (Error handling, logging refinement)

### Pending ⏳
- **Phase 7:** Windows Integration Testing
- **Phase 7:** User Acceptance Testing
- **Phase 7:** Distribution & Release

---

## 12. Next Steps

### Immediate Priorities
1. Complete WebUI testing (canvas tracking, pages, macros, upload, RCU)
2. Windows build and integration testing
3. Error handling refinement
4. Logging optimization for production
5. User acceptance testing

### Future Enhancements
1. CSS Options Manager WebUI page (investigation)
2. Ubuntu version testing and refinement
3. Advanced screen.xml configurations
4. Plugin development for advanced CSS features
5. Auto-update mechanism (optional)

---

**Document Status:** Active Development
**Last Updated:** 2024
**Next Steps:** See TASKS.md for detailed implementation phases
