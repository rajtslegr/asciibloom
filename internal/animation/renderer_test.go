package animation

import (
	"context"
	"testing"
	"time"
)

func TestNewRenderer(t *testing.T) {
	tests := []struct {
		name      string
		width     int
		height    int
		wantPanic bool
	}{
		{
			name:   "valid dimensions",
			width:  80,
			height: 24,
		},
		{
			name:   "small dimensions",
			width:  10,
			height: 10,
		},
		{
			name:   "large dimensions",
			width:  200,
			height: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terminal := &Terminal{
				width:  tt.width,
				height: tt.height,
			}
			generator := &mockGenerator{}

			r := NewRenderer(terminal, generator)

			if r == nil {
				t.Fatal("NewRenderer() returned nil")
			}
			if len(r.buffer) != tt.height {
				t.Errorf("buffer height = %d, want %d", len(r.buffer), tt.height)
			}
			if len(r.buffer) > 0 && len(r.buffer[0]) != tt.width {
				t.Errorf("buffer width = %d, want %d", len(r.buffer[0]), tt.width)
			}
		})
	}
}

func TestRendererRun(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		wantErr bool
	}{
		{
			name:    "context cancellation",
			timeout: 50 * time.Millisecond,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terminal := &Terminal{
				width:  10,
				height: 10,
			}
			generator := &mockGenerator{}
			r := NewRenderer(terminal, generator)

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			err := r.Run(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPickColor(t *testing.T) {
	tests := []struct {
		name      string
		neighbors int
		want      string
	}{
		{
			name:      "high density",
			neighbors: 6,
			want:      "97",
		},
		{
			name:      "medium high density",
			neighbors: 4,
			want:      "96",
		},
		{
			name:      "medium density",
			neighbors: 2,
			want:      "36",
		},
		{
			name:      "low density",
			neighbors: 0,
			want:      "37",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pickColor(tt.neighbors)
			if got != tt.want {
				t.Errorf("pickColor(%d) = %s, want %s", tt.neighbors, got, tt.want)
			}
		})
	}
}

func TestItoa(t *testing.T) {
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
		{
			name: "large number",
			n:    9999,
			want: "9999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := itoa(tt.n)
			if got != tt.want {
				t.Errorf("itoa(%d) = %s, want %s", tt.n, got, tt.want)
			}
		})
	}
}

// mockGenerator is a test double for the Generator interface.
type mockGenerator struct {
	stepCalled   bool
	renderCalled bool
}

func (m *mockGenerator) Step() bool {
	m.stepCalled = true
	return true
}

func (m *mockGenerator) Render(buffer [][]rune) {
	m.renderCalled = true
	// Fill with some content for testing
	for y := range buffer {
		for x := range buffer[y] {
			if x == y {
				buffer[y][x] = '*'
			}
		}
	}
}

func (m *mockGenerator) ClusterSize() int {
	return 0
}

func (m *mockGenerator) PostProcess() {}
