package main

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// FrameBuffer Holds a pixel buffer and renders it to the screen using a texture.
type FrameBuffer struct {
	g            *Game
	texture      *sdl.Texture
	colorBufferA []uint32
}

// NewFrameBuffer Returns a new and initialized FrameBuffer.
func NewFrameBuffer(g *Game) (*FrameBuffer, error) {
	texture, err := g.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, g.width, g.height)
	if err != nil {
		return nil, fmt.Errorf("failed to create texture: %w", err)
	}
	f := &FrameBuffer{
		g:            g,
		colorBufferA: make([]uint32, g.width*g.height),
		texture:      texture,
	}
	f.clearCurrent(g.cellDeadColor)
	return f, nil
}

// clearCurrent Clears the buffer with the provided color.
func (fb *FrameBuffer) clearCurrent(color uint32) {
	for i := range fb.colorBufferA {
		fb.colorBufferA[i] = color
	}
}

// Render Renders current color buffer using a texture.
func (fb *FrameBuffer) Render() error {
	pixels := unsafe.Pointer(&fb.colorBufferA[0])
	err := fb.texture.Update(
		nil,
		pixels,
		(int)(fb.g.width*4),
	)
	if err != nil {
		return err
	}
	return fb.g.renderer.Copy(fb.texture, nil, &sdl.Rect{X: 0, Y: 0, W: fb.g.width, H: fb.g.height})
}

// DrawPixel Draws a pixel with the specified color.
func (fb *FrameBuffer) DrawPixel(x, y int32, color uint32) {
	if x >= 0 && x < fb.g.width && y >= 0 && y < fb.g.height {
		fb.colorBufferA[(y*fb.g.width)+x] = color
	}
}

// DrawRect Draws a rectangle and fills it in with the specified color.
func (fb *FrameBuffer) DrawRect(x, y, width, height int32, color uint32) {
	for j := int32(0); j < height; j++ {
		for i := int32(0); i < width; i++ {
			fb.DrawPixel(x+i, y+j, color)
		}
	}
}

// DrawGrid Draws a grid of cells using the current zoom level and camera offset.
func (fb *FrameBuffer) DrawGrid() error {
	cellPx := fb.g.zoom
	if cellPx < 3 {
		// Grid lines are not useful at very small zoom.
		return nil
	}
	r := uint8((fb.g.gridColor >> 24) & 0xFF)
	g := uint8((fb.g.gridColor >> 16) & 0xFF)
	b := uint8((fb.g.gridColor >> 8) & 0xFF)
	a := uint8(fb.g.gridColor & 0xFF)
	err := fb.g.renderer.SetDrawColor(r, g, b, a)
	if err != nil {
		return err
	}

	// Compute the screen-pixel offset for the first grid line.
	offX := -math.Mod(fb.g.cameraX*cellPx, cellPx)
	if offX > 0 {
		offX -= cellPx
	}
	offY := -math.Mod(fb.g.cameraY*cellPx, cellPx)
	if offY > 0 {
		offY -= cellPx
	}

	for x := offX; x <= float64(fb.g.width); x += cellPx {
		ix := int32(x)
		err = fb.g.renderer.DrawLine(ix, 0, ix, fb.g.height)
		if err != nil {
			return err
		}
	}
	for y := offY; y <= float64(fb.g.height); y += cellPx {
		iy := int32(y)
		err = fb.g.renderer.DrawLine(0, iy, fb.g.width, iy)
		if err != nil {
			return err
		}
	}

	return nil
}
