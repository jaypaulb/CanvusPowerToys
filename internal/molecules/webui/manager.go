package webui

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/config"
	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// Manager handles WebUI integration and local server.
type Manager struct {
	fileService       *services.FileService
	iniParser         *config.INIParser
	server            *http.Server
	serverURL         *widget.Entry
	serverSelect      *widget.Select
	authToken         *widget.Entry
	serverPort        *widget.Entry
	serverStatus      *widget.Label
	enabledPages      map[string]*widget.Check
	selectAllPage     *widget.Check
	suppressSelectAll bool
	startStopBtn      *widget.Button
	canvasService     *CanvasService
	apiRoutes         *APIRoutes
	tokenInstructions *fyne.Container
	tokenLinkButton    *widget.Button
}

type webUIConfiguration struct {
	ServerURL    string          `json:"server_url"`
	AuthToken    string          `json:"auth_token"`
	ServerPort   string          `json:"server_port"`
	EnabledPages map[string]bool `json:"enabled_pages"`
}

// NewManager creates a new WebUI Manager.
func NewManager(fileService *services.FileService) (*Manager, error) {
	return &Manager{
		fileService:       fileService,
		iniParser:         config.NewINIParser(),
		enabledPages:      make(map[string]*widget.Check),
		suppressSelectAll: false,
	}, nil
}

// CreateUI creates the UI for the WebUI tab.
func (m *Manager) CreateUI(window fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("WebUI Integration")
	title.TextStyle = fyne.TextStyle{Bold: true}

	savedConfig := m.loadSavedConfiguration()

	instructions := widget.NewRichTextFromMarkdown(`
**WebUI Integration**

Enable Canvus PowerToys to act as a web server for remote access and control.

**Configuration:**
- Canvus Server URL: Your Canvus server address
- User Auth Token: Access token from Canvus server profile
- WebUI Server Port: Port number for the local WebUI server (default: 8080)
- Enabled Pages: Select which WebUI pages to enable
`)

	// Server URL - Load server names from mt-canvus.ini
	serverURLLabel := widget.NewLabel("Canvus Server URL:")
	serverNames, serverURLs := m.loadServerNames()

	// Create select dropdown with server names
	m.serverSelect = widget.NewSelect(serverNames, func(selected string) {
		// When a server is selected, set the URL (ensure it has https://)
		if url, ok := serverURLs[selected]; ok {
			m.serverURL.SetText(ensureHTTPS(url))
			// Update token instructions with new URL
			m.updateTokenInstructions()
		}
	})
	m.serverSelect.PlaceHolder = "Select server or type URL..."

	// Create entry for typing custom URL
	m.serverURL = widget.NewEntry()
	m.serverURL.SetPlaceHolder("https://your-canvus-server.com or select from dropdown")
	m.loadServerURL()
	if savedConfig != nil && savedConfig.ServerURL != "" {
		m.serverURL.SetText(savedConfig.ServerURL)
	}
	// Update token instructions when URL changes
	m.serverURL.OnChanged = func(_ string) {
		m.updateTokenInstructions()
	}

	// If we have server names, set the first one as default
	if len(serverNames) > 0 {
		m.serverSelect.SetSelected(serverNames[0])
	}

	// Auth Token
	authTokenLabel := widget.NewLabel("User Auth Token:")
	m.authToken = widget.NewEntry()
	m.authToken.SetPlaceHolder("Paste your access token here")
	m.authToken.Password = true // Hide token input
	if savedConfig != nil && savedConfig.AuthToken != "" {
		m.authToken.SetText(savedConfig.AuthToken)
	}

	// Server Port
	serverPortLabel := widget.NewLabel("WebUI Server Port:")
	m.serverPort = widget.NewEntry()
	m.serverPort.SetPlaceHolder("8080")
	m.serverPort.SetText("8080")
	if savedConfig != nil && savedConfig.ServerPort != "" {
		m.serverPort.SetText(savedConfig.ServerPort)
	}

	// Token instructions (dynamic, updates when server URL changes)
	m.tokenInstructions = m.createTokenInstructions()

	// Server Status
	m.serverStatus = widget.NewLabel("Server: Stopped")
	m.serverStatus.Importance = widget.LowImportance

	// Start/Stop button
	m.startStopBtn = widget.NewButton("Start Server", func() {
		m.toggleServer(window)
	})

	// Enabled Pages
	// NOTE: Page selection functionality is not yet implemented. All pages are currently
	// enabled by default and cannot be adjusted by the user. The checkboxes appear enabled
	// (not visually greyed out) but callbacks immediately reset state to prevent user changes.
	// This provides visual feedback that all pages are available while clearly indicating
	// the feature is coming soon via the warning notice.
	pagesLabel := widget.NewLabel("Enabled Pages:")
	pagesLabel.Importance = widget.LowImportance // Grey out the label to indicate it's not functional

	// Coming soon note - orange warning style like custom menu
	comingSoonNote := widget.NewLabel("! Page selection is coming soon")
	comingSoonNote.Importance = widget.DangerImportance
	comingSoonNote.Wrapping = fyne.TextWrapWord

	// WebUI pages (not PowerToys tabs)
	pageOptions := []string{
		"Main",
		"Pages",
		"Macros",
		"Remote Upload",
		"RCU",
	}

	// Select All checkbox
	// NOTE: Checkbox appears enabled but callback immediately resets state to prevent user changes.
	// This maintains visual consistency (all checkboxes look enabled) while preventing state changes.
	m.selectAllPage = widget.NewCheck("Select All", func(checked bool) {
		// Immediately reset to checked: Page selection not yet implemented, prevent user from changing state
		// This ensures the checkbox always appears checked regardless of user clicks
		m.selectAllPage.SetChecked(true)
	})
	m.selectAllPage.SetChecked(true) // Set to checked - all pages are enabled by default

	pageChecks := []fyne.CanvasObject{m.selectAllPage, widget.NewSeparator()}

	// NOTE: All page checkboxes are set to checked (true) by default since page selection
	// is not yet implemented. Checkboxes appear enabled but callbacks immediately reset state
	// to prevent user interaction while maintaining visual consistency.
	for _, page := range pageOptions {
		check := widget.NewCheck(page, func(checked bool) {
			// Immediately reset to checked: Page selection not yet implemented, prevent user from changing state
			// This ensures the checkbox always appears checked regardless of user clicks
			// Note: Each iteration creates a new 'check' variable, so the closure correctly captures it
			check.SetChecked(true)
		})
		m.enabledPages[page] = check
		check.SetChecked(true) // All pages enabled by default until selection is implemented
		pageChecks = append(pageChecks, check)
	}
	// Sync select all state (will be true since all are checked)
	m.syncSelectAllFromChecks()

	// Save configuration button
	saveConfigBtn := widget.NewButton("Save Configuration", func() {
		m.saveConfiguration(window)
	})

	// Test connection button
	testConnectionBtn := widget.NewButton("Test Connection", func() {
		m.testConnection(window)
	})

	// Title bar with Start/Stop button
	titleBar := container.NewBorder(
		nil, nil,
		title,          // Left: Title
		m.startStopBtn, // Right: Start/Stop button
		nil,
	)

	// Left column: Configuration
	leftColumn := container.NewVBox(
		instructions,
		widget.NewSeparator(),
		// Form fields with labels and inputs side-by-side
		container.NewGridWithColumns(2,
			serverURLLabel, container.NewVBox(m.serverSelect, m.serverURL),
		),
		container.NewGridWithColumns(2,
			authTokenLabel, m.authToken,
		),
		container.NewGridWithColumns(2,
			serverPortLabel, m.serverPort,
		),
		m.tokenInstructions,
		widget.NewSeparator(),
		m.serverStatus,
	)

	// Right column: Enabled Pages with buttons
	// Notice at the top, then label, then checkboxes
	rightColumn := container.NewVBox(
		comingSoonNote,
		pagesLabel,
		container.NewVBox(pageChecks...),
		widget.NewSeparator(),
		container.NewVBox(
			saveConfigBtn,
			testConnectionBtn,
		),
	)

	// Two column layout (no individual scrolling - they scroll together)
	twoColumnLayout := container.NewGridWithColumns(2,
		leftColumn,
		rightColumn,
	)

	// Main layout: Title bar on top, two columns below
	mainLayout := container.NewVBox(
		titleBar,
		widget.NewSeparator(),
		twoColumnLayout,
	)

	// Single scroll container for the entire layout
	return container.NewScroll(mainLayout)
}

// ensureHTTPS ensures the URL has https:// prefix if it's missing a protocol.
func ensureHTTPS(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return url
	}
	// Check if URL already has a protocol
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	// Add https:// prefix
	return "https://" + url
}

// loadServerNames loads server names from [server:<name>] sections in mt-canvus.ini.
func (m *Manager) loadServerNames() ([]string, map[string]string) {
	serverNames := []string{}
	serverURLs := make(map[string]string)

	iniPath := m.fileService.DetectMtCanvusIni()
	if iniPath == "" {
		return serverNames, serverURLs
	}

	iniFile, err := m.iniParser.Read(iniPath)
	if err != nil {
		return serverNames, serverURLs
	}

	// Find all [server:<name>] sections
	for _, section := range iniFile.Sections() {
		sectionName := section.Name()
		if strings.HasPrefix(sectionName, "server:") {
			serverName := strings.TrimPrefix(sectionName, "server:")
			if serverName != "" {
				// Get server URL from this section
				serverKey := section.Key("server")
				if serverKey != nil {
					serverURL := serverKey.String()
					if serverURL != "" {
						serverURL = ensureHTTPS(serverURL)
						serverNames = append(serverNames, serverName)
						serverURLs[serverName] = serverURL
					}
				}
			}
		}
	}

	return serverNames, serverURLs
}

// loadServerURL loads the server URL from mt-canvus.ini.
func (m *Manager) loadServerURL() {
	iniPath := m.fileService.DetectMtCanvusIni()
	if iniPath == "" {
		return
	}

	iniFile, err := m.iniParser.Read(iniPath)
	if err != nil {
		return
	}

	// Try to get server URL from [server] or [canvas] section
	sections := []string{"server", "canvas", ""}
	for _, sectionName := range sections {
		section, err := iniFile.GetSection(sectionName)
		if err != nil {
			continue
		}

		// Try common server URL keys
		keys := []string{"server-url", "url", "canvus-server", "server"}
		for _, key := range keys {
			if keyObj := section.Key(key); keyObj != nil {
				url := keyObj.String()
				if url != "" {
					// Ensure URL has https:// prefix
					url = ensureHTTPS(url)
					m.serverURL.SetText(url)
					return
				}
			}
		}
	}
}

// saveConfiguration saves the WebUI configuration.
func (m *Manager) saveConfiguration(window fyne.Window) {
	if err := m.persistConfiguration(); err != nil {
		dialog.ShowError(err, window)
		return
	}

	dialog.ShowInformation("Saved", "Configuration saved successfully", window)
}

// testConnection tests both the local WebUI server and the remote Canvus server connection.
func (m *Manager) testConnection(window fyne.Window) {
	port := m.serverPort.Text
	if port == "" {
		port = "8080"
	}

	serverURL := m.serverURL.Text
	authToken := m.authToken.Text

	// Validate inputs
	if port == "" {
		dialog.ShowError(fmt.Errorf("Port cannot be empty"), window)
		return
	}

	if serverURL == "" {
		dialog.ShowError(fmt.Errorf("Canvus Server URL cannot be empty"), window)
		return
	}

	if authToken == "" {
		dialog.ShowError(fmt.Errorf("Auth token cannot be empty"), window)
		return
	}

	localTestResult, remoteTestResult, localTestSuccess, remoteTestSuccess := m.performConnectionTests(port, serverURL, authToken)
	m.updateStatusFromTestResults(localTestSuccess, remoteTestSuccess)

	// Show results
	if localTestSuccess && remoteTestSuccess {
		dialog.ShowInformation("Connection Test Results",
			fmt.Sprintf("✅ All tests passed!\n\n%s\n\n%s", localTestResult, remoteTestResult),
			window)
	} else if localTestSuccess || remoteTestSuccess {
		dialog.ShowInformation("Connection Test Results",
			fmt.Sprintf("⚠️  Partial success\n\n%s\n\n%s", localTestResult, remoteTestResult),
			window)
	} else {
		dialog.ShowError(fmt.Errorf("❌ All tests failed\n\n%s\n\n%s", localTestResult, remoteTestResult), window)
	}
}

// toggleServer starts or stops the local web server.
func (m *Manager) toggleServer(window fyne.Window) {
	if m.server == nil {
		// Start server
		m.startServer(window)
	} else {
		// Stop server
		m.stopServer()
	}
}

// startServer starts the local web server.
func (m *Manager) startServer(window fyne.Window) {
	port := m.serverPort.Text
	if port == "" {
		port = "8080"
	}

	// Validate port
	if port == "" {
		dialog.ShowError(fmt.Errorf("Port cannot be empty. Please enter a valid port number (e.g., 8080)"), window)
		return
	}

	// Check if port is already in use by checking if server is already running
	if m.server != nil {
		dialog.ShowError(fmt.Errorf("Server is already running. Please stop it first."), window)
		return
	}

	// Check if port is already in use
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Port %s is already in use. Please choose a different port or stop the process using it.", port), window)
		return
	}
	listener.Close() // Close immediately, we'll open it again in the server

	// Get server URL and auth token
	serverURL := m.serverURL.Text
	authToken := m.authToken.Text

	if serverURL == "" {
		dialog.ShowError(fmt.Errorf("Server URL cannot be empty"), window)
		return
	}

	if authToken == "" {
		dialog.ShowError(fmt.Errorf("Auth token cannot be empty"), window)
		return
	}

	if err := m.persistConfiguration(); err != nil {
		dialog.ShowError(fmt.Errorf("failed to save configuration: %w", err), window)
		return
	}

	// Normalize server URL
	apiBaseURL := strings.TrimSuffix(serverURL, "/")
	apiBaseURL = strings.TrimSuffix(apiBaseURL, "/api/v1")
	apiBaseURL = strings.TrimSuffix(apiBaseURL, "/api")

	// Create canvas service (may fail if auto-detection doesn't work, but that's OK)
	canvasService, err := NewCanvasService(m.fileService, apiBaseURL, authToken)
	if err != nil {
		// Show warning but continue - user can override client in WebUI
		dialog.ShowInformation("Canvas Service Warning",
			fmt.Sprintf("Canvas service auto-detection failed: %v\n\n"+
				"The WebUI will still load. You can manually override the client selection in the WebUI.", err), window)
		// Create a minimal canvas service that can be updated later
		// Must initialize canvasTracker to avoid nil pointer panics
		fmt.Printf("[WebUI] Creating minimal CanvasService with apiBaseURL: '%s', authToken length: %d\n", apiBaseURL, len(authToken))
		canvasTracker := webuiatoms.NewCanvasTracker()
		ctx, cancel := context.WithCancel(context.Background())
		canvasService = &CanvasService{
			apiBaseURL:       apiBaseURL,
			authToken:        authToken,
			installationName: "Unknown",
			canvasTracker:    canvasTracker,
			ctx:              ctx,
			cancel:           cancel,
		}
		fmt.Printf("[WebUI] Minimal CanvasService created with apiBaseURL: '%s'\n", canvasService.apiBaseURL)
	}

	// Create API client
	apiClient := webuiatoms.NewAPIClient(apiBaseURL, authToken)

	// Create API routes (uploadDir can be empty for now)
	apiRoutes := NewAPIRoutes(canvasService, apiClient, "")

	// Try to start canvas service, but don't fail if it doesn't work
	// User can override client selection in WebUI
	if err := canvasService.Start(); err != nil {
		fmt.Printf("[WebUI] Canvas service auto-start failed: %v (user can override in WebUI)\n", err)
		// Don't show error dialog - just log it and continue
		// The WebUI will load and user can manually override
	}

	// Store references
	m.canvasService = canvasService
	m.apiRoutes = apiRoutes

	mux := http.NewServeMux()

	// Register API routes first (before static handler)
	// This ensures more specific routes like /api/* take precedence over catch-all /
	apiRoutes.RegisterRoutes(mux)
	fmt.Printf("[WebUI] API routes registered\n")

	// Use StaticHandler to serve actual WebUI pages (not placeholder pages)
	staticHandler := NewStaticHandler()
	staticHandler.ServeFiles(mux)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Debug endpoint to list embedded files
	mux.HandleFunc("/debug/files", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		staticHandler := NewStaticHandler()

		var listFiles func(fs.FS, string, int) string
		listFiles = func(fsys fs.FS, path string, depth int) string {
			if depth > 5 {
				return ""
			}
			result := ""
			entries, err := fs.ReadDir(fsys, path)
			if err != nil {
				return fmt.Sprintf("Error reading %s: %v\n", path, err)
			}
			for _, entry := range entries {
				indent := strings.Repeat("  ", depth)
				result += fmt.Sprintf("%s%s\n", indent, entry.Name())
				if entry.IsDir() {
					subFS, _ := fs.Sub(fsys, path)
					result += listFiles(subFS, entry.Name(), depth+1)
				}
			}
			return result
		}

		fileList := "Embedded filesystem contents:\n"
		fileList += listFiles(staticHandler.fileSystem, ".", 0)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fileList))
	})

	m.server = &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
			// Don't update UI from goroutine - port check already happens before starting
			// If we get here, it's an unexpected error - just log it
			// The server will be nil and user can try starting again
			m.server = nil
		}
	}()

	// Give server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Update UI
	serverURLStr := fmt.Sprintf("http://localhost:%s", port)
	m.serverStatus.SetText(fmt.Sprintf("Server: Running on %s", serverURLStr))
	m.serverStatus.Importance = widget.SuccessImportance
	m.startStopBtn.SetText("Stop Server")

	localTestResult, remoteTestResult, localTestSuccess, remoteTestSuccess := m.performConnectionTests(port, serverURL, authToken)
	m.updateStatusFromTestResults(localTestSuccess, remoteTestSuccess)
	m.showServerStartedDialog(serverURLStr, localTestResult, remoteTestResult, window)
}

// stopServer stops the local web server.
func (m *Manager) stopServer() {
	if m.server == nil {
		return
	}

	// Stop canvas service first to stop workspace subscriptions
	if m.canvasService != nil {
		m.canvasService.Stop()
		m.canvasService = nil
	}

	// Reduced timeout to 5 seconds since SSE handler now checks context every 1 second
	// This should be sufficient for graceful shutdown of all connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.server.Shutdown(ctx); err != nil {
		// Log error but don't fail - server will still stop
		if err == context.DeadlineExceeded {
			fmt.Printf("Server shutdown: Some connections did not close within timeout, forcing close\n")
		} else {
			fmt.Printf("Server shutdown error: %v\n", err)
		}
		// Force close if graceful shutdown failed
		if m.server != nil {
			m.server.Close()
		}
	} else {
		fmt.Printf("Server shutdown: All connections closed gracefully\n")
	}

	m.server = nil
	m.apiRoutes = nil
	m.serverStatus.SetText("Server: Stopped")
	m.serverStatus.Importance = widget.LowImportance
	m.startStopBtn.SetText("Start Server")
}

func (m *Manager) performConnectionTests(port, serverURL, authToken string) (string, string, bool, bool) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var localTestResult, remoteTestResult string
	var localTestSuccess, remoteTestSuccess bool

	m.serverStatus.SetText("Testing local WebUI server...")
	m.serverStatus.Importance = widget.MediumImportance

	if m.server == nil {
		localTestResult = fmt.Sprintf("❌ Local WebUI server is not running\n   Please start the server first.")
		localTestSuccess = false
	} else {
		localTestURL := fmt.Sprintf("http://localhost:%s/health", port)
		req, err := http.NewRequest("GET", localTestURL, nil)
		if err != nil {
			localTestResult = fmt.Sprintf("❌ Failed to create request: %v\n   URL: %s", err, localTestURL)
			localTestSuccess = false
		} else {
			resp, err := client.Do(req)
			if err != nil {
				localTestResult = fmt.Sprintf("❌ Cannot connect to local server on port %s\n   Error: %v\n   URL: %s\n   Make sure the server is running and the port is correct.", port, err, localTestURL)
				localTestSuccess = false
			} else {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					localTestResult = fmt.Sprintf("✅ Local WebUI server responding\n   URL: %s\n   Status: OK", localTestURL)
					localTestSuccess = true
				} else {
					localTestResult = fmt.Sprintf("❌ Server returned HTTP %d\n   URL: %s\n   Server may be running but not responding correctly.", resp.StatusCode, localTestURL)
					localTestSuccess = false
				}
			}
		}
	}

	m.serverStatus.SetText("Testing connection to Canvus server...")
	m.serverStatus.Importance = widget.MediumImportance

	baseURL := strings.TrimSuffix(serverURL, "/")
	baseURL = strings.TrimSuffix(baseURL, "/api/v1")
	baseURL = strings.TrimSuffix(baseURL, "/api")

	remoteTestURL := fmt.Sprintf("%s/api/v1/clients", baseURL)
	req, err := http.NewRequest("GET", remoteTestURL, nil)
	if err != nil {
		remoteTestResult = fmt.Sprintf("❌ Failed to create request: %v\n   URL: %s", err, remoteTestURL)
		remoteTestSuccess = false
	} else {
		req.Header.Set("Private-Token", authToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			remoteTestResult = fmt.Sprintf("❌ Cannot connect to Canvus server\n   Error: %v\n   URL: %s\n   Check your server URL and network connection.", err, remoteTestURL)
			remoteTestSuccess = false
		} else {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				remoteTestResult = fmt.Sprintf("✅ Canvus server connection successful\n   URL: %s\n   Status: OK", remoteTestURL)
				remoteTestSuccess = true
			} else if resp.StatusCode == http.StatusUnauthorized {
				remoteTestResult = fmt.Sprintf("⚠️  Server reachable but authentication failed\n   URL: %s\n   HTTP Status: %d\n   Please check your auth token.", remoteTestURL, resp.StatusCode)
				remoteTestSuccess = false
			} else {
				remoteTestResult = fmt.Sprintf("❌ Server returned HTTP %d\n   URL: %s\n   Server may be reachable but endpoint not available.", resp.StatusCode, remoteTestURL)
				remoteTestSuccess = false
			}
		}
	}

	return localTestResult, remoteTestResult, localTestSuccess, remoteTestSuccess
}

func (m *Manager) updateStatusFromTestResults(localSuccess, remoteSuccess bool) {
	if localSuccess && remoteSuccess {
		m.serverStatus.SetText("All tests passed")
		m.serverStatus.Importance = widget.SuccessImportance
		return
	}

	if localSuccess || remoteSuccess {
		m.serverStatus.SetText("Partial success - see details")
		m.serverStatus.Importance = widget.WarningImportance
		return
	}

	m.serverStatus.SetText("All tests failed")
	m.serverStatus.Importance = widget.DangerImportance
}

func (m *Manager) showServerStartedDialog(serverURL string, localResult, remoteResult string, window fyne.Window) {
	// Create browser widget on canvas with specific size and position
	m.createBrowserWidgetOnCanvas(serverURL, 1024, 768, 4800, 2700)

	resultsLabel := widget.NewLabel(fmt.Sprintf("%s\n\n%s", localResult, remoteResult))
	resultsLabel.Wrapping = fyne.TextWrapWord

	noteLabel := widget.NewLabel("Note: Check the center of the canvas for the WebUI window")
	noteLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("WebUI server is running on %s", serverURL)),
		widget.NewSeparator(),
		widget.NewLabel("Connection Test Results:"),
		resultsLabel,
		widget.NewSeparator(),
		noteLabel,
	)

	dialog.NewCustomConfirm(
		"Server Started & Tested",
		"Close to System Tray",
		"Dismiss",
		content,
		func(closeToTray bool) {
			if closeToTray {
				window.Hide()
			}
		},
		window,
	).Show()
}

func (m *Manager) getWebUIConfigPath() string {
	if m.fileService == nil {
		return ""
	}
	return filepath.Join(m.fileService.GetUserConfigPath(), "CanvusPowerToys", "webui_config.json")
}

func (m *Manager) loadSavedConfiguration() *webUIConfiguration {
	configPath := m.getWebUIConfigPath()
	if configPath == "" {
		return nil
	}

	var cfg webUIConfiguration
	if err := m.fileService.ReadJSONFile(configPath, &cfg); err != nil {
		fmt.Printf("[WebUI] Failed to read saved configuration: %v\n", err)
		return nil
	}

	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	if cfg.EnabledPages == nil {
		cfg.EnabledPages = make(map[string]bool)
	}

	return &cfg
}

func (m *Manager) persistConfiguration() error {
	if m.serverURL == nil || m.authToken == nil || m.serverPort == nil {
		return fmt.Errorf("configuration inputs are not initialized")
	}

	serverURL := strings.TrimSpace(m.serverURL.Text)
	if serverURL == "" {
		return fmt.Errorf("Server URL cannot be empty")
	}

	authToken := strings.TrimSpace(m.authToken.Text)
	if authToken == "" {
		return fmt.Errorf("Auth token cannot be empty")
	}

	port := strings.TrimSpace(m.serverPort.Text)
	if port == "" {
		port = "8080"
	}

	if m.fileService == nil {
		return fmt.Errorf("file service not available")
	}

	configPath := m.getWebUIConfigPath()
	if configPath == "" {
		return fmt.Errorf("unable to determine configuration path")
	}

	cfg := &webUIConfiguration{
		ServerURL:    ensureHTTPS(serverURL),
		AuthToken:    authToken,
		ServerPort:   port,
		EnabledPages: make(map[string]bool),
	}

	for page, check := range m.enabledPages {
		if check != nil {
			cfg.EnabledPages[page] = check.Checked
		}
	}

	return m.fileService.WriteJSONFile(configPath, cfg)
}

func (m *Manager) syncSelectAllFromChecks() {
	if m.selectAllPage == nil {
		return
	}

	allSelected := true
	for _, check := range m.enabledPages {
		if check == nil || !check.Checked {
			allSelected = false
			break
		}
	}
	m.setSelectAllState(allSelected)
}

func (m *Manager) setSelectAllState(checked bool) {
	if m.selectAllPage == nil {
		return
	}
	m.suppressSelectAll = true
	m.selectAllPage.SetChecked(checked)
	m.suppressSelectAll = false
}

// registerPage registers a page handler.
func (m *Manager) registerPage(mux *http.ServeMux, page string) {
	// Convert page name to URL path
	path := "/" + m.pageToPath(page)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<title>%s - Canvus PowerToys</title>
	<meta charset="utf-8">
	<style>
		body { font-family: Arial, sans-serif; margin: 40px; }
		.header { background: #4CAF50; color: white; padding: 20px; border-radius: 5px; }
		.content { margin-top: 20px; }
	</style>
</head>
<body>
	<div class="header">
		<h1>%s</h1>
		<p>Canvus PowerToys WebUI</p>
	</div>
	<div class="content">
		<p>This is a placeholder page for: <strong>%s</strong></p>
		<p>Server URL: %s</p>
		<p>Status: Active</p>
	</div>
</body>
</html>
`, page, page, page, m.serverURL.Text)

		w.Write([]byte(html))
	})
}

// pageToPath converts a page name to a URL path.
func (m *Manager) pageToPath(page string) string {
	// Simple conversion: lowercase, replace spaces with hyphens
	path := strings.ToLower(page)
	path = strings.ReplaceAll(path, " ", "-")
	path = strings.ReplaceAll(path, ".", "")
	return path
}

// createTokenInstructions creates the token instructions widget with dynamic server URL.
func (m *Manager) createTokenInstructions() *fyne.Container {
	title := widget.NewRichTextFromMarkdown("**How to get your Auth Token:**")

	step1 := widget.NewLabel("1. Log in to your Canvus server")
	step2Label := widget.NewLabel("2. Navigate to:")
	step3 := widget.NewLabel("3. Create a new access token")
	step4 := widget.NewLabel("4. Copy and paste it below")

	// Create clickable link button (will be updated with actual URL)
	linkButton := widget.NewButton("", func() {
		serverURL := m.getServerURLForTokenLink()
		if serverURL != "" {
			m.openURL(serverURL)
		}
	})
	linkButton.Importance = widget.LowImportance

	// Container for step 2 with label and link
	step2Container := container.NewHBox(step2Label, linkButton)

	// Store link button reference for updates
	m.tokenLinkButton = linkButton

	// Update link initially
	m.updateTokenInstructions()

	return container.NewVBox(
		title,
		step1,
		step2Container,
		step3,
		step4,
	)
}

// updateTokenInstructions updates the token instructions link with the current server URL.
func (m *Manager) updateTokenInstructions() {
	if m.tokenLinkButton == nil {
		return
	}

	serverURL := m.getServerURLForTokenLink()
	if serverURL != "" {
		m.tokenLinkButton.SetText(serverURL)
		m.tokenLinkButton.OnTapped = func() {
			m.openURL(serverURL)
		}
	} else {
		m.tokenLinkButton.SetText("(Enter server URL above)")
		m.tokenLinkButton.OnTapped = nil
	}
}

// getServerURLForTokenLink gets the server URL for the token link, ensuring it has the profile path.
func (m *Manager) getServerURLForTokenLink() string {
	serverURL := strings.TrimSpace(m.serverURL.Text)
	if serverURL == "" {
		return ""
	}

	// Ensure URL has https:// prefix
	serverURL = ensureHTTPS(serverURL)

	// Remove trailing slash
	serverURL = strings.TrimSuffix(serverURL, "/")

	// Add profile/access-tokens path
	return serverURL + "/profile/access-tokens"
}

// openURL opens a URL in the default browser (cross-platform).
func (m *Manager) openURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // Linux and others
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Run() // Ignore errors - if browser doesn't open, user can copy URL
}

// createBrowserWidgetOnCanvas creates a browser widget on the Canvus canvas via MTCS API.
// Size: width x height, Position: x, y coordinates
func (m *Manager) createBrowserWidgetOnCanvas(url string, width, height, posX, posY int) {
	if m.canvasService == nil {
		fmt.Printf("[createBrowserWidgetOnCanvas] Canvas service not available\n")
		return
	}

	canvasID := m.canvasService.GetCanvasID()
	if canvasID == "" {
		fmt.Printf("[createBrowserWidgetOnCanvas] Canvas ID not available\n")
		return
	}

	// Get server URL and auth token from manager
	serverURL := m.serverURL.Text
	authToken := m.authToken.Text

	if serverURL == "" || authToken == "" {
		fmt.Printf("[createBrowserWidgetOnCanvas] Server URL or auth token not available\n")
		return
	}

	// Normalize server URL (same as in startServer)
	apiBaseURL := strings.TrimSuffix(serverURL, "/")
	apiBaseURL = strings.TrimSuffix(apiBaseURL, "/api/v1")
	apiBaseURL = strings.TrimSuffix(apiBaseURL, "/api")

	// Create API client
	apiClient := webuiatoms.NewAPIClient(apiBaseURL, authToken)

	// Create browser widget payload
	payload := map[string]interface{}{
		"widget_type": "Browser",
		"url":         url,
		"location": map[string]float64{
			"x": float64(posX),
			"y": float64(posY),
		},
		"size": map[string]float64{
			"width":  float64(width),
			"height": float64(height),
		},
	}

	// POST to /canvases/:id/browsers to create browser widget (widgets endpoint is read-only)
	endpoint := fmt.Sprintf("/api/v1/canvases/%s/browsers", canvasID)
	fmt.Printf("[createBrowserWidgetOnCanvas] Creating browser widget at (%d, %d) with size %dx%d, URL: %s\n", posX, posY, width, height, url)

	response, err := apiClient.Post(endpoint, payload)
	if err != nil {
		fmt.Printf("[createBrowserWidgetOnCanvas] ERROR: Failed to create browser widget: %v\n", err)
		return
	}

	fmt.Printf("[createBrowserWidgetOnCanvas] Successfully created browser widget: %s\n", string(response))
}
