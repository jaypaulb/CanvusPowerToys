# Product Requirements Document (PRD)
## Canvus PowerToys

**Version:** 1.0
**Date:** 2024
**Status:** Draft

---

## 1. Problem Statement

AV technicians and Meeting Facilitators managing Multitaction Canvus installations in experience centres need a comprehensive tool to configure and manage Canvus systems efficiently. Currently, configuration requires manual editing of INI files, XML files, and YAML files, which is error-prone and time-consuming. There is no unified interface for managing screen configurations, Canvus settings, CSS customizations, custom menus, and remote access.

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
  - Remote management capabilities
- **Pain Points:** Complex technical configuration, need for IT support

### Use Cases

#### UC1: Configure Multi-Display Setup
**Actor:** AV Technician
**Goal:** Create screen.xml for 10x5 display wall
**Steps:**
1. Open Screen.xml Creator tab
2. See 10x5 grid representation
3. Assign GPU outputs to grid cells (using number drawing method)
4. Define touch area groups via click/drag
5. Generate screen.xml file
6. Optionally update mt-canvus.ini with video-output areas

#### UC2: Configure Canvus Settings
**Actor:** AV Technician or Meeting Facilitator
**Goal:** Modify mt-canvus.ini settings
**Steps:**
1. Open Canvus Config Editor tab
2. Search/filter for specific setting
3. View tooltip explanation
4. Modify setting value
5. Save to user or system config location
6. Backup created automatically

#### UC3: Enable Kiosk Mode
**Actor:** AV Technician
**Goal:** Enable kiosk mode with proper validation
**Steps:**
1. Open CSS Options Manager tab
2. Enable "Kiosk Mode"
3. App validates default-canvas is set (shows error if not)
4. App validates auto-pin=0 (shows error if not)
5. CSS plugin generated automatically
6. mt-canvus.ini updated with plugin-folders

#### UC4: Create Custom Menu
**Actor:** AV Technician
**Goal:** Create custom menu for specific deployment
**Steps:**
1. Open Custom Menu Designer tab
2. Create menu structure (hierarchical)
3. Add menu items with icons from minimal icon set
4. Configure actions (create notes, PDFs, videos, etc.)
5. Set positioning and layout
6. Save menu.yml and update mt-canvus.ini

#### UC5: Remote Management
**Actor:** Meeting Facilitator
**Goal:** Access Canvus management via web interface
**Steps:**
1. Enable WebUI in desktop app
2. Access WebUI from LAN device
3. Use enabled pages/features
4. Manage Canvus remotely

## 4. Functional Requirements

### FR1: Screen.xml Creator
- **FR1.1:** Display 10x5 grid representing video graphics array
- **FR1.2:** Support GPU output assignment via number drawing method
- **FR1.3:** Query OS/GPU for output resolution discovery
- **FR1.4:** Support click/drag for touch area assignment
- **FR1.5:** Support manual touch area index entry
- **FR1.6:** Visual feedback with color coding and pink frames for layout areas
- **FR1.7:** Generate screen.xml based on screen-xml-rules.md
- **FR1.8:** Offer to update mt-canvus.ini with video-output areas

### FR2: Canvus Config Editor
- **FR2.1:** Auto-detect mt-canvus.ini in standard locations
- **FR2.2:** Support manual file path override
- **FR2.3:** Searchable/filterable interface for all ini options
- **FR2.4:** Group options by section or functionality
- **FR2.5:** Hover tooltips for each option
- **FR2.6:** Save to user config (%APPDATA%\MultiTaction\canvus\mt-canvus.ini)
- **FR2.7:** Save to system config (%ProgramData%\MultiTaction\canvus\mt-canvus.ini)
- **FR2.8:** Create backup before saving (smart backup with rotation)

### FR3: CSS Options Manager
- **FR3.1:** Enable/disable rotation (temporary, reverts on canvas close)
- **FR3.2:** Enable/disable video looping (with memory warning)
- **FR3.3:** Enable kiosk mode with validation:
  - Validate default-canvas is set
  - Validate auto-pin=0
  - Hide UI layers appropriately
- **FR3.4:** Enable kiosk plus mode (kiosk + finger menu)
- **FR3.5:** Auto-generate CSS plugin files
- **FR3.6:** Update plugin-folders in mt-canvus.ini
- **FR3.7:** Handle plugin API version compatibility

### FR4: Custom Menu Designer
- **FR4.1:** Visual menu structure editor (hierarchical/tree view)
- **FR4.2:** Form-based item creation (not raw YAML editing)
- **FR4.3:** Icon picker with minimal icon set included
- **FR4.4:** Content file browser for PDFs, videos, images
- **FR4.5:** Coordinate system configuration (viewport vs canvas)
- **FR4.6:** Positioning tools (location, size, origin, offset)
- **FR4.7:** Import existing menu.yml files
- **FR4.8:** Generate YAML format matching CUSTOM_MENU_MANUAL.md
- **FR4.9:** Save menu.yml to same location as mt-canvus.ini
- **FR4.10:** Update custom-menu entry in mt-canvus.ini

### FR5: WebUI Server
- **FR5.1:** Integrated HTTP server in same Go process
- **FR5.2:** On-demand activation (only when enabled)
- **FR5.3:** LAN accessible (not just localhost)
- **FR5.4:** Show only enabled pages/functions
- **FR5.5:** Store Canvus Server Private-Token securely
- **FR5.6:** Display only last 6 digits of saved token
- **FR5.7:** Communicate with MT Canvus Server (online or LAN)
- **FR5.8:** No password/access control for WebUI itself
- **FR5.9:** Migration deferred - needs separate project step

### FR6: Application Core
- **FR6.1:** Tabbed GUI interface (one tab per feature)
- **FR6.2:** System tray integration (minimize to tray only)
- **FR6.3:** Auto-detect configuration files in standard locations
- **FR6.4:** Support Windows and Ubuntu file paths
- **FR6.5:** Form-based validation (prevent invalid configurations)
- **FR6.6:** Real-time validation feedback
- **FR6.7:** Save button disabled until all conflicts resolved
- **FR6.8:** Smart backup creation (only if file changed)
- **FR6.9:** Backup rotation (keep last N backups)
- **FR6.10:** Comprehensive logging in dev mode
- **FR6.11:** Simple error dialogs in production mode
- **FR6.12:** Logs stored in %LOCALAPPDATA%\MultiTaction\Canvus\logs\

## 5. Non-Functional Requirements

### NFR1: Performance
- Standard desktop application performance
- No specific startup time requirements
- No specific memory footprint constraints
- Responsive UI (no noticeable lag)

### NFR2: Reliability
- Zero data loss (backups before all saves)
- Graceful error handling
- No crashes from invalid user input (UI prevents invalid states)

### NFR3: Usability
- Intuitive interface for technical users
- Clear error messages
- Helpful tooltips and documentation
- Consistent UI patterns across tabs

### NFR4: Maintainability
- Single portable executable (no installer)
- Cross-platform codebase (Windows primary, Ubuntu future)
- Well-structured code following atomic design principles
- Comprehensive logging for troubleshooting

### NFR5: Security
- Secure token storage (encrypted)
- No secrets in logs
- Validate all inputs
- Safe file operations (backups, validation)

## 6. Constraints & Assumptions

### Constraints
- **Development Environment:** Linux (developing on Linux, deploying on Windows)
- **Testing:** Iterative approach (build → deploy → test → feedback)
- **WebUI Migration:** Deferred - needs separate project step to analyze canvus-webui
- **Icon Set:** Minimal icon set needed for Custom Menu Designer (to be determined)

### Assumptions
- Users have access to Canvus Server for WebUI token generation
- Users understand basic Canvus concepts (canvases, widgets, etc.)
- Standard Canvus installation paths are accessible
- Users have appropriate file system permissions for configuration files

## 7. Out of Scope

### Phase 1 (Initial Release)
- **WebUI Migration:** Full WebUI feature migration deferred to separate project
- **Auto-update:** No automatic update checking
- **Multi-language:** English only
- **Advanced Features:** Complex screen.xml configurations beyond 10x5 grid
- **Plugin Marketplace:** No plugin distribution system

### Future Considerations
- Ubuntu version (Windows primary for now)
- Advanced WebUI features (after migration analysis)
- Plugin version auto-update
- Additional screen.xml grid sizes/configurations

---

## 8. Dependencies

### External Dependencies
- Multitaction Canvus installation
- Canvus Server (for WebUI features)
- Standard Windows/Linux file system access

### Technical Dependencies
- Go runtime (compiled into executable)
- Fyne GUI framework
- Third-party libraries (see TECH_STACK.md)

---

## 9. Risks & Mitigation

### Risk 1: Cross-Platform Development
- **Risk:** Developing on Linux, deploying on Windows may cause issues
- **Mitigation:** Iterative testing cycle, early Windows builds, user feedback

### Risk 2: WebUI Migration Complexity
- **Risk:** Migrating WebUI from Node.js to Go may be complex
- **Mitigation:** Deferred to separate project step, thorough analysis first

### Risk 3: GPU/Display Detection
- **Risk:** Windows GPU detection may be challenging
- **Mitigation:** Use established libraries, fallback to manual entry

### Risk 4: Configuration File Compatibility
- **Risk:** Changes to Canvus may break configuration file formats
- **Mitigation:** Version validation, comprehensive testing, backup strategy

---

**Document Status:** Ready for development
**Next Steps:** See TASKS.md for implementation phases

