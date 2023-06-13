package main

// RuleB3S23 Implements the Life rule B3/S23.
func (g *Game) RuleB3S23(x, y int32) {
	cellState := g.frameBuffer.GetCellState(x, y, false)
	neighbourStates := g.frameBuffer.GetCellNeighbourStates(x, y)
	var liveCellCount uint8
	for i := 0; i < 8; i++ {
		liveCellCount += uint8(neighbourStates[i])
	}
	switch cellState {
	case ALIVE:
		if liveCellCount == 2 || liveCellCount == 3 {
			g.frameBuffer.SetCellState(ALIVE, x, y, true)
			return
		}
	case DEAD:
		if liveCellCount == 3 {
			g.frameBuffer.SetCellState(ALIVE, x, y, true)
			return
		}
	}
	g.frameBuffer.SetCellState(DEAD, x, y, true)
}
