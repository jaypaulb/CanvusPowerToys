package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jaypaulb/CanvusPowerToys/internal/molecules/screenxml"
	"github.com/jaypaulb/CanvusPowerToys/internal/molecules/tray"
	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// MainWindow represents the main application window with tabs.
type MainWindow struct {
	window fyne.Window
	tabs   *container.AppTabs
	tray   *tray.Manager
}

// NewMainWindow creates a new main window instance.
func NewMainWindow(app fyne.App) *MainWindow {
	window := app.NewWindow("Canvus PowerToys")
	window.Resize(fyne.NewSize(1024, 768))
	window.CenterOnScreen()

	// Create tabs
	// Initialize Screen.xml Creator
	fileService, err := services.NewFileService()
	if err != nil {
		// Fallback if file service fails
		fileService = nil
	}

	var screenXMLCreator fyne.CanvasObject
	if fileService != nil {
		creator, err := screenxml.NewCreator(fileService)
		if err == nil {
			screenXMLCreator = creator.CreateUI(window)
		} else {
			screenXMLCreator = widget.NewLabel("Screen.xml Creator - Error initializing")
		}
	} else {
		screenXMLCreator = widget.NewLabel("Screen.xml Creator - Error initializing file service")
	}

	tabs := container.NewAppTabs(
		&container.TabItem{
			Text:    "Screen.xml Creator",
			Content: screenXMLCreator,
		},
		&container.TabItem{
			Text:    "Config Editor",
			Content: widget.NewLabel("Canvus Config Editor - Coming soon"),
		},
		&container.TabItem{
			Text:    "CSS Options",
			Content: widget.NewLabel("CSS Options Manager - Coming soon"),
		},
		&container.TabItem{
			Text:    "Custom Menu",
			Content: widget.NewLabel("Custom Menu Designer - Coming soon"),
		},
		&container.TabItem{
			Text:    "WebUI",
			Content: widget.NewLabel("WebUI Settings - Coming soon"),
		},
	)

	window.SetContent(tabs)

	// Setup system tray
	trayManager := tray.NewManager(window, app)
	trayManager.Setup()

	// Configure window to minimize to tray
	window.SetCloseIntercept(func() {
		window.Hide()
	})

	return &MainWindow{
		window: window,
		tabs:   tabs,
		tray:   trayManager,
	}
}

// ShowAndRun displays the window and runs the application.
func (mw *MainWindow) ShowAndRun() {
	mw.window.ShowAndRun()
}

// GetWindow returns the underlying Fyne window.
func (mw *MainWindow) GetWindow() fyne.Window {
	return mw.window
}
