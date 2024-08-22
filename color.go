package lipgloss

import (
	"image/color"
	"strconv"

	"github.com/charmbracelet/x/ansi"
	"github.com/lucasb-eyer/go-colorful"
)

// 4-bit color constants.
const (
	Black ansi.BasicColor = iota
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

var noColor = NoColor{}

// NoColor is used to specify the absence of color styling. When this is active
// foreground colors will be rendered with the terminal's default text color,
// and background colors will not be drawn at all.
//
// Example usage:
//
//	var style = someStyle.Background(lipgloss.NoColor{})
type NoColor struct{}

// RGBA returns the RGBA value of this color. Because we have to return
// something, despite this color being the absence of color, we're returning
// black with 100% opacity.
//
// Red: 0x0, Green: 0x0, Blue: 0x0, Alpha: 0xFFFF.
func (n NoColor) RGBA() (r, g, b, a uint32) {
	return 0x0, 0x0, 0x0, 0xFFFF //nolint:gomnd
}

// Color specifies a color by hex or ANSI256 value. For example:
//
//	ansiColor := lipgloss.Color(21)
//	hexColor := lipgloss.Color("#0000ff")
//	uint32Color := lipgloss.Color(0xff0000)
func Color[T string | int](c T) color.Color {
	var col color.Color = noColor
	switch c := any(c).(type) {
	case string:
		if len(c) == 0 {
			return col
		}
		if h, err := colorful.Hex(c); err == nil {
			return h
		} else if i, err := strconv.Atoi(c); err == nil {
			if i < 16 {
				return ansi.BasicColor(i)
			} else if i < 256 {
				return ansi.ExtendedColor(i)
			}
			return ansi.TrueColor(i)
		}
	case int:
		if c < 16 {
			return ansi.BasicColor(c)
		} else if c < 256 {
			return ansi.ExtendedColor(c)
		}
		return ansi.TrueColor(c)
	}
	return col
}

// RGBColor is a color specified by red, green, and blue values.
type RGBColor struct {
	R uint8
	G uint8
	B uint8
}

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface.
func (c RGBColor) RGBA() (r, g, b, a uint32) {
	r |= uint32(c.R) << 8
	g |= uint32(c.G) << 8
	b |= uint32(c.B) << 8
	a = 0xFFFF //nolint:gomnd
	return
}

// ANSIColor is a color specified by an ANSI256 color value.
//
// Example usage:
//
//	colorA := lipgloss.ANSIColor(8)
//	colorB := lipgloss.ANSIColor(134)
type ANSIColor = ansi.ExtendedColor

// IsDarkColor returns whether the given color is dark.
//
// Example usage:
//
//	color := lipgloss.Color("#0000ff")
//	if lipgloss.IsDarkColor(color) {
//		fmt.Println("It's dark!")
//	} else {
//		fmt.Println("It's light!")
//	}
func IsDarkColor(c color.Color) bool {
	col, ok := colorful.MakeColor(c)
	if !ok {
		return true
	}

	_, _, l := col.Hsl()
	return l < 0.5
}
