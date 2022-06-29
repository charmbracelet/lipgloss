package lipgloss

import "image/color"

// TerminalColor is a color intended to be rendered in the terminal.
type TerminalColor interface {
	colorType()
}

// NoColor is used to specify the absence of color styling. When this is active
// foreground colors will be rendered with the terminal's default text color,
// and background colors will not be drawn at all.
//
// Example usage:
//
//	var style = someStyle.Copy().Background(lipgloss.NoColor{})
type NoColor struct{}

func (NoColor) colorType() {}

var noColor = NoColor{}

// Color specifies a color by hex or ANSI value. For example:
//
//	ansiColor := lipgloss.Color("21")
//	hexColor := lipgloss.Color("#0000ff")
type Color string

func (Color) colorType() {}

// ANSIColor is a color specified by an ANSI color value. It's merely syntactic
// sugar for the more general Color function. Invalid colors will render as
// black.
//
// Example usage:
//
//    // These two statements are equivalent.
//    colorA := lipgloss.ANSIColor(21)
//    colorB := lipgloss.Color("21")
//
type ANSIColor uint

func (ANSIColor) colorType() {}

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

func (AdaptiveColor) colorType() {}

// CompleteColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles. Automatic color degredation will not be performed.
type CompleteColor struct {
	TrueColor string
	ANSI256   string
	ANSI      string
}

func (CompleteColor) colorType() {}

// CompleteColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles, with separate options for light and dark backgrounds. Automatic
// color degredation will not be performed.
type CompleteAdaptiveColor struct {
	Light CompleteColor
	Dark  CompleteColor
}

func (CompleteAdaptiveColor) colorType() {}

// hexToColor translates a hex color string (#RRGGBB or #RGB) into a color.RGB,
// which satisfies the color.Color interface. If an invalid string is passed
// black with 100% opacity will be returned: or, in hex format, 0x000000FF.
func hexToColor(hex string) (c color.RGBA) {
	c.A = 0xFF

	if hex == "" || hex[0] != '#' {
		return c
	}

	const (
		fullFormat  = 7 // #RRGGBB
		shortFormat = 4 // #RGB
	)

	switch len(hex) {
	case fullFormat:
		const offset = 4
		c.R = hexToByte(hex[1])<<offset + hexToByte(hex[2])
		c.G = hexToByte(hex[3])<<offset + hexToByte(hex[4])
		c.B = hexToByte(hex[5])<<offset + hexToByte(hex[6])
	case shortFormat:
		const offset = 0x11
		c.R = hexToByte(hex[1]) * offset
		c.G = hexToByte(hex[2]) * offset
		c.B = hexToByte(hex[3]) * offset
	}

	return c
}

func hexToByte(b byte) byte {
	const offset = 10
	switch {
	case b >= '0' && b <= '9':
		return b - '0'
	case b >= 'a' && b <= 'f':
		return b - 'a' + offset
	case b >= 'A' && b <= 'F':
		return b - 'A' + offset
	}
	// Invalid, but just return 0.
	return 0
}
