package generators

import (
	"testing"
)

func TestNewBrownianTree(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "standard size",
			width:  80,
			height: 24,
		},
		{
			name:   "square grid",
			width:  50,
			height: 50,
		},
		{
			name:   "small grid",
			width:  30,
			height: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt := NewBrownianTree(tt.width, tt.height)

			if bt == nil {
				t.Fatal("NewBrownianTree() returned nil")
			}
			if bt.width != tt.width {
				t.Errorf("width = %d, want %d", bt.width, tt.width)
			}
			if bt.height != tt.height {
				t.Errorf("height = %d, want %d", bt.height, tt.height)
			}
			if bt.grid == nil {
				t.Error("grid is nil")
			}
			if len(bt.grid) != tt.height {
				t.Errorf("grid height = %d, want %d", len(bt.grid), tt.height)
			}
		})
	}
}

func TestBrownianTreeStep(t *testing.T) {
	bt := NewBrownianTree(40, 20)

	// Step should return true initially
	if !bt.Step() {
		t.Error("Step() returned false, expected true on first call")
	}

	// After stepping, cluster size should be greater than 0
	if bt.ClusterSize() <= 0 {
		t.Error("ClusterSize() <= 0 after Step()")
	}
}

func TestBrownianTreeReset(t *testing.T) {
	bt := NewBrownianTree(40, 20)

	// Generate some particles
	for i := 0; i < 10; i++ {
		bt.Step()
	}

	initialSize := bt.ClusterSize()
	if initialSize == 0 {
		t.Fatal("No particles generated")
	}

	// Reset - this clears the grid and re-adds seed points
	bt.reset()

	// After reset, cluster should have seed points (not zero)
	if bt.ClusterSize() == 0 {
		t.Error("After reset, ClusterSize() = 0, expected seed points")
	}

	// Grid should have seed points
	hasNonZero := false
	for y := range bt.grid {
		for x := range bt.grid[y] {
			if bt.grid[y][x] > 0 {
				hasNonZero = true
				break
			}
		}
	}
	if !hasNonZero {
		t.Error("Seed points not preserved after reset")
	}
}

func TestBrownianTreeRender(t *testing.T) {
	bt := NewBrownianTree(20, 10)

	// Generate some particles
	bt.Step()

	buffer := make([][]rune, 10)
	for i := range buffer {
		buffer[i] = make([]rune, 20)
		for j := range buffer[i] {
			buffer[i][j] = ' '
		}
	}

	bt.Render(buffer)

	// Buffer should have some non-space characters
	hasContent := false
	for y := range buffer {
		for x := range buffer[y] {
			if buffer[y][x] != ' ' {
				hasContent = true
				break
			}
		}
		if hasContent {
			break
		}
	}

	if !hasContent {
		t.Error("Render() produced empty buffer")
	}
}

func TestBrownianTreeGetSet(t *testing.T) {
	bt := NewBrownianTree(20, 10)

	tests := []struct {
		name string
		x    int
		y    int
		val  int
		want int
	}{
		{
			name: "in bounds",
			x:    5,
			y:    5,
			val:  3,
			want: 3,
		},
		{
			name: "out of bounds negative",
			x:    -1,
			y:    5,
			val:  3,
			want: 0,
		},
		{
			name: "out of bounds positive",
			x:    25,
			y:    5,
			val:  3,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt.set(tt.x, tt.y, tt.val)
			got := bt.get(tt.x, tt.y)
			if got != tt.want {
				t.Errorf("get(%d, %d) = %d, want %d", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestBrownianTreeNeighbors(t *testing.T) {
	bt := NewBrownianTree(20, 10)

	// Set up a pattern
	bt.set(5, 5, 5)
	bt.set(5, 6, 5)
	bt.set(6, 5, 5)

	tests := []struct {
		name string
		x    int
		y    int
		want int
	}{
		{
			name: "center",
			x:    5,
			y:    5,
			want: 2, // 2 neighbors
		},
		{
			name: "neighbor",
			x:    5,
			y:    6,
			want: 2, // neighbors at (5, 5) and (6, 5)
		},
		{
			name: "empty",
			x:    10,
			y:    10,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bt.neighbors(tt.x, tt.y)
			if got != tt.want {
				t.Errorf("neighbors(%d, %d) = %d, want %d", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestBrownianTreeCenter(t *testing.T) {
	bt := NewBrownianTree(40, 20)

	cx, cy := bt.center()

	if cx < 0 || cx > float64(bt.width) {
		t.Errorf("center x = %f, out of bounds [0, %d]", cx, bt.width)
	}
	if cy < 0 || cy > float64(bt.height) {
		t.Errorf("center y = %f, out of bounds [0, %d]", cy, bt.height)
	}
}

func TestBrownianTreeRadius(t *testing.T) {
	bt := NewBrownianTree(40, 20)

	r := bt.radius()

	if r < 10 {
		t.Errorf("radius() = %f, want >= 10", r)
	}
}
