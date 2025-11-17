package screenxml

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// Creator is the main Screen.xml Creator component.
type Creator struct {
	grid           *GridWidget
	gridContainer  *GridContainer
	xmlGenerator   *XMLGenerator
	iniIntegration *INIIntegration
	fileService    *services.FileService
	window         fyne.Window
	areaPerScreen  *widget.Check // Checkbox for "area per screen or per gpu"
}

// NewCreator creates a new Screen.xml Creator.
func NewCreator(fileService *services.FileService) (*Creator, error) {
	grid := NewGridWidget()
	gridContainer := NewGridContainer(grid)

	creator := &Creator{
		grid:           grid,
		gridContainer:  gridContainer,
		iniIntegration: NewINIIntegration(),
		fileService:    fileService,
	}

	// Initialize XML generator (simplified - no longer need separate handlers)
	defaultRes := Resolution{Width: 1920, Height: 1080, Name: "1920x1080 (Full HD)"}
	creator.xmlGenerator = NewXMLGenerator(grid, nil, nil, defaultRes)

	return creator, nil
}

// CreateUI creates the UI for the Screen.xml Creator tab.
func (c *Creator) CreateUI(window fyne.Window) fyne.CanvasObject {
	c.window = window

	// Top row: Checkbox for area generation mode
	// Per Screen = per GPU output (each output is one screen)
	// Per GPU = per window (one window per GPU, matching all outputs on that GPU)
	c.areaPerScreen = widget.NewCheck("Area per GPU (not per Screen)", func(checked bool) {
		// Update XML generator mode
		// checked = true means "per GPU" (per window)
		// checked = false means "per Screen" (per GPU output)
		c.xmlGenerator.SetAreaPerGPU(checked)
	})
	c.areaPerScreen.SetChecked(false) // Default: area per Screen (per GPU output)

	// Top: Checkbox + 3 buttons
	topBar := container.NewHBox(
		c.areaPerScreen,
		widget.NewSeparator(),
		widget.NewButton("Generate Screen.xml", func() {
			c.generateAndPreview(window)
		}),
		widget.NewButton("Save screen.xml", func() {
			c.saveScreenXML(window)
		}),
		widget.NewButton("Update mt-canvus.ini", func() {
			c.updateMtCanvusIni(window)
		}),
	)

	// Main: Grid with cell widgets wrapped in drag selection widget
	dragSelectionWidget := NewDragSelectionWidget(c.gridContainer)

	return container.NewBorder(
		topBar,
		nil, nil, nil,
		dragSelectionWidget,
	)
}

// updateMtCanvusIni updates mt-canvus.ini with video-output configuration.
func (c *Creator) updateMtCanvusIni(window fyne.Window) {
	iniPath := c.fileService.DetectMtCanvusIni()
	if iniPath == "" {
		dialog.ShowInformation("Not Found", "mt-canvus.ini not found in standard locations", window)
		return
	}

	outputCells := c.iniIntegration.DetectVideoOutputs(c.grid)
	if len(outputCells) == 0 {
		dialog.ShowInformation("No Outputs", "No cells without layout detected. Assign GPU outputs to cells (without layer) first.", window)
		return
	}

	if err := c.iniIntegration.UpdateMtCanvusIni(iniPath, outputCells); err != nil {
		dialog.ShowError(err, window)
		return
	}

	dialog.ShowInformation("Success", fmt.Sprintf("mt-canvus.ini updated successfully.\n\nCreated %d output section(s) for cells without layout.", len(outputCells)), window)
}

// generateAndPreview generates screen.xml and shows preview.
func (c *Creator) generateAndPreview(window fyne.Window) {
	// Check if there are any cells with data
	if !c.grid.HasCellsWithData() {
		dialog.ShowError(fmt.Errorf("Please add some displays to your array!"), window)
		return
	}

	xmlData, err := c.xmlGenerator.Generate()
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Validate
	if err := c.xmlGenerator.Validate(xmlData); err != nil {
		dialog.ShowError(fmt.Errorf("validation failed: %w", err), window)
		return
	}

	// Use MultiLineEntry for better performance with large XML content
	// MultiLineEntry handles large content much better than canvas.Text
	// which was causing the freezing issue
	previewEntry := widget.NewMultiLineEntry()
	xmlText := string(xmlData)
	previewEntry.SetText(xmlText)
	// Keep enabled to allow text selection, but prevent editing.
	isResetting := false
	previewEntry.OnChanged = func(text string) {
		if isResetting || text == xmlText {
			return
		}
		isResetting = true
		previewEntry.SetText(xmlText)
		isResetting = false
	}
	previewEntry.Wrapping = fyne.TextWrapOff // Don't wrap XML

	// Copy to clipboard button
	copyBtn := widget.NewButton("Copy to Clipboard", func() {
		window.Clipboard().SetContent(previewEntry.Text)
		dialog.ShowInformation("Copied", "XML content copied to clipboard", window)
	})

	// Create a container with the entry and button
	content := container.NewBorder(
		copyBtn, // Top: Copy button
		nil, nil, nil,
		container.NewScroll(previewEntry), // Center: Scrollable text
	)

	previewDialog := dialog.NewCustom("Generated screen.xml Preview", "Close", content, window)
	previewDialog.Resize(fyne.NewSize(800, 600))
	previewDialog.Show()
}

// saveScreenXML saves the generated screen.xml to file.
func (c *Creator) saveScreenXML(window fyne.Window) {
	// Check if there are any cells with data
	if !c.grid.HasCellsWithData() {
		dialog.ShowError(fmt.Errorf("Please add some displays to your array!"), window)
		return
	}

	xmlData, err := c.xmlGenerator.Generate()
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Validate
	if err := c.xmlGenerator.Validate(xmlData); err != nil {
		dialog.ShowError(fmt.Errorf("validation failed: %w", err), window)
		return
	}

	// Determine default save location (user config first, then system)
	defaultPath := c.fileService.DetectScreenXml()
	if defaultPath == "" {
		// Fallback to user config directory
		defaultPath = filepath.Join(c.fileService.GetUserConfigPath(), "screen.xml")
	}

	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		if _, err := writer.Write(xmlData); err != nil {
			dialog.ShowError(fmt.Errorf("failed to write file: %w", err), window)
			return
		}

		dialog.ShowInformation("Success", "screen.xml saved successfully", window)
	}, window)

	// Set default directory and filename
	if defaultPath != "" {
		dir := filepath.Dir(defaultPath)
		fileName := filepath.Base(defaultPath)

		if dir != "" {
			if dirURI, err := storage.ListerForURI(storage.NewFileURI(dir)); err == nil {
				saveDialog.SetLocation(dirURI)
			}
		}
		saveDialog.SetFileName(fileName)
	}

	saveDialog.Show()
}
