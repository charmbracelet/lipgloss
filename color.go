package lipgloss

import "github.com/muesli/termenv"

// ColorType is an interface used in color specifications.
type ColorType interface {
	value() string
}

// NoColor is used to specify the absence of color styling. When this is active
// foreground colors will be rendered with the terminal's default text color,
// and background colors will not be drawn at all.
//
// Example usage:
//
//     color := NoColor
//
var NoColor = noColor{}

type noColor struct{}

func (n noColor) value() string {
	return ""
}

// Color specifies a color by hex or ANSI value. For example:
//
//     ansiColor := Color("21")
//     hexColor := Color("#0000ff")
//
type Color string

func (c Color) value() string {
	return string(c)
}

// AdaptiveColor provides color alternatives for light and dark backgrounds.
// The appropriate color with be returned based on the darkness of the terminal
// background color determined at runtime.
//
// Example usage:
//
//     color := AdaptiveColor{Light: "#0000ff", Dark: "#000099"}
//
type AdaptiveColor struct {
	Light string
	Dark  string
}

func (a AdaptiveColor) value() string {
	if !darkBackgroundQueried {
		hasDarkBackground = termenv.HasDarkBackground()
	}
	if hasDarkBackground {
		return a.Dark
	}
	return a.Light
}
