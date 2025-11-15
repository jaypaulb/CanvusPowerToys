package screenxml

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// Creator is the main Screen.xml Creator component.
type Creator struct {
	grid            *GridWidget
	gridContainer   *GridContainer
	xmlGenerator    *XMLGenerator
	iniIntegration  *INIIntegration
	fileService     *services.FileService
	window          fyne.Window
}

// NewCreator creates a new Screen.xml Creator.
func NewCreator(fileService *services.FileService) (*Creator, error) {
	grid := NewGridWidget()
	gridContainer := NewGridContainer(grid)

	creator := &Creator{
		grid:            grid,
		gridContainer:      gridContainer,
		iniIntegration:   NewINIIntegration(),
		fileService:      fileService,
	}

	// Initialize XML generator (simplified - no longer need separate handlers)
	defaultRes := Resolution{Width: 1920, Height: 1080, Name: "1920x1080 (Full HD)"}
	creator.xmlGenerator = NewXMLGenerator(grid, nil, nil, defaultRes)

	return creator, nil
}

// CreateUI creates the UI for the Screen.xml Creator tab.
func (c *Creator) CreateUI(window fyne.Window) fyne.CanvasObject {
	c.window = window

	// Top: 3 buttons only
	topBar := container.NewHBox(
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

	videoOutputs := c.iniIntegration.DetectVideoOutputs(c.grid)
	if len(videoOutputs) == 0 {
		dialog.ShowInformation("No Outputs", "No video outputs detected. Assign GPU outputs to cells first.", window)
		return
	}

	if err := c.iniIntegration.UpdateMtCanvusIni(iniPath, videoOutputs); err != nil {
		dialog.ShowError(err, window)
		return
	}

	config := c.iniIntegration.GenerateVideoOutputConfig(videoOutputs)
	dialog.ShowInformation("Success", fmt.Sprintf("mt-canvus.ini updated successfully.\n\nvideo-output=%s", config), window)
}

// generateAndPreview generates screen.xml and shows preview.
func (c *Creator) generateAndPreview(window fyne.Window) {
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

	// Show preview dialog
	preview := widget.NewMultiLineEntry()
	preview.SetText(string(xmlData))
	preview.Disable()

	previewDialog := dialog.NewCustom("Generated screen.xml Preview", "Close", preview, window)
	previewDialog.Resize(fyne.NewSize(800, 600))
	previewDialog.Show()
}

// saveScreenXML saves the generated screen.xml to file.
func (c *Creator) saveScreenXML(window fyne.Window) {
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

	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		if _, err := writer.Write(xmlData); err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowInformation("Success", "screen.xml saved successfully", window)
	}, window)
}


