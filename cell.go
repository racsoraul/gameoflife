package main

type CellState uint8

const (
	DEAD CellState = iota
	ALIVE
)

// setCellState Sets the specified state for the provided cell.
func (g *Game) setCellState(state CellState, x, y int32) {
	switch state {
	case ALIVE:
		g.drawRect(x*g.CellSize, y*g.CellSize, g.CellSize, g.CellSize, g.CellAliveColor)
	case DEAD:
		g.drawRect(x*g.CellSize, y*g.CellSize, g.CellSize, g.CellSize, g.CellDeadColor)
	}
}

// getCellState Get cell's state.
func (g *Game) getCellState(x, y int32) CellState {
	pixelColor := g.getPixelColor(
		(x*g.CellSize)+(g.CellSize/2),
		(y*g.CellSize)+(g.CellSize/2),
	) // Get the color of the pixel at the center of the cell.
	switch pixelColor {
	case g.CellAliveColor:
		return ALIVE
	case g.CellDeadColor:
		return DEAD
	}
	return DEAD
}

// getCellPosFromWindowCoords Returns cell's position located at the window's coordinates.
func (g *Game) getCellPosFromWindowCoords(winX, winY int32) (int32, int32) {
	return winX / g.CellSize, winY / g.CellSize
}

// toggleCellState Toggles cell's state located at the window's coordinates.
func (g *Game) toggleCellState(winX, winY int32) {
	x, y := g.getCellPosFromWindowCoords(winX, winY)
	cellState := g.getCellState(x, y)
	switch cellState {
	case ALIVE:
		g.setCellState(DEAD, x, y)
	case DEAD:
		g.setCellState(ALIVE, x, y)
	}
}

// getCellNeighbourStates Get the state of the eight surrounding cells.
func (g *Game) getCellNeighbourStates(x, y int32) [8]CellState {
	states := [8]CellState{}
	// Left.
	if x-1 >= 0 {
		states[0] = g.getCellState(x-1, y)
	}
	// Right.
	if x+1 < g.width/g.CellSize {
		states[1] = g.getCellState(x+1, y)
	}
	// Top
	if y-1 >= 0 {
		states[2] = g.getCellState(x, y-1)
	}
	// Bottom.
	if y+1 < g.height/g.CellSize {
		states[3] = g.getCellState(x, y+1)
	}
	// Top-Left
	if x-1 >= 0 && y-1 >= 0 {
		states[4] = g.getCellState(x-1, y-1)
	}
	// Top-Right
	if x+1 < g.width/g.CellSize && y-1 >= 0 {
		states[5] = g.getCellState(x+1, y-1)
	}
	// Bottom-Left
	if x-1 >= 0 && y+1 < g.height/g.CellSize {
		states[6] = g.getCellState(x-1, y+1)
	}
	// Bottom-Right
	if x+1 < g.width/g.CellSize && y+1 < g.height/g.CellSize {
		states[7] = g.getCellState(x+1, y+1)
	}
	return states
}
