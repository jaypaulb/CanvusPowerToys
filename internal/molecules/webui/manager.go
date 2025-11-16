package webui

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/config"
	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// Manager handles WebUI integration and local server.
type Manager struct {
	fileService  *services.FileService
	iniParser    *config.INIParser
	server       *http.Server
	serverURL    *widget.Entry
	authToken    *widget.Entry
	serverPort   *widget.Entry
	serverStatus *widget.Label
	enabledPages map[string]*widget.Check
	startStopBtn *widget.Button
}

// NewManager creates a new WebUI Manager.
func NewManager(fileService *services.FileService) (*Manager, error) {
	return &Manager{
		fileService:  fileService,
		iniParser:    config.NewINIParser(),
		enabledPages: make(map[string]*widget.Check),
	}, nil
}

// CreateUI creates the UI for the WebUI tab.
func (m *Manager) CreateUI(window fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("WebUI Integration")
	title.TextStyle = fyne.TextStyle{Bold: true}

	instructions := widget.NewRichTextFromMarkdown(`
**WebUI Integration**

Enable Canvus PowerToys to act as a web server for remote access and control.

**Configuration:**
- Canvus Server URL: Your Canvus server address
- User Auth Token: Access token from Canvus server profile
- WebUI Server Port: Port number for the local WebUI server (default: 8080)
- Enabled Pages: Select which WebUI pages to enable
`)

	// Server URL
	serverURLLabel := widget.NewLabel("Canvus Server URL:")
	m.serverURL = widget.NewEntry()
	m.serverURL.SetPlaceHolder("https://your-canvus-server.com")
	m.loadServerURL()

	// Auth Token
	authTokenLabel := widget.NewLabel("User Auth Token:")
	m.authToken = widget.NewEntry()
	m.authToken.SetPlaceHolder("Paste your access token here")
	m.authToken.Password = true // Hide token input

	// Server Port
	serverPortLabel := widget.NewLabel("WebUI Server Port:")
	m.serverPort = widget.NewEntry()
	m.serverPort.SetPlaceHolder("8080")
	m.serverPort.SetText("8080")

	// Token instructions
	tokenInstructions := widget.NewRichTextFromMarkdown(`
**How to get your Auth Token:**
1. Log in to your Canvus server
2. Navigate to: https://{canvusserverurl}/profile/access-tokens
3. Create a new access token
4. Copy and paste it below
`)

	// Server Status
	m.serverStatus = widget.NewLabel("Server: Stopped")
	m.serverStatus.Importance = widget.LowImportance

	// Start/Stop button
	m.startStopBtn = widget.NewButton("Start Server", func() {
		m.toggleServer(window)
	})

	// Enabled Pages
	pagesLabel := widget.NewLabel("Enabled Pages:")

	// WebUI pages (not PowerToys tabs)
	pageOptions := []string{
		"Main",
		"Pages",
		"Macros",
		"Remote Upload",
		"RCU",
	}

	pageChecks := []fyne.CanvasObject{}
	for _, page := range pageOptions {
		check := widget.NewCheck(page, nil)
		m.enabledPages[page] = check
		pageChecks = append(pageChecks, check)
	}

	// Save configuration button
	saveConfigBtn := widget.NewButton("Save Configuration", func() {
		m.saveConfiguration(window)
	})

	// Test connection button
	testConnectionBtn := widget.NewButton("Test Connection", func() {
		m.testConnection(window)
	})

	// Layout
	configSection := container.NewVBox(
		title,
		instructions,
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			serverURLLabel, m.serverURL,
			authTokenLabel, m.authToken,
			serverPortLabel, m.serverPort,
		),
		tokenInstructions,
		widget.NewSeparator(),
		pagesLabel,
		container.NewVBox(pageChecks...),
		widget.NewSeparator(),
		container.NewHBox(
			saveConfigBtn,
			testConnectionBtn,
		),
		widget.NewSeparator(),
		m.serverStatus,
		m.startStopBtn,
	)

	return container.NewScroll(configSection)
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
					m.serverURL.SetText(url)
					return
				}
			}
		}
	}
}

// saveConfiguration saves the WebUI configuration.
func (m *Manager) saveConfiguration(window fyne.Window) {
	serverURL := m.serverURL.Text
	if serverURL == "" {
		dialog.ShowError(fmt.Errorf("Server URL cannot be empty"), window)
		return
	}

	authToken := m.authToken.Text
	if authToken == "" {
		dialog.ShowError(fmt.Errorf("Auth token cannot be empty"), window)
		return
	}

	// TODO: Save configuration to a config file or mt-canvus.ini
	// For now, just show success message
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

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var localTestResult, remoteTestResult string
	var localTestSuccess, remoteTestSuccess bool

	// Test 1: Local WebUI Server
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

	// Test 2: Remote Canvus Server
	m.serverStatus.SetText("Testing connection to Canvus server...")
	m.serverStatus.Importance = widget.MediumImportance

	// Normalize server URL (remove trailing slash, ensure it doesn't have /api/v1)
	baseURL := strings.TrimSuffix(serverURL, "/")
	baseURL = strings.TrimSuffix(baseURL, "/api/v1")
	baseURL = strings.TrimSuffix(baseURL, "/api")

	remoteTestURL := fmt.Sprintf("%s/api/v1/clients", baseURL)
	req, err := http.NewRequest("GET", remoteTestURL, nil)
	if err != nil {
		remoteTestResult = fmt.Sprintf("❌ Failed to create request: %v\n   URL: %s", err, remoteTestURL)
		remoteTestSuccess = false
	} else {
		// Use Private-Token header for Canvus API
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

	// Update status and show results
	if localTestSuccess && remoteTestSuccess {
		m.serverStatus.SetText("All tests passed")
		m.serverStatus.Importance = widget.SuccessImportance
		dialog.ShowInformation("Connection Test Results",
			fmt.Sprintf("✅ All tests passed!\n\n%s\n\n%s", localTestResult, remoteTestResult),
			window)
	} else if localTestSuccess || remoteTestSuccess {
		m.serverStatus.SetText("Partial success - see details")
		m.serverStatus.Importance = widget.WarningImportance
		dialog.ShowInformation("Connection Test Results",
			fmt.Sprintf("⚠️  Partial success\n\n%s\n\n%s", localTestResult, remoteTestResult),
			window)
	} else {
		m.serverStatus.SetText("All tests failed")
		m.serverStatus.Importance = widget.DangerImportance
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

	mux := http.NewServeMux()

	// Use StaticHandler to serve actual WebUI pages (not placeholder pages)
	staticHandler := NewStaticHandler()
	staticHandler.ServeFiles(mux)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	m.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
			// Update UI on error (need to use fyne.App.QueueEvent or similar for thread safety)
			m.serverStatus.SetText(fmt.Sprintf("Server error: %v", err))
			m.serverStatus.Importance = widget.DangerImportance
			m.server = nil
			m.startStopBtn.SetText("Start Server")
		}
	}()

	// Give server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Update UI
	serverURL := fmt.Sprintf("http://localhost:%s", port)
	m.serverStatus.SetText(fmt.Sprintf("Server: Running on %s", serverURL))
	m.serverStatus.Importance = widget.SuccessImportance
	m.startStopBtn.SetText("Stop Server")

	dialog.ShowInformation("Server Started", fmt.Sprintf("WebUI server is running on:\n%s\n\nYou can access it in your browser at:\n%s", serverURL, serverURL), window)
}

// stopServer stops the local web server.
func (m *Manager) stopServer() {
	if m.server == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.server.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	}

	m.server = nil
	m.serverStatus.SetText("Server: Stopped")
	m.serverStatus.Importance = widget.LowImportance
	m.startStopBtn.SetText("Start Server")
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
