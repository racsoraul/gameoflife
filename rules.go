package main

// directions holds the 8 neighbor offsets for the Moore neighbourhood.
var directions = [8][2]int32{
	{-1, -1}, {0, -1}, {1, -1},
	{-1, 0}, {1, 0},
	{-1, 1}, {0, 1}, {1, 1},
}

// RuleB3S23 Implements the Life rule B3/S23.
func (g *Game) RuleB3S23() {
	// Count live neighbors for every cell that could change the state.
	neighbourCount := make(map[[2]int32]int)
	for pos := range g.cells {
		for _, d := range directions {
			n := [2]int32{pos[0] + d[0], pos[1] + d[1]}
			neighbourCount[n]++
		}
	}

	nextCells := make(map[[2]int32]uint32)
	for pos, count := range neighbourCount {
		age := g.cells[pos]
		if age > 0 && (count == 2 || count == 3) {
			nextCells[pos] = age + 1 // Survived: increment age.
		} else if age == 0 && count == 3 {
			nextCells[pos] = 1 // Born: age starts at 1.
		}
	}

	// Track dead cells for fade trails.
	if g.trailMode {
		const trailFrames = 8
		for pos := range g.cells {
			if nextCells[pos] == 0 {
				g.trails[pos] = trailFrames
			}
		}
	}

	g.cells = nextCells
}
