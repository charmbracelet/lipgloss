package adaptive

import (
	"image/color"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

// HasDarkBackground is true if the terminal has a dark background. It defaults
// to true if the background color cannot be determined.
var HasDarkBackground = func() bool {
	bg, err := QueryBackgroundColor(os.Stdin, os.Stdout)
	if err != nil {
		return true
	}

	return lipgloss.IsDarkColor(bg)
}()

// QueryBackgroundColor queries the terminal for the background color. If the
// terminal does not support querying the background color, nil is returned.
//
// Example usage:
//
//	bg, _ := adaptive.QueryBackgroundColor(os.Stdin, os.Stdout)
//	fmt.Printf("Background color: %v, isDark: %v\n", bg, lipgloss.IsDarkColor(bg))
func QueryBackgroundColor(stdin term.File, stdout term.File) (color.Color, error) {
	state, err := term.MakeRaw(stdin.Fd())
	if err != nil {
		return nil, err
	}

	defer term.Restore(stdin.Fd(), state) //nolint:errcheck

	return queryBackgroundColor(stdin, stdout)
}
