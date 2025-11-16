package screenxml

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// GridContainer is a container that displays a grid of cell widgets with drag selection.
type GridContainer struct {
	grid        *GridWidget
	cellWidgets [GridRows][GridCols]*CellWidget
	container   *fyne.Container
	dragStart   *fyne.Position
	dragEnd     *fyne.Position
	isDragging  bool
	selectedCells map[string]bool // Map of "row:col" to track selected cells
}

// NewGridContainer creates a new grid container with cell widgets.
func NewGridContainer(grid *GridWidget) *GridContainer {
	gc := &GridContainer{
		grid:          grid,
		cellWidgets:   [GridRows][GridCols]*CellWidget{},
		selectedCells: make(map[string]bool),
	}

	// Create cell widgets
	cellContainer := container.NewGridWithColumns(GridCols)
	for row := 0; row < GridRows; row++ {
		for col := 0; col < GridCols; col++ {
			cellWidget := NewCellWidget(grid, row, col)
			gc.cellWidgets[row][col] = cellWidget
			cellContainer.Add(cellWidget)
		}
	}

	gc.container = cellContainer
	return gc
}

// GetContainer returns the container for this grid.
func (gc *GridContainer) GetContainer() *fyne.Container {
	return gc.container
}

// HandleDragStart handles the start of a drag operation.
func (gc *GridContainer) HandleDragStart(pos fyne.Position) {
	// Check if click is on whitespace (not on a form control)
	// For now, we'll assume any click starts drag - form controls will handle their own events
	gc.dragStart = &pos
	gc.isDragging = false
	gc.selectedCells = make(map[string]bool)
}

// HandleDragUpdate handles drag updates.
func (gc *GridContainer) HandleDragUpdate(pos fyne.Position) {
	if gc.dragStart == nil {
		return
	}

	gc.isDragging = true
	gc.dragEnd = &pos

	// Calculate which cells are in the drag rectangle
	gc.updateSelectedCells()
}

// HandleDragEnd handles the end of a drag operation.
func (gc *GridContainer) HandleDragEnd() {
	if gc.dragStart == nil || !gc.isDragging {
		gc.dragStart = nil
		gc.dragEnd = nil
		gc.isDragging = false
		return
	}

	// Assign all selected cells to lowest available layout index
	gc.assignSelectedCellsToIndex()

	gc.dragStart = nil
	gc.dragEnd = nil
	gc.isDragging = false
	gc.selectedCells = make(map[string]bool)
}

// updateSelectedCells updates which cells are selected based on drag rectangle.
func (gc *GridContainer) updateSelectedCells() {
	if gc.dragStart == nil || gc.dragEnd == nil {
		return
	}

	// Convert positions to cell coordinates
	startRow, startCol := gc.positionToCell(*gc.dragStart)
	endRow, endCol := gc.positionToCell(*gc.dragEnd)

	// Ensure start is top-left, end is bottom-right
	if startRow > endRow {
		startRow, endRow = endRow, startRow
	}
	if startCol > endCol {
		startCol, endCol = endCol, startCol
	}

	// Mark all cells in rectangle as selected
	gc.selectedCells = make(map[string]bool)
	for row := startRow; row <= endRow; row++ {
		for col := startCol; col <= endCol; col++ {
			if row >= 0 && row < GridRows && col >= 0 && col < GridCols {
				key := fmt.Sprintf("%d:%d", row, col)
				gc.selectedCells[key] = true
			}
		}
	}
}

// positionToCell converts a position to cell coordinates.
// Accounts for spacing between cells (2px gap).
func (gc *GridContainer) positionToCell(pos fyne.Position) (row, col int) {
	size := gc.container.Size()
	if size.Width == 0 || size.Height == 0 {
		// Container not laid out yet, return invalid
		return -1, -1
	}

	// Account for spacing (2px gap between cells)
	spacing := float32(2)
	cellWidth := (size.Width - float32(GridCols-1)*spacing) / float32(GridCols)
	cellHeight := (size.Height - float32(GridRows-1)*spacing) / float32(GridRows)

	// Calculate which cell the position is in, accounting for spacing
	col = int(pos.X / (cellWidth + spacing))
	row = int(pos.Y / (cellHeight + spacing))

	// Clamp to valid range
	if row < 0 {
		row = 0
	} else if row >= GridRows {
		row = GridRows - 1
	}
	if col < 0 {
		col = 0
	} else if col >= GridCols {
		col = GridCols - 1
	}

	return row, col
}

// assignSelectedCellsToIndex assigns all selected cells to the lowest available layout index.
func (gc *GridContainer) assignSelectedCellsToIndex() {
	if len(gc.selectedCells) == 0 {
		return
	}

	// Find lowest available index (0-3)
	usedIndices := make(map[int]bool)
	for row := 0; row < GridRows; row++ {
		for col := 0; col < GridCols; col++ {
			cell := gc.grid.GetCell(row, col)
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
	lowestIndex := 0
	for i := 0; i <= 3; i++ {
		if !usedIndices[i] {
			lowestIndex = i
			break
		}
		// If all indices 0-3 are used, use 0 (wrap around)
		if i == 3 {
			lowestIndex = 0
		}
	}

	// Assign selected cells to this index
	for key := range gc.selectedCells {
		var row, col int
		fmt.Sscanf(key, "%d:%d", &row, &col)
		if row >= 0 && row < GridRows && col >= 0 && col < GridCols {
			gc.grid.SetCellIndex(row, col, strconv.Itoa(lowestIndex))
			// Update cell widget's layer checkbox
			if gc.cellWidgets[row][col] != nil {
				gc.cellWidgets[row][col].Refresh()
			}
		}
	}

	gc.grid.Refresh()
}

