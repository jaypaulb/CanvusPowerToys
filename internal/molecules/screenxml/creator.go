package screenxml

import (
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// Creator is the main Screen.xml Creator component.
type Creator struct {
	grid              *GridWidget
	gpuAssignment     *GPUAssignment
	touchAreaHandler  *TouchAreaHandler
	resolutionHandler *ResolutionHandler
	xmlGenerator      *XMLGenerator
	iniIntegration    *INIIntegration
	fileService       *services.FileService
}

// NewCreator creates a new Screen.xml Creator.
func NewCreator(fileService *services.FileService) (*Creator, error) {
	grid := NewGridWidget()
	gpuAssignment := NewGPUAssignment(grid)
	touchAreaHandler := NewTouchAreaHandler(grid)
	resolutionHandler := NewResolutionHandler()

	creator := &Creator{
		grid:              grid,
		gpuAssignment:     gpuAssignment,
		touchAreaHandler:  touchAreaHandler,
		resolutionHandler: resolutionHandler,
		iniIntegration:    NewINIIntegration(),
		fileService:       fileService,
	}

	// Initialize XML generator
	resolution := resolutionHandler.GetCurrentResolution()
	creator.xmlGenerator = NewXMLGenerator(grid, gpuAssignment, touchAreaHandler, resolution)

	return creator, nil
}

// CreateUI creates the UI for the Screen.xml Creator tab.
func (c *Creator) CreateUI(window fyne.Window) fyne.CanvasObject {
	// Left panel: Controls
	controls := container.NewVBox(
		widget.NewLabel("Screen.xml Creator"),
		widget.NewSeparator(),
		c.gpuAssignment.CreateUI(),
		widget.NewSeparator(),
		c.touchAreaHandler.CreateUI(),
		widget.NewSeparator(),
		c.resolutionHandler.CreateUI(),
		widget.NewSeparator(),
		c.createActionButtons(window),
	)

	// Right panel: Grid
	gridContainer := container.NewBorder(
		widget.NewLabel("10x5 Grid - Click cells to assign GPU outputs, drag to assign touch areas"),
		nil, nil, nil,
		c.grid,
	)

	// Split view
	split := container.NewHSplit(controls, gridContainer)
	split.SetOffset(0.3) // 30% for controls, 70% for grid

	return split
}

// createActionButtons creates the action buttons (Generate, Save, etc.).
func (c *Creator) createActionButtons(window fyne.Window) fyne.CanvasObject {
	generateBtn := widget.NewButton("Generate screen.xml", func() {
		c.generateAndPreview(window)
	})

	saveBtn := widget.NewButton("Save screen.xml", func() {
		c.saveScreenXML(window)
	})

	detectBtn := widget.NewButton("Auto-detect mt-canvus.ini", func() {
		c.autoDetectConfig(window)
	})

	return container.NewVBox(
		generateBtn,
		saveBtn,
		detectBtn,
	)
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

		// Offer to update mt-canvus.ini
		if c.iniIntegration.ShouldUpdateIni(c.grid) {
			c.offerUpdateIni(window)
		} else {
			dialog.ShowInformation("Success", "screen.xml saved successfully", window)
		}
	}, window)
}

// offerUpdateIni offers to update mt-canvus.ini with video-output.
func (c *Creator) offerUpdateIni(window fyne.Window) {
	videoOutputs := c.iniIntegration.DetectVideoOutputs(c.grid)
	config := c.iniIntegration.GenerateVideoOutputConfig(videoOutputs)

	content := widget.NewRichTextFromMarkdown(fmt.Sprintf(`
**Update mt-canvus.ini?**

Detected video outputs: %s

Would you like to update mt-canvus.ini with:
video-output=%s
`, strings.Join(videoOutputs, ", "), config))

	updateBtn := widget.NewButton("Update", func() {
		iniPath := c.fileService.DetectMtCanvusIni()
		if iniPath == "" {
			dialog.ShowError(fmt.Errorf("mt-canvus.ini not found"), window)
			return
		}

		if err := c.iniIntegration.UpdateMtCanvusIni(iniPath, videoOutputs); err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowInformation("Success", "mt-canvus.ini updated successfully", window)
	})

	cancelBtn := widget.NewButton("Cancel", func() {})

	dialog.ShowCustom("Update mt-canvus.ini", "Cancel", container.NewVBox(
		content,
		container.NewHBox(updateBtn, cancelBtn),
	), window)
}

// autoDetectConfig attempts to auto-detect mt-canvus.ini.
func (c *Creator) autoDetectConfig(window fyne.Window) {
	iniPath := c.fileService.DetectMtCanvusIni()
	if iniPath == "" {
		dialog.ShowInformation("Not Found", "mt-canvus.ini not found in standard locations", window)
	} else {
		dialog.ShowInformation("Found", fmt.Sprintf("Found mt-canvus.ini at:\n%s", iniPath), window)
	}
}
