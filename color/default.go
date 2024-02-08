package color

import "github.com/charmbracelet/lipgloss"

var osProfile = NewProfile(nil)

// Color returns a TerminalColor based on the color profile.
func Color(c string) lipgloss.TerminalColor {
	return osProfile.Color(c)
}
