package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestMainFunction(t *testing.T) {
	// main() calls cmd.Execute() which requires a terminal,
	// so we can only verify it compiles and imports work correctly.
	// Integration tests would need a pseudo-terminal setup.
}
