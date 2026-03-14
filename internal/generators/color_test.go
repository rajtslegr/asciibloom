package generators

import (
	"strings"
	"testing"
)

func TestNewColorGrid(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := NewColorGrid(tt.width, tt.height)

			if cg == nil {
				t.Fatal("NewColorGrid() returned nil")
			}
			if cg.Width != tt.width {
				t.Errorf("Width = %d, want %d", cg.Width, tt.width)
			}
			if cg.Height != tt.height {
				t.Errorf("Height = %d, want %d", cg.Height, tt.height)
			}
			if len(cg.Chars) != tt.height {
				t.Errorf("Chars height = %d, want %d", len(cg.Chars), tt.height)
			}
			if len(cg.Intensities) != tt.height {
				t.Errorf("Intensities height = %d, want %d", len(cg.Intensities), tt.height)
			}
		})
	}
}

func TestBrownianPalette(t *testing.T) {
	palette := BrownianPalette()

	if len(palette) == 0 {
		t.Error("BrownianPalette() returned empty palette")
	}

	// Check that all colors are valid grayscale-ish
	for i, c := range palette {
		if c.R < 0 || c.R > 255 {
			t.Errorf("palette[%d].R = %d, want [0, 255]", i, c.R)
		}
		if c.G < 0 || c.G > 255 {
			t.Errorf("palette[%d].G = %d, want [0, 255]", i, c.G)
		}
		if c.B < 0 || c.B > 255 {
			t.Errorf("palette[%d].B = %d, want [0, 255]", i, c.B)
		}
	}
}

func TestGrayCode(t *testing.T) {
	tests := []struct {
		name       string
		brightness int
		wantStart  string
	}{
		{
			name:       "minimum",
			brightness: 0,
			wantStart:  "232",
		},
		{
			name:       "maximum",
			brightness: 255,
			wantStart:  "255",
		},
		{
			name:       "negative (clamped)",
			brightness: -10,
			wantStart:  "232",
		},
		{
			name:       "over max (clamped)",
			brightness: 300,
			wantStart:  "255",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := grayCode(tt.brightness)
			if got != tt.wantStart {
				t.Errorf("grayCode(%d) = %s, want %s", tt.brightness, got, tt.wantStart)
			}
		})
	}
}

func TestIntToStr(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{
			name: "single digit",
			n:    5,
			want: "5",
		},
		{
			name: "double digit",
			n:    42,
			want: "42",
		},
		{
			name: "triple digit",
			n:    123,
			want: "123",
		},
		{
			name: "zero",
			n:    0,
			want: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := intToStr(tt.n)
			if got != tt.want {
				t.Errorf("intToStr(%d) = %s, want %s", tt.n, got, tt.want)
			}
		})
	}
}

func TestColorizedOutput(t *testing.T) {
	cg := NewColorGrid(10, 5)
	palette := BrownianPalette()

	// Add some content
	cg.Chars[2][3] = '*'
	cg.Intensities[2][3] = 0.5

	output := cg.ColorizedOutput(palette)

	// Output should contain ANSI escape sequences
	if !strings.Contains(output, "\x1b[") {
		t.Error("ColorizedOutput() missing ANSI escape sequences")
	}

	// Should contain the character we set
	if !strings.Contains(output, "*") {
		t.Error("ColorizedOutput() missing expected character")
	}

	// Should reset colors at the end
	if !strings.HasSuffix(output, "\x1b[39m") {
		t.Error("ColorizedOutput() missing color reset")
	}
}
