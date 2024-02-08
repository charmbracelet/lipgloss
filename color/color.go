package color

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// NoColor is used to specify the absence of color styling. When this is active
// foreground colors will be rendered with the terminal's default text color,
// and background colors will not be drawn at all.
//
// Example usage:
//
//	var style = someStyle.Copy().Background(lipgloss.NoColor{})
type NoColor = termenv.NoColor

var _ lipgloss.TerminalColor = NoColor{}

// Color specifies a color by hex or ANSI value. For example:
//
//	ansiColor := lipgloss.Color("21")
//	hexColor := lipgloss.Color("#0000ff")
// type Color string
//
// var _ lipgloss.TerminalColor = Color("")
//
// func (c Color) color(r *Renderer) termenv.Color {
// 	return r.ColorProfile().Color(string(c))
// }

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface. Note that on error we return black with 100% opacity, or:
//
// Red: 0x0, Green: 0x0, Blue: 0x0, Alpha: 0xFFFF.
//
// Deprecated.
// func (c Color) RGBA() (r, g, b, a uint32) {
// 	return termenv.ConvertToRGB(c.color(renderer)).RGBA()
// }

// ANSIColor is a color specified by an ANSI color value. It's merely syntactic
// sugar for the more general Color function. Invalid colors will render as
// black.
//
// Example usage:
//
//	// These two statements are equivalent.
//	colorA := lipgloss.ANSIColor(21)
//	colorB := lipgloss.Color("21")
type ANSIColor uint

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface. Note that on error we return black with 100% opacity, or:
//
// Red: 0x0, Green: 0x0, Blue: 0x0, Alpha: 0xFFFF.
func (ac ANSIColor) RGBA() (r, g, b, a uint32) {
	cf := termenv.ANSI256Color(int(ac))
	return termenv.ConvertToRGB(cf).RGBA()
}

// AdaptiveColor provides color options for light and dark backgrounds. The
// appropriate color will be returned at runtime based on the darkness of the
// terminal background color.
//
// Example usage:
//
//	color := lipgloss.AdaptiveColor{Light: "#0000ff", Dark: "#000099"}
type AdaptiveColor struct {
	Light lipgloss.TerminalColor
	Dark  lipgloss.TerminalColor
}

// Color returns the color that should be used.
func (ac AdaptiveColor) Color(isDark bool) lipgloss.TerminalColor {
	if isDark {
		return ac.Dark
	}
	return ac.Light
}

// CompleteColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed.
type CompleteColor struct {
	TrueColor string
	ANSI256   string
	ANSI      string
}

func (c CompleteColor) Color(p Profile) lipgloss.TerminalColor {
	switch termenv.Profile(p) { //nolint:exhaustive
	case termenv.TrueColor:
		return p.Color(c.TrueColor)
	case termenv.ANSI256:
		return p.Color(c.ANSI256)
	case termenv.ANSI:
		return p.Color(c.ANSI)
	default:
		return termenv.NoColor{}
	}
}

// CompleteAdaptiveColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles, with separate options for light and dark backgrounds. Automatic
// color degradation will not be performed.
type CompleteAdaptiveColor struct {
	Light CompleteColor
	Dark  CompleteColor
}

func (cac CompleteAdaptiveColor) Color(p Profile, isDark bool) lipgloss.TerminalColor {
	if isDark {
		return cac.Dark.Color(p)
	}
	return cac.Light.Color(p)
}
