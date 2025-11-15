package screenxml

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// DragSelectionWidget is a widget that wraps the grid container and handles drag selection.
type DragSelectionWidget struct {
	widget.BaseWidget
	gridContainer *GridContainer
	dragStart     *fyne.Position
	dragEnd       *fyne.Position
	isDragging    bool
	selectionRect *canvas.Rectangle
}

// NewDragSelectionWidget creates a new drag selection widget.
func NewDragSelectionWidget(gridContainer *GridContainer) *DragSelectionWidget {
	ds := &DragSelectionWidget{
		gridContainer: gridContainer,
		selectionRect: canvas.NewRectangle(color.RGBA{R: 100, G: 150, B: 255, A: 100}),
	}
	ds.selectionRect.Hide()
	ds.ExtendBaseWidget(ds)
	return ds
}

// CreateRenderer creates the renderer for the drag selection widget.
func (ds *DragSelectionWidget) CreateRenderer() fyne.WidgetRenderer {
	return &dragSelectionRenderer{
		widget:         ds,
		gridContainer:  ds.gridContainer.GetContainer(),
		selectionRect:  ds.selectionRect,
	}
}

// MouseDown handles mouse down events.
func (ds *DragSelectionWidget) MouseDown(e *fyne.PointEvent) {
	// Check if click is on a form control (we'll detect this by checking if it's on a cell widget)
	// For now, we'll start drag on any mouse down - form controls will consume their own events
	ds.dragStart = &e.Position
	ds.isDragging = false
	ds.selectionRect.Hide()
	ds.Refresh()
}

// MouseUp handles mouse up events.
func (ds *DragSelectionWidget) MouseUp(e *fyne.PointEvent) {
	if ds.dragStart != nil && ds.isDragging {
		// Convert positions to container coordinates
		containerPos := ds.getContainerPosition(*ds.dragStart)
		endPos := ds.getContainerPosition(e.Position)

		// Update grid container with drag positions
		ds.gridContainer.HandleDragStart(containerPos)
		ds.gridContainer.HandleDragUpdate(endPos)
		ds.gridContainer.HandleDragEnd()
	}

	ds.dragStart = nil
	ds.dragEnd = nil
	ds.isDragging = false
	ds.selectionRect.Hide()
	ds.Refresh()
}

// MouseDragged handles mouse drag events.
func (ds *DragSelectionWidget) MouseDragged(e *fyne.DragEvent) {
	if ds.dragStart == nil {
		return
	}

	ds.isDragging = true
	ds.dragEnd = &e.Position

	// Update selection rectangle
	ds.updateSelectionRect()

	// Update grid container with current drag position
	containerPos := ds.getContainerPosition(e.Position)
	ds.gridContainer.HandleDragStart(ds.getContainerPosition(*ds.dragStart))
	ds.gridContainer.HandleDragUpdate(containerPos)

	ds.Refresh()
}

// getContainerPosition converts a widget-relative position to container-relative position.
func (ds *DragSelectionWidget) getContainerPosition(pos fyne.Position) fyne.Position {
	// The position is already relative to the widget, which contains the grid container
	// So we can use it directly
	return pos
}

// updateSelectionRect updates the selection rectangle based on drag start and end.
func (ds *DragSelectionWidget) updateSelectionRect() {
	if ds.dragStart == nil || ds.dragEnd == nil {
		ds.selectionRect.Hide()
		return
	}

	// Calculate rectangle bounds
	startX := ds.dragStart.X
	startY := ds.dragStart.Y
	endX := ds.dragEnd.X
	endY := ds.dragEnd.Y

	// Ensure start is top-left, end is bottom-right
	if startX > endX {
		startX, endX = endX, startX
	}
	if startY > endY {
		startY, endY = endY, startY
	}

	// Set rectangle position and size
	ds.selectionRect.Move(fyne.NewPos(startX, startY))
	ds.selectionRect.Resize(fyne.NewSize(endX-startX, endY-startY))
	ds.selectionRect.Show()
}

// dragSelectionRenderer renders the drag selection widget.
type dragSelectionRenderer struct {
	widget        *DragSelectionWidget
	gridContainer *fyne.Container
	selectionRect *canvas.Rectangle
}

func (r *dragSelectionRenderer) Layout(size fyne.Size) {
	// Layout grid container to fill the widget
	r.gridContainer.Resize(size)
	r.gridContainer.Move(fyne.NewPos(0, 0))

	// Selection rectangle is positioned dynamically
}

func (r *dragSelectionRenderer) MinSize() fyne.Size {
	return r.gridContainer.MinSize()
}

func (r *dragSelectionRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.gridContainer, r.selectionRect}
}

func (r *dragSelectionRenderer) Refresh() {
	// Update selection rectangle if dragging
	if r.widget.isDragging {
		r.widget.updateSelectionRect()
	}
	canvas.Refresh(r.widget)
}

func (r *dragSelectionRenderer) Destroy() {
}

