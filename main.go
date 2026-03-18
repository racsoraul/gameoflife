package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gameoflife/mcell"
)

func main() {
	var pattern *Pattern
	var err error

	if len(os.Args) > 1 {
		path := os.Args[1]
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".mc":
			pattern, err = loadMacrocell(path)
		default:
			pattern, err = LoadPattern(path)
		}
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Default to first catalog entry.
		pattern, err = LoadCatalogEntry(&catalog[0])
		if err != nil {
			log.Fatal(err)
		}
	}

	gameOfLife, err := NewGame("Conway's Game of Life", 800, 800, 10)
	if err != nil {
		log.Fatal(err)
	}
	// If using a default catalog pattern, sync the catalog index.
	if len(os.Args) <= 1 {
		gameOfLife.catalogIndex = 0
	}
	if err = gameOfLife.Run(pattern); err != nil {
		log.Fatal(err)
	}
}

// loadMacrocell loads a Golly Macrocell (.mc) file into a Pattern.
func loadMacrocell(path string) (*Pattern, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cells, err := mcell.Parse(f)
	if err != nil {
		return nil, err
	}

	// Width/Height are 0 so loadPattern's centering offset is zero —
	// mcell.Parse already centers coordinates around (0,0).
	return &Pattern{
		Rule:  "B3/S23",
		Cells: cells,
	}, nil
}
