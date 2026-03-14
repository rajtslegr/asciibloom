// Package core provides shared types and utilities for generators.
package core

// GeneratorType represents the type of generator to use.
type GeneratorType string

const (
	TypeBrownian  GeneratorType = "brownian"
	TypeFlowField GeneratorType = "flow"
)

// Grid provides common grid operations for generators.
type Grid struct {
	width  int
	height int
	cells  [][]int
}

// NewGrid creates a new grid with the given dimensions.
func NewGrid(width, height int) *Grid {
	cells := make([][]int, height)
	for y := range cells {
		cells[y] = make([]int, width)
	}
	return &Grid{
		width:  width,
		height: height,
		cells:  cells,
	}
}

// Width returns the grid width.
func (g *Grid) Width() int { return g.width }

// Height returns the grid height.
func (g *Grid) Height() int { return g.height }

// Get returns the value at position (x, y), or 0 if out of bounds.
func (g *Grid) Get(x, y int) int {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return 0
	}
	return g.cells[y][x]
}

// Set sets the value at position (x, y) if within bounds and new value is higher.
func (g *Grid) Set(x, y, val int) {
	if x >= 0 && x < g.width && y >= 0 && y < g.height && g.cells[y][x] < val {
		g.cells[y][x] = val
	}
}

// Clear resets all cells to 0.
func (g *Grid) Clear() {
	for y := range g.cells {
		for x := range g.cells[y] {
			g.cells[y][x] = 0
		}
	}
}

// CountNeighbors counts non-zero neighbors around (x, y).
func (g *Grid) CountNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			if g.Get(x+dx, y+dy) > 0 {
				count++
			}
		}
	}
	return count
}

// CharForIntensity returns an ASCII character based on neighbor count and intensity.
func CharForIntensity(neighbors, intensity int) rune {
	switch {
	case neighbors >= 6:
		return 'O'
	case neighbors >= 4:
		return 'o'
	case neighbors >= 3:
		return '*'
	case neighbors >= 2:
		return ':'
	case intensity >= 5:
		return '+'
	default:
		return '.'
	}
}
