package animation

import (
	"context"
	"time"
)

const (
	framesPerSecond = 20
	frameDuration   = time.Second / framesPerSecond
)

// Generator defines the interface for content generation in the animation.
type Generator interface {
	Step() bool
	Render(buffer [][]rune)
	ClusterSize() int
	PostProcess()
}

// Renderer handles the animation loop and rendering to the terminal.
type Renderer struct {
	terminal  *Terminal
	generator Generator
	buffer    [][]rune
}

// NewRenderer creates a new renderer with the given terminal and generator.
func NewRenderer(terminal *Terminal, generator Generator) *Renderer {
	height := terminal.Height()
	width := terminal.Width()

	buffer := make([][]rune, height)
	for i := range buffer {
		buffer[i] = make([]rune, width)
		for j := range buffer[i] {
			buffer[i][j] = ' '
		}
	}

	return &Renderer{
		terminal:  terminal,
		generator: generator,
		buffer:    buffer,
	}
}

// Run starts the animation loop and blocks until ctx is canceled.
// Returns an error if rendering fails.
func (r *Renderer) Run(ctx context.Context) error {
	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.terminal.showCursor()
			r.terminal.clearScreen()
			r.terminal.moveCursor(1, 1)
			return nil
		case <-ticker.C:
			r.generator.Step()
			r.render()
		}
	}
}

func (r *Renderer) render() {
	height := len(r.buffer)
	if height == 0 {
		return
	}
	width := len(r.buffer[0])

	for y := range r.buffer {
		for x := range r.buffer[y] {
			r.buffer[y][x] = ' '
		}
	}

	r.generator.Render(r.buffer)

	var output []rune
	output = append(output, []rune("\x1b[H")...)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			char := r.buffer[y][x]
			if char == ' ' {
				continue
			}

			output = append(output, []rune("\x1b["+itoa(y+1)+";"+itoa(x+1)+"H")...)

			neighbors := r.countNeighbors(x, y)
			maxNeighbors := r.maxPossibleNeighbors(x, y)
			// Normalize neighbor count to 0-8 scale
			normalized := neighbors * 8 / maxNeighbors
			color := pickColor(normalized)

			output = append(output, []rune("\x1b["+color+"m")...)
			output = append(output, char)
		}
	}

	output = append(output, []rune("\x1b[0m")...)
	r.terminal.Write(string(output))
}

func pickColor(neighbors int) string {
	switch {
	case neighbors >= 6:
		return "97"
	case neighbors >= 4:
		return "96"
	case neighbors >= 2:
		return "36"
	default:
		return "37"
	}
}

func (r *Renderer) countNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if ny >= 0 && ny < len(r.buffer) && nx >= 0 && nx < len(r.buffer[0]) {
				if r.buffer[ny][nx] != ' ' {
					count++
				}
			}
		}
	}
	return count
}

func (r *Renderer) maxPossibleNeighbors(x, y int) int {
	maxNeighbors := 0
	height := len(r.buffer)
	width := len(r.buffer[0])
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if ny >= 0 && ny < height && nx >= 0 && nx < width {
				maxNeighbors++
			}
		}
	}
	return maxNeighbors
}

func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return itoa(n/10) + string(rune('0'+n%10))
}
