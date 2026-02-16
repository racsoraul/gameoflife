package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/ttf"
)

// loadFont tries Windows-specific font paths and returns the first one that works.
func loadFont(size int) (*ttf.Font, error) {
	paths := []string{
		"C:\\Windows\\Fonts\\arial.ttf",
	}
	for _, p := range paths {
		f, err := ttf.OpenFont(p, size)
		if err == nil {
			return f, nil
		}
	}
	return nil, fmt.Errorf("no suitable font found for windows")
}
