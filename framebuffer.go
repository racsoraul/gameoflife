package main

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// FrameBuffer It's composed of two color buffers. It allows computing the next generation of cells based on the
// previous generation, by alternating both color buffers. It renders to screen using a texture.
type FrameBuffer struct {
	g            *Game
	texture      *sdl.Texture
	section      *sdl.Rect // Section of the texture to render.
	colorBufferA []uint32  // Index 0.
	colorBufferB []uint32  // Index 1.
	index        uint8     // Tracks what color buffer to use.
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
		colorBufferB: make([]uint32, g.width*g.height),
		texture:      texture,
		section:      nil,
	}
	f.clearCurrent(g.CellDeadColor)
	return f, nil
}

// clearNext Clears next buffer with the provided color. It clears the color buffer that will render
// the next generation of cells.
func (fb *FrameBuffer) clearNext(color uint32) {
	for i := 0; i < len(fb.colorBufferA); i++ {
		if fb.index == 0 {
			fb.colorBufferB[i] = color
		} else {
			fb.colorBufferA[i] = color
		}
	}
}

// clearCurrent Clears current buffer with the provided color.
func (fb *FrameBuffer) clearCurrent(color uint32) {
	for i := 0; i < len(fb.colorBufferA); i++ {
		if fb.index == 0 {
			fb.colorBufferA[i] = color
		} else {
			fb.colorBufferB[i] = color
		}
	}
}

// Render Renders current color buffer using a texture. It increments or decrements the fb.index
// to alternate the color buffers to render the next generation.
func (fb *FrameBuffer) Render() error {
	var pixels unsafe.Pointer
	if fb.index == 0 {
		pixels = unsafe.Pointer(&fb.colorBufferA[0])
		// Swap buffer if the game is not paused.
		if fb.g.playing {
			fb.index++
		}
	} else {
		pixels = unsafe.Pointer(&fb.colorBufferB[0])
		// Swap buffer if the game is not paused.
		if fb.g.playing {
			fb.index--
		}
	}
	err := fb.texture.Update(
		nil,
		pixels,
		(int)(fb.g.width*4),
	)
	if err != nil {
		return err
	}
	err = fb.g.renderer.Copy(fb.texture, fb.section, &sdl.Rect{X: 0, Y: 0, W: fb.g.width, H: fb.g.height})
	if err != nil {
		return err
	}
	if fb.g.playing {
		// Clear the next buffer if the game is not paused.
		fb.clearNext(fb.g.CellDeadColor)
	}
	return nil
}

// DrawPixel Draws a pixel with the specified color. If next is false, it uses the current color buffer.
func (fb *FrameBuffer) DrawPixel(x, y int32, color uint32, next bool) {
	if x >= 0 && x < fb.g.width && y >= 0 && y < fb.g.height {
		var index uint8
		if next {
			index = 1
		}
		if fb.index == index {
			fb.colorBufferA[(y*fb.g.width)+x] = color
		} else {
			fb.colorBufferB[(y*fb.g.width)+x] = color
		}
	}
}

// GetPixelColor Returns color of the specified pixel. If next is false, it uses the current color buffer.
func (fb *FrameBuffer) GetPixelColor(x, y int32, next bool) uint32 {
	if x >= 0 && x < fb.g.width && y >= 0 && y < fb.g.height {
		var index uint8
		if next {
			index = 1
		}
		if fb.index == index {
			return fb.colorBufferA[(y*fb.g.width)+x]
		}
		return fb.colorBufferB[(y*fb.g.width)+x]
	}
	return fb.g.CellDeadColor
}

// DrawRect Draws a rectangle and fills it in with the specified color. If next is false, it uses the current color buffer.
func (fb *FrameBuffer) DrawRect(x, y, width, height int32, color uint32, next bool) {
	for j := int32(0); j < height; j++ {
		for i := int32(0); i < width; i++ {
			fb.DrawPixel(x+i, y+j, color, next)
		}
	}
}

// DrawGrid Draws a grid of cells with size cellSize and color GridColor using the renderer (overlay).
func (fb *FrameBuffer) DrawGrid() error {
	if fb.g.cellSize <= 1 {
		return nil
	}
	r := uint8((fb.g.GridColor >> 24) & 0xFF)
	g := uint8((fb.g.GridColor >> 16) & 0xFF)
	b := uint8((fb.g.GridColor >> 8) & 0xFF)
	a := uint8(fb.g.GridColor & 0xFF)
	err := fb.g.renderer.SetDrawColor(r, g, b, a)
	if err != nil {
		return err
	}

	for y := int32(0); y <= fb.g.height; y += fb.g.cellSize {
		err = fb.g.renderer.DrawLine(0, y, fb.g.width, y)
		if err != nil {
			return err
		}
	}
	for x := int32(0); x <= fb.g.width; x += fb.g.cellSize {
		err = fb.g.renderer.DrawLine(x, 0, x, fb.g.height)
		if err != nil {
			return err
		}
	}

	return nil
}
