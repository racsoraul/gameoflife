package main

// translateAliveCells Translates every alive cell along every row of the grid.
func (g *Game) translateAliveCells(x, y int32) {
	state := g.frameBuffer.GetCellState(x, y, false)
	if state == ALIVE {
		if x+1 < g.width/g.CellSize {
			g.frameBuffer.SetCellState(ALIVE, x+1, y, true)
		} else if y+1 < g.height/g.CellSize {
			g.frameBuffer.SetCellState(ALIVE, 0, y+1, true)
		} else {
			g.frameBuffer.SetCellState(ALIVE, 0, 0, true)
		}
	}
}

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
