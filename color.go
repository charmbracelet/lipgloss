package lipgloss

import (
	"github.com/muesli/termenv"
)

var noColor = termenv.NoColor{}

// TerminalColor is a color intended to be rendered in the terminal.
type TerminalColor interface {
	termenv.Color
}
