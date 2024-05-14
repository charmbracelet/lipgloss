package lipgloss

import (
	"image/color"
	"strconv"

	"github.com/charmbracelet/x/ansi"
	"github.com/lucasb-eyer/go-colorful"
)

// 4-bit color constants.
const (
	Black basicColor = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White

	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

// TerminalColor is a color intended to be rendered in the terminal.
type TerminalColor interface {
	color(p Profile, hasLightBg bool) ansi.Color
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

func (NoColor) color(Profile, bool) ansi.Color {
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

// Color returns a color.Color from a hex or ANSI value (0-16, 16-255).
func (c Color) Color() (col ansi.Color) {
	s := string(c)
	if len(s) == 0 {
		return noColor
	}

	if h, err := colorful.Hex(s); err == nil {
		tc := uint32(h.R*255)<<16 + uint32(h.G*255)<<8 + uint32(h.B*255)
		col = ansi.TrueColor(tc)
	} else if i, err := strconv.Atoi(s); err == nil {
		if i < 16 {
			col = ansi.BasicColor(i)
		} else if i < 256 {
			col = ansi.ExtendedColor(i)
		} else {
			col = ansi.TrueColor(i)
		}
	}

	return noColor
}

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface. Note that on error we return black with 100% opacity, or:
//
// Red: 0x0, Green: 0x0, Blue: 0x0, Alpha: 0xFFFF.
func (c Color) RGBA() (r, g, b, a uint32) {
	return c.Color().RGBA()
}

func (c Color) color(p Profile, _ bool) ansi.Color {
	return p.Convert(c)
}

// RGBColor is a color specified by red, green, and blue values.
type RGBColor struct {
	R uint8
	G uint8
	B uint8
}

// Color returns a color.Color from an RGB value.
func (c RGBColor) Color() ansi.Color {
	return ansi.TrueColor(uint32(c.R)<<16 + uint32(c.G)<<8 + uint32(c.B))
}

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface. Note that on error we return black with 100% opacity, or:
//
// Red: 0x0, Green: 0x0, Blue: 0x0, Alpha: 0xFFFF.
func (c RGBColor) RGBA() (r, g, b, a uint32) {
	return c.Color().RGBA()
}

func (c RGBColor) color(p Profile, _ bool) ansi.Color {
	return p.Convert(c)
}

// ANSIColor is a color specified by an ANSI256 color value.
//
// Example usage:
//
//	colorA := lipgloss.ANSIColor(8)
//	colorB := lipgloss.ANSIColor(134)
type ANSIColor uint8

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface. Note that on error we return black with 100% opacity, or:
//
// Red: 0x0, Green: 0x0, Blue: 0x0, Alpha: 0xFFFF.
func (c ANSIColor) RGBA() (r, g, b, a uint32) {
	return ansi.ExtendedColor(c).RGBA()
}

func (ac ANSIColor) color(p Profile, _ bool) ansi.Color {
	return p.Convert(ansi.ExtendedColor(ac))
}

// basicColor is a color specified by an ANSI 4-bit color value.
type basicColor uint8

// Color returns the color.Color from a basic color value (0-16).
func (bc basicColor) Color() ansi.Color {
	return ansi.BasicColor(bc)
}

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface. Note that on error we return black with 100% opacity, or:
//
// Red: 0x0, Green: 0x0, Blue: 0x0, Alpha: 0xFFFF.
func (c basicColor) RGBA() (r, g, b, a uint32) {
	return c.Color().RGBA()
}

func (bc basicColor) color(p Profile, _ bool) ansi.Color {
	return p.Convert(ansi.BasicColor(bc))
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

func (ac AdaptiveColor) color(p Profile, hasLightBg bool) ansi.Color {
	if hasLightBg {
		return Color(ac.Light).color(p, hasLightBg)
	}
	return Color(ac.Dark).color(p, hasLightBg)
}

// CompleteColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed.
type CompleteColor struct {
	TrueColor string
	ANSI256   string
	ANSI      string
}

func (c CompleteColor) color(p Profile, hasLightBg bool) ansi.Color {
	switch p { //nolint:exhaustive
	case TrueColor:
		return Color(c.TrueColor).color(p, hasLightBg)
	case ANSI256:
		return Color(c.ANSI256).color(p, hasLightBg)
	case ANSI:
		return Color(c.ANSI).color(p, hasLightBg)
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

func (cac CompleteAdaptiveColor) color(p Profile, hasLightBg bool) ansi.Color {
	if hasLightBg {
		return cac.Light.color(p, hasLightBg)
	}
	return cac.Light.color(p, hasLightBg)
}

// ConvertToRGB converts a Color to a colorful.Color.
func ConvertToRGB(c ansi.Color) colorful.Color {
	ch, _ := colorful.MakeColor(c)
	return ch
}
