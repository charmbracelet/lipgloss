package lipgloss

import (
	"github.com/muesli/termenv"
)

var (
	ColorProfile      termenv.Profile            = termenv.ColorProfile()
	color             func(string) termenv.Color = ColorProfile.Color
	hasDarkBackground bool                       = termenv.HasDarkBackground()
)

// TerminalColor is a color intended to be rendered in the terminal. It
// satisfies the Go color.Color interface.
type TerminalColor interface {
	value() string
	color() termenv.Color
	RGBA() (r, g, b, a uint32)
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

// RGBA returns the RGBA value of this color. Because we have to return
// something, despite this color being the absence of color, we're returning
// defaults for all values.
func (n NoColor) RGBA() (r, g, b, a uint32) {
	return 0, 0, 0, 0
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

// RGBA returns the RGBA value of this color.
func (c Color) RGBA() (r, g, b, a uint32) {
	return c.RGBA()
}

// AdaptiveColor provides color options for light and dark backgrounds. The
// appropriate color with be returned at runtime based on the darkness of the
// terminal background color.
//
// Example usage:
//
//     color := lipgloss.AdaptiveColor{Light: "#0000ff", Dark: "#000099"}
//
type AdaptiveColor struct {
	Light string
	Dark  string
}

func (ac AdaptiveColor) value() string {
	if hasDarkBackground {
		return ac.Dark
	}
	return ac.Light
}

func (ac AdaptiveColor) color() termenv.Color {
	return color(ac.value())
}

// RGBA returns the RGBA value of this color.
func (ac AdaptiveColor) RGBA() (r, g, b, a uint32) {
	return ac.RGBA()
}
