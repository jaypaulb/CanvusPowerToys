package app

import (
	"bytes"
	"image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/jaypaulb/CanvusPowerToys/assets"
	"github.com/jaypaulb/CanvusPowerToys/internal/molecules/configeditor"
	"github.com/jaypaulb/CanvusPowerToys/internal/molecules/custommenu"
	"github.com/jaypaulb/CanvusPowerToys/internal/molecules/cssoptions"
	"github.com/jaypaulb/CanvusPowerToys/internal/molecules/screenxml"
	"github.com/jaypaulb/CanvusPowerToys/internal/molecules/tray"
	"github.com/jaypaulb/CanvusPowerToys/internal/molecules/webui"
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

	// Set window icon from embedded assets
	if iconData, err := assets.Icons.ReadFile("CanvusPowerToysIcon.png"); err == nil {
		if _, err := png.Decode(bytes.NewReader(iconData)); err == nil {
			icon := fyne.NewStaticResource("CanvusPowerToysIcon.png", iconData)
			window.SetIcon(icon)
		}
	}

	// Create tabs
	// Initialize Screen.xml Creator
	fileService, err := services.NewFileService()
	var screenXMLCreator fyne.CanvasObject
	if err == nil {
		creator, err := screenxml.NewCreator(fileService)
		if err == nil {
			screenXMLCreator = creator.CreateUI(window)
		} else {
			screenXMLCreator = widget.NewLabel("Screen.xml Creator - Error initializing")
		}
	} else {
		screenXMLCreator = widget.NewLabel("Screen.xml Creator - Error initializing file service")
	}

	// Initialize Config Editor
	var configEditor fyne.CanvasObject
	if fileService != nil {
		editor, err := configeditor.NewEditor(fileService)
		if err == nil {
			configEditor = editor.CreateUI(window)
		} else {
			configEditor = widget.NewLabel("Config Editor - Error initializing")
		}
	} else {
		configEditor = widget.NewLabel("Config Editor - Error initializing file service")
	}

	// Initialize CSS Options Manager
	var cssOptions fyne.CanvasObject
	if fileService != nil {
		cssMgr, err := cssoptions.NewManager(fileService)
		if err == nil {
			cssOptions = cssMgr.CreateUI(window)
		} else {
			cssOptions = widget.NewLabel("CSS Options Manager - Error initializing")
		}
	} else {
		cssOptions = widget.NewLabel("CSS Options Manager - Error initializing file service")
	}

	// Initialize Custom Menu Designer
	var customMenu fyne.CanvasObject
	if fileService != nil {
		menuDesigner, err := custommenu.NewDesigner(fileService)
		if err == nil {
			customMenu = menuDesigner.CreateUI(window)
		} else {
			customMenu = widget.NewLabel("Custom Menu Designer - Error initializing")
		}
	} else {
		customMenu = widget.NewLabel("Custom Menu Designer - Error initializing file service")
	}

	// Initialize WebUI Manager
	var webUI fyne.CanvasObject
	if fileService != nil {
		webUIMgr, err := webui.NewManager(fileService)
		if err == nil {
			webUI = webUIMgr.CreateUI(window)
		} else {
			webUI = widget.NewLabel("WebUI - Error initializing")
		}
	} else {
		webUI = widget.NewLabel("WebUI - Error initializing file service")
	}

	tabs := container.NewAppTabs(
		&container.TabItem{
			Text:    "Screen.xml Creator",
			Content: screenXMLCreator,
		},
		&container.TabItem{
			Text:    "Config Editor",
			Content: configEditor,
		},
		&container.TabItem{
			Text:    "CSS Options",
			Content: cssOptions,
		},
		&container.TabItem{
			Text:    "Custom Menu",
			Content: customMenu,
		},
		&container.TabItem{
			Text:    "WebUI",
			Content: webUI,
		},
	)

	// Create a note label with warning triangle, aligned inline with tabs
	// Using Unicode characters: ⚠ (warning triangle) and ↑ (upward arrow)
	trayNote := widget.NewLabel("⚠ Close to System Tray ↑")
	trayNote.TextStyle = fyne.TextStyle{Italic: true}

	// Create a horizontal container to place note inline with tab bar
	// The tabs widget will render its tab bar, and we position the note to the right
	// Using Border layout with tabs in center and note on right side
	contentWithNote := container.NewBorder(
		nil,   // Top: nothing
		nil,   // Bottom: nothing
		nil,   // Left: nothing
		trayNote, // Right: note (will appear to the right of tabs)
		tabs,     // Center: tabs widget
	)

	window.SetContent(contentWithNote)

	// Setup system tray (this also sets up close intercept)
	trayManager := tray.NewManager(window, app)
	trayManager.Setup()

	// Store reference to window for minimize handling
	mw := &MainWindow{
		window: window,
		tabs:   tabs,
		tray:   trayManager,
	}

	// Minimize to tray is handled by the tray manager
	// The window will hide when minimized (handled in tray.Setup)

	return mw
}

// ShowAndRun displays the window and runs the application.
func (mw *MainWindow) ShowAndRun() {
	mw.window.ShowAndRun()
}

// GetWindow returns the underlying Fyne window.
func (mw *MainWindow) GetWindow() fyne.Window {
	return mw.window
}

