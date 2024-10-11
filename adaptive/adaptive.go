package adaptive

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

// Stdin, Stdout, and HasDarkBackground are the standard input, output, and
// default background color value to use. They can be overridden by the
// importing program.
var (
	Stdin             = os.Stdin
	Stdout            = os.Stdout
	HasDarkBackground = true
)

// colorFn is the light-dark Lip Gloss color function to use to determine the
// appropriate color based on the terminal's background color.
// When a program imports this package, it will query the terminal's background
// color and use it to determine whether to use the light or dark color.
var colorFn lipgloss.Adapt

func init() {
	Query()
}

// Query queries the terminal's background color and updates the color function
// accordingly.
func Query() {
	colorFn = lipgloss.Adapt(func() bool {
		state, err := term.MakeRaw(Stdin.Fd())
		if err != nil {
			return HasDarkBackground
		}

		defer term.Restore(Stdin.Fd(), state) //nolint:errcheck

		bg, err := queryBackgroundColor(Stdin, Stdout)
		if err == nil {
			return lipgloss.IsDarkColor(bg)
		}

		return HasDarkBackground
	}())
}
