package main

import "log"

func main() {
	gameOfLife, err := NewGame("Conway's Game of Life", 800, 800)
	if err != nil {
		log.Fatal(err)
	}

	gameOfLife.CellSize = 6

	err = gameOfLife.Run()
	if err != nil {
		log.Fatal(err)
	}
}
