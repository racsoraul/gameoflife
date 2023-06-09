package main

import (
	"unsafe"
)

// clearColorBuffer Clears the game's color buffer with the provided color.
func (g *Game) clearColorBuffer(color uint32) {
	for i := 0; i < len(g.colorBuffer); i++ {
		g.colorBuffer[i] = color
	}
}

// renderColorBuffer Renders the current game's color buffer to the current rendering target using a texture.
func (g *Game) renderColorBuffer() error {
	err := g.colorBufferTexture.Update(
		nil,
		unsafe.Pointer(&g.colorBuffer[0]),
		(int)(g.width*(int32)(unsafe.Sizeof(g.width))),
	)
	if err != nil {
		return err
	}
	err = g.renderer.Copy(g.colorBufferTexture, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// drawPixel Draws a pixel with the specified color in the colorBuffer.
func (g *Game) drawPixel(x, y int32, color uint32) {
	if x >= 0 && x < g.width && y >= 0 && y < g.height {
		g.colorBuffer[(y*g.width)+x] = color
	}
}

// drawRect Draws a rectangle and fills it in with the specified color.
func (g *Game) drawRect(x, y, width, height int32, color uint32) {
	for j := int32(0); j < height; j++ {
		for i := int32(0); i < width; i++ {
			g.drawPixel(x+i, y+j, color)
		}
	}
}

// drawGrid Draws a grid of cells with size CellSize and color GridColor.
func (g *Game) drawGrid() {
	for y := int32(0); y < g.height; y++ {
		for x := int32(0); x < g.width; x++ {
			if y%g.CellSize == 0 || x%g.CellSize == 0 {
				g.drawPixel(x, y, g.GridColor)
			}
		}
	}
}
