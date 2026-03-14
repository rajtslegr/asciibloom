package core

import (
	"testing"
)

func TestNewGrid(t *testing.T) {
	g := NewGrid(20, 10)

	if g.Width() != 20 {
		t.Errorf("Width() = %d, want 20", g.Width())
	}
	if g.Height() != 10 {
		t.Errorf("Height() = %d, want 10", g.Height())
	}
}

func TestGridGetSet(t *testing.T) {
	g := NewGrid(20, 10)

	// Test Get/Set
	g.Set(5, 5, 3)
	if g.Get(5, 5) != 3 {
		t.Errorf("Get(5, 5) = %d, want 3", g.Get(5, 5))
	}

	// Test Set only increases value
	g.Set(5, 5, 2)
	if g.Get(5, 5) != 3 {
		t.Errorf("Get(5, 5) after lower set = %d, want 3", g.Get(5, 5))
	}

	g.Set(5, 5, 5)
	if g.Get(5, 5) != 5 {
		t.Errorf("Get(5, 5) after higher set = %d, want 5", g.Get(5, 5))
	}
}

func TestGridOutOfBounds(t *testing.T) {
	g := NewGrid(20, 10)

	if g.Get(-1, 5) != 0 {
		t.Error("Get(-1, 5) should return 0")
	}
	if g.Get(25, 5) != 0 {
		t.Error("Get(25, 5) should return 0")
	}
	if g.Get(5, -1) != 0 {
		t.Error("Get(5, -1) should return 0")
	}
	if g.Get(5, 15) != 0 {
		t.Error("Get(5, 15) should return 0")
	}

	// Set should not panic on out of bounds
	g.Set(-1, 5, 5)
	g.Set(25, 5, 5)
}

func TestGridCountNeighbors(t *testing.T) {
	g := NewGrid(20, 10)

	g.Set(5, 5, 5)
	g.Set(4, 5, 1)
	g.Set(6, 5, 1)
	g.Set(5, 4, 1)

	count := g.CountNeighbors(5, 5)
	if count != 3 {
		t.Errorf("CountNeighbors(5, 5) = %d, want 3", count)
	}
}

func TestGridClear(t *testing.T) {
	g := NewGrid(20, 10)

	g.Set(5, 5, 5)
	g.Set(10, 10, 3)

	g.Clear()

	if g.Get(5, 5) != 0 {
		t.Error("Clear() did not reset cell at (5, 5)")
	}
	if g.Get(10, 10) != 0 {
		t.Error("Clear() did not reset cell at (10, 10)")
	}
}

func TestCharForIntensity(t *testing.T) {
	tests := []struct {
		neighbors int
		intensity int
		want      rune
	}{
		{6, 5, 'O'},
		{4, 5, 'o'},
		{3, 5, '*'},
		{2, 5, ':'},
		{1, 5, '+'},
		{0, 5, '+'},
		{0, 4, '.'},
		{0, 1, '.'},
	}

	for _, tt := range tests {
		got := CharForIntensity(tt.neighbors, tt.intensity)
		if got != tt.want {
			t.Errorf("CharForIntensity(%d, %d) = %c, want %c", tt.neighbors, tt.intensity, got, tt.want)
		}
	}
}
