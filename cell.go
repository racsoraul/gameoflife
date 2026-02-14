package main

// CellState It's either Dead or Alive.
type CellState uint8

const (
	DEAD CellState = iota
	ALIVE
)

// SetCellState Sets the specified state for the provided cell. If next is false, it uses the current buffer.
func (fb *FrameBuffer) SetCellState(state CellState, x, y int32, next bool) {
	var csColor uint32
	switch state {
	case ALIVE:
		csColor = fb.g.CellAliveColor
	case DEAD:
		csColor = fb.g.CellDeadColor
	}
	fb.DrawRect(x*fb.g.cellSize, y*fb.g.cellSize, fb.g.cellSize, fb.g.cellSize, csColor, next)
}

// GetCellState Get cell's state. If next is false, it uses the current buffer.
func (fb *FrameBuffer) GetCellState(x, y int32, next bool) CellState {
	pixelColor := fb.GetPixelColor(
		(x*fb.g.cellSize)+(fb.g.cellSize/2),
		(y*fb.g.cellSize)+(fb.g.cellSize/2),
		next,
	) // Get the color of the pixel at the center of the cell.
	switch pixelColor {
	case fb.g.CellAliveColor:
		return ALIVE
	case fb.g.CellDeadColor:
		return DEAD
	}
	return DEAD
}

// GetCellPosFromWindowCoords Returns cell's position from the given relative window's coordinates.
func (fb *FrameBuffer) GetCellPosFromWindowCoords(winX, winY int32) (int32, int32) {
	return winX / fb.g.cellSize, winY / fb.g.cellSize
}

// ToggleCellState Toggles cell's state located at the window's coordinates in current color buffer.
// If next is false, it uses the current buffer.
func (fb *FrameBuffer) ToggleCellState(winX, winY int32, next bool) {
	x, y := fb.GetCellPosFromWindowCoords(winX, winY)
	cellState := fb.GetCellState(x, y, next)
	switch cellState {
	case ALIVE:
		fb.SetCellState(DEAD, x, y, next)
	case DEAD:
		fb.SetCellState(ALIVE, x, y, next)
	}
}

// GetCellNeighbourStates Get the state of the eight surrounding cells from current color buffer.
func (fb *FrameBuffer) GetCellNeighbourStates(x, y int32) [8]CellState {
	states := [8]CellState{}
	// Left.
	if x-1 >= 0 {
		states[0] = fb.GetCellState(x-1, y, false)
	}
	// Right.
	if x+1 < fb.g.width/fb.g.cellSize {
		states[1] = fb.GetCellState(x+1, y, false)
	}
	// Top
	if y-1 >= 0 {
		states[2] = fb.GetCellState(x, y-1, false)
	}
	// Bottom.
	if y+1 < fb.g.height/fb.g.cellSize {
		states[3] = fb.GetCellState(x, y+1, false)
	}
	// Top-Left
	if x-1 >= 0 && y-1 >= 0 {
		states[4] = fb.GetCellState(x-1, y-1, false)
	}
	// Top-Right
	if x+1 < fb.g.width/fb.g.cellSize && y-1 >= 0 {
		states[5] = fb.GetCellState(x+1, y-1, false)
	}
	// Bottom-Left
	if x-1 >= 0 && y+1 < fb.g.height/fb.g.cellSize {
		states[6] = fb.GetCellState(x-1, y+1, false)
	}
	// Bottom-Right
	if x+1 < fb.g.width/fb.g.cellSize && y+1 < fb.g.height/fb.g.cellSize {
		states[7] = fb.GetCellState(x+1, y+1, false)
	}
	return states
}
