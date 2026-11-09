package compat

import (
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
)

var (
	// HasDarkBackground is true if the terminal has a dark background.
	//
	// When the NO_COLOR environment variable is set (per
	// https://no-color.org), this remains false and the background
	// query is skipped to avoid the init-time stdin/stdout side
	// effects (raw-mode toggle, OSC 11 + DA1 query writes, blocking
	// stdin read) that would otherwise leak ANSI bytes into the
	// terminal and steal pre-buffered stdin input under a PTY.
	HasDarkBackground bool

	// Profile is the color profile of the terminal.
	Profile colorprofile.Profile
)

func init() {
	// colorprofile.Detect reads only the environment and the stdout
	// file descriptor; it has no stdin side effects, so it is always
	// safe to run at init time. It honors NO_COLOR internally.
	Profile = colorprofile.Detect(os.Stdout, os.Environ())

	// Skip the background-color query when NO_COLOR is set. Per
	// https://no-color.org, any non-empty value disables color
	// behavior. HasDarkBackground therefore defaults to false, which
	// callers interpret as a light background — the safer default
	// when querying is disallowed.
	if os.Getenv("NO_COLOR") != "" {
		return
	}
	HasDarkBackground = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
}

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
