//go:build windows
// +build windows

package tray

import (
	"fyne.io/fyne/v2"
	"github.com/getlantern/systray"
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

	go func() {
		systray.Run(m.onReady, m.onExit)
	}()
}

func (m *Manager) onReady() {
	systray.SetTitle("Canvus PowerToys")
	systray.SetTooltip("Canvus PowerToys")

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
