package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Pattern represents a parsed RLE pattern with its alive cell coordinates.
type Pattern struct {
	Width  int
	Height int
	Rule   string
	Cells  [][2]int32 // Alive cell positions as (x, y).
}

// LoadPattern loads an RLE pattern from a file path.
func LoadPattern(path string) (*Pattern, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open pattern file: %w", err)
	}
	defer f.Close()
	return ParseRLE(f)
}

// ParseRLE parses an RLE-encoded pattern from the given reader.
func ParseRLE(r io.Reader) (*Pattern, error) {
	scanner := bufio.NewScanner(r)
	var headerParsed bool
	var pattern Pattern
	var rleData strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Skip comment lines.
		if line[0] == '#' {
			continue
		}
		// Parse header line.
		if !headerParsed && strings.HasPrefix(line, "x") {
			if err := parseHeader(line, &pattern); err != nil {
				return nil, fmt.Errorf("failed to parse header: %w", err)
			}
			headerParsed = true
			continue
		}
		rleData.WriteString(line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}
	if !headerParsed {
		return nil, fmt.Errorf("missing RLE header line")
	}

	if err := decodeRLE(rleData.String(), &pattern); err != nil {
		return nil, fmt.Errorf("failed to decode RLE data: %w", err)
	}
	return &pattern, nil
}

// parseHeader parses the header line of an RLE pattern and updates the Pattern struct.
func parseHeader(line string, p *Pattern) error {
	// Expected format: x = N, y = N[, rule = ...]
	parts := strings.Split(line, ",")
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		switch key {
		case "x":
			n, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid x value: %w", err)
			}
			p.Width = n
		case "y":
			n, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid y value: %w", err)
			}
			p.Height = n
		case "rule":
			p.Rule = value
		}
	}
	return nil
}

// decodeRLE decodes a Run-Length Encoded (RLE) pattern string into a Pattern struct.
func decodeRLE(data string, p *Pattern) error {
	var x, y int32
	var runCount int

	for _, ch := range data {
		switch {
		case ch == '!':
			return nil
		case ch == '$':
			count := max(runCount, 1)
			y += int32(count)
			x = 0
			runCount = 0
		case ch == 'b':
			count := max(runCount, 1)
			x += int32(count)
			runCount = 0
		case ch == 'o':
			count := max(runCount, 1)
			for range count {
				p.Cells = append(p.Cells, [2]int32{x, y})
				x++
			}
			runCount = 0
		case unicode.IsDigit(ch):
			runCount = runCount*10 + int(ch-'0')
		default:
			// Ignore whitespace and unknown characters.
		}
	}
	return nil
}
