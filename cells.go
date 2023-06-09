package main

type CellState uint8

const (
	DEAD CellState = iota
	ALIVE
)

// setCellState Sets the state of the specified cell.
func (g *Game) setCellState(state CellState, x, y int32) {
	switch state {
	case ALIVE:
		g.drawRect(x*g.CellSize, y*g.CellSize, g.CellSize, g.CellSize, g.CellAliveColor)
	case DEAD:
		g.drawRect(x*g.CellSize, y*g.CellSize, g.CellSize, g.CellSize, g.CellDeadColor)
	}
}
