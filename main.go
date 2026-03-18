package main

import (
	"log"
	"os"
)

func main() {
	var pattern *Pattern
	var err error

	if len(os.Args) > 1 {
		// Load a pattern from a file argument.
		pattern, err = LoadPattern(os.Args[1])
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
