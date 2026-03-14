package generators

import (
	"crypto/rand"
	"math"
	mathrand "math/rand"

	"asciibloom/internal/core"
)

// Mandelbrot generates a static Mandelbrot set with animated edges.
type Mandelbrot struct {
	grid   *core.Grid
	width  int
	height int
	rng    *mathrand.Rand

	centerRe, centerIm float64
	viewSize           float64
	maxIter            int
	time               float64
}

// NewMandelbrot creates a new Mandelbrot set generator.
func NewMandelbrot(width, height int) *Mandelbrot {
	var seed int64
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err == nil {
		for i := 0; i < 8; i++ {
			seed = (seed << 8) | int64(buf[i])
		}
	}

	rng := mathrand.New(mathrand.NewSource(seed))

	m := &Mandelbrot{
		width:   width,
		height:  height,
		grid:    core.NewGrid(width, height),
		rng:     rng,
		maxIter: 60,
		time:    0,
	}
	m.reset()
	return m
}

func (m *Mandelbrot) mandelbrotIter(cr, ci float64) int {
	zr, zi := 0.0, 0.0

	for i := 0; i < m.maxIter; i++ {
		zr2 := zr * zr
		zi2 := zi * zi
		if zr2+zi2 > 4.0 {
			return i
		}
		zi = 2*zr*zi + ci
		zr = zr2 - zi2 + cr
	}
	return m.maxIter
}

func (m *Mandelbrot) reset() {
	// Different regions of the mandelbrot, all fitting on screen
	regions := []struct {
		re, im float64
		zoom   float64
		desc   string
	}{
		{-0.5, 0.0, 2.6, "main"},              // Main cardioid
		{-1.0, 0.0, 1.2, "leftbulb"},          // Left period-2 bulb
		{-0.12, 0.74, 0.8, "topbulb"},         // Top bulbs
		{-0.745, 0.105, 0.15, "seahorse"},     // Seahorse valley
		{-0.16, 1.035, 0.2, "elephant"},       // Elephant valley
		{-0.7747, 0.1242, 0.03, "minimandel"}, // Mini mandelbrot
		{-0.235, 0.827, 0.4, "dendrite"},      // Dendrite structure
		{0.28, 0.01, 0.6, "rightside"},        // Right side
		{-1.25, 0.0, 0.5, "tip"},              // Far left tip
	}

	region := regions[m.rng.Intn(len(regions))]
	m.centerRe = region.re
	m.centerIm = region.im
	m.viewSize = region.zoom
	m.time = 0
}

func (m *Mandelbrot) Step() bool {
	m.grid.Clear()
	m.time += 0.08

	aspectRatio := float64(m.width) / float64(m.height) * 2.0
	spanRe := m.viewSize
	spanIm := spanRe / aspectRatio

	currMinRe := m.centerRe - spanRe/2
	currMaxRe := m.centerRe + spanRe/2
	currMinIm := m.centerIm - spanIm/2
	currMaxIm := m.centerIm + spanIm/2

	// Compute iteration counts
	iters := make([][]int, m.height)
	for y := 0; y < m.height; y++ {
		iters[y] = make([]int, m.width)
		for x := 0; x < m.width; x++ {
			re := currMinRe + (currMaxRe-currMinRe)*float64(x)/float64(m.width-1)
			im := currMinIm + (currMaxIm-currMinIm)*float64(y)/float64(m.height-1)
			iters[y][x] = m.mandelbrotIter(re, im)
		}
	}

	// Find the boundary and animate it
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			iter := iters[y][x]

			// Inside the set - use 'O' (intensity 13)
			if iter >= m.maxIter {
				m.grid.Set(x, y, 13)
				continue
			}

			// Check for boundary: inside neighbor
			isBoundary := false
			for dy := -1; dy <= 1 && !isBoundary; dy++ {
				for dx := -1; dx <= 1 && !isBoundary; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx, ny := x+dx, y+dy
					if nx >= 0 && nx < m.width && ny >= 0 && ny < m.height {
						if iters[ny][nx] >= m.maxIter {
							isBoundary = true
						}
					}
				}
			}

			if isBoundary {
				// Animate boundary with random twinkling using full char range
				sparkle := m.rng.Float64()
				wave := math.Sin(float64(x)*0.3 + float64(y)*0.2 + m.time*2)

				// Map random+wave to full intensity range 1-13 for maximum sparkle variety
				// Excludes 0 (outside '.') and 13 (inside 'O')
				intensity := 1 + int(sparkle*10)
				if wave > 0.3 {
					intensity += 2
				}
				if intensity > 12 {
					intensity = 12
				}
				if intensity < 1 {
					intensity = 1
				}
				m.grid.Set(x, y, intensity)
			} else {
				// Outside the set - use '.' character (intensity 0)
				m.grid.Set(x, y, 0)
			}
		}
	}

	return true
}

func (m *Mandelbrot) Render(buffer [][]rune) {
	// Extended character set for more sparkle variety
	// Inside uses index 14 ('O'), outside uses index 0 ('.')
	// Boundary sparkles use full range
	chars := []rune{'.', ':', '-', '~', '=', '+', '*', '^', '%', '#', '&', '$', 'o', 'O', '@'}

	for y := 0; y < len(buffer) && y < m.height; y++ {
		for x := 0; x < len(buffer[y]) && x < m.width; x++ {
			v := m.grid.Get(x, y)
			switch {
			case v < 0:
				buffer[y][x] = ' '
			case v >= len(chars):
				buffer[y][x] = chars[len(chars)-1]
			default:
				buffer[y][x] = chars[v]
			}
		}
	}
}

func (m *Mandelbrot) ClusterSize() int {
	return int(m.time * 100)
}

func (m *Mandelbrot) PostProcess() {}
