package generators

import (
	"testing"
)

func TestNewFlowField(t *testing.T) {
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
			ff := NewFlowField(tt.width, tt.height)

			if ff == nil {
				t.Fatal("NewFlowField() returned nil")
			}
			if ff.width != tt.width {
				t.Errorf("width = %d, want %d", ff.width, tt.width)
			}
			if ff.height != tt.height {
				t.Errorf("height = %d, want %d", ff.height, tt.height)
			}
			if ff.grid == nil {
				t.Error("grid is nil")
			}
			if ff.grid.Width() != tt.width {
				t.Errorf("grid width = %d, want %d", ff.grid.Width(), tt.width)
			}
		})
	}
}

func TestFlowFieldStep(t *testing.T) {
	ff := NewFlowField(40, 20)

	// Step should return true initially
	if !ff.Step() {
		t.Error("Step() returned false, expected true on first call")
	}

	// After stepping, should have particles or trails
	if ff.ClusterSize() < 0 {
		t.Error("ClusterSize() < 0 after Step()")
	}
}

func TestFlowFieldReset(t *testing.T) {
	ff := NewFlowField(40, 20)

	// Generate some particles
	for i := 0; i < 10; i++ {
		ff.Step()
	}

	ff.reset()

	// After reset, particles should be empty but grid should have seeds
	if len(ff.particles) != 0 {
		t.Error("Particles not cleared after reset")
	}
}

func TestFlowFieldRender(t *testing.T) {
	ff := NewFlowField(20, 10)

	// Generate some particles
	ff.Step()

	buffer := make([][]rune, 10)
	for i := range buffer {
		buffer[i] = make([]rune, 20)
		for j := range buffer[i] {
			buffer[i][j] = ' '
		}
	}

	ff.Render(buffer)

	// Buffer should have some non-space characters (from seed points at minimum)
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

func TestFlowFieldNoise2D(t *testing.T) {
	ff := NewFlowField(40, 20)

	// Noise should return consistent values for same inputs
	n1 := ff.noise2D(10.0, 10.0)
	n2 := ff.noise2D(10.0, 10.0)

	if n1 != n2 {
		t.Error("noise2D not deterministic for same inputs")
	}

	// Noise should return different values for different inputs (usually)
	n3 := ff.noise2D(20.0, 20.0)
	if n1 == n3 {
		t.Error("noise2D returned same value for different inputs")
	}
}

func TestFlowFieldGetFlowVector(t *testing.T) {
	ff := NewFlowField(40, 20)

	vx, vy := ff.getFlowVector(10.0, 10.0)

	// Vector should be normalized-ish
	magnitude := vx*vx + vy*vy
	if magnitude < 0.5 || magnitude > 2.0 {
		t.Errorf("flow vector magnitude = %f, expected around 1.0", magnitude)
	}
}
