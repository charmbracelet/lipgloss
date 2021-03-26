package lipgloss

import (
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/termenv"
)

// whitespace is a whitespace renderer.
type whitespace struct {
	style termenv.Style
	chars string
}

// Render spacespace.
func (w whitespace) render(width int) string {
	if w.chars == "" {
		w.chars = " "
	}

	charWidth := ansi.PrintableRuneWidth(w.chars)
	line := strings.Repeat(w.chars, width/charWidth)

	// Fill in gaps with spaces
	short := width - ansi.PrintableRuneWidth(line)
	if short > 0 {
		line += strings.Repeat(" ", short)
	}

	return w.style.Styled(line)
}

// WhiteSpaceOption sets a styling rule for rendering whitespace.
type WhitespaceOption func(*whitespace)

// WithWhiteSpaceBackground sets the background color of the whitespace.
func WithWhitespaceBackground(c Color) WhitespaceOption {
	return func(w *whitespace) {
		w.style = w.style.Background(c.color())
	}
}
