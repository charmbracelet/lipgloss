package adaptive

import (
	"image/color"
)

// Color returns the color that should be used based on the terminal's
// background color.
func Color(light, dark any) color.Color {
	return colorFn.Color(light, dark)
}
