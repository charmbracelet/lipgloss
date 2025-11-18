package compat

import (
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
)

var (
	// HasDarkBackground is true if the terminal has a dark background.
	HasDarkBackground = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)

	// Profile is the color profile of the terminal.
	Profile = colorprofile.Detect(os.Stdout, os.Environ())
)

// AdaptiveColor provides color options for light and dark backgrounds. The
// appropriate color will be returned at runtime based on the darkness of the
// terminal background color.
//
// Example usage:
//
//	color := lipgloss.AdaptiveColor{Light: "#0000ff", Dark: "#000099"}
type AdaptiveColor struct {
	Light color.Color
	Dark  color.Color
}

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface.
func (c AdaptiveColor) RGBA() (uint32, uint32, uint32, uint32) {
	if HasDarkBackground {
		return c.Dark.RGBA()
	}
	return c.Light.RGBA()
}

// CompleteColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed.
type CompleteColor struct {
	TrueColor color.Color
	ANSI256   color.Color
	ANSI      color.Color
}

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface.
func (c CompleteColor) RGBA() (uint32, uint32, uint32, uint32) {
	switch Profile { //nolint:exhaustive
	case colorprofile.TrueColor:
		return c.TrueColor.RGBA()
	case colorprofile.ANSI256:
		return c.ANSI256.RGBA()
	case colorprofile.ANSI:
		return c.ANSI.RGBA()
	}
	return lipgloss.NoColor{}.RGBA()
}

// CompleteAdaptiveColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles, with separate options for light and dark backgrounds. Automatic
// color degradation will not be performed.
type CompleteAdaptiveColor struct {
	Light CompleteColor
	Dark  CompleteColor
}

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface.
func (c CompleteAdaptiveColor) RGBA() (uint32, uint32, uint32, uint32) {
	if HasDarkBackground {
		return c.Dark.RGBA()
	}
	return c.Light.RGBA()
}
