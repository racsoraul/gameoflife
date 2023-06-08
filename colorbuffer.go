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
