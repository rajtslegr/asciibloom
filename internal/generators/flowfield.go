package generators

import (
	"crypto/rand"
	"math"
	mathrand "math/rand"

	"asciibloom/internal/core"
)

// FlowField simulates particles flowing through a Perlin-like noise field,
// creating organic stream patterns similar to water currents or wind flows.
type FlowField struct {
	grid          *core.Grid
	particles     []particle
	particleCount int
	rng           *mathrand.Rand
	width         int
	height        int
	noiseScale    float64
	timeOffset    float64
	maxParticles  int
}

// particle represents a single flowing particle with position and velocity.
type particle struct {
	x, y    float64
	vx, vy  float64
	life    int
	maxLife int
}

// NewFlowField creates a new flow field generator with the given dimensions.
func NewFlowField(width, height int) *FlowField {
	var seed int64
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err == nil {
		for i := 0; i < 8; i++ {
			seed = (seed << 8) | int64(buf[i])
		}
	}

	f := &FlowField{
		width:        width,
		height:       height,
		grid:         core.NewGrid(width, height),
		rng:          mathrand.New(mathrand.NewSource(seed)),
		noiseScale:   0.05,
		maxParticles: width * height / 8,
		particles:    make([]particle, 0, 100),
	}
	f.init()
	return f
}

func (f *FlowField) init() {
	// Initialize with some seed particles in interesting patterns
	numSeeds := 3 + f.rng.Intn(4)
	for i := 0; i < numSeeds; i++ {
		cx := f.width / 2
		cy := f.height / 2
		if f.width > 20 {
			cx = f.rng.Intn(f.width-20) + 10
		}
		if f.height > 10 {
			cy = f.rng.Intn(f.height-10) + 5
		}

		// Create seed clusters
		for j := 0; j < 8; j++ {
			offsetX := f.rng.Intn(5) - 2
			offsetY := f.rng.Intn(5) - 2
			f.grid.Set(cx+offsetX, cy+offsetY, 7)
		}
	}
}

// noise2D generates a simplified noise value for flow field calculation.
// Uses sine wave combinations for organic, deterministic patterns.
func (f *FlowField) noise2D(x, y float64) float64 {
	sx := x * f.noiseScale
	sy := y * f.noiseScale
	t := f.timeOffset * 0.01

	// Combine multiple sine waves for organic patterns
	n1 := math.Sin(sx*1.5+sy*0.5+t) * math.Cos(sy*1.2-sx*0.3+t)
	n2 := math.Sin(sx*0.7-sy*1.3+t*0.5) * 0.5
	n3 := math.Cos(sx*2.0+sy*0.8-t*0.3) * 0.25

	return (n1 + n2 + n3) * 2 * math.Pi
}

// getFlowVector returns the flow direction at position (x, y).
func (f *FlowField) getFlowVector(x, y float64) (float64, float64) {
	angle := f.noise2D(x, y)
	vx := math.Cos(angle)
	vy := math.Sin(angle)
	return vx, vy
}

// spawnParticle creates a new particle at the edge or randomly.
func (f *FlowField) spawnParticle() particle {
	var x, y float64

	// Spawn from edges based on flow direction
	side := f.rng.Intn(4)
	switch side {
	case 0: // top
		x = float64(f.rng.Intn(f.width))
		y = 0
	case 1: // right
		x = float64(f.width - 1)
		y = float64(f.rng.Intn(f.height))
	case 2: // bottom
		x = float64(f.rng.Intn(f.width))
		y = float64(f.height - 1)
	case 3: // left
		x = 0
		y = float64(f.rng.Intn(f.height))
	}

	vx, vy := f.getFlowVector(x, y)

	return particle{
		x:       x,
		y:       y,
		vx:      vx,
		vy:      vy,
		life:    0,
		maxLife: 50 + f.rng.Intn(100),
	}
}

// Step performs one iteration of flow field simulation.
// Returns true if simulation is active.
func (f *FlowField) Step() bool {
	if f.particleCount >= f.maxParticles {
		f.reset()
	}

	// Update time for animated flow
	f.timeOffset += 0.5

	// Spawn new particles
	numNew := 3 + f.rng.Intn(5)
	for i := 0; i < numNew && len(f.particles) < 50; i++ {
		f.particles = append(f.particles, f.spawnParticle())
	}

	// Update existing particles
	newParticles := make([]particle, 0, len(f.particles))

	for _, p := range f.particles {
		// Update velocity based on flow field
		vx, vy := f.getFlowVector(p.x, p.y)

		// Smooth interpolation between current and target velocity
		p.vx = p.vx*0.7 + vx*0.3
		p.vy = p.vy*0.7 + vy*0.3

		// Move particle
		speed := 0.8 + f.rng.Float64()*0.4
		p.x += p.vx * speed
		p.y += p.vy * speed
		p.life++

		// Check if particle is still valid
		ix, iy := int(p.x), int(p.y)
		if p.life >= p.maxLife || ix < 0 || ix >= f.width || iy < 0 || iy >= f.height {
			continue
		}

		// Leave trail on grid
		if f.rng.Float64() < 0.3 {
			neighbors := f.grid.CountNeighbors(ix, iy)
			intensity := 3 + neighbors
			if intensity > 7 {
				intensity = 7
			}
			f.grid.Set(ix, iy, intensity)
			f.particleCount++
		}

		newParticles = append(newParticles, p)
	}

	f.particles = newParticles

	return f.particleCount > 0
}

func (f *FlowField) reset() {
	f.grid.Clear()
	f.particles = f.particles[:0]
	f.particleCount = 0
	f.timeOffset = 0
	f.init()
}

// Render writes the current flow field state to the provided buffer.
func (f *FlowField) Render(buffer [][]rune) {
	for y := 0; y < len(buffer) && y < f.height; y++ {
		for x := 0; x < len(buffer[y]) && x < f.width; x++ {
			v := f.grid.Get(x, y)
			if v == 0 {
				buffer[y][x] = ' '
				continue
			}

			n := f.grid.CountNeighbors(x, y)
			buffer[y][x] = core.CharForIntensity(n, v)
		}
	}
}

// ClusterSize returns the current number of trail particles.
func (f *FlowField) ClusterSize() int {
	return f.particleCount
}

// PostProcess performs any post-processing after rendering.
func (f *FlowField) PostProcess() {}
