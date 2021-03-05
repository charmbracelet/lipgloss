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

// Style contains formatting instructions for a given string.
type Style struct {
	bold          *bool
	italic        *bool
	underline     *bool
	strikethrough *bool
	reverse       *bool
	blink         *bool
	faint         *bool
	foreground    *ColorType
	background    *ColorType

	// If the string contains multiple lines, they'll wrap at this value.
	// Lines will also be padded with spaces so they'll all be the same width.
	width *int

	// Text alignment.
	align *Align

	// Padding. This will be colored according to the Background value if
	// StylePadding is true.
	leftPadding   *int
	rightPadding  *int
	topPadding    *int
	bottomPadding *int

	// Whether or not to apply styling to the padding. Most notably, this
	// determines whether or not the indent background color is styled.
	stylePadding *bool

	// Margins. These will never be colored.
	leftMargin   *int
	rightMargin  *int
	topMargin    *int
	bottomMargin *int

	// When set, render as a single line.
	inline *bool

	// If set, we truncate lines at this value after all other style has been
	// applied. That is to say, the physical width of strings will not exceed
	// this value, so this can be handy when building user interfaces.
	maxWidth *int

	// Whether or not to remove trailing spaces with no background color. By
	// default we leave them in.
	drawClearTrailingSpaces *bool

	// Whether to apply underlines and strikes to whitespace like margins and
	// padding. We don't do this by default as it's likely not what people
	// want, but you can turn it on if you so desire.
	underlineWhitespace     *bool
	strikethroughWhitespace *bool
}

// NewStyle returns a new, empty Style. It's syntatic sugar for the literal
// Style{}.
func NewStyle() Style {
	return Style{}
}

// Inherit creates a new style using a given style as the starting point. It's
// syntatic sugar for Style{}.Inherit().
func Inherit(s Style) Style {
	return s
}

// Inherit takes values from another style and applies them to this style. Only
// values explicitly set on the style in argument will be applied. Values on
// the style of parent of this method will be overwritten.
func (o Style) Inherit(i Style) Style {
	// We could use reflection here, but it's slow, so we're doing things the
	// long way.

	// Inline
	if i.bold != nil {
		o.bold = i.bold
	}
	if i.italic != nil {
		o.italic = i.italic
	}
	if i.underline != nil {
		o.underline = i.underline
	}
	if i.strikethrough != nil {
		o.strikethrough = i.strikethrough
	}
	if i.reverse != nil {
		o.reverse = i.reverse
	}
	if i.blink != nil {
		o.blink = i.blink
	}
	if i.faint != nil {
		o.faint = i.faint
	}

	// Colors
	if i.foreground != nil {
		o.foreground = i.foreground
	}
	if i.background != nil {
		o.background = i.background
	}

	// Width
	if i.width != nil {
		o.width = i.width
	}

	// Alignment
	if i.align != nil {
		o.align = i.align
	}

	// Padding
	if i.leftPadding != nil {
		o.leftPadding = i.leftPadding
	}
	if i.rightPadding != nil {
		o.rightPadding = i.rightPadding
	}
	if i.topPadding != nil {
		o.rightPadding = i.topPadding
	}
	if i.bottomPadding != nil {
		o.bottomPadding = i.bottomPadding
	}
	if i.stylePadding != nil {
		o.stylePadding = i.stylePadding
	}

	// Margins
	if i.leftMargin != nil {
		o.leftMargin = i.leftMargin
	}
	if i.rightMargin != nil {
		o.rightMargin = i.rightMargin
	}
	if i.topMargin != nil {
		o.topMargin = i.topMargin
	}
	if i.bottomMargin != nil {
		o.bottomMargin = i.bottomMargin
	}

	// Etc
	if i.maxWidth != nil {
		o.maxWidth = i.maxWidth
	}
	if i.inline != nil {
		o.inline = i.inline
	}
	if i.drawClearTrailingSpaces != nil {
		o.drawClearTrailingSpaces = i.drawClearTrailingSpaces
	}
	if i.underlineWhitespace != nil {
		o.underlineWhitespace = i.underlineWhitespace
	}
	if i.strikethroughWhitespace != nil {
		o.strikethroughWhitespace = i.strikethroughWhitespace
	}

	return o
}

// Render applies formatting to a given string.
func (s Style) Render(str string) string {
	var (
		singleLine = s.inline != nil && *s.inline

		styler termenv.Style

		// Additional styling helpers for spaces, which won't always be
		// applicable
		styleSpaces bool
		spaceStyler termenv.Style
	)

	// Helper conditions
	underline := s.underline != nil && *s.underline
	underlineWhitespace := s.underlineWhitespace != nil && *s.underlineWhitespace
	strike := s.strikethrough != nil && *s.strikethrough
	strikeWhitespace := s.strikethroughWhitespace != nil && *s.strikethroughWhitespace

	if (underline && !underlineWhitespace) || (strike && !strikeWhitespace) {
		styleSpaces = true
	}

	if s.bold != nil && *s.bold {
		styler = styler.Bold()
	}
	if s.italic != nil && *s.italic {
		styler = styler.Italic()
	}
	if s.reverse != nil && *s.reverse {
		styler = styler.Reverse()
	}
	if s.blink != nil && *s.blink {
		styler = styler.Blink()
	}
	if s.faint != nil && *s.faint {
		styler = styler.Faint()
	}

	if s.foreground != nil {
		switch c := (*s.foreground).(type) {
		case Color, AdaptiveColor:
			styler = styler.Foreground(color(c.value()))

			if styleSpaces {
				spaceStyler = spaceStyler.Foreground(color(c.value()))
			}
		}
	}

	if s.background != nil {
		switch c := (*s.background).(type) {
		case Color, AdaptiveColor:
			styler = styler.Background(color(c.value()))

			if styleSpaces {
				spaceStyler = spaceStyler.Background(color(c.value()))
			}
		}
	}

	if styleSpaces {
		if underlineWhitespace {
			spaceStyler = spaceStyler.Underline()
		}
		if strikeWhitespace {
			spaceStyler = spaceStyler.CrossOut()
		}
	}

	if underline {
		styler = styler.Underline()
	}

	if strike {
		styler = styler.CrossOut()
	}

	// Strip newlines in single line mode
	if singleLine {
		str = strings.Replace(str, "\n", "", -1)
	}

	if s.stylePadding != nil && !*s.stylePadding {
		str = styler.Styled(str)
	}

	// Word wrap
	if !singleLine && s.width != nil && *s.width > 0 {
		var leftPadding, rightPadding int
		if s.leftPadding != nil {
			leftPadding = *s.leftPadding
		}
		if s.rightPadding != nil {
			rightPadding = *s.rightPadding
		}
		str = wordwrap.String(str, *s.width-leftPadding-rightPadding)
	}

	// Is a background color set?
	var hasBackgroundColor bool
	if s.background != nil {
		hasBackgroundColor = true
		switch (*s.background).(type) {
		case noColor:
			hasBackgroundColor = false
		}
	}

	// Left/right padding
	if s.leftPadding != nil {
		str = padLeft(str, *s.leftPadding)
	}

	drawClearTrailingSpaces := true
	if s.drawClearTrailingSpaces != nil {
		drawClearTrailingSpaces = *s.drawClearTrailingSpaces
	}

	if drawClearTrailingSpaces || hasBackgroundColor {
		var rightPadding int
		if s.rightPadding != nil {
			rightPadding = *s.rightPadding
		}
		var stylePadding bool
		if s.stylePadding != nil {
			stylePadding = *s.stylePadding
		}
		str = padRight(str, rightPadding, stylePadding)
	}

	// Top/bottom padding
	if s.topPadding != nil && *s.topPadding > 0 && !singleLine {
		str = strings.Repeat("\n", *s.topPadding) + str
	}
	if s.bottomPadding != nil && *s.bottomPadding > 0 && !singleLine {
		str += strings.Repeat("\n", *s.bottomPadding)
	}

	numLines := strings.Count(str, "\n")

	// Set alignment. This will also pad short lines with spaces so that all
	// lines are the same length, so we run it under a few different conditions
	// beyond alignment.
	{
		align := AlignLeft
		if s.align != nil {
			align = *s.align
		}

		if numLines > 0 && (align != AlignLeft || drawClearTrailingSpaces || hasBackgroundColor) {
			str = alignText(str, align)
		}
	}

	if s.stylePadding != nil && *s.stylePadding {

		// We have to do some extra work to not render underlines and/or
		// strikes on spaces.
		if styleSpaces {
			var b strings.Builder

			for _, c := range str {
				if unicode.IsSpace(c) {
					b.WriteString(spaceStyler.Styled(string(c)))
					continue
				}
				b.WriteString(styler.Styled(string(c)))
			}

			str = b.String()

		} else {
			str = styler.Styled(str)
		}

	} else {
		str = styler.Styled(str)
	}

	// Add left margin
	if s.leftMargin != nil {
		str = padLeft(str, *s.leftMargin)
	}

	// Add right margin
	if s.rightMargin != nil && drawClearTrailingSpaces {
		str = padRight(str, *s.rightMargin, false)
	}

	// Top/bottom margin
	if !singleLine {
		var maybeSpaces string

		if drawClearTrailingSpaces {
			_, width := getLines(str)
			maybeSpaces = strings.Repeat(" ", width)
		}

		if s.topMargin != nil && *s.topMargin > 0 {
			str = strings.Repeat(maybeSpaces+"\n", *s.topMargin) + str
		}
		if s.bottomMargin != nil && *s.bottomMargin > 0 {
			str += strings.Repeat("\n"+maybeSpaces, *s.bottomMargin) + "\n"
		}
	}

	// Truncate accoridng to MaxWidth
	if s.maxWidth != nil && *s.maxWidth > 0 {
		lines := strings.Split(str, "\n")

		for i := range lines {
			lines[i] = truncate.String(lines[i], uint(*s.maxWidth))
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
