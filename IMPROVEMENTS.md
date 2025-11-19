# Recent Improvements - Canvus PowerToys

This document summarizes the improvements made to address CSS Options and WebUI issues.

## 1. CSS Options Preview - Improved UX

### Problem
The "Preview CSS" dialog was showing a cramped 1-line text box, making it very difficult to read the generated CSS content.

### Solution
- Upgraded the preview dialog to use a proper scrollable container with minimum size: **750×400px**
- Increased dialog size to **900×700px** for better visibility
- Made the text entry read-only (by ignoring onChange) so users can select/copy but not edit
- Added more spacious layout for comfortable reading of CSS code

### Impact
Users can now easily preview CSS content before applying it, improving the feature's usability.

---

## 2. Generate Plugin Button - Marked as Coming Soon

### Problem
The "Generate Plugin" button was active and clickable, but the feature wasn't fully implemented, which could confuse users.

### Solution
- Changed button text to **"Generate Plugin (Coming Soon)"**
- Disabled the button using `Disable()`
- Replaced the callback with a user-friendly info dialog explaining alternatives:
  - Preview CSS using Preview CSS button
  - Manually create plugin files following Canvus documentation
  - Use "Launch Canvus with Current Config" for temporary testing
- Added comprehensive TODO comments explaining what needs to be implemented:
  - Create .canvusplugin manifest
  - Generate styles.css
  - Update mt-canvus.ini plugin-folders
  - Add documentation for plugin format
  - Update UI to show plugin status

### Impact
Clearer UX that prevents users from expecting a non-working feature while still showing what the full feature will do.

---

## 3. WebUI Server Logging - File-Based Logging System

### Problem
When WebUI server encounters errors (connection issues, canvas tracking failures, etc.), the error messages only print to stdout via `fmt.Printf()`. Users running the GUI app without a console window cannot see these logs, making debugging impossible.

### Solution
Implemented a comprehensive file-based logging system:

#### New Component: `internal/atoms/webui/logger.go`
- **WebUILogger** struct for managing WebUI-specific logging
- Logs to **both console and file** (dual output)
- Log file path: `%LOCALAPPDATA%\MultiTaction\Canvus\logs\webui\webui_YYYY-MM-DD_HH-MM-SS.log`
- Automatic log rotation: keeps logs organized by timestamp
- **Automatic cleanup**: removes log files older than 7 days to prevent disk space issues
- Methods:
  - `NewWebUILogger()` - Creates logger and initializes log directory
  - `Log(message string)` - Log formatted message
  - `Logf(format string, args ...)` - Log with format string
  - `GetLogPath()` - Retrieve current log file path
  - `Close()` - Close log file handle
  - `CleanupOldLogs()` - Clean up old logs (called during server initialization)

#### Updated Components:
1. **`internal/organisms/webui/server.go`**
   - Initialized WebUILogger during server creation
   - Added logging to `Start()` method:
     - Server initialization
     - Canvas service startup
     - HTTP server startup
     - Error conditions
   - Added logging to `Stop()` method:
     - Server shutdown
     - Graceful shutdown progress
     - Connection closing status
   - Added `GetLogPath()` method to retrieve log file location

#### Log Locations
**Windows:**
- `C:\Users\<username>\AppData\Local\MultiTaction\Canvus\logs\webui\webui_YYYY-MM-DD_HH-MM-SS.log`

**Linux:**
- `~/.local/share/MultiTaction/Canvus/logs/webui/webui_YYYY-MM-DD_HH-MM-SS.log`

#### Usage
When WebUI server starts, logs are automatically written to the timestamped log file. Users can:
1. Check logs for debugging when WebUI isn't working
2. Share logs with support for issue resolution
3. Find canvas tracking errors, connection issues, API failures, etc.

### Impact
- **Enables debugging** without console access
- **Reduces support friction** - users can capture detailed logs
- **Automatic cleanup** prevents disk space issues
- **Dual output** maintains console visibility for developers

---

## 4. WebUI Canvas Tracking - Architecture Documentation

### Issue Analysis
Users reported that WebUI wasn't displaying the active canvas name in the top bar.

### Root Causes Identified
The issue has multiple potential causes:

1. **Canvas Event Reception**: WorkspaceSubscriber may not be receiving SSE events from Canvus API
2. **Canvas Name Fetching**: If canvas_id exists but canvas_name is empty, CanvasService attempts async fetch but returns empty string immediately
3. **API Connectivity**: Canvus Server API may be unreachable or auth token invalid
4. **Client ID Resolution**: Installation name to client ID resolution may fail

### Architecture (for future debugging)

```
┌─────────────────────────────────────────────────────────────────┐
│ Canvas Tracking Data Flow                                       │
└─────────────────────────────────────────────────────────────────┘

1. CanvasService.Start()
   ├─ Resolves client_id from installation_name (via ClientResolver)
   ├─ Creates WorkspaceSubscriber
   └─ Subscribes to workspace updates (SSE stream)

2. WorkspaceSubscriber receives canvas events
   ├─ Event: canvas_id changes
   ├─ Event: workspace updates
   └─ Updates CanvasTracker with canvas_id and canvas_name

3. CanvasTracker stores current state
   └─ Methods: GetCanvasID(), GetCanvasName(), UpdateCanvas()

4. SSEHandler.HandleSubscribe() - WebUI Server-Sent Events
   ├─ Initial send: GetCanvasID() + GetCanvasName()
   ├─ Polls every 1 second
   ├─ Sends canvas_update event if values changed
   └─ Sends keepalive comment (: keepalive\n) to maintain connection

5. Frontend (common.js) receives SSE updates
   ├─ Connected to /api/sse endpoint
   ├─ Listens for canvas_update events
   ├─ Updates navbar: navbarCanvasName.textContent
   └─ Updates sessionStorage: canvasName

6. Fallback: /api/installation/info endpoint
   ├─ Called on page load
   ├─ Returns current state from CanvasService
   ├─ Used if SSE connection has issues
   └─ Updates navbar: navbarCanvasName.textContent
```

### Debugging WebUI Issues

With the new logging system, users and developers can now:

1. **Enable WebUI server** in the PowerToys GUI
2. **Wait 5-10 seconds** for initialization
3. **Check WebUI log file** at:
   - Windows: `%LOCALAPPDATA%\MultiTaction\Canvus\logs\webui\webui_*.log`
   - Linux: `~/.local/share/MultiTaction/Canvus/logs/webui/webui_*.log`
4. **Look for errors in log file:**
   - `[Server.Start]` - Server initialization issues
   - `[CanvasService]` - Canvas tracking/connection issues
   - `ERROR` - Any error conditions
   - `Failed to resolve client_id` - Installation name not found
   - `Failed to fetch canvas` - API connectivity issue

### Recommended Next Steps

For production debugging of canvas display issues:

1. Check WebUI logs (now available)
2. Verify Canvus Server URL and auth token are correct
3. Verify network connectivity between PowerToys and Canvus Server
4. Check if installation_name from screen.xml matches a client on Canvus Server
5. Consider implementing manual client override in WebUI (already partially implemented)

---

## Testing the Changes

### Build
```bash
make build              # Build for current platform
make build-windows      # Build for Windows
```

### Test CSS Preview
1. Open CSS Options Manager tab
2. Enable any CSS option
3. Click "Preview CSS"
4. Verify the dialog shows a large, readable text area (not a 1-line box)

### Test Generate Plugin Button
1. Open CSS Options Manager tab
2. Look at "Generate Plugin" button
3. Verify:
   - Button text says "(Coming Soon)"
   - Button is greyed out/disabled
   - Clicking shows info dialog with alternatives

### Test WebUI Logging
1. Open WebUI Settings tab
2. Enter valid Canvus Server URL and auth token
3. Click "Start Server"
4. Check log file:
   - Windows: `%LOCALAPPDATA%\MultiTaction\Canvus\logs\webui\`
   - Linux: `~/.local/share/MultiTaction/Canvus/logs/webui/`
5. Verify log file contains startup messages and any errors

---

## Files Modified

1. **internal/molecules/cssoptions/manager.go**
   - Updated `previewCSS()` method with larger dialog and better UX
   - Updated `generateBtn` initialization with "Coming Soon" message and disabled state
   - Added comprehensive TODO comment

2. **internal/atoms/webui/logger.go** (NEW)
   - Complete WebUI logging implementation
   - File-based logging with automatic cleanup

3. **internal/organisms/webui/server.go**
   - Added WebUILogger field to Server struct
   - Updated `NewServer()` to initialize logger
   - Updated `Start()` with logging
   - Updated `Stop()` with logging
   - Added `GetLogPath()` method

---

## Known Limitations & Future Improvements

### CSS Options
- Generate Plugin feature is marked as coming soon
- Implementation needed for plugin manifest generation
- Documentation needed for plugin format

### WebUI Logging
- Logs are currently console + file (development + debug)
- Could add log level filtering (INFO, WARN, ERROR) in future
- Could add structured logging (JSON) for parsing

### Canvas Tracking
- Manual client override partially implemented
- Could add auto-discovery of active canvas from Canvus API
- Could add canvas history tracking

---

---

## 5. Version Display in Application UI

### Problem
Users couldn't easily identify which version of the application they were running without checking file properties.

### Solution
- Added version number display in the top-right corner of the main window
- Shows current version (e.g., "v1.0.74")
- Version updates automatically with each build

### Impact
Users can immediately identify which version they're running, making it easier to verify updates and report issues.

---

## 6. Application Control Buttons

### Problem
- No way to quit the application from within the UI - had to use system tray menu
- Closing window minimized to tray, but users couldn't easily find the tray icon

### Solution
Added two buttons in the top-right corner:
- **"Minimize to Tray"** button - hides window, app stays running in system tray
- **"Quit"** button (red) - exits application completely

### Impact
Users now have explicit, visible controls for both minimizing and quitting, improving usability.

---

## 7. Full Connection Test Flow

### Problem
The "Test Connection" button only checked if the API was reachable, not if the full canvas tracking pipeline would work.

### Solution
Implemented a comprehensive 4-step connection test:

1. **Step 1: Get clients list** - Fetches all clients from `/api/v1/clients`
2. **Step 2: Find client** - Matches installation_name from mt-canvus.ini or device name
3. **Step 3: Get workspace** - Fetches `/api/v1/clients/{id}/workspaces/0` to get canvas ID
4. **Step 4: Get canvas** - Fetches `/api/v1/canvases/{id}` to get canvas name

On success, displays:
```
Status: Ready to start server using ClientName : CanvasName
```

### Impact
- Users get immediate feedback that their configuration will work
- Shows exactly which client and canvas will be used
- Lists available clients if match fails (helps debugging)

---

## 8. Self-Signed Certificate Support

### Problem
Users connecting to Canvus servers with self-signed SSL certificates got x509 certificate errors.

### Solution
Added TLS certificate verification skip for all HTTP clients:
- `internal/atoms/webui/api_client.go` - Main API client
- `internal/atoms/webui/client_resolver.go` - Client ID resolution
- `internal/molecules/webui/manager.go` - Connection tests

Uses:
```go
TLSClientConfig: &tls.Config{
    InsecureSkipVerify: true,
}
```

### Impact
- Supports internal/development Canvus servers with self-signed certs
- Works with expired certificates
- Works with certificates from untrusted CAs
- Essential for enterprise deployments using internal PKI

---

## 9. Version Management System

### Problem
Build versions stayed the same (e.g., always 1.0.66), making it impossible to identify which binary contained which changes.

### Solution
Created an automated version increment system:

- **`scripts/increment-version.sh`** - Auto-increments patch version (1.0.66 → 1.0.67)
- **`VERSION_WORKFLOW.md`** - Complete workflow guide
- **`BUILD_WITH_VERSION.md`** - Quick reference

Usage:
```bash
./scripts/increment-version.sh && make build-windows
```

### Impact
- Each build gets a unique version number
- Binary filename shows version (e.g., canvus-powertoys.1.0.74.exe)
- Easy to identify which binary has the latest changes

---

## Summary

These improvements make the application more user-friendly and debuggable:

1. ✅ **CSS Preview** - Much better UX with readable content
2. ✅ **Generate Plugin** - Clear status and user expectations
3. ✅ **WebUI Logging** - Users can now debug WebUI issues without console access
4. 📝 **Canvas Tracking** - Architecture documented for future debugging
5. ✅ **Version Display** - Version visible in application UI
6. ✅ **Control Buttons** - Explicit Minimize to Tray and Quit buttons
7. ✅ **Full Connection Test** - 4-step API flow test with ClientName : CanvasName status
8. ✅ **Self-Signed Certs** - Support for servers with self-signed SSL certificates
9. ✅ **Version Management** - Automated version increment system for builds

All changes maintain backward compatibility and don't require configuration changes.
