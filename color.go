package lipgloss

import (
	"math"
	"strconv"

	ansi "github.com/charmbracelet/x/exp/term/ansi/style"
	"github.com/lucasb-eyer/go-colorful"
)

// TerminalColor is a color intended to be rendered in the terminal.
type TerminalColor interface {
	Color(*Renderer) ansi.Color
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

var _ TerminalColor = NoColor{}

// Color implements the TerminalColor interface.
func (NoColor) Color(*Renderer) ansi.Color {
	return nil
}

// Color specifies a color by hex or ANSI value. For example:
//
//	ansiColor := lipgloss.Color("21")
//	hexColor := lipgloss.Color("#0000ff")
type Color string

var _ TerminalColor = Color("")

// Color implements the TerminalColor interface.
func (c Color) Color(r *Renderer) ansi.Color {
	return r.ColorProfile().Color(string(c))
}

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

var _ TerminalColor = ANSIColor(0)

// Color implements the TerminalColor interface.
func (ac ANSIColor) Color(r *Renderer) ansi.Color {
	return Color(strconv.FormatUint(uint64(ac), 10)).Color(r)
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

var _ TerminalColor = AdaptiveColor{}

// Color implements the TerminalColor interface.
func (ac AdaptiveColor) Color(r *Renderer) ansi.Color {
	if r.HasDarkBackground() {
		return r.ColorProfile().Color(ac.Dark)
	}
	return r.ColorProfile().Color(ac.Light)
}

// CompleteColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed.
type CompleteColor struct {
	TrueColor string
	ANSI256   string
	ANSI      string
}

var _ TerminalColor = CompleteColor{}

// Color implements the TerminalColor interface.
func (c CompleteColor) Color(r *Renderer) ansi.Color {
	p := r.ColorProfile()
	switch p { //nolint:exhaustive
	case TrueColor:
		return p.Color(c.TrueColor)
	case ANSI256:
		return p.Color(c.ANSI256)
	case ANSI:
		return p.Color(c.ANSI)
	default:
		return noColor.Color(r)
	}
}

// CompleteAdaptiveColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles, with separate options for light and dark backgrounds. Automatic
// color degradation will not be performed.
type CompleteAdaptiveColor struct {
	Light CompleteColor
	Dark  CompleteColor
}

var _ TerminalColor = CompleteAdaptiveColor{}

// Color implements the TerminalColor interface.
func (cac CompleteAdaptiveColor) Color(r *Renderer) ansi.Color {
	if r.HasDarkBackground() {
		return cac.Dark.Color(r)
	}
	return cac.Light.Color(r)
}

// rgbToHex converts red, green, and blue values to a hexadecimal value.
//
//	hex := RgbToHex(0, 0, 255) // 0x0000FF
func rgbToHex(r, g, b uint32) uint32 {
	return r<<16 + g<<8 + b
}

func ansi256ToANSIColor(c ansi.ExtendedColor) ansi.BasicColor {
	var r int
	md := math.MaxFloat64

	h, _ := colorful.Hex(ansiHex[c])
	for i := 0; i <= 15; i++ {
		hb, _ := colorful.Hex(ansiHex[i])
		d := h.DistanceHSLuv(hb)

		if d < md {
			md = d
			r = i
		}
	}

	return ansi.BasicColor(r)
}

func hexToANSI256Color(c colorful.Color) ansi.ExtendedColor {
	v2ci := func(v float64) int {
		if v < 48 {
			return 0
		}
		if v < 115 {
			return 1
		}
		return int((v - 35) / 40)
	}

	// Calculate the nearest 0-based color index at 16..231
	r := v2ci(c.R * 255.0) // 0..5 each
	g := v2ci(c.G * 255.0)
	b := v2ci(c.B * 255.0)
	ci := 36*r + 6*g + b /* 0..215 */

	// Calculate the represented colors back from the index
	i2cv := [6]int{0, 0x5f, 0x87, 0xaf, 0xd7, 0xff}
	cr := i2cv[r] // r/g/b, 0..255 each
	cg := i2cv[g]
	cb := i2cv[b]

	// Calculate the nearest 0-based gray index at 232..255
	var grayIdx int
	average := (r + g + b) / 3
	if average > 238 {
		grayIdx = 23
	} else {
		grayIdx = (average - 3) / 10 // 0..23
	}
	gv := 8 + 10*grayIdx // same value for r/g/b, 0..255

	// Return the one which is nearer to the original input rgb value
	c2 := colorful.Color{R: float64(cr) / 255.0, G: float64(cg) / 255.0, B: float64(cb) / 255.0}
	g2 := colorful.Color{R: float64(gv) / 255.0, G: float64(gv) / 255.0, B: float64(gv) / 255.0}
	colorDist := c.DistanceHSLuv(c2)
	grayDist := c.DistanceHSLuv(g2)

	if colorDist <= grayDist {
		return ansi.ExtendedColor(16 + ci)
	}
	return ansi.ExtendedColor(232 + grayIdx)
}
