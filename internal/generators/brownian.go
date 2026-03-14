package generators

import (
	"crypto/rand"
	"math"
	mathrand "math/rand"
)

// BrownianTree generates a Brownian tree pattern using particle aggregation.
type BrownianTree struct {
	grid          [][]int
	width         int
	height        int
	particleCount int
	rng           *mathrand.Rand

	minX, maxX int
	minY, maxY int

	seedPoints [][2]int
	maxSize    int
}

// NewBrownianTree creates a new Brownian tree generator with the given dimensions.
func NewBrownianTree(width, height int) *BrownianTree {
	var seed int64
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err == nil {
		for i := 0; i < 8; i++ {
			seed = (seed << 8) | int64(buf[i])
		}
	}

	b := &BrownianTree{
		width:  width,
		height: height,
		rng:    mathrand.New(mathrand.NewSource(seed)),
	}
	b.init()
	return b
}

func (b *BrownianTree) init() {
	b.grid = make([][]int, b.height)
	for y := range b.grid {
		b.grid[y] = make([]int, b.width)
	}

	b.seedPoints = make([][2]int, 0)
	seeds := 1 + b.rng.Intn(3)
	for i := 0; i < seeds; i++ {
		cx := b.width / 2
		cy := b.height / 2
		if b.width > 20 {
			cx = b.rng.Intn(b.width-20) + 10
		}
		if b.height > 10 {
			cy = b.rng.Intn(b.height-10) + 5
		}
		b.seedPoints = append(b.seedPoints, [2]int{cx, cy})
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				b.set(cx+dx, cy+dy, 5)
			}
		}
	}

	b.minX, b.maxX = b.width, 0
	b.minY, b.maxY = b.height, 0
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			if b.grid[y][x] > 0 {
				if x < b.minX {
					b.minX = x
				}
				if x > b.maxX {
					b.maxX = x
				}
				if y < b.minY {
					b.minY = y
				}
				if y > b.maxY {
					b.maxY = y
				}
			}
		}
	}

	b.maxSize = b.width * b.height / 12
}

func (b *BrownianTree) reset() {
	for y := range b.grid {
		for x := range b.grid[y] {
			b.grid[y][x] = 0
		}
	}
	b.particleCount = 0
	b.minX, b.maxX = b.width, 0
	b.minY, b.maxY = b.height, 0

	for _, seed := range b.seedPoints {
		cx, cy := seed[0], seed[1]
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				b.set(cx+dx, cy+dy, 5)
			}
		}
	}
}

func (b *BrownianTree) set(x, y, val int) {
	if x >= 0 && x < b.width && y >= 0 && y < b.height && b.grid[y][x] < val {
		b.grid[y][x] = val
		b.particleCount++
	}
}

func (b *BrownianTree) get(x, y int) int {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return 0
	}
	return b.grid[y][x]
}

func (b *BrownianTree) neighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if (dx != 0 || dy != 0) && b.get(x+dx, y+dy) > 0 {
				count++
			}
		}
	}
	return count
}

func (b *BrownianTree) center() (float64, float64) {
	cx := float64(b.minX+b.maxX) / 2
	cy := float64(b.minY+b.maxY) / 2
	if cx < 1 {
		cx = float64(b.width) / 2
	}
	if cy < 1 {
		cy = float64(b.height) / 2
	}
	return cx, cy
}

func (b *BrownianTree) radius() float64 {
	dx := float64(b.maxX - b.minX)
	dy := float64(b.maxY - b.minY)
	r := math.Max(dx, dy) / 2
	if r < 10 {
		r = 10
	}
	return r
}

// Step performs one iteration of particle simulation.
// Returns true if particles are still being generated.
func (b *BrownianTree) Step() bool {
	if b.particleCount >= b.maxSize {
		b.reset()
	}

	cx, cy := b.center()
	r := b.radius()

	for i := 0; i < 50+b.rng.Intn(30); i++ {
		spawnR := r + 15 + b.rng.Float64()*25
		angle := b.rng.Float64() * 2 * math.Pi
		x := int(cx + math.Cos(angle)*spawnR)
		y := int(cy + math.Sin(angle)*spawnR)

		for step := 0; step < 3000; step++ {
			switch b.rng.Intn(8) {
			case 0:
				x--
			case 1:
				x++
			case 2:
				y--
			case 3:
				y++
			case 4:
				x--
				y--
			case 5:
				x++
				y--
			case 6:
				x--
				y++
			case 7:
				x++
				y++
			}

			if x < 0 || x >= b.width || y < 0 || y >= b.height {
				break
			}

			n := b.neighbors(x, y)
			if n > 0 {
				p := 0.12 + float64(n)*0.06
				if b.rng.Float64() < p {
					val := n + 2
					if val > 7 {
						val = 7
					}
					b.set(x, y, val)

					if x < b.minX {
						b.minX = x
					}
					if x > b.maxX {
						b.maxX = x
					}
					if y < b.minY {
						b.minY = y
					}
					if y > b.maxY {
						b.maxY = y
					}
					break
				}
			}
		}
	}

	return b.particleCount > 0
}

// Render writes the current tree state to the provided buffer.
func (b *BrownianTree) Render(buffer [][]rune) {
	for y := 0; y < len(buffer) && y < b.height; y++ {
		for x := 0; x < len(buffer[y]) && x < b.width; x++ {
			v := b.grid[y][x]
			if v == 0 {
				buffer[y][x] = ' '
				continue
			}

			n := b.neighbors(x, y)

			var ch rune
			switch {
			case n >= 6:
				ch = 'O'
			case n >= 4:
				ch = 'o'
			case n >= 3:
				ch = '*'
			case n >= 2:
				ch = ':'
			case v >= 5:
				ch = '+'
			default:
				ch = '.'
			}
			buffer[y][x] = ch
		}
	}
}

// ClusterSize returns the current number of particles in the tree.
func (b *BrownianTree) ClusterSize() int {
	return b.particleCount
}

// PostProcess performs any post-processing after rendering.
func (b *BrownianTree) PostProcess() {}
