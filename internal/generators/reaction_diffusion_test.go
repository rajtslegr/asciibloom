package generators

import (
	"testing"
)

func TestNewReactionDiffusion(t *testing.T) {
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
			rd := NewReactionDiffusion(tt.width, tt.height)

			if rd == nil {
				t.Fatal("NewReactionDiffusion() returned nil")
			}
			if rd.width != tt.width {
				t.Errorf("width = %d, want %d", rd.width, tt.width)
			}
			if rd.height != tt.height {
				t.Errorf("height = %d, want %d", rd.height, tt.height)
			}
			if rd.grid == nil {
				t.Error("grid is nil")
			}
			if rd.u == nil || rd.v == nil {
				t.Error("concentration grids are nil")
			}
		})
	}
}

func TestReactionDiffusionStep(t *testing.T) {
	rd := NewReactionDiffusion(40, 20)

	// Step should return true
	if !rd.Step() {
		t.Error("Step() returned false, expected true")
	}

	// Steps counter should increment
	if rd.steps != 1 {
		t.Errorf("steps = %d, want 1", rd.steps)
	}

	// ClusterSize should return steps
	if rd.ClusterSize() != 1 {
		t.Errorf("ClusterSize() = %d, want 1", rd.ClusterSize())
	}
}

func TestReactionDiffusionReset(t *testing.T) {
	rd := NewReactionDiffusion(40, 20)

	// Run some steps
	for i := 0; i < 10; i++ {
		rd.Step()
	}

	rd.reset()

	// After reset, steps should be 0
	if rd.steps != 0 {
		t.Errorf("steps = %d after reset, want 0", rd.steps)
	}
}

func TestReactionDiffusionRender(t *testing.T) {
	rd := NewReactionDiffusion(20, 10)

	// Run some steps to generate pattern
	for i := 0; i < 50; i++ {
		rd.Step()
	}

	buffer := make([][]rune, 10)
	for i := range buffer {
		buffer[i] = make([]rune, 20)
		for j := range buffer[i] {
			buffer[i][j] = ' '
		}
	}

	rd.Render(buffer)

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
		t.Error("Render() produced empty buffer after 50 steps")
	}
}

func TestReactionDiffusionLaplacian(t *testing.T) {
	rd := NewReactionDiffusion(20, 10)

	// Test laplacian at center with uniform grid (should be ~0)
	for y := range rd.u {
		for x := range rd.u[y] {
			rd.u[y][x] = 1.0
		}
	}

	lap := rd.laplacian(rd.u, 10, 5)
	if lap < -0.001 || lap > 0.001 {
		t.Errorf("laplacian of uniform grid = %f, expected ~0", lap)
	}
}

func TestReactionDiffusionMaxSteps(t *testing.T) {
	rd := NewReactionDiffusion(20, 10)
	rd.maxSteps = 5

	// Run to max steps
	for i := 0; i < 5; i++ {
		rd.Step()
	}

	// Next step should trigger reset
	rd.Step()

	if rd.steps != 1 {
		t.Errorf("steps = %d after reset, expected reset to 1", rd.steps)
	}
}
