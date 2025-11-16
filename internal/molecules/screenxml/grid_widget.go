package screenxml

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

const (
	GridCols = 6
	GridRows = 4
)

// CellState represents the state of a grid cell.
type CellState struct {
	GPUOutput    string     // Format: "gpu#:output#" (e.g., "1:2")
	TouchArea    int        // Touch area index (-1 if not assigned)
	IsLayoutArea bool       // True if part of layout area (pink frame)
	Resolution   Resolution // Resolution for this cell (default: 1920x1080)
	Index        string     // Index string (empty by default)
}

// GridWidget is a custom widget for displaying a 10x5 grid.
type GridWidget struct {
	widget.BaseWidget
	cells       [GridRows][GridCols]*CellState
	onCellClick func(row, col int)
	onCellDrag  func(startRow, startCol, endRow, endCol int)
	dragStart   *fyne.Position
	isDragging  bool
}

// NewGridWidget creates a new grid widget.
func NewGridWidget() *GridWidget {
	grid := &GridWidget{
		cells: [GridRows][GridCols]*CellState{},
	}

	// Initialize all cells
	for row := 0; row < GridRows; row++ {
		for col := 0; col < GridCols; col++ {
			grid.cells[row][col] = &CellState{
				TouchArea:    -1,
				IsLayoutArea: false,
				Resolution:   Resolution{Width: 1920, Height: 1080, Name: "1920x1080 (Full HD)"},
				Index:        "",
			}
		}
	}

	grid.ExtendBaseWidget(grid)
	return grid
}

// SetOnCellClick sets the callback for cell click events.
func (g *GridWidget) SetOnCellClick(fn func(row, col int)) {
	g.onCellClick = fn
}

// SetOnCellDrag sets the callback for cell drag events.
func (g *GridWidget) SetOnCellDrag(fn func(startRow, startCol, endRow, endCol int)) {
	g.onCellDrag = fn
}

// GetCell returns the state of a cell.
func (g *GridWidget) GetCell(row, col int) *CellState {
	if row < 0 || row >= GridRows || col < 0 || col >= GridCols {
		return nil
	}
	return g.cells[row][col]
}

// SetCellGPUOutput sets the GPU output for a cell.
func (g *GridWidget) SetCellGPUOutput(row, col int, gpuOutput string) {
	if row < 0 || row >= GridRows || col < 0 || col >= GridCols {
		return
	}
	g.cells[row][col].GPUOutput = gpuOutput
	g.Refresh()
}

// ClearCellGPUOutput clears the GPU output for a cell.
func (g *GridWidget) ClearCellGPUOutput(row, col int) {
	if row < 0 || row >= GridRows || col < 0 || col >= GridCols {
		return
	}
	g.cells[row][col].GPUOutput = ""
	g.Refresh()
}

// SetCellResolution sets the resolution for a cell.
func (g *GridWidget) SetCellResolution(row, col int, res Resolution) {
	if row < 0 || row >= GridRows || col < 0 || col >= GridCols {
		return
	}
	g.cells[row][col].Resolution = res
	g.Refresh()
}

// SetCellIndex sets the index for a cell.
func (g *GridWidget) SetCellIndex(row, col int, index string) {
	if row < 0 || row >= GridRows || col < 0 || col >= GridCols {
		return
	}
	g.cells[row][col].Index = index
	g.Refresh()
}

// Tapped handles tap/click events on the grid.
func (g *GridWidget) Tapped(e *fyne.PointEvent) {
	row, col := g.getCellFromPosition(e.Position)
	if row >= 0 && col >= 0 && g.onCellClick != nil {
		g.onCellClick(row, col)
	}
}

// MouseDown handles mouse down events for drag start.
func (g *GridWidget) MouseDown(e *fyne.PointEvent) {
	row, col := g.getCellFromPosition(e.Position)
	if row >= 0 && col >= 0 {
		g.dragStart = &e.Position
		g.isDragging = false
	}
}

// MouseUp handles mouse up events for drag end.
func (g *GridWidget) MouseUp(e *fyne.PointEvent) {
	if g.dragStart != nil && g.isDragging && g.onCellDrag != nil {
		startRow, startCol := g.getCellFromPosition(*g.dragStart)
		endRow, endCol := g.getCellFromPosition(e.Position)
		if startRow >= 0 && startCol >= 0 && endRow >= 0 && endCol >= 0 {
			g.onCellDrag(startRow, startCol, endRow, endCol)
		}
	}
	g.dragStart = nil
	g.isDragging = false
}

// MouseDragged handles mouse drag events.
func (g *GridWidget) MouseDragged(e *fyne.DragEvent) {
	if g.dragStart != nil {
		g.isDragging = true
		g.Refresh()
	}
}

// getCellFromPosition converts a screen position to grid cell coordinates.
func (g *GridWidget) getCellFromPosition(pos fyne.Position) (row, col int) {
	size := g.Size()
	cellWidth := size.Width / GridCols
	cellHeight := size.Height / GridRows

	col = int(pos.X / cellWidth)
	row = int(pos.Y / cellHeight)

	if row < 0 || row >= GridRows || col < 0 || col >= GridCols {
		return -1, -1
	}
	return row, col
}

// CreateRenderer creates the renderer for the grid widget.
func (g *GridWidget) CreateRenderer() fyne.WidgetRenderer {
	return &gridRenderer{
		grid:   g,
		cells:  make([][]fyne.CanvasObject, GridRows),
		border: canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255}),
	}
}

// gridRenderer renders the grid widget.
type gridRenderer struct {
	grid   *GridWidget
	cells  [][]fyne.CanvasObject
	border *canvas.Rectangle
}

func (r *gridRenderer) Layout(size fyne.Size) {
	// Add spacing between cells (2px gap)
	spacing := float32(2)
	cellWidth := (size.Width - float32(GridCols-1)*spacing) / GridCols
	cellHeight := (size.Height - float32(GridRows-1)*spacing) / GridRows

	for row := 0; row < GridRows; row++ {
		if r.cells[row] == nil {
			r.cells[row] = make([]fyne.CanvasObject, GridCols)
		}
		for col := 0; col < GridCols; col++ {
			if r.cells[row][col] == nil {
				rect := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 255})
				rect.StrokeColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
				rect.StrokeWidth = 1
				r.cells[row][col] = rect
			}
			// Calculate position with spacing
			x := float32(col) * (cellWidth + spacing)
			y := float32(row) * (cellHeight + spacing)
			r.cells[row][col].Resize(fyne.NewSize(cellWidth, cellHeight))
			r.cells[row][col].Move(fyne.NewPos(x, y))
		}
	}

	r.border.Resize(size)
	r.border.Move(fyne.NewPos(0, 0))
}

func (r *gridRenderer) MinSize() fyne.Size {
	return fyne.NewSize(600, 300) // Minimum size for 10x5 grid
}

func (r *gridRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.border}
	for row := 0; row < GridRows; row++ {
		if r.cells[row] == nil {
			continue
		}
		for col := 0; col < GridCols; col++ {
			if r.cells[row][col] != nil {
				objects = append(objects, r.cells[row][col])
			}
		}
	}
	return objects
}

func (r *gridRenderer) Refresh() {
	cellState := r.grid.cells
	for row := 0; row < GridRows; row++ {
		// Ensure row slice is initialized
		if r.cells[row] == nil {
			r.cells[row] = make([]fyne.CanvasObject, GridCols)
		}
		for col := 0; col < GridCols; col++ {
			// Ensure cell is initialized
			if r.cells[row][col] == nil {
				rect := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 255})
				rect.StrokeColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
				rect.StrokeWidth = 1
				r.cells[row][col] = rect
			}
			if r.cells[row][col] != nil {
				rect := r.cells[row][col].(*canvas.Rectangle)
				state := cellState[row][col]

				// Set cell color based on state
				if state.GPUOutput != "" {
					// Cell has GPU output assigned - light blue
					rect.FillColor = color.RGBA{R: 173, G: 216, B: 230, A: 255}
				} else if state.IsLayoutArea {
					// Layout area - light pink border
					rect.FillColor = color.RGBA{R: 255, G: 240, B: 245, A: 255}
					rect.StrokeColor = color.RGBA{R: 255, G: 192, B: 203, A: 255}
					rect.StrokeWidth = 2
				} else {
					// Default - white
					rect.FillColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
					rect.StrokeColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
					rect.StrokeWidth = 1
				}
			}
		}
	}
	canvas.Refresh(r.grid)
}

func (r *gridRenderer) Destroy() {
}

