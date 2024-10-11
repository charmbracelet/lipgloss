package adaptive

import (
	"image/color"

	"github.com/charmbracelet/lipgloss"
)

// Color returns the color that should be used based on the terminal's
// background color.
func Color(light, dark any) color.Color {
	return lipgloss.Adapt(HasDarkBackground).Color(light, dark)
}
