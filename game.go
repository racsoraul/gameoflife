package main

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const zoomFactor = 1.2

// Game Holds the configs and state of the Conway's Game of Life.
type Game struct {
	running              bool
	title                string
	width                int32
	height               int32
	window               *sdl.Window
	renderer             *sdl.Renderer
	frameBuffer          *FrameBuffer // Holds every generation of cells.
	playing              bool         // Acts as Play/Pause for the game.
	activeTransitionRule string
	leftClickPressed     bool
	rightClickPressed    bool
	step                 bool // Progresses one generation while on pause.
	cellAliveColor       uint32
	cellDeadColor        uint32
	gridColor            uint32
	enableGrid           bool
	heatmapMode          bool
	trailMode            bool
	fullscreen           bool
	trails               map[[2]int32]uint8 // Recently dead cells with remaining fade frames.
	fps                  uint32
	generation           uint32 // Current generation number.
	font                 *ttf.Font
	infoHeight           int32
	genPerSec            int      // Target generations per second (1-60).
	genAccumulator       uint64   // Accumulated ms for generation timing.
	catalogIndex         int      // Current index in the pattern catalog (-1 = custom/file pattern).
	currentPattern       *Pattern // Last loaded pattern, used for reset.
	// Infinite grid. Maps cell coordinate to age (generations alive). Only live cells are stored.
	cells          map[[2]int32]uint32
	nextCells      map[[2]int32]uint32 // Scratch buffer for next generation (reused to avoid alloc).
	neighbourCount map[[2]int32]int    // Scratch buffer for neighbor counting (reused to avoid alloc).
	// Camera/viewport. World coordinate of the top-left corner of the screen.
	cameraX float64
	cameraY float64
	zoom    float64 // Pixels per cell.
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
		heatmapMode:    true,
		trailMode:      true,
		cellAliveColor: 0xFFFFFFFF,
		cellDeadColor:  0x000000FF,
		gridColor:      0xA9A9A9FF,
		fps:            60,
		infoHeight:     40,
		genPerSec:      25,
		catalogIndex:   -1,
		cells:          make(map[[2]int32]uint32),
		nextCells:      make(map[[2]int32]uint32),
		neighbourCount: make(map[[2]int32]int),
		trails:         make(map[[2]int32]uint8),
		zoom:           float64(cellSize),
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
func (g *Game) Run(pattern *Pattern) error {
	g.loadPattern(pattern)
	g.renderWorld()
	err := g.frameBuffer.Render()
	if err != nil {
		return fmt.Errorf("failed to render initial configuration")
	}
	g.renderer.Present()

	g.running = true
	tpf := uint64(1000 / g.fps)
	lastTicks := sdl.GetTicks64()
	for g.running {
		g.processInput()

		now := sdl.GetTicks64()
		delta := now - lastTicks
		if delta < tpf {
			sdl.Delay(1) // Yield CPU instead of busy-waiting.
			continue
		}
		lastTicks = now

		// Accumulate time for generation updates (decoupled from render FPS).
		// Cap delta to prevent burst of generations after suspend/minimize.
		g.genAccumulator += min(delta, 500)
		genInterval := uint64(1000 / g.genPerSec)
		for g.playing && g.genAccumulator >= genInterval {
			g.genAccumulator -= genInterval
			g.generation++
			g.RuleB3S23()
			if g.step {
				g.step = false
				g.playing = false
				g.genAccumulator = 0
				break
			}
		}

		fps := uint32(0)
		if delta > 0 {
			fps = uint32(1000 / delta)
		}
		g.render(fps)
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

// screenToWorld converts screen pixel coordinates to world cell coordinates.
func (g *Game) screenToWorld(screenX, screenY int32) (int32, int32) {
	wx := int32(math.Floor(g.cameraX + float64(screenX)/g.zoom))
	wy := int32(math.Floor(g.cameraY + float64(screenY)/g.zoom))
	return wx, wy
}

// updateTitle updates the window title with grid info.
func (g *Game) updateTitle() {
	g.window.SetTitle(fmt.Sprintf("%s [ZOOM:%.0f]", g.title, g.zoom))
}

// loadPattern clears the grid and places the given pattern centered at origin.
func (g *Game) loadPattern(pattern *Pattern) {
	g.currentPattern = pattern
	g.cells = make(map[[2]int32]uint32)
	g.nextCells = make(map[[2]int32]uint32)
	g.neighbourCount = make(map[[2]int32]int)
	g.trails = make(map[[2]int32]uint8)
	g.generation = 0
	g.genAccumulator = 0
	g.playing = false
	g.step = false
	g.activeTransitionRule = pattern.Rule
	offsetX := -int32(pattern.Width) / 2
	offsetY := -int32(pattern.Height) / 2
	for _, cell := range pattern.Cells {
		g.cells[[2]int32{offsetX + cell[0], offsetY + cell[1]}] = 1
	}
	g.cameraX = -float64(g.width) / (2 * g.zoom)
	g.cameraY = -float64(g.height) / (2 * g.zoom)
	g.updateTitle()
}

// toggleFullscreen switches between windowed and fullscreen desktop mode.
func (g *Game) toggleFullscreen() {
	g.fullscreen = !g.fullscreen
	var flag uint32
	if g.fullscreen {
		flag = sdl.WINDOW_FULLSCREEN_DESKTOP
	}
	if err := g.window.SetFullscreen(flag); err != nil {
		slog.Error("failed to toggle fullscreen", "error", err)
		g.fullscreen = !g.fullscreen
		return
	}
	// Remember what the camera is centered on before resizing.
	centerX := g.cameraX + float64(g.width)/(2*g.zoom)
	centerY := g.cameraY + float64(g.height)/(2*g.zoom)

	// Update dimensions to match the new window size.
	w, h := g.window.GetSize()
	g.width = w
	g.height = h - g.infoHeight

	// Recenter camera on the same world point.
	g.cameraX = centerX - float64(g.width)/(2*g.zoom)
	g.cameraY = centerY - float64(g.height)/(2*g.zoom)

	// Destroy the old texture and rebuild the framebuffer with new dimensions.
	if err := g.frameBuffer.texture.Destroy(); err != nil {
		slog.Error("failed to destroy old texture", "error", err)
	}
	buffer, err := NewFrameBuffer(g)
	if err != nil {
		slog.Error("failed to recreate framebuffer", "error", err)
		return
	}
	g.frameBuffer = buffer
}

// applyZoom multiplies the zoom by factor, keeping the center of the screen fixed.
func (g *Game) applyZoom(factor float64) {
	centerX := g.cameraX + float64(g.width)/(2*g.zoom)
	centerY := g.cameraY + float64(g.height)/(2*g.zoom)
	g.zoom *= factor
	if g.zoom < 1 {
		g.zoom = 1
	}
	g.cameraX = centerX - float64(g.width)/(2*g.zoom)
	g.cameraY = centerY - float64(g.height)/(2*g.zoom)
	g.updateTitle()
}

// processInput Processes the system events from the event queue.
func (g *Game) processInput() {
	panSpeed := 10.0 / g.zoom
	for nextEvent := sdl.PollEvent(); nextEvent != nil; nextEvent = sdl.PollEvent() {
		switch event := nextEvent.(type) {
		case *sdl.QuitEvent:
			g.running = false
		case *sdl.KeyboardEvent:
			if event.Type == sdl.KEYDOWN {
				switch event.Keysym.Sym {
				case sdl.K_ESCAPE:
					g.running = false
					return
				case sdl.K_p, sdl.K_SPACE:
					g.playing = !g.playing
					if g.playing {
						g.genAccumulator = 0
					}
					continue
				case sdl.K_g:
					g.enableGrid = !g.enableGrid
				case sdl.K_s:
					if !g.playing {
						g.playing = true
						g.step = true
					}
					continue
				case sdl.K_LEFT:
					g.cameraX -= panSpeed
				case sdl.K_RIGHT:
					g.cameraX += panSpeed
				case sdl.K_UP:
					g.cameraY -= panSpeed
				case sdl.K_DOWN:
					g.cameraY += panSpeed
				case sdl.K_EQUALS, sdl.K_PLUS:
					g.applyZoom(zoomFactor)
				case sdl.K_MINUS:
					g.applyZoom(1 / zoomFactor)
				case sdl.K_RIGHTBRACKET:
					g.genPerSec = min(g.genPerSec+1, 60)
				case sdl.K_LEFTBRACKET:
					g.genPerSec = max(g.genPerSec-1, 1)
				case sdl.K_f:
					g.toggleFullscreen()
				case sdl.K_h:
					g.heatmapMode = !g.heatmapMode
				case sdl.K_t:
					g.trailMode = !g.trailMode
					if !g.trailMode {
						g.trails = make(map[[2]int32]uint8)
					}
				case sdl.K_c:
					// Clear the grid.
					g.cells = make(map[[2]int32]uint32)
					g.nextCells = make(map[[2]int32]uint32)
					g.neighbourCount = make(map[[2]int32]int)
					g.trails = make(map[[2]int32]uint8)
					g.generation = 0
					g.genAccumulator = 0
					g.playing = false
					g.step = false
				case sdl.K_r:
					// Reset to the current pattern.
					if g.currentPattern != nil {
						g.loadPattern(g.currentPattern)
					}
				case sdl.K_n:
					g.catalogIndex = (g.catalogIndex + 1) % len(catalog)
					entry := &catalog[g.catalogIndex]
					p, err := LoadCatalogEntry(entry)
					if err != nil {
						slog.Error("failed to load catalog pattern", "name", entry.Name, "error", err)
						continue
					}
					g.loadPattern(p)
				}
			}
		case *sdl.MouseButtonEvent:
			switch event.Button {
			case sdl.BUTTON_LEFT:
				if event.Type == sdl.MOUSEBUTTONDOWN {
					wx, wy := g.screenToWorld(event.X, event.Y)
					pos := [2]int32{wx, wy}
					if g.cells[pos] > 0 {
						delete(g.cells, pos)
					} else {
						g.cells[pos] = 1
					}
					g.leftClickPressed = true
				} else {
					g.leftClickPressed = false
				}
			case sdl.BUTTON_RIGHT:
				g.rightClickPressed = event.Type == sdl.MOUSEBUTTONDOWN
			}
		case *sdl.MouseMotionEvent:
			if g.leftClickPressed {
				wx, wy := g.screenToWorld(event.X, event.Y)
				g.cells[[2]int32{wx, wy}] = 1
			}
			if g.rightClickPressed {
				g.cameraX -= float64(event.XRel) / g.zoom
				g.cameraY -= float64(event.YRel) / g.zoom
			}
		case *sdl.MouseWheelEvent:
			if event.Y > 0 {
				g.applyZoom(zoomFactor)
			} else if event.Y < 0 {
				g.applyZoom(1 / zoomFactor)
			}
		}
	}
}

// ageToColor returns an RGBA8888 color for a cell based on its age.
// New cells are bright green, aging through yellow, orange, red, then blue/purple.
func ageToColor(age uint32) uint32 {
	switch {
	case age <= 5:
		// Bright green → yellow (green stays 255, red ramps up).
		r := min(age*50, 255)
		return (r << 24) | (0xFF << 16) | (0x00 << 8) | 0xFF
	case age <= 15:
		// Yellow → orange → red (green ramps down).
		g := 255 - min((age-5)*25, 255)
		return (0xFF << 24) | (g << 16) | (0x00 << 8) | 0xFF
	case age <= 30:
		// Red → purple (blue ramps up).
		b := min((age-15)*17, 255)
		return (0xFF << 24) | (0x00 << 16) | (b << 8) | 0xFF
	default:
		// Deep blue/purple for ancient cells.
		return (0x80 << 24) | (0x00 << 16) | (0xFF << 8) | 0xFF
	}
}

// renderWorld draws the current cell state from the infinite grid into the pixel buffer.
func (g *Game) renderWorld() {
	g.frameBuffer.clearCurrent(g.cellDeadColor)

	cellPx := max(int32(g.zoom), 1)

	// Compute the visible world-cell bounding box once.
	// Any cell outside this range is guaranteed off-screen — skip it without
	// doing per-cell float math.
	visMinX := int32(math.Floor(g.cameraX)) - 1
	visMinY := int32(math.Floor(g.cameraY)) - 1
	visMaxX := int32(math.Ceil(g.cameraX+float64(g.width)/g.zoom)) + 1
	visMaxY := int32(math.Ceil(g.cameraY+float64(g.height)/g.zoom)) + 1

	// Render fade trails (behind alive cells).
	if g.trailMode {
		for pos, remaining := range g.trails {
			if g.cells[pos] > 0 {
				delete(g.trails, pos)
				continue
			}
			if pos[0] < visMinX || pos[0] > visMaxX || pos[1] < visMinY || pos[1] > visMaxY {
				// Still decay even when off-screen.
				g.trails[pos] = remaining - 1
				if remaining <= 1 {
					delete(g.trails, pos)
				}
				continue
			}
			screenX := int32(float64(pos[0])*g.zoom - g.cameraX*g.zoom)
			screenY := int32(float64(pos[1])*g.zoom - g.cameraY*g.zoom)
			brightness := uint32(remaining) * 25
			trailColor := (brightness << 24) | (brightness << 16) | (brightness << 8) | 0xFF
			g.frameBuffer.DrawRect(screenX, screenY, cellPx, cellPx, trailColor)
			g.trails[pos] = remaining - 1
			if remaining <= 1 {
				delete(g.trails, pos)
			}
		}
	}

	for pos, age := range g.cells {
		// Fast bounds check in world coordinates — avoids float math for off-screen cells.
		if pos[0] < visMinX || pos[0] > visMaxX || pos[1] < visMinY || pos[1] > visMaxY {
			continue
		}
		screenX := int32(float64(pos[0])*g.zoom - g.cameraX*g.zoom)
		screenY := int32(float64(pos[1])*g.zoom - g.cameraY*g.zoom)
		color := g.cellAliveColor
		if g.heatmapMode {
			color = ageToColor(age)
		}
		g.frameBuffer.DrawRect(screenX, screenY, cellPx, cellPx, color)
	}
}

// render the game state to screen.
func (g *Game) render(fps uint32) {
	g.renderWorld()

	err := g.frameBuffer.Render()
	if err != nil {
		slog.Error("failed to render frame buffer", "error", err)
	}

	if g.enableGrid {
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
	if g.enableGrid {
		grid = "On"
	}
	infoTextColor := sdl.Color{R: 200, G: 200, B: 200, A: 255}

	// Shortcuts.
	heatmap := "Off"
	if g.heatmapMode {
		heatmap = "On"
	}
	trail := "Off"
	if g.trailMode {
		trail = "On"
	}
	shortcuts := fmt.Sprintf("Exit: ESC | %s: P/Space | Step: S | Grid(%s): G | Heatmap(%s): H | Trail(%s): T | Fullscreen: F | Speed: [/] | Clear: C | Reset: R | Next: N", pausePlay, grid, heatmap, trail)
	err = g.drawText(shortcuts, 10, g.height+5, infoTextColor)
	if err != nil {
		return err
	}

	// Generation/Step and FPS.
	patternName := "Custom"
	if g.catalogIndex >= 0 {
		patternName = catalog[g.catalogIndex].Name
	}
	stats := fmt.Sprintf("%s | %s | Gen: %d | Speed: %d gen/s | FPS: %d | Click: toggle | Right-drag: pan", patternName, g.activeTransitionRule, g.generation, g.genPerSec, fps)
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
