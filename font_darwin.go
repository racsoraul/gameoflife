package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/ttf"
)

// loadFont tries macOS-specific font paths and returns the first one that works.
func loadFont(size int) (*ttf.Font, error) {
	paths := []string{
		"/System/Library/Fonts/Supplemental/Arial.ttf",
		"/System/Library/Fonts/Helvetica.ttc",
		"/Library/Fonts/Arial.ttf",
	}
	for _, p := range paths {
		f, err := ttf.OpenFont(p, size)
		if err == nil {
			return f, nil
		}
	}
	return nil, fmt.Errorf("no suitable font found for darwin")
}
