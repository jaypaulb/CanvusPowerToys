//go:build windows
// +build windows

package tray

import (
	"bytes"
	"image/png"

	"fyne.io/fyne/v2"
	"github.com/getlantern/systray"
	"github.com/jaypaulb/CanvusPowerToys/assets"
)

// Manager handles system tray integration.
type Manager struct {
	window fyne.Window
	app    fyne.App
}

// NewManager creates a new system tray manager.
func NewManager(window fyne.Window, app fyne.App) *Manager {
	return &Manager{
		window: window,
		app:    app,
	}
}

// Setup initializes the system tray.
func (m *Manager) Setup() {
	// Hide window when close button is clicked (minimize to tray)
	m.window.SetCloseIntercept(func() {
		m.window.Hide()
	})

	// Setup minimize-to-tray hook (Windows-specific)
	// This intercepts the minimize button to hide window instead
	if err := setupMinimizeToTray(m); err != nil {
		// If hook setup fails, continue without it
		// Close button will still work to hide to tray
	}

	// Start systray in background
	go func() {
		systray.Run(m.onReady, m.onExit)
	}()
}

func (m *Manager) onReady() {
	systray.SetTitle("Canvus PowerToys")
	systray.SetTooltip("Canvus PowerToys")

	// Load and set tray icon from embedded assets
	// systray.SetIcon expects PNG image data as []byte
	if iconData, err := assets.Icons.ReadFile("icons/CanvusPowerToysIcon.png"); err == nil {
		// Verify it's valid PNG by attempting to decode
		if _, err := png.Decode(bytes.NewReader(iconData)); err == nil {
			systray.SetIcon(iconData)
		}
	}

	showItem := systray.AddMenuItem("Show", "Show window")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "Exit application")

	go func() {
		for {
			select {
			case <-showItem.ClickedCh:
				m.window.Show()
				m.window.RequestFocus()
			case <-quitItem.ClickedCh:
				systray.Quit()
				m.app.Quit()
			}
		}
	}()
}

func (m *Manager) onExit() {
}
