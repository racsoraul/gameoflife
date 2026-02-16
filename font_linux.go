package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/ttf"
)

// loadFont tries Linux-specific font paths and returns the first one that works.
func loadFont(size int) (*ttf.Font, error) {
	paths := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
		"/usr/share/fonts/dejavu-sans-fonts/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/liberation-sans/LiberationSans-Regular.ttf",
		"/usr/share/fonts/truetype/freefont/FreeSans.ttf",
		"/usr/share/fonts/gnu-free/FreeSans.ttf",
	}
	for _, p := range paths {
		f, err := ttf.OpenFont(p, size)
		if err == nil {
			return f, nil
		}
	}
	return nil, fmt.Errorf("no suitable font found for linux")
}
