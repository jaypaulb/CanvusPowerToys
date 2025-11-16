package screenxml

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
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
	indexSpin   fyne.CanvasObject // Container with entry and up/down buttons
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

	// Layout Index: number input with up/down buttons (0-3, step 1)
	indexEntry := widget.NewEntry()
	indexEntry.SetText(cell.Index)
	indexEntry.Validator = func(text string) error {
		if text == "" {
			return nil
		}
		val, err := strconv.Atoi(text)
		if err != nil {
			return fmt.Errorf("must be a number")
		}
		if val < 0 || val > 3 {
			return fmt.Errorf("must be between 0 and 3")
		}
		return nil
	}
	indexEntry.OnChanged = func(text string) {
		if text != "" {
			cw.grid.SetCellIndex(cw.row, cw.col, text)
			cw.grid.Refresh()
		}
	}

	// Up/Down buttons for index
	upBtn := widget.NewButton("▲", func() {
		current := 0
		if cell.Index != "" {
			if idx, err := strconv.Atoi(cell.Index); err == nil {
				current = idx
			}
		}
		if current < 3 {
			current++
			indexEntry.SetText(strconv.Itoa(current))
		}
	})
	downBtn := widget.NewButton("▼", func() {
		current := 0
		if cell.Index != "" {
			if idx, err := strconv.Atoi(cell.Index); err == nil {
				current = idx
			}
		}
		if current > 0 {
			current--
			indexEntry.SetText(strconv.Itoa(current))
		}
	})
	indexContainer := container.NewBorder(nil, nil, nil, container.NewVBox(upBtn, downBtn), indexEntry)
	cw.indexSpin = indexContainer

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

	// Set index to 1
	cw.grid.SetCellIndex(cw.row, cw.col, "1")
	// Update index entry if it exists
	if indexContainer, ok := cw.indexSpin.(*fyne.Container); ok {
		for _, obj := range indexContainer.Objects {
			if entry, ok := obj.(*widget.Entry); ok {
				entry.SetText("1")
				break
			}
		}
	}

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

// CreateRenderer creates the renderer for the cell widget.
func (cw *CellWidget) CreateRenderer() fyne.WidgetRenderer {
	// Col1: Labels
	gpuLabel := widget.NewLabel("GPU.Output:")
	resLabel := widget.NewLabel("Resolution:")
	indexLabel := widget.NewLabel("Layer:")

	// Col2: Controls
	col1 := container.NewVBox(gpuLabel, resLabel, indexLabel)
	col2 := container.NewVBox(cw.gpuSelect, cw.resSelect, cw.indexSpin)

	// Main layout: 2 cols, 4 rows
	// Use VBox for rows, then GridWithColumns for cols
	row1 := container.NewGridWithColumns(2, col1, col2)
	row2 := widget.NewLabel("") // Empty
	row3 := widget.NewLabel("") // Empty
	// Row 4: Auto Fill button spans full width (both columns)
	// Put button directly in VBox - it will take full width
	row4 := cw.autoFillBtn

	content := container.NewVBox(row1, row2, row3, row4)

	return &cellRenderer{
		cell:    cw,
		content: content,
	}
}

// cellRenderer renders the cell widget.
type cellRenderer struct {
	cell    *CellWidget
	content *fyne.Container
}

func (r *cellRenderer) Layout(size fyne.Size) {
	r.content.Resize(size)
}

func (r *cellRenderer) MinSize() fyne.Size {
	return r.content.MinSize()
}

func (r *cellRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
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
		if cell.Index != "" {
			if indexContainer, ok := r.cell.indexSpin.(*fyne.Container); ok {
				for _, obj := range indexContainer.Objects {
					if entry, ok := obj.(*widget.Entry); ok {
						if entry.Text != cell.Index {
							entry.SetText(cell.Index)
						}
						break
					}
				}
			}
		}
	}
}

func (r *cellRenderer) Destroy() {
}

