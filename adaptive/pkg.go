package adaptive

import (
	"image/color"

	"github.com/charmbracelet/lipgloss"
)

// Print writes the given text to the default Lip Gloss writer.
func Print(text string) (n int, err error) {
	return Default().Print(text)
}

// Println writes the given text to the default Lip Gloss writer followed by a
// newline.
func Println(text string) (n int, err error) {
	return Default().Println(text)
}

// Printf writes the given text to the default Lip Gloss writer with the given
// format.
func Printf(format string, a ...interface{}) (n int, err error) {
	return Default().Printf(format, a...)
}

// AdaptiveColor returns a color appropriate for the terminal background color.
//
// This function is a convenience function that uses the default adaptive color
// output which is based on the standard input and output and environment
// variables.
//
// Example usage:
//
//	light := lipgloss.Color("#0000ff")
//	dark := lipgloss.Color("#000099")
//	color := adaptive.AdaptiveColor(light, dark)
func AdaptiveColor(light, dark lipgloss.TerminalColor) lipgloss.TerminalColor {
	return Default().AdaptiveColor(light, dark)
}

// CompleteColor specifies exact values for TrueColor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed, instead the
// appropriate color will be returned at based on the terminal color profile.
//
// This function is a convenience function that uses the default adaptive color
// output which is based on the standard input and output and environment
// variables.
func CompleteColor(trueColor, ansi256, ansi lipgloss.TerminalColor) lipgloss.TerminalColor {
	return Default().CompleteColor(trueColor, ansi256, ansi)
}

// HasDarkBackground returns whether the default output has a dark background.
//
// If the background color is not set, it will default to true.
//
// This function is a convenience function that uses the default adaptive color
// output which is based on the standard input and output and environment
// variables.
func HasDarkBackground() bool {
	return Default().HasDarkBackground()
}

// QueryBackgroundColor queries, caches, and returns the terminal background
// color for the default output.
//
// This function is a convenience function that uses the default adaptive color
// output which is based on the standard input and output and environment
// variables.
func QueryBackgroundColor() (c color.Color, ok bool) {
	return Default().QueryBackgroundColor()
}
