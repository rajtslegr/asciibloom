# asciibloom

Terminal ASCII art with organic patterns. Made with Go.

## Build

```bash
make build
# or
go build .
```

## Usage

```bash
./asciibloom                # Random mode
./asciibloom -m brownian    # Brownian tree
./asciibloom -m flowfield   # Flow field
./asciibloom -m mandelbrot  # Mandelbrot set
./asciibloom -m reaction    # Gray-Scott reaction-diffusion
```

Press `Ctrl+C` to exit.

## Modes

- **brownian** - Diffusion-limited aggregation creating tree-like structures
- **flowfield** - Flow field simulation with particles following vector fields  
- **mandelbrot** - Mandelbrot set fractal with animated boundary edges
- **reaction** - Gray-Scott reaction-diffusion creating organic patterns (coral, spots, maze, waves)
