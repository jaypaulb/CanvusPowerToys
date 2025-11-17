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
	// Hide window when close button (X) is clicked - this hides to tray
	// Note: Minimize button minimizes to taskbar (standard Windows behavior)
	m.window.SetCloseIntercept(func() {
		m.window.Hide()
	})

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
	// Since embed is "//go:embed icons", path is relative to icons directory
	if iconData, err := assets.Icons.ReadFile("CanvusPowerToysIcon.png"); err == nil {
		// Verify it's valid PNG by attempting to decode
		if _, err := png.Decode(bytes.NewReader(iconData)); err == nil {
			systray.SetIcon(iconData)
		}
	} else {
		// Log error for debugging (only in development)
		// Try alternative path in case embed structure is different
		if iconData, err := assets.Icons.ReadFile("icons/CanvusPowerToysIcon.png"); err == nil {
			if _, err := png.Decode(bytes.NewReader(iconData)); err == nil {
				systray.SetIcon(iconData)
			}
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
