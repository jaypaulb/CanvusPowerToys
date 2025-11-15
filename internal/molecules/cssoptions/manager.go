package cssoptions

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/ini.v1"

	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/backup"
	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/config"
	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// PluginManifest represents a .canvusplugin manifest file.
type PluginManifest struct {
	APIVersion  string `json:"api-version"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

// Manager handles CSS options and plugin generation.
type Manager struct {
	iniParser            *config.INIParser
	fileService          *services.FileService
	backupManager        *backup.Manager
	movingEnabled        *widget.Check
	scalingEnabled       *widget.Check
	rotationEnabled      *widget.Check
	videoLoopEnabled     *widget.Check
	kioskModeEnabled     *widget.Check
	kioskPlusEnabled     *widget.Check
	hideTitleBarsEnabled *widget.Check
	hideResizeHandlesEnabled *widget.Check
	statusLabel          *widget.Label
}

// NewManager creates a new CSS Options Manager.
func NewManager(fileService *services.FileService) (*Manager, error) {
	return &Manager{
		iniParser:     config.NewINIParser(),
		fileService:   fileService,
		backupManager: backup.NewManager(""),
	}, nil
}

// CreateUI creates the UI for the CSS Options Manager tab.
func (m *Manager) CreateUI(window fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("CSS Options Manager")
	title.TextStyle = fyne.TextStyle{Bold: true}

	instructions := widget.NewRichTextFromMarkdown(`
**CSS Options Manager**

Enable CSS-based features for Canvus. These options create plugins that modify Canvus behavior.

**Requirements:**
- default-canvas must be set in mt-canvus.ini
- auto-pin must be 0 for kiosk modes
`)

	// Widget Options section
	widgetOptionsLabel := widget.NewLabel("Widget Options")
	widgetOptionsLabel.TextStyle = fyne.TextStyle{Bold: true}
	widgetOptionsTooltip := widget.NewLabel("Control widget interaction capabilities (temporary, reverts on canvas close)")

	movingLabel := widget.NewLabel("Enable Moving")
	movingTooltip := widget.NewLabel("Allow users to move canvas items")
	m.movingEnabled = widget.NewCheck("", nil)

	scalingLabel := widget.NewLabel("Enable Scaling")
	scalingTooltip := widget.NewLabel("Allow users to resize canvas items")
	m.scalingEnabled = widget.NewCheck("", nil)

	rotationLabel := widget.NewLabel("Enable Rotation")
	rotationTooltip := widget.NewLabel("Allow users to rotate canvas items")
	m.rotationEnabled = widget.NewCheck("", nil)

	// Video Loop option
	videoLoopLabel := widget.NewLabel("Enable Video Looping")
	videoLoopTooltip := widget.NewLabel("Warning: May use significant memory with large videos")
	videoLoopWarning := widget.NewLabel("⚠ Memory Warning: Large videos may consume significant resources")
	videoLoopWarning.Importance = widget.WarningImportance
	m.videoLoopEnabled = widget.NewCheck("", nil)

	// UI Visibility Options
	uiVisibilityLabel := widget.NewLabel("UI Visibility Options")
	uiVisibilityLabel.TextStyle = fyne.TextStyle{Bold: true}

	hideTitleBarsLabel := widget.NewLabel("Hide Title Bars")
	hideTitleBarsTooltip := widget.NewLabel("Hide title bars on canvas items")
	m.hideTitleBarsEnabled = widget.NewCheck("", nil)

	hideResizeHandlesLabel := widget.NewLabel("Hide Resize Handles")
	hideResizeHandlesTooltip := widget.NewLabel("Hide resize handles on canvas items")
	m.hideResizeHandlesEnabled = widget.NewCheck("", nil)

	// Kiosk Mode option
	kioskModeLabel := widget.NewLabel("Enable Kiosk Mode")
	kioskModeTooltip := widget.NewLabel("Requires: default-canvas set, auto-pin=0")
	m.kioskModeEnabled = widget.NewCheck("", nil)

	// Kiosk Plus Mode option
	kioskPlusLabel := widget.NewLabel("Enable Kiosk Plus Mode")
	kioskPlusTooltip := widget.NewLabel("Requires: default-canvas set, auto-pin=0")
	m.kioskPlusEnabled = widget.NewCheck("", nil)

	// Status label
	m.statusLabel = widget.NewLabel("Ready")

	// Buttons
	generateBtn := widget.NewButton("Generate Plugin", func() {
		m.generatePlugin(window)
	})

	validateBtn := widget.NewButton("Validate Requirements", func() {
		m.validateRequirements(window)
	})

	launchBtn := widget.NewButton("Launch Canvus with Current Config", func() {
		m.launchCanvusWithConfig(window)
	})

	form := container.NewVBox(
		title,
		instructions,
		widget.NewSeparator(),
		widgetOptionsLabel,
		widgetOptionsTooltip,
		container.NewGridWithColumns(2,
			movingLabel, m.movingEnabled,
			movingTooltip, widget.NewLabel(""),
			scalingLabel, m.scalingEnabled,
			scalingTooltip, widget.NewLabel(""),
			rotationLabel, m.rotationEnabled,
			rotationTooltip, widget.NewLabel(""),
		),
		widget.NewSeparator(),
		videoLoopLabel, m.videoLoopEnabled,
		videoLoopTooltip,
		videoLoopWarning,
		widget.NewSeparator(),
		uiVisibilityLabel,
		container.NewGridWithColumns(2,
			hideTitleBarsLabel, m.hideTitleBarsEnabled,
			hideTitleBarsTooltip, widget.NewLabel(""),
			hideResizeHandlesLabel, m.hideResizeHandlesEnabled,
			hideResizeHandlesTooltip, widget.NewLabel(""),
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			kioskModeLabel, m.kioskModeEnabled,
			kioskModeTooltip, widget.NewLabel(""),
			kioskPlusLabel, m.kioskPlusEnabled,
			kioskPlusTooltip, widget.NewLabel(""),
		),
		widget.NewSeparator(),
		m.statusLabel,
		container.NewHBox(
			validateBtn,
			generateBtn,
			launchBtn,
		),
	)

	return container.NewScroll(form)
}

// validateRequirements validates that requirements are met for enabled options.
func (m *Manager) validateRequirements(window fyne.Window) {
	iniPath := m.fileService.DetectMtCanvusIni()
	if iniPath == "" {
		dialog.ShowError(fmt.Errorf("mt-canvus.ini not found"), window)
		return
	}

	iniFile, err := m.iniParser.Read(iniPath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to read mt-canvus.ini: %w", err), window)
		return
	}

	var errors []string
	var warnings []string

	// Check default-canvas
	defaultCanvas := ""
	if section, err := iniFile.GetSection(""); err == nil {
		if key := section.Key("default-canvas"); key != nil {
			defaultCanvas = key.String()
		}
	}
	if defaultCanvas == "" {
		errors = append(errors, "default-canvas is not set")
	}

	// Check auto-pin
	autoPin := ""
	if section, err := iniFile.GetSection(""); err == nil {
		if key := section.Key("auto-pin"); key != nil {
			autoPin = key.String()
		}
	}

	// Validate kiosk mode requirements
	if m.kioskModeEnabled.Checked || m.kioskPlusEnabled.Checked {
		if defaultCanvas == "" {
			errors = append(errors, "default-canvas must be set for kiosk modes")
		}
		if autoPin != "0" {
			errors = append(errors, "auto-pin must be 0 for kiosk modes")
		}
	}

	if len(errors) > 0 {
		m.statusLabel.SetText(fmt.Sprintf("Validation failed: %s", errors[0]))
		dialog.ShowError(fmt.Errorf("Validation failed:\n%s", fmt.Sprintf("%s", errors[0])), window)
		return
	}

	if len(warnings) > 0 {
		m.statusLabel.SetText(fmt.Sprintf("Warnings: %s", warnings[0]))
	} else {
		m.statusLabel.SetText("All requirements met")
		dialog.ShowInformation("Valid", "All requirements are met", window)
	}
}

// generatePlugin generates the CSS plugin files.
func (m *Manager) generatePlugin(window fyne.Window) {
	// Validate first
	iniPath := m.fileService.DetectMtCanvusIni()
	if iniPath == "" {
		dialog.ShowError(fmt.Errorf("mt-canvus.ini not found"), window)
		return
	}

	iniFile, err := m.iniParser.Read(iniPath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to read mt-canvus.ini: %w", err), window)
		return
	}

	// Check requirements
	if m.kioskModeEnabled.Checked || m.kioskPlusEnabled.Checked {
		defaultCanvas := ""
		if section, err := iniFile.GetSection(""); err == nil {
			if key := section.Key("default-canvas"); key != nil {
				defaultCanvas = key.String()
			}
		}
		if defaultCanvas == "" {
			dialog.ShowError(fmt.Errorf("default-canvas must be set for kiosk modes"), window)
			return
		}
	}

	// Get plugin directory
	pluginDir := m.getPluginDirectory()
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		dialog.ShowError(fmt.Errorf("failed to create plugin directory: %w", err), window)
		return
	}

	// Generate plugin manifest
	manifest := m.createPluginManifest()
	manifestPath := filepath.Join(pluginDir, ".canvusplugin")
	if err := m.writeManifest(manifestPath, manifest); err != nil {
		dialog.ShowError(fmt.Errorf("failed to write manifest: %w", err), window)
		return
	}

	// Generate CSS
	css := m.generateCSS()
	cssPath := filepath.Join(pluginDir, "styles.css")
	if err := os.WriteFile(cssPath, []byte(css), 0644); err != nil {
		dialog.ShowError(fmt.Errorf("failed to write CSS: %w", err), window)
		return
	}

	// Update mt-canvus.ini with plugin folder
	if err := m.updatePluginFolders(iniFile, iniPath, pluginDir); err != nil {
		dialog.ShowError(fmt.Errorf("failed to update mt-canvus.ini: %w", err), window)
		return
	}

	m.statusLabel.SetText("Plugin generated successfully")
	dialog.ShowInformation("Success", fmt.Sprintf("CSS plugin generated at:\n%s", pluginDir), window)
}

// createPluginManifest creates the plugin manifest.
func (m *Manager) createPluginManifest() *PluginManifest {
	return &PluginManifest{
		APIVersion:  "1.0",
		Name:        "CanvusPowerToys CSS Options",
		Version:     "1.0.0",
		Description: "CSS options generated by Canvus PowerToys",
	}
}

// writeManifest writes the plugin manifest to file.
func (m *Manager) writeManifest(path string, manifest *PluginManifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// generateCSS generates the CSS content based on enabled options.
func (m *Manager) generateCSS() string {
	var css strings.Builder

	// Widget Options: Control moving, scaling, and rotation
	// Note: These are temporary and revert when canvas is closed
	// Always generate CSS for widget options (users can disable them)
	css.WriteString("/* Widget Options CSS - Control widget interaction capabilities */\n")
	css.WriteString("/* Note: These are temporary and revert when canvas is closed */\n")
	css.WriteString("CanvusCanvasWidget > * {\n")
	if m.movingEnabled.Checked {
		css.WriteString("  input-translate: true !important;\n")
	} else {
		css.WriteString("  input-translate: false !important;\n")
	}
	if m.scalingEnabled.Checked {
		css.WriteString("  input-scale: true !important;\n")
	} else {
		css.WriteString("  input-scale: false !important;\n")
	}
	if m.rotationEnabled.Checked {
		css.WriteString("  input-rotation: true !important;\n")
	} else {
		css.WriteString("  input-rotation: false !important;\n")
	}
	css.WriteString("}\n\n")

	// Video Looping: Enable looping for video widgets
	// Note: playmode: loop may need testing - similar to playmode: no-loop used for animated icons
	if m.videoLoopEnabled.Checked {
		css.WriteString("/* Video Looping CSS - Enables automatic looping for videos */\n")
		css.WriteString("/* Warning: May consume significant memory with large videos */\n")
		css.WriteString("CanvusCanvasVideo {\n")
		css.WriteString("  playmode: loop !important;\n")
		css.WriteString("}\n\n")
		css.WriteString("/* Note: A toggle button for video looping would require plugin code (C++), not just CSS */\n")
		css.WriteString("/* This is a future enhancement that would need a custom button widget */\n\n")
	}

	// Kiosk Mode: Hide UI layers (SidebarWidget and MainMenu), hide finger menu
	if m.kioskModeEnabled.Checked {
		css.WriteString("/* Kiosk Mode CSS - Hides UI layers for kiosk presentation mode */\n")
		css.WriteString("/* Hides sidebar menu that appears by widgets */\n")
		css.WriteString("SidebarWidget {\n")
		css.WriteString("  display: none !important;\n")
		css.WriteString("}\n\n")
		css.WriteString("/* Hides main menu at side/bottom of screen */\n")
		css.WriteString("MainMenu {\n")
		css.WriteString("  display: none !important;\n")
		css.WriteString("}\n\n")
		css.WriteString("/* Hides finger menu (CanvusCanvasMenu) */\n")
		css.WriteString("CanvusCanvasMenu {\n")
		css.WriteString("  display: none !important;\n")
		css.WriteString("}\n\n")
		css.WriteString("/* Note: Third-party touch menus are not affected by these rules */\n\n")
	}

	// Kiosk Plus Mode: Hide UI layers but keep finger menu visible
	if m.kioskPlusEnabled.Checked {
		css.WriteString("/* Kiosk Plus Mode CSS - Hides UI layers but keeps finger menu enabled */\n")
		css.WriteString("/* Hides sidebar menu that appears by widgets */\n")
		css.WriteString("SidebarWidget {\n")
		css.WriteString("  display: none !important;\n")
		css.WriteString("}\n\n")
		css.WriteString("/* Hides main menu at side/bottom of screen */\n")
		css.WriteString("MainMenu {\n")
		css.WriteString("  display: none !important;\n")
		css.WriteString("}\n\n")
		css.WriteString("/* Keep finger menu (CanvusCanvasMenu) visible for content creation */\n")
		css.WriteString("/* CanvusCanvasMenu is NOT hidden in Plus mode */\n\n")
		css.WriteString("/* Note: Third-party touch menus are not affected by these rules */\n\n")
	}

	// Hide Title Bars
	if m.hideTitleBarsEnabled.Checked {
		css.WriteString("/* Hide Title Bars CSS - Hides title bars on canvas items */\n")
		css.WriteString("TitleBarWidget {\n")
		css.WriteString("  display: none !important;\n")
		css.WriteString("}\n\n")
	}

	// Hide Resize Handles
	if m.hideResizeHandlesEnabled.Checked {
		css.WriteString("/* Hide Resize Handles CSS - Hides resize handles on canvas items */\n")
		css.WriteString("CanvusResizeHandleWidget {\n")
		css.WriteString("  display: none !important;\n")
		css.WriteString("}\n\n")
	}

	// Ensure third-party touch menu is not hidden (applies to both kiosk modes)
	if m.kioskModeEnabled.Checked || m.kioskPlusEnabled.Checked {
		css.WriteString("/* Ensure third-party touch menus remain visible */\n")
		css.WriteString("/* Third-party menus should not be hidden by kiosk modes */\n")
		css.WriteString("/* Note: Adjust selector if third-party menus use different class names */\n\n")
	}

	return css.String()
}

// getPluginDirectory returns the plugin directory path.
func (m *Manager) getPluginDirectory() string {
	configDir := m.fileService.GetUserConfigPath()
	return filepath.Join(configDir, "plugins", "canvus-powertoys-css")
}

// updatePluginFolders updates the plugin-folders setting in mt-canvus.ini.
// Plugin folders are registered in the [content] section per mt-canvus-plugin-api.md
func (m *Manager) updatePluginFolders(iniFile *ini.File, iniPath, pluginDir string) error {
	section, err := iniFile.GetSection("content")
	if err != nil {
		section, _ = iniFile.NewSection("content")
	}

	// Normalize path for INI file: replace backslashes with forward slashes
	// Windows paths need to use / or \\ in INI files, forward slashes are safer
	normalizedPluginDir := strings.ReplaceAll(pluginDir, "\\", "/")

	// Get existing plugin-folders
	existingFolders := ""
	if key := section.Key("plugin-folders"); key != nil {
		existingFolders = key.String()
	}

	// Normalize existing folders too (in case they have backslashes)
	normalizedExisting := strings.ReplaceAll(existingFolders, "\\", "/")

	// Add our plugin directory if not already present
	if !strings.Contains(normalizedExisting, normalizedPluginDir) {
		if normalizedExisting != "" {
			normalizedExisting += "," + normalizedPluginDir
		} else {
			normalizedExisting = normalizedPluginDir
		}
		section.Key("plugin-folders").SetValue(normalizedExisting)
	}

	// Create backup before updating
	if _, err := os.Stat(iniPath); err == nil {
		if _, err := m.backupManager.CreateBackup(iniPath); err != nil {
			fmt.Printf("Warning: Failed to create backup: %v\n", err)
		}
	}

	return m.iniParser.Write(iniFile, iniPath)
}

// launchCanvusWithConfig launches Canvus with the current CSS configuration.
func (m *Manager) launchCanvusWithConfig(window fyne.Window) {
	// Check if any options are enabled
	if !m.movingEnabled.Checked && !m.scalingEnabled.Checked && !m.rotationEnabled.Checked &&
		!m.videoLoopEnabled.Checked && !m.kioskModeEnabled.Checked && !m.kioskPlusEnabled.Checked &&
		!m.hideTitleBarsEnabled.Checked && !m.hideResizeHandlesEnabled.Checked {
		dialog.ShowError(fmt.Errorf("please enable at least one CSS option before launching"), window)
		return
	}

	// Find Canvus executable
	canvusExe := m.findCanvusExecutable()
	if canvusExe == "" {
		// Ask user to select Canvus executable
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			uri := reader.URI()
			var exePath string
			if uri.Scheme() == "file" {
				exePath = uri.Path()
			} else {
				exePath = uri.String()
			}
			m.launchCanvusWithPath(window, exePath)
		}, window)
		return
	}

	m.launchCanvusWithPath(window, canvusExe)
}

// launchCanvusWithPath launches Canvus with the specified executable path.
func (m *Manager) launchCanvusWithPath(window fyne.Window, canvusExe string) {
	// Generate CSS file for current options
	cssDir := m.getCSSDirectory()
	if err := os.MkdirAll(cssDir, 0755); err != nil {
		dialog.ShowError(fmt.Errorf("failed to create CSS directory: %w", err), window)
		return
	}

	// Generate unique filename based on enabled options
	cssFileName := m.generateCSSFileName()
	cssPath := filepath.Join(cssDir, cssFileName)

	// Generate CSS content
	css := m.generateCSS()
	if err := os.WriteFile(cssPath, []byte(css), 0644); err != nil {
		dialog.ShowError(fmt.Errorf("failed to write CSS file: %w", err), window)
		return
	}

	// Launch Canvus with CSS file
	workingDir := filepath.Dir(canvusExe)
	cmd := exec.Command(canvusExe, "--css", cssPath)
	cmd.Dir = workingDir

	if err := cmd.Start(); err != nil {
		dialog.ShowError(fmt.Errorf("failed to launch Canvus: %w", err), window)
		return
	}

	m.statusLabel.SetText(fmt.Sprintf("Canvus launched with CSS: %s", cssFileName))
	dialog.ShowInformation("Launched", fmt.Sprintf("Canvus launched with CSS configuration:\n%s", cssPath), window)
}

// findCanvusExecutable attempts to find the Canvus executable in common locations.
func (m *Manager) findCanvusExecutable() string {
	// Common Windows installation locations
	possiblePaths := []string{
		`C:\Program Files\MultiTaction\Canvus\Canvus.exe`,
		`C:\Program Files (x86)\MultiTaction\Canvus\Canvus.exe`,
		`C:\MultiTaction\Canvus\Canvus.exe`,
	}

	// Check if Canvus.exe exists in any of these locations
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Check environment variables
	if programFiles := os.Getenv("ProgramFiles"); programFiles != "" {
		path := filepath.Join(programFiles, "MultiTaction", "Canvus", "Canvus.exe")
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	if programFilesX86 := os.Getenv("ProgramFiles(x86)"); programFilesX86 != "" {
		path := filepath.Join(programFilesX86, "MultiTaction", "Canvus", "Canvus.exe")
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// getCSSDirectory returns the directory for CSS files.
func (m *Manager) getCSSDirectory() string {
	configDir := m.fileService.GetUserConfigPath()
	return filepath.Join(configDir, "css")
}

// generateCSSFileName generates a unique filename based on enabled options.
func (m *Manager) generateCSSFileName() string {
	var parts []string

	if m.movingEnabled.Checked {
		parts = append(parts, "move")
	}
	if m.scalingEnabled.Checked {
		parts = append(parts, "scale")
	}
	if m.rotationEnabled.Checked {
		parts = append(parts, "rotate")
	}
	if m.videoLoopEnabled.Checked {
		parts = append(parts, "loop")
	}
	if m.kioskModeEnabled.Checked {
		parts = append(parts, "kiosk")
	}
	if m.kioskPlusEnabled.Checked {
		parts = append(parts, "kioskplus")
	}
	if m.hideTitleBarsEnabled.Checked {
		parts = append(parts, "notitle")
	}
	if m.hideResizeHandlesEnabled.Checked {
		parts = append(parts, "noresize")
	}

	if len(parts) == 0 {
		parts = append(parts, "default")
	}

	filename := "canvus-" + strings.Join(parts, "-") + ".css"
	return filename
}

// generateShortcutName generates a unique shortcut name based on enabled options.
func (m *Manager) generateShortcutName() string {
	var parts []string

	if m.movingEnabled.Checked {
		parts = append(parts, "Move")
	}
	if m.scalingEnabled.Checked {
		parts = append(parts, "Scale")
	}
	if m.rotationEnabled.Checked {
		parts = append(parts, "Rotate")
	}
	if m.videoLoopEnabled.Checked {
		parts = append(parts, "Loop")
	}
	if m.kioskModeEnabled.Checked {
		parts = append(parts, "Kiosk")
	}
	if m.kioskPlusEnabled.Checked {
		parts = append(parts, "KioskPlus")
	}
	if m.hideTitleBarsEnabled.Checked {
		parts = append(parts, "NoTitle")
	}
	if m.hideResizeHandlesEnabled.Checked {
		parts = append(parts, "NoResize")
	}

	if len(parts) == 0 {
		return "Canvus Custom CSS"
	}

	return "Canvus " + strings.Join(parts, " ")
}
