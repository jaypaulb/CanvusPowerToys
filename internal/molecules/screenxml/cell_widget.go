package screenxml

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CellWidget is a widget that displays a cell with embedded form controls.
type CellWidget struct {
	widget.BaseWidget
	grid        *GridWidget
	row, col    int
	gpuSelect   *widget.Select
	resSelect   *widget.Select
	layerCheck  *widget.Check // Checkbox for layer (in layer or not)
	autoFillBtn *widget.Button
	onChanged   func()
}

// NewCellWidget creates a new cell widget with embedded form.
func NewCellWidget(grid *GridWidget, row, col int) *CellWidget {
	cw := &CellWidget{
		grid: grid,
		row:  row,
		col:  col,
	}
	cw.ExtendBaseWidget(cw)
	cw.buildUI()
	return cw
}

// buildUI builds the UI for the cell (2 cols, 4 rows).
func (cw *CellWidget) buildUI() {
	cell := cw.grid.GetCell(cw.row, cw.col)
	if cell == nil {
		return
	}

	// Generate GPU output options (1:1, 1:2, ... up to 6:4)
	var gpuOptions []string
	gpuOptions = append(gpuOptions, "") // Empty option
	for gpu := 1; gpu <= 6; gpu++ {
		for output := 1; output <= 4; output++ {
			gpuOptions = append(gpuOptions, fmt.Sprintf("%d:%d", gpu, output))
		}
	}

	// GPU Output dropdown
	cw.gpuSelect = widget.NewSelect(gpuOptions, func(selected string) {
		if selected == "" {
			cw.grid.ClearCellGPUOutput(cw.row, cw.col)
		} else {
			cw.grid.SetCellGPUOutput(cw.row, cw.col, selected)
		}
		cw.grid.Refresh()
	})
	if cell.GPUOutput != "" {
		cw.gpuSelect.SetSelected(cell.GPUOutput)
	}

	// Resolution dropdown
	resOptions := make([]string, len(CommonResolutions))
	for i, res := range CommonResolutions {
		resOptions[i] = res.Name
	}
	cw.resSelect = widget.NewSelect(resOptions, func(selected string) {
		for _, res := range CommonResolutions {
			if res.Name == selected {
				cw.grid.SetCellResolution(cw.row, cw.col, res)
				break
			}
		}
		cw.grid.Refresh()
	})
	currentResName := cell.Resolution.Name
	if currentResName == "" {
		currentResName = "1920x1080 (Full HD)"
	}
	cw.resSelect.SetSelected(currentResName)

	// Layer: checkbox (in layer or not)
	cw.layerCheck = widget.NewCheck("", func(checked bool) {
		if checked {
			// Assign to lowest available layout index
			lowestIndex := cw.findLowestAvailableIndex()
			cw.grid.SetCellIndex(cw.row, cw.col, strconv.Itoa(lowestIndex))
		} else {
			// Remove from layer
			cw.grid.SetCellIndex(cw.row, cw.col, "")
		}
		cw.grid.Refresh()
	})
	// Set checked if cell has an index (is in a layer)
	cw.layerCheck.SetChecked(cell.Index != "")

	// Auto Fill button (spans both cols in row 4)
	cw.autoFillBtn = widget.NewButton("Auto Fill", func() {
		cw.handleAutoFill()
	})
}

// handleAutoFill applies next GPU index, sets res to 1080p, sets index 1.
func (cw *CellWidget) handleAutoFill() {
	// Find next available GPU output
	nextGPU := cw.findNextGPUOutput()

	// Set GPU output
	cw.grid.SetCellGPUOutput(cw.row, cw.col, nextGPU)
	cw.gpuSelect.SetSelected(nextGPU)

	// Set resolution to 1080p
	for _, res := range CommonResolutions {
		if res.Name == "1920x1080 (Full HD)" {
			cw.grid.SetCellResolution(cw.row, cw.col, res)
			cw.resSelect.SetSelected(res.Name)
			break
		}
	}

	// Set index to 1 (check the layer checkbox)
	cw.grid.SetCellIndex(cw.row, cw.col, "1")
	cw.layerCheck.SetChecked(true)

	cw.grid.Refresh()
}

// findNextGPUOutput finds the next available GPU output.
func (cw *CellWidget) findNextGPUOutput() string {
	// Check all cells to find the highest assigned GPU output
	maxGpu := 0
	maxOutput := 0

	for row := 0; row < GridRows; row++ {
		for col := 0; col < GridCols; col++ {
			cell := cw.grid.GetCell(row, col)
			if cell != nil && cell.GPUOutput != "" {
				parts := strings.Split(cell.GPUOutput, ":")
				if len(parts) == 2 {
					var gpu, output int
					fmt.Sscanf(parts[0], "%d", &gpu)
					fmt.Sscanf(parts[1], "%d", &output)
					if gpu > maxGpu || (gpu == maxGpu && output > maxOutput) {
						maxGpu = gpu
						maxOutput = output
					}
				}
			}
		}
	}

	// Start from 1:1 if nothing assigned, otherwise continue from next
	if maxGpu == 0 {
		return "1:1"
	}

	// Increment output (1:1 -> 1:2 -> 1:3 -> 1:4 -> 2:1 -> 2:2, etc.)
	if maxOutput < 4 {
		return fmt.Sprintf("%d:%d", maxGpu, maxOutput+1)
	}
	return fmt.Sprintf("%d:1", maxGpu+1)
}

// findLowestAvailableIndex finds the lowest available layout index (0-3).
func (cw *CellWidget) findLowestAvailableIndex() int {
	usedIndices := make(map[int]bool)
	for row := 0; row < GridRows; row++ {
		for col := 0; col < GridCols; col++ {
			cell := cw.grid.GetCell(row, col)
			if cell != nil && cell.Index != "" {
				if idx, err := strconv.Atoi(cell.Index); err == nil {
					if idx >= 0 && idx <= 3 {
						usedIndices[idx] = true
					}
				}
			}
		}
	}

	// Find lowest available index (0-3)
	for i := 0; i <= 3; i++ {
		if !usedIndices[i] {
			return i
		}
	}
	// If all indices are used, return 0
	return 0
}

// CreateRenderer creates the renderer for the cell widget.
func (cw *CellWidget) CreateRenderer() fyne.WidgetRenderer {
	// Col1: Labels
	gpuLabel := widget.NewLabel("GPU.Output:")
	resLabel := widget.NewLabel("Resolution:")
	layerLabel := widget.NewLabel("Layer:")

	// Col2: Controls
	col1 := container.NewVBox(gpuLabel, resLabel, layerLabel)
	col2 := container.NewVBox(cw.gpuSelect, cw.resSelect, cw.layerCheck)

	// Row 1: Labels and controls
	row1 := container.NewGridWithColumns(2, col1, col2)
	
	// Auto Fill button immediately below Layer (spans both columns)
	row2 := cw.autoFillBtn

	content := container.NewVBox(row1, row2)

	// MT Blue border: #36A9E1 (RGB: 54, 169, 225)
	border := canvas.NewRectangle(color.RGBA{})
	border.StrokeColor = color.RGBA{R: 54, G: 169, B: 225, A: 255} // MT Blue
	border.StrokeWidth = 2
	border.FillColor = color.RGBA{R: 255, G: 255, B: 255, A: 0} // Transparent fill

	return &cellRenderer{
		cell:    cw,
		content: content,
		border:  border,
	}
}

// cellRenderer renders the cell widget.
type cellRenderer struct {
	cell    *CellWidget
	content *fyne.Container
	border *canvas.Rectangle
}

func (r *cellRenderer) Layout(size fyne.Size) {
	r.content.Resize(size)
	r.content.Move(fyne.NewPos(0, 0))

	// Border fills the entire cell
	r.border.Resize(size)
	r.border.Move(fyne.NewPos(0, 0))
}

func (r *cellRenderer) MinSize() fyne.Size {
	return r.content.MinSize()
}

func (r *cellRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.border, r.content}
}

func (r *cellRenderer) Refresh() {
	// Update controls based on cell state
	cell := r.cell.grid.GetCell(r.cell.row, r.cell.col)
	if cell != nil {
		if cell.GPUOutput != "" && r.cell.gpuSelect.Selected != cell.GPUOutput {
			r.cell.gpuSelect.SetSelected(cell.GPUOutput)
		}
		if cell.Resolution.Name != "" && r.cell.resSelect.Selected != cell.Resolution.Name {
			r.cell.resSelect.SetSelected(cell.Resolution.Name)
		}
		// Update layer checkbox based on cell index
		isInLayer := cell.Index != ""
		if r.cell.layerCheck.Checked != isInLayer {
			r.cell.layerCheck.SetChecked(isInLayer)
		}
	}
}

func (r *cellRenderer) Destroy() {
}

