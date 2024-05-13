package lipgloss

import (
	"image/color"

	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/lucasb-eyer/go-colorful"
)

// TerminalColor is a color intended to be rendered in the terminal.
type TerminalColor interface {
	color(*Renderer) ansi.Color
}

var noColor = NoColor{}

// NoColor is used to specify the absence of color styling. When this is active
// foreground colors will be rendered with the terminal's default text color,
// and background colors will not be drawn at all.
//
// Example usage:
//
//	var style = someStyle.Copy().Background(lipgloss.NoColor{})
type NoColor struct{}

func (NoColor) color(*Renderer) ansi.Color {
	return color.Black
}

// RGBA returns the RGBA value of this color. Because we have to return
// something, despite this color being the absence of color, we're returning
// black with 100% opacity.
//
// Red: 0x0, Green: 0x0, Blue: 0x0, Alpha: 0xFFFF.
func (n NoColor) RGBA() (r, g, b, a uint32) {
	return 0x0, 0x0, 0x0, 0xFFFF //nolint:gomnd
}

// Color specifies a color by hex or ANSI value. For example:
//
//	ansiColor := lipgloss.Color("21")
//	hexColor := lipgloss.Color("#0000ff")
type Color string

func (c Color) color(r *Renderer) ansi.Color {
	return r.ColorProfile().Color(string(c))
}

// ANSIColor is a color specified by an ANSI256 color value.
//
// Example usage:
//
//	colorA := lipgloss.ANSIColor(8)
//	colorB := lipgloss.ANSIColor(134)
type ANSIColor uint8

func (ac ANSIColor) color(r *Renderer) ansi.Color {
	return r.ColorProfile().Convert(ansi.ExtendedColor(ac))
}

// AdaptiveColor provides color options for light and dark backgrounds. The
// appropriate color will be returned at runtime based on the darkness of the
// terminal background color.
//
// Example usage:
//
//	color := lipgloss.AdaptiveColor{Light: "#0000ff", Dark: "#000099"}
type AdaptiveColor struct {
	Light string
	Dark  string
}

func (ac AdaptiveColor) color(r *Renderer) ansi.Color {
	if r.HasDarkBackground() {
		return Color(ac.Dark).color(r)
	}
	return Color(ac.Light).color(r)
}

// CompleteColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed.
type CompleteColor struct {
	TrueColor string
	ANSI256   string
	ANSI      string
}

func (c CompleteColor) color(r *Renderer) ansi.Color {
	p := r.ColorProfile()
	switch p { //nolint:exhaustive
	case TrueColor:
		return p.Color(c.TrueColor)
	case ANSI256:
		return p.Color(c.ANSI256)
	case ANSI:
		return p.Color(c.ANSI)
	default:
		return noColor
	}
}

// CompleteAdaptiveColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles, with separate options for light and dark backgrounds. Automatic
// color degradation will not be performed.
type CompleteAdaptiveColor struct {
	Light CompleteColor
	Dark  CompleteColor
}

func (cac CompleteAdaptiveColor) color(r *Renderer) ansi.Color {
	if r.HasDarkBackground() {
		return cac.Dark.color(r)
	}
	return cac.Light.color(r)
}

// ConvertToRGB converts a Color to a colorful.Color.
func ConvertToRGB(c ansi.Color) colorful.Color {
	ch, _ := colorful.MakeColor(c)
	return ch
}
