package screenxml

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// GPUAssignment handles GPU output assignment workflow.
type GPUAssignment struct {
	grid        *GridWidget
	gpuInput    *widget.Entry
	outputInput *widget.Entry
	statusLabel *widget.Label
}

// NewGPUAssignment creates a new GPU assignment handler.
func NewGPUAssignment(grid *GridWidget) *GPUAssignment {
	return &GPUAssignment{
		grid: grid,
	}
}

// CreateUI creates the UI for GPU assignment.
func (ga *GPUAssignment) CreateUI() fyne.CanvasObject {
	// Instructions label
	instructions := widget.NewRichTextFromMarkdown(`
**GPU Output Assignment Instructions:**
1. Draw a large number on each physical screen
2. Enter the GPU number (e.g., 0, 1, 2)
3. Enter the output number (e.g., 1, 2, 3)
4. Click the corresponding grid cell to assign
5. Click again to remove assignment
`)

	// GPU and output input fields
	gpuLabel := widget.NewLabel("GPU Number:")
	ga.gpuInput = widget.NewEntry()
	ga.gpuInput.SetPlaceHolder("0")

	outputLabel := widget.NewLabel("Output Number:")
	ga.outputInput = widget.NewEntry()
	ga.outputInput.SetPlaceHolder("1")

	// Assign button
	assignBtn := widget.NewButton("Assign to Clicked Cell", func() {
		ga.enableAssignmentMode()
	})

	// Status label
	ga.statusLabel = widget.NewLabel("Ready - Click a cell to assign GPU output")

	// Layout
	form := container.NewVBox(
		instructions,
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			gpuLabel, ga.gpuInput,
			outputLabel, ga.outputInput,
		),
		assignBtn,
		ga.statusLabel,
	)

	return form
}

// enableAssignmentMode enables click-to-assign mode.
func (ga *GPUAssignment) enableAssignmentMode() {
	gpuStr := strings.TrimSpace(ga.gpuInput.Text)
	outputStr := strings.TrimSpace(ga.outputInput.Text)

	if gpuStr == "" || outputStr == "" {
		ga.statusLabel.SetText("Error: Please enter both GPU and output numbers")
		return
	}

	gpuNum, err := strconv.Atoi(gpuStr)
	if err != nil || gpuNum < 0 {
		ga.statusLabel.SetText("Error: Invalid GPU number")
		return
	}

	outputNum, err := strconv.Atoi(outputStr)
	if err != nil || outputNum < 0 {
		ga.statusLabel.SetText("Error: Invalid output number")
		return
	}

	gpuOutput := fmt.Sprintf("%d.%d", gpuNum, outputNum)
	ga.statusLabel.SetText(fmt.Sprintf("Assignment mode: %s - Click a cell to assign or remove", gpuOutput))

	// Set up click handler
	ga.grid.SetOnCellClick(func(row, col int) {
		cell := ga.grid.GetCell(row, col)
		if cell == nil {
			return
		}

		// Toggle assignment: if already assigned to this GPU output, remove it
		if cell.GPUOutput == gpuOutput {
			ga.grid.ClearCellGPUOutput(row, col)
			ga.statusLabel.SetText(fmt.Sprintf("Removed %s from cell (%d,%d)", gpuOutput, row, col))
		} else {
			ga.grid.SetCellGPUOutput(row, col, gpuOutput)
			ga.statusLabel.SetText(fmt.Sprintf("Assigned %s to cell (%d,%d)", gpuOutput, row, col))
		}
	})
}

// GetGPUOutputs returns all assigned GPU outputs as a map of cell coordinates to GPU output strings.
func (ga *GPUAssignment) GetGPUOutputs() map[string]string {
	outputs := make(map[string]string)
	for row := 0; row < GridRows; row++ {
		for col := 0; col < GridCols; col++ {
			cell := ga.grid.GetCell(row, col)
			if cell != nil && cell.GPUOutput != "" {
				key := fmt.Sprintf("%d,%d", row, col)
				outputs[key] = cell.GPUOutput
			}
		}
	}
	return outputs
}
