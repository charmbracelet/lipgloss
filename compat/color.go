package compat

import (
	"os"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/lipgloss/v2"
)

var (
	// HasDarkBackground is true if the terminal has a dark background.
	HasDarkBackground = func() bool {
		hdb, _ := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
		return hdb
	}()

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
	Light any
	Dark  any
}

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface.
func (c AdaptiveColor) RGBA() (uint32, uint32, uint32, uint32) {
	if HasDarkBackground {
		return lipgloss.Color(c.Dark).RGBA()
	}
	return lipgloss.Color(c.Light).RGBA()
}

// CompleteColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed.
type CompleteColor struct {
	TrueColor any
	ANSI256   any
	ANSI      any
}

// RGBA returns the RGBA value of this color. This satisfies the Go Color
// interface.
func (c CompleteColor) RGBA() (uint32, uint32, uint32, uint32) {
	switch Profile {
	case colorprofile.TrueColor:
		return lipgloss.Color(c.TrueColor).RGBA()
	case colorprofile.ANSI256:
		return lipgloss.Color(c.ANSI256).RGBA()
	case colorprofile.ANSI:
		return lipgloss.Color(c.ANSI).RGBA()
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
