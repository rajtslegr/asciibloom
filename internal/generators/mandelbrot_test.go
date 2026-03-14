package generators

import (
	"testing"
)

func TestNewMandelbrot(t *testing.T) {
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
			m := NewMandelbrot(tt.width, tt.height)

			if m == nil {
				t.Fatal("NewMandelbrot() returned nil")
			}
			if m.width != tt.width {
				t.Errorf("width = %d, want %d", m.width, tt.width)
			}
			if m.height != tt.height {
				t.Errorf("height = %d, want %d", m.height, tt.height)
			}
			if m.grid == nil {
				t.Error("grid is nil")
			}
			if m.grid.Height() != tt.height {
				t.Errorf("grid height = %d, want %d", m.grid.Height(), tt.height)
			}
			if m.grid.Width() != tt.width {
				t.Errorf("grid width = %d, want %d", m.grid.Width(), tt.width)
			}
			if m.rng == nil {
				t.Error("rng is nil")
			}
		})
	}
}

func TestMandelbrotStep(t *testing.T) {
	m := NewMandelbrot(40, 20)

	// Step should always return true for Mandelbrot
	if !m.Step() {
		t.Error("Step() returned false, expected true")
	}

	// After stepping, cluster size should be greater than 0 (time based)
	if m.ClusterSize() < 0 {
		t.Error("ClusterSize() < 0 after Step()")
	}

	// Grid should have some non-zero values
	hasNonZero := false
	for y := 0; y < m.grid.Height(); y++ {
		for x := 0; x < m.grid.Width(); x++ {
			if m.grid.Get(x, y) > 0 {
				hasNonZero = true
				break
			}
		}
		if hasNonZero {
			break
		}
	}

	if !hasNonZero {
		t.Error("Step() did not set any grid values")
	}
}

func TestMandelbrotRender(t *testing.T) {
	m := NewMandelbrot(20, 10)

	// Generate some content
	m.Step()

	buffer := make([][]rune, 10)
	for i := range buffer {
		buffer[i] = make([]rune, 20)
		for j := range buffer[i] {
			buffer[i][j] = ' '
		}
	}

	m.Render(buffer)

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

	// Check that we have the expected character set
	validChars := map[rune]bool{
		'.': true, ':': true, '-': true, '~': true, '=': true,
		'+': true, '*': true, '^': true, '%': true, '#': true,
		'&': true, '$': true, 'o': true, 'O': true, '@': true,
	}

	for y := range buffer {
		for x := range buffer[y] {
			char := buffer[y][x]
			if char != ' ' && !validChars[char] {
				t.Errorf("Render() produced invalid character: %c", char)
			}
		}
	}
}

func TestMandelbrotReset(t *testing.T) {
	m := NewMandelbrot(40, 20)

	// Step a few times
	for i := 0; i < 5; i++ {
		m.Step()
	}

	initialTime := m.time

	// Reset
	m.reset()

	// After reset, time should be 0
	if m.time != 0 {
		t.Errorf("After reset, time = %f, want 0", m.time)
	}

	// View size should be reset
	if m.viewSize == 0 {
		t.Error("After reset, viewSize = 0")
	}

	// Time should be different from before (if we stepped)
	if initialTime == 0 {
		t.Error("Initial time was 0 before steps")
	}
}

func TestMandelbrotMandelbrotIter(t *testing.T) {
	m := NewMandelbrot(40, 20)

	tests := []struct {
		name     string
		cr, ci   float64
		expected int
	}{
		{
			name:     "center of main cardioid",
			cr:       0.0,
			ci:       0.0,
			expected: 60, // Should reach maxIter (inside set)
		},
		{
			name:     "outside set",
			cr:       2.0,
			ci:       2.0,
			expected: 1, // Should escape quickly
		},
		{
			name:     "left bulb",
			cr:       -1.0,
			ci:       0.0,
			expected: 60, // Inside period-2 bulb
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.mandelbrotIter(tt.cr, tt.ci)
			// For points inside, we expect maxIter
			// For points outside, we expect less than maxIter
			if tt.expected == 60 && got != tt.expected {
				t.Errorf("mandelbrotIter(%f, %f) = %d, want %d", tt.cr, tt.ci, got, tt.expected)
			}
			if tt.expected == 1 && got >= 60 {
				t.Errorf("mandelbrotIter(%f, %f) = %d, expected to escape quickly", tt.cr, tt.ci, got)
			}
		})
	}
}

func TestMandelbrotPostProcess(t *testing.T) {
	m := NewMandelbrot(40, 20)

	// PostProcess should not panic
	m.PostProcess()
}

func TestMandelbrotDifferentRegions(t *testing.T) {
	// Run multiple times to potentially hit different regions
	for i := 0; i < 20; i++ {
		m := NewMandelbrot(40, 20)
		m.Step()

		// Each run should produce a valid mandelbrot
		// (either filled or detailed depending on region)
		hasInside := false
		hasOutside := false

		for y := 0; y < m.grid.Height(); y++ {
			for x := 0; x < m.grid.Width(); x++ {
				val := m.grid.Get(x, y)
				if val == 13 { // Inside set
					hasInside = true
				}
				if val == 0 { // Outside set
					hasOutside = true
				}
			}
		}

		// Should have both inside and outside
		if !hasInside {
			t.Error("Step() produced grid with no inside set values")
		}
		if !hasOutside {
			t.Error("Step() produced grid with no outside set values")
		}
	}
}
