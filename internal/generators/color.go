package generators

// ColorGrid stores character and intensity data for colorized output.
type ColorGrid struct {
	Width       int
	Height      int
	Chars       [][]rune
	Intensities [][]float64
}

func NewColorGrid(width, height int) *ColorGrid {
	chars := make([][]rune, height)
	intensities := make([][]float64, height)
	for y := 0; y < height; y++ {
		chars[y] = make([]rune, width)
		intensities[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			chars[y][x] = ' '
		}
	}
	return &ColorGrid{
		Width:       width,
		Height:      height,
		Chars:       chars,
		Intensities: intensities,
	}
}

// Color represents an RGB color.
type Color struct {
	R, G, B int
}

// BrownianPalette returns a grayscale palette for Brownian tree visualization.
func BrownianPalette() []Color {
	return []Color{
		{232, 232, 242},
		{220, 220, 235},
		{200, 200, 225},
		{180, 180, 215},
		{160, 160, 205},
		{140, 140, 195},
		{120, 120, 185},
		{100, 100, 175},
		{80, 80, 165},
		{60, 60, 155},
	}
}

func grayCode(brightness int) string {
	if brightness < 0 {
		brightness = 0
	}
	if brightness > 255 {
		brightness = 255
	}
	for i := 232; i <= 255; i++ {
		gray := 8 + 10*(i-232)
		if gray >= brightness {
			return intToStr(i)
		}
	}
	return "255"
}

func intToStr(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return intToStr(n/10) + string(rune('0'+n%10))
}

// ColorizedOutput generates ANSI escape sequences for colored terminal output.
func (cg *ColorGrid) ColorizedOutput(palette []Color) string {
	var result []rune
	result = append(result, []rune("\x1b[H")...)

	for y := 0; y < cg.Height; y++ {
		for x := 0; x < cg.Width; x++ {
			char := cg.Chars[y][x]
			if char == ' ' {
				continue
			}

			result = append(result, []rune("\x1b["+intToStr(y+1)+";"+intToStr(x+1)+"H")...)

			intensity := cg.Intensities[y][x]
			if intensity < 0 {
				intensity = 0
			}
			if intensity > 1 {
				intensity = 1
			}

			gray := int(150 + intensity*100)
			if gray > 255 {
				gray = 255
			}

			result = append(result, []rune("\x1b[38;5;"+grayCode(gray)+"m")...)
			result = append(result, char)
		}
	}

	result = append(result, []rune("\x1b[39m")...)
	return string(result)
}
