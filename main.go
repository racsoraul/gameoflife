package main

import "log"

func main() {
	gameOfLife, err := NewGame("Conway's Game of Life", 500, 500, 10)
	if err != nil {
		log.Fatal(err)
	}

	err = gameOfLife.Run()
	if err != nil {
		log.Fatal(err)
	}
}
