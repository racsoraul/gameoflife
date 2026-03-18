// Package mcell parses Golly Macrocell (.mc) files into a list of alive cell
// coordinates. The format is a quadtree where each node is either a 2x2 leaf
// (level 1) with literal cell states, or an internal node referencing four
// child quadrants (NW, NE, SW, SE). Node 0 is the implicit empty node at any
// level.
package mcell

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// node is a single macrocell quadtree node.
type node struct {
	level int
	// For level 1: nw/ne/sw/se are cell states (0 or 1).
	// For level 2+: nw/ne/sw/se are 1-indexed references into the node list.
	nw, ne, sw, se int
}

// Parse reads a Golly Macrocell file and returns the alive cell coordinates.
// Coordinates are centered so that (0,0) is roughly the middle of the pattern.
func Parse(r io.Reader) ([][2]int32, error) {
	nodes, err := parseNodes(r)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, nil
	}

	// Walk the quadtree to collect raw coordinates.
	root := nodes[len(nodes)-1]
	var raw [][2]int64
	walk(nodes, len(nodes), root.level, 0, 0, func(x, y int64) {
		raw = append(raw, [2]int64{x, y})
	})
	if len(raw) == 0 {
		return nil, nil
	}

	// Find bounding box and center the pattern around (0,0).
	minX, minY := raw[0][0], raw[0][1]
	maxX, maxY := minX, minY
	for _, c := range raw[1:] {
		minX = min(minX, c[0])
		maxX = max(maxX, c[0])
		minY = min(minY, c[1])
		maxY = max(maxY, c[1])
	}
	cx := (minX + maxX) / 2
	cy := (minY + maxY) / 2

	cells := make([][2]int32, len(raw))
	for i, c := range raw {
		cells[i] = [2]int32{int32(c[0] - cx), int32(c[1] - cy)}
	}
	return cells, nil
}

// parseNodes reads all node definitions from the macrocell file.
func parseNodes(r io.Reader) ([]node, error) {
	scanner := bufio.NewScanner(r)
	// Macrocell files can have long lines in extended formats.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var nodes []node
	headerSeen := false

	for scanner.Scan() {
		line := scanner.Text()

		// Skip the [M2] header line.
		if !headerSeen {
			if strings.HasPrefix(line, "[M2]") {
				headerSeen = true
				continue
			}
		}

		// Skip comments and empty lines.
		if line == "" || line[0] == '#' {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 5 {
			continue
		}

		vals := make([]int, 5)
		for i, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil {
				return nil, fmt.Errorf("invalid node value %q: %w", p, err)
			}
			vals[i] = v
		}

		nodes = append(nodes, node{
			level: vals[0],
			nw:    vals[1],
			ne:    vals[2],
			sw:    vals[3],
			se:    vals[4],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read macrocell file: %w", err)
	}
	return nodes, nil
}

// walk recursively traverses the quadtree and calls emit for each alive cell.
// idx is 1-indexed into nodes. x, y is the top-left corner of this node's region.
func walk(nodes []node, idx, level int, x, y int64, emit func(int64, int64)) {
	if idx <= 0 || idx > len(nodes) {
		return
	}

	n := nodes[idx-1] // Convert to 0-indexed.

	if level == 1 {
		// Leaf: 2x2 grid of literal cell states.
		if n.nw != 0 {
			emit(x, y)
		}
		if n.ne != 0 {
			emit(x+1, y)
		}
		if n.sw != 0 {
			emit(x, y+1)
		}
		if n.se != 0 {
			emit(x+1, y+1)
		}
		return
	}

	half := int64(1) << (level - 1)
	walk(nodes, n.nw, level-1, x, y, emit)
	walk(nodes, n.ne, level-1, x+half, y, emit)
	walk(nodes, n.sw, level-1, x, y+half, emit)
	walk(nodes, n.se, level-1, x+half, y+half, emit)
}
