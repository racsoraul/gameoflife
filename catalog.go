package main

import "strings"

// CatalogEntry holds a named RLE pattern.
type CatalogEntry struct {
	Name string
	RLE  string
}

// catalog contains classic Game of Life patterns.
var catalog = []CatalogEntry{
	{Name: "R-pentomino", RLE: "x = 3, y = 3, rule = B3/S23\nb2o$2ob$bo!"},
	{Name: "Glider", RLE: "x = 3, y = 3, rule = B3/S23\nbo$2bo$3o!"},
	{Name: "Gosper Glider Gun", RLE: "x = 36, y = 9, rule = B3/S23\n24bo$22bobo$12b2o6b2o12b2o$11bo3bo4b2o12b2o$2o8bo5bo3b2o$2o8bo3bob2o4bobo$10bo5bo7bo$11bo3bo$12b2o!"},
	{Name: "Lightweight Spaceship", RLE: "x = 5, y = 4, rule = B3/S23\nbo2bo$o$o3bo$b4o!"},
	{Name: "Pulsar", RLE: "x = 13, y = 13, rule = B3/S23\n2b3o3b3o2b$o4bobo4bo$o4bobo4bo$o4bobo4bo$2b3o3b3o2b2$2b3o3b3o2b$o4bobo4bo$o4bobo4bo$o4bobo4bo$2b3o3b3o!"},
	{Name: "Pentadecathlon", RLE: "x = 10, y = 3, rule = B3/S23\n2bo4bo$2ob4ob2o$2bo4bo!"},
	{Name: "Acorn", RLE: "x = 7, y = 3, rule = B3/S23\nbo$3bo$2o2b3o!"},
	{Name: "Diehard", RLE: "x = 8, y = 3, rule = B3/S23\n6bo$2o$bo3b3o!"},
	{Name: "Beacon", RLE: "x = 4, y = 4, rule = B3/S23\n2o$2o$2b2o$2b2o!"},
	{Name: "Blinker", RLE: "x = 3, y = 1, rule = B3/S23\n3o!"},
	{Name: "Toad", RLE: "x = 4, y = 2, rule = B3/S23\nb3o$3ob!"},
	{Name: "Block", RLE: "x = 2, y = 2, rule = B3/S23\n2o$2o!"},
	{Name: "Beehive", RLE: "x = 4, y = 3, rule = B3/S23\nb2o$o2bo$b2o!"},
	{Name: "Heavyweight Spaceship", RLE: "x = 7, y = 5, rule = B3/S23\n3bo$bo4bo$o$o5bo$b6o!"},
	{Name: "Infinite Growth", RLE: "x = 39, y = 1, rule = B3/S23\n8ob5o3b3o6b7ob5o!"},
}

// LoadCatalogEntry parses a catalog entry into a Pattern.
func LoadCatalogEntry(entry *CatalogEntry) (*Pattern, error) {
	return ParseRLE(strings.NewReader(entry.RLE))
}
