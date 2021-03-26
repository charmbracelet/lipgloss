package lipgloss

import (
	"github.com/muesli/termenv"
)

var (
	ColorProfile      termenv.Profile            = termenv.ColorProfile()
	color             func(string) termenv.Color = ColorProfile.Color
	hasDarkBackground bool                       = termenv.HasDarkBackground()
)

// ColorType is an interface used in color specifications.
type ColorType interface {
	value() string
	color() termenv.Color
}

// NoColor is used to specify the absence of color styling. When this is active
// foreground colors will be rendered with the terminal's default text color,
// and background colors will not be drawn at all.
//
// Example usage:
//
//     var style = someStyle.Copy().Background(lipgloss.NoColor{})
//
type NoColor struct{}

func (n NoColor) value() string {
	return ""
}

func (n NoColor) color() termenv.Color {
	return color("")
}

var noColor = NoColor{}

// Color specifies a color by hex or ANSI value. For example:
//
//     ansiColor := lipgloss.Color("21")
//     hexColor := lipgloss.Color("#0000ff")
//
type Color string

func (c Color) value() string {
	return string(c)
}

func (c Color) color() termenv.Color {
	return color(string(c))
}

// AdaptiveColor provides color options for light and dark backgrounds. The
// appropriate color with be returned based on the darkness of the terminal
// background color determined at runtime.
//
// Example usage:
//
//     color := lipgloss.AdaptiveColor{Light: "#0000ff", Dark: "#000099"}
//
type AdaptiveColor struct {
	Light string
	Dark  string
}

func (a AdaptiveColor) value() string {
	if hasDarkBackground {
		return a.Dark
	}
	return a.Light
}

func (a AdaptiveColor) color() termenv.Color {
	return color(a.value())
}
