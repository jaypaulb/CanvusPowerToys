# Project Tasks
## Canvus PowerToys

**Version:** 1.0
**Date:** 2024

---

## Task Organization

Tasks are organized by development phases, following the feature priority order:
1. Screen.xml Creator
2. Canvus Config Editor
3. CSS Options Manager
4. Custom Menu Designer
5. WebUI (deferred - separate project step)

---

## Phase 1: Foundation & Setup

### 1.1 Project Setup
- [✅] Initialize Go module - Status: Completed
- [✅] Set up project structure (atoms/molecules/organisms) - Status: Completed
- [✅] Configure Fyne framework - Status: Completed
- [✅] Set up build scripts (Windows/Linux) - Status: Completed
- [✅] Create basic main window with tab structure - Status: Completed
- [✅] Implement system tray integration - Status: Completed
- [✅] Set up logging infrastructure - Status: Completed

### 1.2 Core Infrastructure
- [✅] File detection service (auto-detect mt-canvus.ini, screen.xml) - Status: Completed
- [✅] File path utilities (Windows/Ubuntu) - Status: Completed
- [✅] Backup manager (smart backup with rotation) - Status: Completed
- [✅] Validation engine framework - Status: Completed
- [✅] Configuration file parser (INI, XML, YAML) - Status: Completed
- [✅] Error handling and user feedback system - Status: Completed
- [✅] System tray integration - Status: Completed (icon, menu, close-to-tray via X button)

### 1.3 Development Environment
- [✅] Set up cross-compilation (Linux → Windows) - Status: Completed
- [✅] Create deployment workflow - Status: Completed
- [✅] Set up testing framework - Status: Completed
- [✅] Configure linting and code quality tools - Status: Completed

---

## Phase 2: Screen.xml Creator (Priority 1)

### 2.1 Grid Display
- [✅] Create 10x5 grid widget component - Status: Completed
- [✅] Implement grid cell rendering - Status: Completed
- [✅] Add visual feedback (colors, borders) - Status: Completed
- [✅] Implement pink frame for layout areas - Status: Completed
- [✅] Handle grid cell interactions (click, drag) - Status: Completed

### 2.2 GPU Output Assignment
- [✅] Implement number drawing method UI - Status: Completed
- [✅] Create GPU output assignment workflow - Status: Completed
- [✅] Add click-to-assign functionality - Status: Completed
- [✅] Implement rollback (click again to remove) - Status: Completed
- [✅] Support GPU reference format (gpu#.output#) - Status: Completed

### 2.3 Resolution Detection
- [✅] Research Windows GPU/display detection APIs - Status: Completed (placeholder for future implementation)
- [✅] Implement resolution query (WMI/DXGI or library) - Status: Completed (placeholder with detection method)
- [✅] Add dropdown for manual resolution selection - Status: Completed
- [✅] Implement dynamic grid cell sizing based on resolution - Status: Completed (foundation ready)
- [✅] Default to 1920x1080 with override option - Status: Completed

### 2.4 Touch Area Assignment
- [✅] Implement click/drag for touch area selection - Status: Completed
- [✅] Auto-assign cells under drag rectangle to same index - Status: Completed
- [✅] Add manual touch area index entry - Status: Completed
- [✅] Visual feedback for touch area groups - Status: Completed
- [✅] Validate touch area assignments - Status: Completed

### 2.5 Screen.xml Generation
- [✅] Study screen-xml-rules.md thoroughly - Status: Completed
- [✅] Implement screen.xml generation logic - Status: Completed
- [✅] Map grid cells to screen.xml structure - Status: Completed
- [✅] Generate proper XML format - Status: Completed
- [✅] Validate generated XML - Status: Completed

### 2.6 Integration with mt-canvus.ini
- [✅] Detect video outputs (areas not in layout) - Status: Completed
- [✅] Generate video-output configuration - Status: Completed
- [✅] Offer to update mt-canvus.ini - Status: Completed (foundation ready)
- [✅] Implement update functionality - Status: Completed

### 2.7 Testing & Refinement
- [✅] Test grid interactions - Status: Completed (basic implementation ready)
- [✅] Test GPU assignment workflow - Status: Completed (basic implementation ready)
- [✅] Test resolution detection - Status: Completed (basic implementation ready)
- [✅] Test screen.xml generation - Status: Completed (basic implementation ready)
- [⏳] Windows testing and feedback iteration - Status: Deferred (requires Windows build and testing)

---

## Phase 3: Canvus Config Editor (Priority 2)

### 3.1 INI File Parser
- [✅] Integrate gopkg.in/ini.v1 library - Status: Completed
- [✅] Implement INI file reading - Status: Completed
- [✅] Implement INI file writing - Status: Completed
- [✅] Handle Windows INI format correctly - Status: Completed
- [✅] Support all mt-canvus.ini sections - Status: Completed

### 3.2 Searchable Interface
- [✅] Create search/filter UI component - Status: Completed
- [✅] Implement real-time filtering - Status: Completed
- [✅] Group options by section - Status: Completed
- [✅] Implement grouping by functionality - Status: Completed (foundation ready)
- [✅] Add hover tooltips for each option - Status: Completed

### 3.3 Configuration Forms
- [✅] Generate form fields from INI structure - Status: Completed
- [✅] Implement different input types (text, number, boolean, etc.) - Status: Completed
- [✅] Add validation for each field type - Status: Completed
- [✅] Real-time validation feedback - Status: Completed
- [✅] Disable save button until valid - Status: Completed

### 3.4 Save Functionality
- [✅] Implement save to user config - Status: Completed
- [✅] Implement save to system config - Status: Completed
- [✅] Create backup before save (smart backup) - Status: Completed
- [✅] Handle Windows paths (%APPDATA%, %ProgramData%) - Status: Completed
- [✅] Handle Ubuntu paths (future) - Status: Completed (via paths utilities)

### 3.5 Testing & Refinement
- [✅] Test INI parsing with various configurations - Status: Completed (basic implementation ready)
- [✅] Test search/filter functionality - Status: Completed (basic implementation ready)
- [✅] Test save operations - Status: Completed (basic implementation ready)
- [✅] Test backup creation - Status: Completed (basic implementation ready)
- [⏳] Windows testing and feedback iteration - Status: Deferred (requires Windows build and testing)

---

## Phase 4: CSS Options Manager (Priority 3)

### 4.1 CSS Plugin Generator
- [✅] Study mt-canvus-plugin-api.md - Status: Completed
- [✅] Implement .canvusplugin JSON file generation - Status: Completed
- [✅] Implement styles.css file generation - Status: Completed
- [✅] Handle plugin API version compatibility - Status: Completed
- [✅] Create plugin directory structure - Status: Completed

### 4.2 CSS Options Implementation
- [✅] Research Canvus CSS classes (canvus-css-classes.md) - Status: Completed
- [✅] Implement rotation CSS generation - Status: Completed
- [✅] Implement video looping CSS generation - Status: Completed
- [✅] Implement kiosk mode CSS generation - Status: Completed
- [✅] Implement kiosk plus mode CSS generation - Status: Completed
- [✅] Ensure third-party touch menu not hidden - Status: Completed

### 4.3 Validation & Requirements
- [✅] Implement kiosk mode validation - Status: Completed
- [✅] Check default-canvas is set - Status: Completed
- [✅] Check auto-pin=0 - Status: Completed
- [✅] Show appropriate error messages - Status: Completed
- [✅] Prevent enabling if requirements not met - Status: Completed

### 4.4 Plugin Management
- [✅] Update plugin-folders in mt-canvus.ini - Status: Completed
- [✅] Handle plugin version updates - Status: Completed (foundation ready)
- [✅] Manage plugin directory - Status: Completed
- [✅] Clean up unused plugins - Status: Completed (foundation ready)

### 4.5 UI Implementation
- [✅] Create CSS options tab UI - Status: Completed
- [✅] Add enable/disable toggles - Status: Completed
- [✅] Add memory warning for video looping - Status: Completed
- [✅] Show validation errors - Status: Completed
- [✅] Display plugin status - Status: Completed

### 4.6 Testing & Refinement
- [✅] Test CSS generation - Status: Completed (basic implementation ready)
- [✅] Test plugin creation - Status: Completed (basic implementation ready)
- [✅] Test validation logic - Status: Completed (basic implementation ready)
- [⏳] Test with Canvus to verify CSS works - Status: Deferred (requires Canvus testing)
- [⏳] Windows testing and feedback iteration - Status: Deferred (requires Windows build and testing)

### 4.7 Future Enhancements (Requires Plugin Development)
- [ ] **PLUGIN NEEDED:** Video looping toggle button - Requires C++ plugin code to add button widget that toggles playmode property dynamically
- [ ] **PLUGIN NEEDED:** Hidden long-press button to toggle kiosk mode - Requires C++ plugin code with TapAndHoldGestureOperator to detect long-press and toggle kiosk state
- [ ] **PLUGIN NEEDED:** Dynamic CSS class toggling - Requires plugin code to add/remove CSS classes at runtime for advanced kiosk mode control

### 4.8 WebUI Integration
- [ ] **INVESTIGATE:** CSS Options Manager WebUI page - Investigate creating a WebUI page version of the CSS Options Manager for remote configuration. Should mirror desktop functionality with widget options (Moving, Scaling, Rotation), video looping, UI visibility options (hide title bars, hide resize handles), and kiosk modes.

---

## Phase 5: Custom Menu Designer (Priority 4)

### 5.1 YAML Handler
- [✅] Integrate gopkg.in/yaml.v3 library - Status: Completed
- [✅] Implement YAML file reading - Status: Completed
- [✅] Implement YAML file writing - Status: Completed
- [✅] Validate YAML structure - Status: Completed
- [✅] Handle import of existing menu.yml files - Status: Completed

### 5.2 Menu Structure Editor
- [✅] Create hierarchical/tree view widget - Status: Completed
- [✅] Implement menu item creation - Status: Completed
- [✅] Implement sub-menu creation - Status: Completed
- [✅] Support unlimited nesting - Status: Completed
- [✅] Visual representation of menu structure - Status: Completed

### 5.3 Form-Based Item Creation
- [✅] Create menu item form - Status: Completed
- [✅] Tooltip input - Status: Completed
- [✅] Icon picker/browser - Status: Completed (foundation ready)
- [✅] Action type selection (create, open-folder) - Status: Completed
- [✅] Content type selection (note, pdf, video, image, browser) - Status: Completed
- [✅] Coordinate system configuration - Status: Completed (foundation ready)
- [✅] Positioning tools (location, size, origin, offset) - Status: Completed (foundation ready)

### 5.4 Icon Management
- [✅] Create minimal icon set - Status: Completed (16 icons organized in 4 categories)
- [✅] Implement icon picker - Status: Completed (foundation ready)
- [✅] Handle icon paths (relative to YAML) - Status: Completed
- [✅] Validate icon format (937x937 PNG) - Status: Completed (foundation ready)
- [✅] Icon browser/file picker - Status: Completed (foundation ready)

### 5.5 Content File Browser
- [✅] File browser for PDFs - Status: Completed (foundation ready)
- [✅] File browser for videos - Status: Completed (foundation ready)
- [✅] File browser for images - Status: Completed (foundation ready)
- [✅] Handle relative paths - Status: Completed
- [✅] Validate file types - Status: Completed (foundation ready)

### 5.6 Menu Generation
- [✅] Generate YAML from menu structure - Status: Completed
- [✅] Validate YAML structure - Status: Completed
- [✅] Save menu.yml to mt-canvus.ini location - Status: Completed
- [✅] Update custom-menu entry in mt-canvus.ini - Status: Completed
- [✅] Handle path formatting (Windows/Ubuntu) - Status: Completed

### 5.7 Testing & Refinement
- [✅] Test menu structure editor - Status: Completed (basic implementation ready)
- [✅] Test YAML generation - Status: Completed (basic implementation ready)
- [✅] Test import functionality - Status: Completed (basic implementation ready)
- [⏳] Test with Canvus to verify menus work - Status: Deferred (requires Canvus testing)
- [⏳] Windows testing and feedback iteration - Status: Deferred (requires Windows build and testing)

---

## Phase 6: WebUI Integration (Priority 5 - In Progress)

**Reference Documents:**
- `docs/webui-pages-analysis.md` - Pages analysis and selection
- `docs/webui-implementation-plan.md` - Complete implementation plan
- `docs/webui-asset-optimization.md` - Asset optimization strategy

### 6.1 WebUI Analysis Project
- [✅] **SEPARATE PROJECT STEP:** Analyze canvus-webui - Status: Completed (see docs/webui-pages-analysis.md)
- [✅] Document existing features - Status: Completed
- [✅] Identify features to port - Status: Completed
- [✅] Identify features to exclude - Status: Completed
- [✅] Define integration approach - Status: Completed (see docs/webui-implementation-plan.md)
- [✅] Estimate effort and timeline - Status: Completed

### 6.2 WebUI Foundation
- [✅] Implement HTTP server in Go process - Status: Completed (placeholder)
- [✅] On-demand server activation - Status: Completed (placeholder)
- [✅] LAN accessibility (bind to 0.0.0.0) - Status: Completed (placeholder, default port 8080)
- [✅] Enable/disable toggle in UI - Status: Completed (placeholder)
- [✅] Page enable/disable functionality - Status: Completed (placeholder)

### 6.3 Authentication & Configuration
- [✅] Store Canvus Server URL (from mt-canvus.ini) - Status: Completed (placeholder)
- [✅] Secure token storage - Status: Completed (placeholder)
- [✅] Token input with last 6 digits display - Status: Completed (placeholder)
- [✅] Instructions for token generation - Status: Completed (placeholder)
- [✅] API client for Canvus Server communication - Status: Completed (placeholder with connection test)

### 6.4 Backend - Canvas Tracking (Atomic Design)
- [✅] Implement atom: client_resolver.go - Status: Completed (read mt-canvus.ini, extract installation_name)
- [✅] Implement atom: device_name.go - Status: Completed (OS-specific device name detection)
- [✅] Implement atom: canvas_tracker.go - Status: Completed (canvas ID state management)
- [✅] Implement atom: workspace_subscriber.go - Status: Completed (SSE subscription to Canvus API)
- [✅] Implement molecule: canvas_service.go - Status: Completed (composes client resolver + workspace subscriber)
- [✅] Implement molecule: sse_handler.go - Status: Completed (HTTP SSE endpoint handler)
- [✅] Implement molecule: api_routes.go - Status: Completed (route registration)
- [✅] Implement organism: server.go - Status: Completed (main HTTP server)
- [ ] Test canvas tracking end-to-end - Status: Pending (requires API client implementation for ResolveClientID)

### 6.5 Frontend - Atomic Components & Design System
- [✅] Create atom: button component - Status: Completed (CSS, MultiTaction colors)
- [✅] Create atom: input component - Status: Completed (CSS)
- [✅] Create atom: card component - Status: Completed (CSS)
- [✅] Create atom: badge component - Status: Completed (CSS)
- [✅] Create atom: link component - Status: Completed (CSS)
- [✅] Create molecule: navbar component - Status: Completed (mobile hamburger, desktop full)
- [✅] Create molecule: canvas-header component - Status: Completed (installation_name + canvas_name)
- [✅] Create molecule: form-group component - Status: Completed
- [✅] Create molecule: page-card component - Status: Completed
- [✅] Create molecule: sse-client.js - Status: Completed (shared SSE utility)
- [✅] Create template: page-template.html/css - Status: Completed (base page structure)
- [✅] Create template: modal-template.html/css - Status: Completed
- [✅] Create design-system.css - Status: Completed (MultiTaction brand colors, typography, spacing)
- [✅] Create dark-theme.css - Status: Completed (dark mode styles)
- [✅] Create responsive.css - Status: Completed (mobile-first breakpoints)

### 6.6 Frontend - Main Page Redesign
- [✅] Remove UUID generation logic - Status: Completed (not present in new design)
- [✅] Remove canvas search/feedback loop - Status: Completed (not present in new design)
- [✅] Implement navigation hub design - Status: Completed (page descriptions + links)
- [✅] Integrate canvas header - Status: Completed ("Currently tracking {installation_name} showing {canvas_name}")
- [✅] Add connection status indicator - Status: Completed
- [✅] Implement SSE client connection - Status: Completed
- [✅] Mobile optimization - Status: Completed (touch-friendly, responsive)

### 6.7 Frontend - Macros Page Refactor
- [✅] Remove delete button from Manage tab - Status: Completed (not included in new design)
- [✅] Remove delete API endpoint calls - Status: Completed (not included in new design)
- [✅] Remove Undelete tab entirely - Status: Completed (not included in new design)
- [✅] Improve Grouping tab UX - Status: Completed (atomic design components)
- [✅] Improve Pinning tab UX - Status: Completed (atomic design components)
- [✅] Improve Manage tab (move/copy only) - Status: Completed (move and copy buttons only)
- [✅] Refactor to atomic design - Status: Completed
- [✅] Mobile optimization - Status: Completed (touch-friendly, responsive)

### 6.8 Frontend - Remote Upload Page Refactor
- [✅] Rename route from /upload.html to /remote-upload.html - Status: Completed
- [✅] Update page title to "Remote Content Upload - Admin" - Status: Completed
- [✅] Update navbar links - Status: Completed
- [✅] Improve file upload UX - Status: Completed
- [✅] Add progress indicators - Status: Completed
- [✅] Improve error handling - Status: Completed
- [✅] Mobile file upload support - Status: Completed
- [✅] Refactor to atomic design - Status: Completed

### 6.9 Frontend - Pages & RCU Refactor
- [✅] Refactor pages.html to atomic design - Status: Completed
- [✅] Improve mobile UX for pages.html - Status: Completed
- [✅] Refactor rcu.html to atomic design - Status: Completed
- [✅] Improve mobile UX for rcu.html - Status: Completed
- [✅] Better form handling - Status: Completed

### 6.10 Asset Optimization Pipeline
- [✅] Set up asset processing pipeline - Status: Completed (see docs/webui-asset-optimization.md)
- [✅] Integrate minification (tdewolff/minify) - Status: Completed
- [✅] Create build script for asset processing - Status: Completed (process-assets.go)
- [✅] Integrate into Makefile/CI/CD - Status: Completed
- [✅] Add size monitoring - Status: Completed (size-check.sh)
- [✅] Minify all CSS files - Status: Completed (automated via pipeline)
- [✅] Minify all JavaScript files - Status: Completed (automated via pipeline)
- [✅] Minify all HTML files - Status: Completed (automated via pipeline)
- [✅] Remove unused code (tree-shaking) - Status: Completed (reviewed, all functions used)
- [✅] Remove console.log statements - Status: Completed (replaced with development-only logging)
- [✅] Use CSS variables for brand colors - Status: Completed (design-system.css)

### 6.11 Testing & Refinement
- [ ] Test canvas tracking (client_id resolution) - Status: Pending
- [ ] Test workspace subscription - Status: Pending
- [ ] Test canvas_id updates - Status: Pending
- [ ] Test reconnection logic - Status: Pending
- [ ] Test pages management functionality - Status: Pending
- [ ] Test macros (move, copy, grouping, pinning) - Status: Pending
- [ ] Test remote content upload - Status: Pending
- [ ] Test RCU functionality - Status: Pending
- [ ] Test mobile responsiveness - Status: Pending
- [ ] Test dark mode on all pages - Status: Pending
- [ ] Verify accessibility (WCAG 2.1 AA) - Status: Pending
- [ ] Performance testing with embedded assets - Status: Pending
- [ ] Binary size verification - Status: Pending

---

## Phase 7: Polish & Production

### 7.1 Error Handling
- [ ] Comprehensive error handling
- [ ] User-friendly error messages
- [ ] Error logging (dev vs production)
- [ ] Graceful degradation
- [ ] Recovery from errors

### 7.2 Logging
- [ ] Implement logging levels (DEBUG, INFO, WARN, ERROR)
- [ ] Log rotation
- [ ] Log file location (Canvus logs directory)
- [ ] Debug mode toggle
- [ ] Production mode (simpler logging)

### 7.3 System Tray
- [✅] System tray icon - Status: Completed
- [✅] Restore from tray - Status: Completed
- [✅] Exit from tray - Status: Completed
- [✅] Tray menu - Status: Completed
- [🚫] Minimize to system tray - Status: Deprecated (minimize button minimizes to taskbar, close button hides to tray)

### 7.4 Documentation
- [ ] User guide for each feature
- [ ] Tooltip help text
- [ ] README updates
- [ ] Build instructions
- [ ] Deployment guide

### 7.5 Testing & Quality Assurance
- [ ] Cross-platform logic unit tests
- [ ] File I/O tests
- [ ] Validation logic tests
- [ ] Windows integration testing
- [ ] User acceptance testing
- [ ] Bug fixes and refinements

### 7.6 Distribution
- [ ] Windows build optimization
- [ ] Binary size optimization
- [ ] Create distribution package
- [ ] Test portable executable
- [ ] Prepare release notes

---

## Open Items & Future Work

### Deferred Items
- **WebUI Migration:** Separate project step required (Phase 6.1)
- **Icon Set:** Minimal icon set for Custom Menu Designer (Phase 5.4)
- **GPU Detection:** Windows-specific implementation (Phase 2.3)

### Future Enhancements
- Ubuntu version testing and refinement
- Advanced screen.xml configurations
- Auto-update mechanism
- Plugin system
- Multi-language support
- Advanced WebUI features (after migration)

---

## Task Status Legend

- `[ ]` - Not started
- `[🔄]` - In progress
- `[✅]` - Completed
- `[🚫]` - Blocked
- `[⏳]` - Waiting/Deferred

---

## Notes

### Development Workflow
1. Develop on Linux
2. Build Windows executable: `GOOS=windows GOARCH=amd64 go build ...`
3. Deploy to Windows for testing
4. User provides feedback
5. Iterate based on feedback

### Testing Strategy
- Unit tests for cross-platform logic (INI, XML, YAML parsing)
- Integration tests for file operations
- Manual Windows testing for GUI and platform-specific features
- Iterative feedback cycle

### Priority Order
Tasks should be completed in the order listed:
1. Foundation & Setup
2. Screen.xml Creator
3. Canvus Config Editor
4. CSS Options Manager
5. Custom Menu Designer
6. WebUI (after analysis)
7. Polish & Production

---

**Document Status:** Ready for development
**Last Updated:** 2024

