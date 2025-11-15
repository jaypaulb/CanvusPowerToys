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
- [ ] Implement click/drag for touch area selection
- [ ] Auto-assign cells under drag rectangle to same index
- [ ] Add manual touch area index entry
- [ ] Visual feedback for touch area groups
- [ ] Validate touch area assignments

### 2.5 Screen.xml Generation
- [ ] Study screen-xml-rules.md thoroughly
- [ ] Implement screen.xml generation logic
- [ ] Map grid cells to screen.xml structure
- [ ] Generate proper XML format
- [ ] Validate generated XML

### 2.6 Integration with mt-canvus.ini
- [ ] Detect video outputs (areas not in layout)
- [ ] Generate video-output configuration
- [ ] Offer to update mt-canvus.ini
- [ ] Implement update functionality

### 2.7 Testing & Refinement
- [ ] Test grid interactions
- [ ] Test GPU assignment workflow
- [ ] Test resolution detection
- [ ] Test screen.xml generation
- [ ] Windows testing and feedback iteration

---

## Phase 3: Canvus Config Editor (Priority 2)

### 3.1 INI File Parser
- [ ] Integrate gopkg.in/ini.v1 library
- [ ] Implement INI file reading
- [ ] Implement INI file writing
- [ ] Handle Windows INI format correctly
- [ ] Support all mt-canvus.ini sections

### 3.2 Searchable Interface
- [ ] Create search/filter UI component
- [ ] Implement real-time filtering
- [ ] Group options by section
- [ ] Implement grouping by functionality
- [ ] Add hover tooltips for each option

### 3.3 Configuration Forms
- [ ] Generate form fields from INI structure
- [ ] Implement different input types (text, number, boolean, etc.)
- [ ] Add validation for each field type
- [ ] Real-time validation feedback
- [ ] Disable save button until valid

### 3.4 Save Functionality
- [ ] Implement save to user config
- [ ] Implement save to system config
- [ ] Create backup before save (smart backup)
- [ ] Handle Windows paths (%APPDATA%, %ProgramData%)
- [ ] Handle Ubuntu paths (future)

### 3.5 Testing & Refinement
- [ ] Test INI parsing with various configurations
- [ ] Test search/filter functionality
- [ ] Test save operations
- [ ] Test backup creation
- [ ] Windows testing and feedback iteration

---

## Phase 4: CSS Options Manager (Priority 3)

### 4.1 CSS Plugin Generator
- [ ] Study mt-canvus-plugin-api.md
- [ ] Implement .canvusplugin JSON file generation
- [ ] Implement styles.css file generation
- [ ] Handle plugin API version compatibility
- [ ] Create plugin directory structure

### 4.2 CSS Options Implementation
- [ ] Research Canvus CSS classes (canvus-css-classes.md)
- [ ] Implement rotation CSS generation
- [ ] Implement video looping CSS generation
- [ ] Implement kiosk mode CSS generation
- [ ] Implement kiosk plus mode CSS generation
- [ ] Ensure third-party touch menu not hidden

### 4.3 Validation & Requirements
- [ ] Implement kiosk mode validation
- [ ] Check default-canvas is set
- [ ] Check auto-pin=0
- [ ] Show appropriate error messages
- [ ] Prevent enabling if requirements not met

### 4.4 Plugin Management
- [ ] Update plugin-folders in mt-canvus.ini
- [ ] Handle plugin version updates
- [ ] Manage plugin directory
- [ ] Clean up unused plugins

### 4.5 UI Implementation
- [ ] Create CSS options tab UI
- [ ] Add enable/disable toggles
- [ ] Add memory warning for video looping
- [ ] Show validation errors
- [ ] Display plugin status

### 4.6 Testing & Refinement
- [ ] Test CSS generation
- [ ] Test plugin creation
- [ ] Test validation logic
- [ ] Test with Canvus to verify CSS works
- [ ] Windows testing and feedback iteration

---

## Phase 5: Custom Menu Designer (Priority 4)

### 5.1 YAML Handler
- [ ] Integrate gopkg.in/yaml.v3 library
- [ ] Implement YAML file reading
- [ ] Implement YAML file writing
- [ ] Validate YAML structure
- [ ] Handle import of existing menu.yml files

### 5.2 Menu Structure Editor
- [ ] Create hierarchical/tree view widget
- [ ] Implement menu item creation
- [ ] Implement sub-menu creation
- [ ] Support unlimited nesting
- [ ] Visual representation of menu structure

### 5.3 Form-Based Item Creation
- [ ] Create menu item form
- [ ] Tooltip input
- [ ] Icon picker/browser
- [ ] Action type selection (create, open-folder)
- [ ] Content type selection (note, pdf, video, image, browser)
- [ ] Coordinate system configuration
- [ ] Positioning tools (location, size, origin, offset)

### 5.4 Icon Management
- [ ] Create minimal icon set
- [ ] Implement icon picker
- [ ] Handle icon paths (relative to YAML)
- [ ] Validate icon format (937x937 PNG)
- [ ] Icon browser/file picker

### 5.5 Content File Browser
- [ ] File browser for PDFs
- [ ] File browser for videos
- [ ] File browser for images
- [ ] Handle relative paths
- [ ] Validate file types

### 5.6 Menu Generation
- [ ] Generate YAML from menu structure
- [ ] Validate YAML structure
- [ ] Save menu.yml to mt-canvus.ini location
- [ ] Update custom-menu entry in mt-canvus.ini
- [ ] Handle path formatting (Windows/Ubuntu)

### 5.7 Testing & Refinement
- [ ] Test menu structure editor
- [ ] Test YAML generation
- [ ] Test import functionality
- [ ] Test with Canvus to verify menus work
- [ ] Windows testing and feedback iteration

---

## Phase 6: WebUI Integration (Priority 5 - Deferred)

### 6.1 WebUI Analysis Project
- [ ] **SEPARATE PROJECT STEP:** Analyze canvus-webui
- [ ] Document existing features
- [ ] Identify features to port
- [ ] Identify features to exclude
- [ ] Define integration approach
- [ ] Estimate effort and timeline

### 6.2 WebUI Foundation (After Analysis)
- [ ] Implement HTTP server in Go process
- [ ] On-demand server activation
- [ ] LAN accessibility (bind to 0.0.0.0)
- [ ] Enable/disable toggle in UI
- [ ] Page enable/disable functionality

### 6.3 Authentication & Configuration
- [ ] Store Canvus Server URL (from mt-canvus.ini)
- [ ] Secure token storage
- [ ] Token input with last 6 digits display
- [ ] Instructions for token generation
- [ ] API client for Canvus Server communication

### 6.4 WebUI Pages (TBD After Analysis)
- [ ] Define pages to implement
- [ ] Implement page routing
- [ ] Create page templates
- [ ] Implement API endpoints
- [ ] Static file serving

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
- [ ] Minimize to system tray only
- [ ] System tray icon
- [ ] Restore from tray
- [ ] Exit from tray
- [ ] Tray menu (if needed)

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

