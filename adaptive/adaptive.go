package adaptive

import (
	"image/color"
	"io"

	"github.com/charmbracelet/x/term"
)

// queryBackgroundColor returns the terminal background color.
func queryBackgroundColor(in io.Reader, out io.Writer) (color.Color, bool) {
	if output, ok := out.(term.File); ok {
		if !term.IsTerminal(output.Fd()) {
			return nil, false
		}
	}

	if input, ok := in.(term.File); ok {
		state, err := term.MakeRaw(input.Fd())
		if err != nil {
			return nil, false
		}

		defer term.Restore(input.Fd(), state) //nolint:errcheck
	}

	c, err := term.QueryBackgroundColor(in, out)
	if err != nil || c == nil {
		return nil, false
	}

	return c, true
}
