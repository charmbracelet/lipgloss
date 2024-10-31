package lipgloss

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/charmbracelet/colorprofile"
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
	return 0x0, 0x0, 0x0, 0xFFFF //nolint:mnd
}

// Color specifies a color by hex or ANSI256 value. For example:
//
//	ansiColor := lipgloss.Color(21)
//	hexColor := lipgloss.Color("#0000ff")
//	uint32Color := lipgloss.Color(0xff0000)
func Color(c any) color.Color {
	switch c := c.(type) {
	case nil:
		return noColor
	case ansi.BasicColor:
		return c
	case ansi.ExtendedColor:
		return c
	case ansi.TrueColor:
		return c
	case string:
		if len(c) == 0 {
			return noColor
		}
		if h, err := colorful.Hex(c); err == nil {
			return h
		} else if i, err := strconv.Atoi(c); err == nil {
			if i < 16 { //nolint:mnd
				return ansi.BasicColor(i) //nolint:gosec
			} else if i < 256 { //nolint:mnd
				return ansi.ExtendedColor(i) //nolint:gosec
			}
			return ansi.TrueColor(i) //nolint:gosec
		}
		return noColor
	case int:
		if c < 16 { //nolint:mnd
			return ansi.BasicColor(c) //nolint:gosec
		} else if c < 256 { //nolint:mnd
			return ansi.ExtendedColor(c) //nolint:gosec
		}
		return ansi.TrueColor(c) //nolint:gosec
	case color.Color:
		return c
	}
	return Color(fmt.Sprint(c))
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
	const shift = 8
	r |= uint32(c.R) << shift
	g |= uint32(c.G) << shift
	b |= uint32(c.B) << shift
	a = 0xFFFF
	return
}

// ANSIColor is a color specified by an ANSI256 color value.
//
// Example usage:
//
//	colorA := lipgloss.ANSIColor(8)
//	colorB := lipgloss.ANSIColor(134)
type ANSIColor = ansi.ExtendedColor

// LightDarkFunc is a function that returns a color based on whether the
// terminal has a light or dark background. You can create one of these with
// [LightDark].
//
// Example:
//
//	lightDark := lipgloss.LightDark(hasDarkBackground)
//	myHotColor := lightDark("#ff0000", "#0000ff")
//
// For more info see [LightDark].
type LightDarkFunc func(light, dark any) color.Color

// LightDark is a simple helper type that can be used to choose the appropriate
// color based on whether the terminal has a light or dark background.
//
//	lightDark := lipgloss.LightDark(hasDarkBackground)
//	theRightColor := lightDark("#0000ff", "#ff0000")
//
// In practice, there are slightly different workflows between Bubble Tea and
// Lip Gloss standalone.
//
// In Bubble Tea listen for tea.BackgroundColorMsg, which automatically
// flows through Update on start, and whenever the background color changes:
//
//	case tea.BackgroundColorMsg:
//	    m.hasDarkBackground = msg.IsDark()
//
// Later, when you're rendering:
//
//	lightDark := lipgloss.LightDark(m.hasDarkBackground)
//	myHotColor := lightDark("#ff0000", "#0000ff")
//
// In standalone Lip Gloss, the workflow is simpler:
//
//	hasDarkBG, _ := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
//	lightDark := lipgloss.LightDark(hasDarkBG)
//	myHotColor := lightDark("#ff0000", "#0000ff")
func LightDark(isDark bool) LightDarkFunc {
	return func(light, dark any) color.Color {
		if isDark {
			return Color(dark)
		}
		return Color(light)
	}
}

// isDarkColor returns whether the given color is dark (based on the luminance
// portion of the color as interpreted as HSL).
//
// Example usage:
//
//	color := lipgloss.Color("#0000ff")
//	if lipgloss.isDarkColor(color) {
//		fmt.Println("It's dark! I love darkness!")
//	} else {
//		fmt.Println("It's light! Cover your eyes!")
//	}
func isDarkColor(c color.Color) bool {
	col, ok := colorful.MakeColor(c)
	if !ok {
		return true
	}

	_, _, l := col.Hsl()
	return l < 0.5 //nolint:mnd
}

// CompleteFunc is a function that returns the appropriate color based on the
// given color profile.
//
// Example usage:
//
//	p := colorprofile.Detect(os.Stderr, os.Environ())
//	complete := lipgloss.Complete(p)
//	color := complete(
//		lipgloss.Color(1), // ANSI
//		lipgloss.Color(124), // ANSI256
//		lipgloss.Color("#ff34ac"), // TrueColor
//	)
//	fmt.Println("Ooh, pretty color: ", color)
//
// For more info see [Complete].
type CompleteFunc func(ansi, ansi256, truecolor color.Color) color.Color

// Complete returns a function that will return the appropriate color based on
// the given color profile.
//
// Example usage:
//
//	p := colorprofile.Detect(os.Stderr, os.Environ())
//	complete := lipgloss.Complete(p)
//	color := complete(
//	    lipgloss.Color(1), // ANSI
//	    lipgloss.Color(124), // ANSI256
//	    lipgloss.Color("#ff34ac"), // TrueColor
//	)
//	fmt.Println("Ooh, pretty color: ", color)
func Complete(p colorprofile.Profile) CompleteFunc {
	return func(ansi, ansi256, truecolor color.Color) color.Color {
		switch p {
		case colorprofile.ANSI:
			return ansi
		case colorprofile.ANSI256:
			return ansi256
		case colorprofile.TrueColor:
			return truecolor
		}
		return noColor
	}
}
