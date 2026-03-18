package main

// directions holds the 8 neighbor offsets for the Moore neighbourhood.
var directions = [8][2]int32{
	{-1, -1}, {0, -1}, {1, -1},
	{-1, 0}, {1, 0},
	{-1, 1}, {0, 1}, {1, 1},
}

// RuleB3S23 Implements the Life rule B3/S23.
// Reuses g.neighbourCount and g.nextCells maps to avoid allocation per generation.
func (g *Game) RuleB3S23() {
	// Clear scratch maps (retains allocated bucket memory).
	clear(g.neighbourCount)
	clear(g.nextCells)

	// Count live neighbors for every cell that could change state.
	for pos := range g.cells {
		for _, d := range directions {
			n := [2]int32{pos[0] + d[0], pos[1] + d[1]}
			g.neighbourCount[n]++
		}
	}

	for pos, count := range g.neighbourCount {
		age := g.cells[pos]
		if age > 0 && (count == 2 || count == 3) {
			g.nextCells[pos] = age + 1 // Survived: increment age.
		} else if age == 0 && count == 3 {
			g.nextCells[pos] = 1 // Born: age starts at 1.
		}
	}

	// Track dead cells for fade trails.
	if g.trailMode {
		const trailFrames = 8
		for pos := range g.cells {
			if g.nextCells[pos] == 0 {
				g.trails[pos] = trailFrames
			}
		}
	}

	// Swap cells and nextCells so the current buffer becomes the scratch for next generation.
	g.cells, g.nextCells = g.nextCells, g.cells
}
