package lipgloss

import (
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/termenv"
)

// Whitespace is a whitespace renderer.
type Whitespace struct {
	re    *Renderer
	style termenv.Style
	chars string
}

// NewWhitespace creates a new whitespace renderer. The order of the options
// matters, it you'r using WithWhitespaceRenderer, make sure it comes first as
// other options might depend on it.
func NewWhitespace(opts ...WhitespaceOption) *Whitespace {
	w := &Whitespace{re: renderer}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Render whitespaces.
func (w Whitespace) render(width int) string {
	if w.chars == "" {
		w.chars = " "
	}

	r := []rune(w.chars)
	j := 0
	b := strings.Builder{}

	// Cycle through runes and print them into the whitespace.
	for i := 0; i < width; {
		b.WriteRune(r[j])
		j++
		if j >= len(r) {
			j = 0
		}
		i += ansi.PrintableRuneWidth(string(r[j]))
	}

	// Fill any extra gaps white spaces. This might be necessary if any runes
	// are more than one cell wide, which could leave a one-rune gap.
	short := width - ansi.PrintableRuneWidth(b.String())
	if short > 0 {
		b.WriteString(strings.Repeat(" ", short))
	}

	return w.style.Styled(b.String())
}

// WhitespaceOption sets a styling rule for rendering whitespace.
type WhitespaceOption func(*Whitespace)

// WithWhitespaceRenderer sets the lipgloss renderer to be used for rendering.
func WithWhitespaceRenderer(r *Renderer) WhitespaceOption {
	return func(w *Whitespace) {
		w.re = r
	}
}

// WithWhitespaceForeground sets the color of the characters in the whitespace.
func WithWhitespaceForeground(c TerminalColor) WhitespaceOption {
	return func(w *Whitespace) {
		w.style = w.style.Foreground(c.color(w.re))
	}
}

// WithWhitespaceBackground sets the background color of the whitespace.
func WithWhitespaceBackground(c TerminalColor) WhitespaceOption {
	return func(w *Whitespace) {
		w.style = w.style.Background(c.color(w.re))
	}
}

// WithWhitespaceChars sets the characters to be rendered in the whitespace.
func WithWhitespaceChars(s string) WhitespaceOption {
	return func(w *Whitespace) {
		w.chars = s
	}
}
