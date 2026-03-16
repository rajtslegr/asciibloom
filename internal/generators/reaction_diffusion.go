package generators

import (
	"crypto/rand"
	"math"
	mathrand "math/rand"

	"asciibloom/internal/core"
)

// ReactionDiffusion implements the Gray-Scott reaction-diffusion model.
// It simulates the interaction between two chemical substances U and V,
// creating organic patterns like spots, stripes, and coral-like structures.
type ReactionDiffusion struct {
	grid   *core.Grid
	width  int
	height int
	rng    *mathrand.Rand

	// Chemical concentrations
	u [][]float64
	v [][]float64

	// Diffusion rates
	du float64
	dv float64

	// Feed and kill rates (determine pattern type)
	f float64
	k float64

	// Timestep
	dt float64

	// Pattern iteration counter
	steps    int
	maxSteps int
}

// NewReactionDiffusion creates a new reaction-diffusion generator.
func NewReactionDiffusion(width, height int) *ReactionDiffusion {
	var seed int64
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err == nil {
		for i := 0; i < 8; i++ {
			seed = (seed << 8) | int64(buf[i])
		}
	}

	r := &ReactionDiffusion{
		width:    width,
		height:   height,
		grid:     core.NewGrid(width, height),
		rng:      mathrand.New(mathrand.NewSource(seed)),
		du:       0.16,
		dv:       0.08,
		dt:       1.0,
		steps:    0,
		maxSteps: 2000,
	}
	r.reset()
	return r
}

func (r *ReactionDiffusion) reset() {
	// Initialize concentration grids
	r.u = make([][]float64, r.height)
	r.v = make([][]float64, r.height)
	for y := range r.u {
		r.u[y] = make([]float64, r.width)
		r.v[y] = make([]float64, r.width)
	}

	// Set up different pattern parameters (F, k pairs)
	patterns := []struct {
		f, k float64
		desc string
	}{
		{0.035, 0.060, "coral"},       // Coral growth
		{0.025, 0.060, "spots"},       // Spots
		{0.055, 0.062, "maze"},        // Maze-like
		{0.030, 0.062, "waves"},       // Waves
		{0.039, 0.058, "fingerprint"}, // Fingerprint
		{0.042, 0.061, "chaos"},       // Chaotic
	}

	pattern := patterns[r.rng.Intn(len(patterns))]
	r.f = pattern.f
	r.k = pattern.k

	// Initialize with uniform U and small V perturbations
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			r.u[y][x] = 1.0
			r.v[y][x] = 0.0
		}
	}

	// Create seed regions for V with varied shapes
	numSeeds := 4 + r.rng.Intn(4)
	for i := 0; i < numSeeds; i++ {
		cx := r.rng.Intn(r.width-10) + 5
		cy := r.rng.Intn(r.height-6) + 3
		size := 2 + r.rng.Intn(4)

		// Random shape: circle, square, or cross
		shape := r.rng.Intn(3)
		for dy := -size; dy <= size; dy++ {
			for dx := -size; dx <= size; dx++ {
				nx, ny := cx+dx, cy+dy
				if nx < 0 || nx >= r.width || ny < 0 || ny >= r.height {
					continue
				}

				var place bool
				switch shape {
				case 0: // Circle
					place = dx*dx+dy*dy <= size*size
				case 1: // Square
					place = true
				case 2: // Cross
					place = dx == 0 || dy == 0 || (dx == dy) || (dx == -dy)
				}

				if place {
					// Add some noise for organic feel
					noise := r.rng.Float64() * 0.3
					r.v[ny][nx] = 0.8 + noise
					r.u[ny][nx] = 0.2 - noise
				}
			}
		}
	}

	// Add background noise for texture
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			if r.rng.Float64() < 0.02 {
				r.v[y][x] = r.rng.Float64() * 0.3
				r.u[y][x] = 1.0 - r.v[y][x]
			}
		}
	}

	r.steps = 0
	r.grid.Clear()
}

// laplacian computes the discrete Laplacian at position (x, y) for the given grid.
func (r *ReactionDiffusion) laplacian(grid [][]float64, x, y int) float64 {
	// Use a 3x3 kernel with weights:
	// 0.05 0.2 0.05
	// 0.2  -1  0.2
	// 0.05 0.2 0.05
	var sum float64
	sum += grid[y][x] * -1.0

	// Direct neighbors (weight 0.2)
	if x > 0 {
		sum += grid[y][x-1] * 0.2
	}
	if x < r.width-1 {
		sum += grid[y][x+1] * 0.2
	}
	if y > 0 {
		sum += grid[y-1][x] * 0.2
	}
	if y < r.height-1 {
		sum += grid[y+1][x] * 0.2
	}

	// Diagonal neighbors (weight 0.05)
	if x > 0 && y > 0 {
		sum += grid[y-1][x-1] * 0.05
	}
	if x < r.width-1 && y > 0 {
		sum += grid[y-1][x+1] * 0.05
	}
	if x > 0 && y < r.height-1 {
		sum += grid[y+1][x-1] * 0.05
	}
	if x < r.width-1 && y < r.height-1 {
		sum += grid[y+1][x+1] * 0.05
	}

	return sum
}

// Step performs one iteration of the reaction-diffusion simulation.
func (r *ReactionDiffusion) Step() bool {
	if r.steps >= r.maxSteps {
		r.reset()
	}

	// Create temporary grids for the next state
	un := make([][]float64, r.height)
	vn := make([][]float64, r.height)
	for y := range un {
		un[y] = make([]float64, r.width)
		vn[y] = make([]float64, r.width)
	}

	// Compute next state
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			u := r.u[y][x]
			v := r.v[y][x]

			// Gray-Scott reaction terms
			uv2 := u * v * v

			// Laplacian for diffusion
			lu := r.laplacian(r.u, x, y)
			lv := r.laplacian(r.v, x, y)

			// Update equations
			un[y][x] = u + r.dt*(r.du*lu-uv2+r.f*(1.0-u))
			vn[y][x] = v + r.dt*(r.dv*lv+uv2-(r.f+r.k)*v)

			// Clamp values
			if un[y][x] < 0 {
				un[y][x] = 0
			}
			if un[y][x] > 1 {
				un[y][x] = 1
			}
			if vn[y][x] < 0 {
				vn[y][x] = 0
			}
			if vn[y][x] > 1 {
				vn[y][x] = 1
			}
		}
	}

	// Swap grids
	r.u = un
	r.v = vn

	r.steps++

	// Update the grid for rendering (map V concentration to intensity 0-14)
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			v := r.v[y][x]
			// Map V concentration (0-1) to intensity (0-14) for more variety
			intensity := int(v * 14)
			if intensity > 14 {
				intensity = 14
			}
			if intensity < 0 {
				intensity = 0
			}
			r.grid.Set(x, y, intensity)
		}
	}

	return true
}

// Render writes the current state to the provided buffer with enhanced visuals.
func (r *ReactionDiffusion) Render(buffer [][]rune) {
	// Extended character palette for richer visualization
	// Ordered from low to high intensity
	chars := []rune{' ', '.', ':', '-', '~', '=', '+', '*', '^', '%', '#', '&', '$', 'O', '@'}

	for y := 0; y < len(buffer) && y < r.height; y++ {
		for x := 0; x < len(buffer[y]) && x < r.width; x++ {
			v := r.v[y][x]

			if v < 0.15 {
				buffer[y][x] = ' '
				continue
			}

			// Calculate gradient for edge detection
			gradX, gradY := r.gradient(x, y)
			gradientMag := math.Sqrt(gradX*gradX + gradY*gradY)

			// Base intensity from V concentration
			baseIntensity := int(v * float64(len(chars)-1))
			if baseIntensity >= len(chars) {
				baseIntensity = len(chars) - 1
			}

			// Add sparkle on edges (high gradient areas)
			if gradientMag > 0.3 && r.rng.Float64() < 0.15 {
				// Random high-intensity character for sparkle effect
				sparkleIdx := len(chars) - 1 - r.rng.Intn(3)
				buffer[y][x] = chars[sparkleIdx]
			} else {
				buffer[y][x] = chars[baseIntensity]
			}
		}
	}
}

// gradient computes the gradient magnitude at position (x, y).
func (r *ReactionDiffusion) gradient(x, y int) (float64, float64) {
	var dx, dy float64

	if x > 0 {
		dx = r.v[y][x] - r.v[y][x-1]
	}
	if x < r.width-1 {
		dx = (dx + r.v[y][x+1] - r.v[y][x]) / 2
	}

	if y > 0 {
		dy = r.v[y][x] - r.v[y-1][x]
	}
	if y < r.height-1 {
		dy = (dy + r.v[y+1][x] - r.v[y][x]) / 2
	}

	return dx, dy
}

// ClusterSize returns the current step count.
func (r *ReactionDiffusion) ClusterSize() int {
	return r.steps
}

// PostProcess performs any post-processing after rendering.
func (r *ReactionDiffusion) PostProcess() {}
