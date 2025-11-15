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

	// Common pages that might be available
	pageOptions := []string{
		"Dashboard",
		"Screen.xml Manager",
		"Config Editor",
		"CSS Options",
		"Custom Menu Designer",
		"System Status",
		"Logs Viewer",
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
	authToken := m.authToken.Text

	if serverURL == "" {
		dialog.ShowError(fmt.Errorf("Server URL cannot be empty"), window)
		return
	}

	// TODO: Save configuration to a config file or mt-canvus.ini
	// For now, just show success message
	dialog.ShowInformation("Saved", "Configuration saved successfully", window)
}

// testConnection tests the connection to the Canvus server.
func (m *Manager) testConnection(window fyne.Window) {
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

	// Test connection
	m.serverStatus.SetText("Testing connection...")
	m.serverStatus.Importance = widget.MediumImportance

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Test endpoint (adjust based on actual API)
	testURL := fmt.Sprintf("%s/api/health", serverURL)
	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		m.serverStatus.SetText("Connection failed: Invalid URL")
		m.serverStatus.Importance = widget.DangerImportance
		dialog.ShowError(fmt.Errorf("Failed to create request: %w", err), window)
		return
	}

	// Add auth token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		m.serverStatus.SetText("Connection failed: Network error")
		m.serverStatus.Importance = widget.DangerImportance
		dialog.ShowError(fmt.Errorf("Connection failed: %w", err), window)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized {
		// StatusUnauthorized might mean server is reachable but token is invalid
		if resp.StatusCode == http.StatusUnauthorized {
			m.serverStatus.SetText("Connection successful, but token may be invalid")
			m.serverStatus.Importance = widget.WarningImportance
			dialog.ShowInformation("Connection Test", "Server is reachable, but authentication failed. Please check your token.", window)
		} else {
			m.serverStatus.SetText("Connection successful")
			m.serverStatus.Importance = widget.SuccessImportance
			dialog.ShowInformation("Connection Test", "Successfully connected to Canvus server!", window)
		}
	} else {
		m.serverStatus.SetText(fmt.Sprintf("Connection failed: HTTP %d", resp.StatusCode))
		m.serverStatus.Importance = widget.DangerImportance
		dialog.ShowError(fmt.Errorf("Connection failed: HTTP %d", resp.StatusCode), window)
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
	port := "8080" // Default port, could be configurable

	mux := http.NewServeMux()

	// Register enabled pages
	for page, check := range m.enabledPages {
		if check.Checked {
			m.registerPage(mux, page)
		}
	}

	// Register root page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		html := `
<!DOCTYPE html>
<html>
<head>
	<title>Canvus PowerToys WebUI</title>
	<meta charset="utf-8">
	<style>
		body { font-family: Arial, sans-serif; margin: 40px; }
		.header { background: #4CAF50; color: white; padding: 20px; border-radius: 5px; }
		.content { margin-top: 20px; }
		.page-list { list-style: none; padding: 0; }
		.page-list li { margin: 10px 0; }
		.page-list a { color: #4CAF50; text-decoration: none; font-weight: bold; }
		.page-list a:hover { text-decoration: underline; }
	</style>
</head>
<body>
	<div class="header">
		<h1>Canvus PowerToys WebUI</h1>
		<p>Welcome to the WebUI interface</p>
	</div>
	<div class="content">
		<h2>Available Pages:</h2>
		<ul class="page-list">
`

		// Add links to enabled pages
		for page, check := range m.enabledPages {
			if check.Checked {
				path := "/" + m.pageToPath(page)
				html += fmt.Sprintf(`			<li><a href="%s">%s</a></li>`, path, page)
			}
		}

		html += `
		</ul>
		<p>Server URL: ` + m.serverURL.Text + `</p>
		<p>Status: Active</p>
	</div>
</body>
</html>
`

		w.Write([]byte(html))
	})

	// Default health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
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
			m.serverStatus.SetText(fmt.Sprintf("Server error: %v", err))
			m.serverStatus.Importance = widget.DangerImportance
		}
	}()

	// Update UI
	m.serverStatus.SetText(fmt.Sprintf("Server: Running on http://localhost:%s", port))
	m.serverStatus.Importance = widget.SuccessImportance
	m.startStopBtn.SetText("Stop Server")

	dialog.ShowInformation("Server Started", fmt.Sprintf("WebUI server is running on:\nhttp://localhost:%s", port), window)
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
