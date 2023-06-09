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
