// Package cmd provides the command execution logic for asciibloom.
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"asciibloom/internal/animation"
	"asciibloom/internal/generators"
)

// Execute runs the main animation loop with proper signal handling and cleanup.
// It returns an error if terminal initialization or animation fails.
func Execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	term, err := animation.NewTerminal()
	if err != nil {
		return fmt.Errorf("initialize terminal: %w", err)
	}
	defer term.Restore()

	generator := generators.NewBrownianTree(term.Width(), term.Height())
	renderer := animation.NewRenderer(term, generator)

	done := make(chan error, 1)
	go func() {
		done <- renderer.Run(ctx)
	}()

	select {
	case <-sigChan:
		cancel()
		term.Restore()
		return nil
	case <-term.InterruptChan():
		cancel()
		term.Restore()
		return nil
	case err := <-done:
		return err
	}
}
