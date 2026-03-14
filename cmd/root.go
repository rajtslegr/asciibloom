// Package cmd provides the command execution logic for asciibloom.
package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"asciibloom/internal/animation"
	"asciibloom/internal/core"
	"asciibloom/internal/generators"

	"github.com/spf13/cobra"
)

var mode string

var rootCmd = &cobra.Command{
	Use:   "asciibloom",
	Short: "Terminal ASCII art generator with organic growth patterns",
	Long: `asciibloom generates organic ASCII art animations directly in your terminal.

It supports multiple visualization modes:
  - brownian: Diffusion-limited aggregation creating tree-like structures
  - flow: Flow field simulation with particles following vector fields
  - mandelbrot: Mandelbrot set fractal visualization`,
	RunE: runAnimation,
}

func init() {
	rootCmd.Flags().StringVarP(&mode, "mode", "m", "", "Animation mode: brownian, flow, or mandelbrot (random if not specified)")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func runAnimation(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	term, err := animation.NewTerminal()
	if err != nil {
		return fmt.Errorf("initialize terminal: %w", err)
	}
	defer term.Restore()

	genType := parseMode(mode)
	generator := createGenerator(genType, term.Width(), term.Height())
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

func parseMode(m string) core.GeneratorType {
	switch m {
	case "flow", "flowfield":
		return core.TypeFlowField
	case "brownian", "brown":
		return core.TypeBrownian
	case "mandelbrot", "mandel":
		return core.TypeMandelbrot
	default:
		// Random mode when not specified
		r := rand.Intn(3)
		switch r {
		case 0:
			return core.TypeBrownian
		case 1:
			return core.TypeFlowField
		}
		return core.TypeMandelbrot
	}
}

func createGenerator(genType core.GeneratorType, width, height int) interface {
	Step() bool
	Render(buffer [][]rune)
	ClusterSize() int
	PostProcess()
} {
	switch genType {
	case core.TypeFlowField:
		return generators.NewFlowField(width, height)
	case core.TypeBrownian:
		return generators.NewBrownianTree(width, height)
	case core.TypeMandelbrot:
		return generators.NewMandelbrot(width, height)
	default:
		return generators.NewBrownianTree(width, height)
	}
}
