package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Game Holds the configs and state of the Conway's Game of Life.
type Game struct {
	running          bool
	title            string
	width            int32
	height           int32
	window           *sdl.Window
	renderer         *sdl.Renderer
	frameBuffer      *FrameBuffer // Holds every generation of cells.
	playing          bool         // Acts as Play/Pause for the game.
	leftClickPressed bool
	step             bool // Progresses one generation while on pause.
	cellSize         int32
	CellAliveColor   uint32
	CellDeadColor    uint32
	GridColor        uint32
	EnableGrid       bool
	FPS              uint32
	generation       uint32 // Current generation number.
	font             *ttf.Font
	infoHeight       int32
}

// NewGame Returns a new initialized game.
func NewGame(title string, width, height, cellSize int32) (*Game, error) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})))

	game := &Game{
		running:        false,
		title:          title,
		width:          width,
		height:         height,
		playing:        false, // Paused by default.
		cellSize:       cellSize,
		CellAliveColor: 0xFFFFFFFF,
		CellDeadColor:  0x000000FF,
		GridColor:      0xA9A9A9FF,
		FPS:            60,
		infoHeight:     40,
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
	g.window.SetTitle(fmt.Sprintf("%s [%dx%d]", g.title, g.width/g.cellSize, g.height/g.cellSize))

	// Initial configuration.
	posX := ((g.width / g.cellSize) / 2) - 1
	posY := ((g.height / g.cellSize) / 2) - 1
	g.frameBuffer.SetCellState(ALIVE, posX, posY, false)
	g.frameBuffer.SetCellState(ALIVE, posX+1, posY, false)
	g.frameBuffer.SetCellState(ALIVE, posX, posY+1, false)
	g.frameBuffer.SetCellState(ALIVE, posX-1, posY+1, false)
	g.frameBuffer.SetCellState(ALIVE, posX, posY+2, false)
	err := g.frameBuffer.Render()
	if err != nil {
		return fmt.Errorf("failed to render initial configuration")
	}
	g.renderer.Present()

	g.running = true
	tpf := uint64(1000 / g.FPS)
	lastTicks := sdl.GetTicks64()
	for g.running {
		g.processInput()

		delta := sdl.GetTicks64() - lastTicks
		if delta < tpf {
			continue
		}
		lastTicks = sdl.GetTicks64()

		g.update()
		fps := uint32(0)
		if delta > 0 {
			fps = uint32(1000 / delta)
		}
		g.render(fps)

		if g.step {
			g.step = false
			g.playing = false
		}
	}
	return g.shutdown()
}

// init Initialize game resources.
func (g *Game) init() error {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		return fmt.Errorf("failed to initialize subsystems: %w", err)
	}

	window, err := sdl.CreateWindow(
		g.title,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		g.width,
		g.height+g.infoHeight,
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

	if err = ttf.Init(); err != nil {
		return fmt.Errorf("failed to initialize TTF: %w", err)
	}

	font, err := loadFont(14)
	if err != nil {
		slog.Warn("failed to load font", "error", err)
	}
	g.font = font

	return nil
}

// shutdown Clean up, free, and destroy resources.
func (g *Game) shutdown() error {
	if g.font != nil {
		g.font.Close()
	}
	rendErr := g.renderer.Destroy()
	winErr := g.window.Destroy()
	errs := errors.Join(rendErr, winErr)
	ttf.Quit()
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
			if event.Type == sdl.KEYDOWN {
				if event.Keysym.Sym == sdl.K_ESCAPE {
					g.running = false
					return
				}
				if event.Keysym.Sym == sdl.K_p {
					g.playing = !g.playing
					continue
				}
				if event.Keysym.Sym == sdl.K_g {
					g.EnableGrid = !g.EnableGrid
				}
				if event.Keysym.Sym == sdl.K_s {
					if !g.playing {
						g.playing = true
						g.step = true
					}
					continue
				}
			}
		case *sdl.MouseButtonEvent:
			if event.Type == sdl.MOUSEBUTTONDOWN {
				g.frameBuffer.ToggleCellState(event.X, event.Y, false)
				g.leftClickPressed = true
			} else {
				g.leftClickPressed = false
			}
		case *sdl.MouseMotionEvent:
			if g.leftClickPressed {
				g.frameBuffer.ToggleCellState(event.X, event.Y, false)
			}
		}
	}
}

// update game state.
func (g *Game) update() {
	if !g.playing {
		return
	}
	g.generation++
	for y := int32(0); y < g.height/g.cellSize; y++ {
		for x := int32(0); x < g.width/g.cellSize; x++ {
			g.RuleB3S23(x, y)
		}
	}
}

// render the game state to screen.
func (g *Game) render(fps uint32) {
	err := g.frameBuffer.Render()
	if err != nil {
		slog.Error("failed to render frame buffer", "error", err)
	}

	if g.EnableGrid {
		err = g.frameBuffer.DrawGrid()
		if err != nil {
			slog.Error("failed to draw grid", "error", err)
		}
	}

	err = g.renderInfoSection(fps)
	if err != nil {
		slog.Error("failed to render info section", "error", err)
	}

	g.renderer.Present()
}

// renderInfoSection Render an info section with the shortcuts and stats of the game.
func (g *Game) renderInfoSection(fps uint32) error {
	if g.font == nil {
		return nil
	}

	// Draw background for the info section.
	rect := sdl.Rect{
		X: 0,
		Y: g.height,
		W: g.width,
		H: g.infoHeight,
	}

	// Dark gray.
	err := g.renderer.SetDrawColor(32, 32, 32, 255)
	if err != nil {
		return err
	}
	err = g.renderer.FillRect(&rect)
	if err != nil {
		return err
	}

	// Shortcuts
	pausePlay := "Pause"
	if !g.playing {
		pausePlay = "Play"
	}
	grid := "Off"
	if g.EnableGrid {
		grid = "On"
	}
	infoTextColor := sdl.Color{R: 200, G: 200, B: 200, A: 255}

	// Shortcuts.
	shortcuts := fmt.Sprintf("Exit: ESC | %s: P | Step: S | Grid(%s): G", pausePlay, grid)
	err = g.drawText(shortcuts, 10, g.height+5, infoTextColor)
	if err != nil {
		return err
	}

	// Generation/Step and FPS.
	stats := fmt.Sprintf("Gen: %d | FPS: %d | Click on any cell to toggle its state", g.generation, fps)
	err = g.drawText(stats, 10, g.height+22, infoTextColor)
	if err != nil {
		return err
	}

	return nil
}

// drawText Renders the provided text.
func (g *Game) drawText(text string, x, y int32, color sdl.Color) error {
	if g.font == nil || text == "" {
		return nil
	}
	surface, err := g.font.RenderUTF8Blended(text, color)
	if err != nil {
		return err
	}
	defer surface.Free()

	texture, err := g.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	dst := sdl.Rect{
		X: x,
		Y: y,
		W: surface.W,
		H: surface.H,
	}
	return g.renderer.Copy(texture, nil, &dst)
}
