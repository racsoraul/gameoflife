package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

// Game Holds the configs and state of the Conway's Game of Life.
type Game struct {
	running        bool
	title          string
	width          int32
	height         int32
	window         *sdl.Window
	renderer       *sdl.Renderer
	frameBuffer    *FrameBuffer // Holds every generation of cells.
	CellSize       int32
	CellAliveColor uint32
	CellDeadColor  uint32
	GridColor      uint32
	FPS            uint32
}

// NewGame Returns a new initialized game.
func NewGame(title string, width, height int32) (*Game, error) {
	game := &Game{
		running:        false,
		title:          title,
		width:          width,
		height:         height,
		CellSize:       20,
		CellAliveColor: 0xFFFFFFFF,
		CellDeadColor:  0x00000000,
		GridColor:      0xFFFFFFFF,
		FPS:            60,
	}
	err := game.init()
	if err != nil {
		return nil, fmt.Errorf("failed to create init game resources: %w", err)
	}
	buffer, err := NewFrameBuffer(game)
	if err != nil {
		return nil, fmt.Errorf("failed to create new frame buffer: %w", err)
	}
	game.frameBuffer = buffer
	return game, nil
}

// Run the game. Calling this will block until exiting the game.
func (g *Game) Run() error {
	g.title = fmt.Sprintf("%s [%dx%d]", g.title, g.width/g.CellSize, g.height/g.CellSize)
	g.window.SetTitle(g.title)
	g.running = true
	tpf := uint64(1000 / g.FPS)
	lastTicks := sdl.GetTicks64()
	for g.running {
		delta := sdl.GetTicks64() - lastTicks
		if delta < tpf {
			continue
		}
		lastTicks = sdl.GetTicks64()
		fmt.Println("FPS:", 1000/delta)

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
	winErr := g.window.Destroy()
	errs := errors.Join(rendErr, winErr)
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
		case *sdl.MouseButtonEvent:
			if event.Type == sdl.MOUSEBUTTONDOWN {
				g.frameBuffer.ToggleCellState(event.X, event.Y)
			}
		}
	}
}

// update game state.
func (g *Game) update() {
	for y := int32(0); y < g.width/g.CellSize; y++ {
		for x := int32(0); x < g.height/g.CellSize; x++ {
			// TODO: Add rules here...
		}
	}
}

// render the game state to screen.
func (g *Game) render() {
	g.frameBuffer.DrawGrid()
	err := g.frameBuffer.Render()
	if err != nil {
		log.Println(err)
	}
	g.renderer.Present()
}
