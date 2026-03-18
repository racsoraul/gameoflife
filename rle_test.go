package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseRLE_RPentomino(t *testing.T) {
	input := `#N R-pentomino
x = 3, y = 3, rule = B3/S23
b2o$2ob$bo!
`
	pattern, err := ParseRLE(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pattern.Rule != "B3/S23" {
		t.Fatalf("expected rule B3/S23, got %s", pattern.Rule)
	}
	if pattern.Width != 3 || pattern.Height != 3 {
		t.Fatalf("expected 3x3, got %dx%d", pattern.Width, pattern.Height)
	}
	// R-pentomino cells: (1,0),(2,0),(0,1),(1,1),(1,2)
	expected := [][2]int32{{1, 0}, {2, 0}, {0, 1}, {1, 1}, {1, 2}}
	if len(pattern.Cells) != len(expected) {
		t.Fatalf("expected %d cells, got %d", len(expected), len(pattern.Cells))
	}
	for i, cell := range pattern.Cells {
		if cell != expected[i] {
			t.Errorf("cell %d: expected %v, got %v", i, expected[i], cell)
		}
	}
}

func TestParseRLE_Glider(t *testing.T) {
	input := `x = 3, y = 3
bo$2bo$3o!
`
	pattern, err := ParseRLE(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Glider: (1,0),(2,1),(0,2),(1,2),(2,2)
	expected := [][2]int32{{1, 0}, {2, 1}, {0, 2}, {1, 2}, {2, 2}}
	if len(pattern.Cells) != len(expected) {
		t.Fatalf("expected %d cells, got %d", len(expected), len(pattern.Cells))
	}
	for i, cell := range pattern.Cells {
		if cell != expected[i] {
			t.Errorf("cell %d: expected %v, got %v", i, expected[i], cell)
		}
	}
}

func TestParseRLE_MissingHeader(t *testing.T) {
	input := `bo$2bo$3o!`
	_, err := ParseRLE(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for missing header")
	}
}

func TestParseRLE_AllPatternFiles(t *testing.T) {
	files, err := filepath.Glob("patterns/*.rle")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no pattern files found in patterns/")
	}
	for _, f := range files {
		name := filepath.Base(f)
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}
			p, err := ParseRLE(strings.NewReader(string(data)))
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			if len(p.Cells) == 0 {
				t.Error("pattern has no alive cells")
			}
			if p.Width == 0 || p.Height == 0 {
				t.Errorf("invalid dimensions: %dx%d", p.Width, p.Height)
			}
			t.Logf("%dx%d, %d cells, rule=%s", p.Width, p.Height, len(p.Cells), p.Rule)
		})
	}
}
