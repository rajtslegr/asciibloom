package animation

import (
	"testing"
)

func TestTerminalDimensions(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "standard terminal",
			width:  80,
			height: 24,
		},
		{
			name:   "wide terminal",
			width:  120,
			height: 30,
		},
		{
			name:   "small terminal",
			width:  40,
			height: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terminal := &Terminal{
				width:  tt.width,
				height: tt.height,
			}

			if got := terminal.Width(); got != tt.width {
				t.Errorf("Width() = %d, want %d", got, tt.width)
			}
			if got := terminal.Height(); got != tt.height {
				t.Errorf("Height() = %d, want %d", got, tt.height)
			}
		})
	}
}

func TestStringWidth(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{
			name: "ascii characters",
			s:    "hello",
			want: 5,
		},
		{
			name: "unicode characters",
			s:    "こんにちは",
			want: 5,
		},
		{
			name: "with spaces",
			s:    "  hello  ",
			want: 5,
		},
		{
			name: "empty string",
			s:    "",
			want: 0,
		},
		{
			name: "only spaces",
			s:    "   ",
			want: 0,
		},
	}

	terminal := &Terminal{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := terminal.StringWidth(tt.s)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}
