package lipgloss

import (
	"strings"
	"unicode"

	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
)

// ANSI reset sequence.
const resetSeq = termenv.CSI + termenv.ResetSeq + "m"

// Style describes formatting instructions for a given string.
type Style struct {
	Bold          bool
	Italic        bool
	Underline     bool
	Strikethrough bool
	Reverse       bool
	Blink         bool
	Faint         bool
	Foreground    ColorType
	Background    ColorType

	// If the string contains multiple lines, they'll wrap at this value.
	// Lines will also be padded with spaces so they'll all be the same width.
	Width int

	// Text alignment.
	Align Align

	// Padding. This will be colored according to the Background value if
	// StylePadding is true.
	LeftPadding   int
	RightPadding  int
	TopPadding    int
	BottomPadding int

	// Whether or not to apply styling to the padding. Most notably, this
	// determines whether or not the indent background color is styled.
	StylePadding bool

	// Margins. These will never be colored.
	LeftMargin   int
	RightMargin  int
	TopMargin    int
	BottomMargin int

	// If set, we truncate lines at this value after all other style has been
	// applied. That is to say, the physical width of strings will not exceed
	// this value, so this can be handy when building user interfaces.
	MaxWidth int

	// Whether or not to remove trailing spaces with no background color. By
	// default we leave them in.
	RenderClearTrailingSpaces bool

	// Whether to draw underlines on spaces (like padding). We don't do this by
	// default as it's likely not what people want, but you can turn it on if
	// you so desire.
	RenderUnderlinesOnSpaces bool
}

// Apply applies formatting to a given string.
func (s Style) Apply(str string) string {
	return s.apply(str, false)
}

// Apply as applies formatting to the given string, removing newlines and
// skipping over block-level rules. Use this internally if you require an
// inline-style only in your library or component.
func (s Style) ApplyInline(str string) string {
	return s.apply(str, true)
}

// WithMaxWidth applies a max width to a given style. This is useful in
// enforcing a certain width at render time, particularly with aribtrary
// strings and styles.
//
// Example:
//
//     var userInput string = "..."
//     var userStyle = text.Style{ /* ... */ }
//     fmt.Println(userStyle.WithMaxWidth(16).Apply(userInput))
//
func (s Style) WithMaxWidth(n int) Style {
	s.MaxWidth = n
	return s
}

func (s Style) apply(str string, singleLine bool) string {
	var (
		styler = termenv.Style{}

		// A copy of the main termenv styler, but without underlines. Used to
		// not render underlines on spaces, if applicable. It's a pointer so
		// we can treat it like a maybe monad, since it won't always be
		// applicable.
		noUnderlineStyler *termenv.Style
	)

	if s.Bold {
		styler = styler.Bold()
	}
	if s.Italic {
		styler = styler.Italic()
	}
	if s.Strikethrough {
		styler = styler.CrossOut()
	}
	if s.Reverse {
		styler = styler.Reverse()
	}
	if s.Blink {
		styler = styler.Blink()
	}
	if s.Faint {
		styler = styler.Faint()
	}

	switch c := s.Foreground.(type) {
	case Color, AdaptiveColor:
		styler = styler.Foreground(color(c.value()))
	}

	switch c := s.Background.(type) {
	case Color, AdaptiveColor:
		styler = styler.Background(color(c.value()))
	}

	if s.Underline {
		if s.Underline && !s.RenderUnderlinesOnSpaces {
			stylerCopy := styler
			noUnderlineStyler = &stylerCopy
		}
		styler = styler.Underline()
	}

	// Strip spaces in single line mode
	if singleLine {
		str = strings.Replace(str, "\n", "", -1)
	}

	if !s.StylePadding {
		str = styler.Styled(str)
	}

	// Word wrap
	if !singleLine && s.Width > 0 {
		str = wordwrap.String(str, s.Width-s.LeftPadding-s.RightPadding)
	}

	// Is a background color set?
	backgroundColorSet := true
	switch s.Background.(type) {
	case nil, noColor:
		backgroundColorSet = false
	}

	// Left/right padding
	str = padLeft(str, s.LeftPadding)
	if !s.RenderClearTrailingSpaces || backgroundColorSet {
		str = padRight(str, s.RightPadding, s.StylePadding)
	}

	// Top/bottom padding
	if !singleLine && s.TopPadding > 0 {
		str = strings.Repeat("\n", s.TopPadding) + str
	}
	if !singleLine && s.BottomPadding > 0 {
		str += strings.Repeat("\n", s.BottomPadding)
	}

	numLines := strings.Count(str, "\n")

	// Set alignment. This will also pad short lines with spaces so that all
	// lines are the same length.
	if numLines > 0 && (!s.RenderClearTrailingSpaces || backgroundColorSet || s.Align != AlignLeft) {
		str = align(str, s.Align)
	}

	if s.StylePadding {
		// We have to do some extra work to not render underlines on spaces
		if noUnderlineStyler != nil {
			var b strings.Builder

			for _, c := range str {
				if unicode.IsSpace(c) {
					b.WriteString(noUnderlineStyler.Styled(string(c)))
					continue
				}
				b.WriteString(styler.Styled(string(c)))
			}

			str = b.String()

		} else {
			str = styler.Styled(str)
		}
	}

	// Add left margin
	str = padLeft(str, s.LeftMargin)

	// Add right margin
	if !s.RenderClearTrailingSpaces {
		str = padRight(str, s.RightMargin, false)
	}

	// Top/bottom margin
	if !singleLine {
		var maybeSpaces string

		if s.RenderClearTrailingSpaces {
			_, width := getLines(str)
			maybeSpaces = strings.Repeat(" ", width)
		}

		if s.TopMargin > 0 {
			str = strings.Repeat(maybeSpaces+"\n", s.TopMargin) + str
		}
		if s.BottomMargin > 0 {
			str += strings.Repeat("\n"+maybeSpaces, s.BottomMargin) + "\n"
		}
	}

	// Truncate accoridng to MaxWidth
	if s.MaxWidth > 0 {
		lines := strings.Split(str, "\n")

		for i := range lines {
			lines[i] = truncate.String(lines[i], uint(s.MaxWidth))
		}

		str = strings.Join(lines, "\n")
	}

	return str
}

// Apply left padding.
func padLeft(str string, n int) string {
	if n == 0 || str == "" {
		return str
	}
	return indent.String(str, uint(n))
}

// Apply right right padding.
func padRight(str string, n int, stylePadding bool) string {
	if n == 0 || str == "" {
		return str
	}

	padding := strings.Repeat(" ", n)

	var maybeReset string
	if !stylePadding {
		maybeReset = resetSeq
	}

	lines := strings.Split(str, "\n")
	for i := range lines {
		lines[i] += maybeReset + padding
	}

	return strings.Join(lines, "\n")
}
