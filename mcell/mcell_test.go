package mcell

import (
	"strings"
	"testing"
)

func TestParse_SingleGlider(t *testing.T) {
	// A hand-crafted macrocell encoding of a glider.
	// Level 1 nodes are 2x2 leaves, level 2 is 4x4, level 3 is 8x8.
	// Glider (3x3):
	//   .#.
	//   ..#
	//   ###
	// Placed in the NW quadrant of an 8x8 grid:
	//   NW of NW (2x2): .#  → nw=0, ne=1, sw=0, se=0
	//                    ..
	//   NE of NW (2x2): ..  → all 0
	//                    ..
	//   SW of NW (2x2): ..  → nw=0, ne=1, sw=1, se=1
	//                    ##
	//   SE of NW (2x2): ..  → nw=1, ne=0, sw=0, se=0
	//                    ..
	input := `[M2] (golly 2.8)
#R B3/S23
1 0 1 0 0
1 0 0 0 1
1 1 1 1 0
2 1 0 3 2
3 4 0 0 0
`
	cells, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cells) != 5 {
		t.Fatalf("expected 5 cells, got %d", len(cells))
	}
	t.Logf("glider cells: %v", cells)
}

func TestParse_Empty(t *testing.T) {
	input := `[M2] (golly 2.8)
#R B3/S23
1 0 0 0 0
2 0 0 0 0
`
	cells, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cells) != 0 {
		t.Fatalf("expected 0 cells, got %d", len(cells))
	}
}

func TestParse_Block(t *testing.T) {
	// 2x2 block — all alive.
	input := `[M2] (golly 2.8)
#R B3/S23
1 1 1 1 1
`
	cells, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cells) != 4 {
		t.Fatalf("expected 4 cells, got %d", len(cells))
	}
}

func TestParse_InvalidLine(t *testing.T) {
	input := `[M2] (golly 2.8)
1 abc 0 0 0
`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid input")
	}
}
