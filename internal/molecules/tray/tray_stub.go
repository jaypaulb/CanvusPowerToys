//go:build !windows
// +build !windows

package tray

import (
	"fyne.io/fyne/v2"
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
	m.window.SetCloseIntercept(func() {
		m.window.Hide()
	})
}
