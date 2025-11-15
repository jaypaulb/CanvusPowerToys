package cssoptions

import (
	"encoding/json"
	"fmt"
	"os"
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
	iniParser        *config.INIParser
	fileService      *services.FileService
	backupManager    *backup.Manager
	rotationEnabled  *widget.Check
	videoLoopEnabled *widget.Check
	kioskModeEnabled *widget.Check
	kioskPlusEnabled *widget.Check
	statusLabel      *widget.Label
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

	// Rotation option
	rotationLabel := widget.NewLabel("Enable Rotation")
	rotationTooltip := widget.NewLabel("Allow rotation of screen items (temporary, reverts on close)")
	m.rotationEnabled = widget.NewCheck("", nil)

	// Video Loop option
	videoLoopLabel := widget.NewLabel("Enable Video Looping")
	videoLoopTooltip := widget.NewLabel("Warning: May use significant memory with large videos")
	videoLoopWarning := widget.NewLabel("⚠ Memory Warning: Large videos may consume significant resources")
	videoLoopWarning.Importance = widget.WarningImportance
	m.videoLoopEnabled = widget.NewCheck("", nil)

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

	form := container.NewVBox(
		title,
		instructions,
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			rotationLabel, m.rotationEnabled,
			rotationTooltip, widget.NewLabel(""),
			videoLoopLabel, m.videoLoopEnabled,
			videoLoopTooltip, widget.NewLabel(""),
			videoLoopWarning, widget.NewLabel(""),
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

	if m.rotationEnabled.Checked {
		css.WriteString("/* Rotation CSS */\n")
		css.WriteString(".canvus-item {\n")
		css.WriteString("  transform-origin: center;\n")
		css.WriteString("}\n\n")
	}

	if m.videoLoopEnabled.Checked {
		css.WriteString("/* Video Looping CSS */\n")
		css.WriteString("video {\n")
		css.WriteString("  loop: true;\n")
		css.WriteString("}\n\n")
	}

	if m.kioskModeEnabled.Checked {
		css.WriteString("/* Kiosk Mode CSS */\n")
		css.WriteString(".canvus-kiosk-mode {\n")
		css.WriteString("  /* Kiosk mode styles */\n")
		css.WriteString("}\n\n")
	}

	if m.kioskPlusEnabled.Checked {
		css.WriteString("/* Kiosk Plus Mode CSS */\n")
		css.WriteString(".canvus-kiosk-plus-mode {\n")
		css.WriteString("  /* Kiosk plus mode styles */\n")
		css.WriteString("}\n\n")
	}

	// Ensure third-party touch menu is not hidden
	css.WriteString("/* Ensure third-party touch menu is visible */\n")
	css.WriteString(".third-party-menu {\n")
	css.WriteString("  display: block !important;\n")
	css.WriteString("  visibility: visible !important;\n")
	css.WriteString("}\n")

	return css.String()
}

// getPluginDirectory returns the plugin directory path.
func (m *Manager) getPluginDirectory() string {
	configDir := m.fileService.GetUserConfigPath()
	return filepath.Join(configDir, "plugins", "canvus-powertoys-css")
}

// updatePluginFolders updates the plugin-folders setting in mt-canvus.ini.
func (m *Manager) updatePluginFolders(iniFile *ini.File, iniPath, pluginDir string) error {
	section, err := iniFile.GetSection("")
	if err != nil {
		section, _ = iniFile.NewSection("")
	}

	// Get existing plugin-folders
	existingFolders := ""
	if key := section.Key("plugin-folders"); key != nil {
		existingFolders = key.String()
	}

	// Add our plugin directory if not already present
	if !strings.Contains(existingFolders, pluginDir) {
		if existingFolders != "" {
			existingFolders += "," + pluginDir
		} else {
			existingFolders = pluginDir
		}
		section.Key("plugin-folders").SetValue(existingFolders)
	}

	// Create backup before updating
	if _, err := os.Stat(iniPath); err == nil {
		if _, err := m.backupManager.CreateBackup(iniPath); err != nil {
			fmt.Printf("Warning: Failed to create backup: %v\n", err)
		}
	}

	return m.iniParser.Write(iniFile, iniPath)
}
