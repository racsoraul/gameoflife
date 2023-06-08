package main

import (
	"errors"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// Game Holds the configs and state of the Conway's Game of Life.
type Game struct {
	running  bool
	title    string
	width    int32
	height   int32
	window   *sdl.Window
	renderer *sdl.Renderer
}

// NewGame Returns a new initialized game.
func NewGame(title string, width, height int32) (*Game, error) {
	game := &Game{
		title:  title,
		width:  width,
		height: height,
	}
	err := game.init()
	if err != nil {
		return nil, fmt.Errorf("failed to create new game: %w", err)
	}
	return game, nil
}

// Run the game. Calling this will block until exiting the game.
func (g *Game) Run() error {
	g.running = true
	for g.running {
		g.processInput()
		g.update()
		g.render()
	}
	return g.shutdown()
}

// init Initialize game resources.
func (g *Game) init() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("failed to initialize subsystems: %w", err)
	}

	window, err := sdl.CreateWindow(
		g.title,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		g.width,
		g.height,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return fmt.Errorf("failed to create Window: %w", err)
	}
	g.window = window

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return fmt.Errorf("failed to create renderer: %w", err)
	}
	g.renderer = renderer

	return nil
}

// shutdown Clean up, free and destroy resources.
func (g *Game) shutdown() error {
	rendErr := g.renderer.Destroy()
	winError := g.window.Destroy()
	errs := errors.Join(rendErr, winError)
	sdl.Quit()
	return errs
}

// processInput Processes the system events from the event queue.
func (g *Game) processInput() {
	for nextEvent := sdl.PollEvent(); nextEvent != nil; nextEvent = sdl.PollEvent() {
		switch event := nextEvent.(type) {
		case *sdl.QuitEvent:
			g.running = false
		case *sdl.KeyboardEvent:
			if event.Keysym.Sym == sdl.K_ESCAPE {
				g.running = false
			}
		}
	}
}

// update game state.
func (g *Game) update() {}

// render the game state to screen.
func (g *Game) render() {
	g.renderer.Present()
}
