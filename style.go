package lipgloss

import (
	"strings"

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
	// colorWhitespace is true.
	leftPadding   *int
	rightPadding  *int
	topPadding    *int
	bottomPadding *int

	// Whether or not to apply a background color to the whitespace surrounding
	// text blocks. Most notably, this determines whether or not the indent
	// background color will be styled.
	colorWhitespace *bool

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

	// Whether to apply underlines and strikes to whitespace like padding. We
	// don't do this by default as it's likely not what people usually want,
	// but it can can turned it on for certain graphic effects.
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
	if i.colorWhitespace != nil {
		o.colorWhitespace = i.colorWhitespace
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
		singleLine       = s.inline != nil && *s.inline
		styler           termenv.Style
		whitespaceStyler termenv.Style // won't always be applicable
	)

	// Is a background color set?
	var hasBackgroundColor bool
	if s.background != nil {
		_, ok := (*s.background).(noColor)
		hasBackgroundColor = !ok
	}

	// By default, we color padding and space around paragraphs by default if
	// a background color is set.
	colorWhitespace := hasBackgroundColor && (s.colorWhitespace == nil || *s.colorWhitespace)

	// Helper conditions. Niladic types make our conditions rather long and
	// convoluted.
	underline := s.underline != nil && *s.underline
	underlineWhitespace := underline && s.underlineWhitespace != nil && *s.underlineWhitespace
	strike := s.strikethrough != nil && *s.strikethrough
	strikeWhitespace := strike && s.strikethroughWhitespace != nil && *s.strikethroughWhitespace

	// Whether or not to apply foreground styling to whitespace for things like
	// strikethroughs and underlines.
	styleWhitespace := underlineWhitespace || strikeWhitespace

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

			if styleWhitespace {
				whitespaceStyler = whitespaceStyler.Foreground(color(c.value()))
			}
		}
	}

	if s.background != nil {
		switch c := (*s.background).(type) {
		case Color, AdaptiveColor:
			styler = styler.Background(color(c.value()))

			if colorWhitespace {
				whitespaceStyler = whitespaceStyler.Background(color(c.value()))
			}
		}
	}

	if underlineWhitespace {
		whitespaceStyler = whitespaceStyler.Underline()
	}
	if strikeWhitespace {
		whitespaceStyler = whitespaceStyler.CrossOut()
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

	// We draw clear trailing spaces by default
	drawClearTrailingSpaces := s.drawClearTrailingSpaces == nil || *s.drawClearTrailingSpaces

	// Render core text
	{
		var b strings.Builder
		l := strings.Split(str, "\n")
		for i := range l {
			b.WriteString(styler.Styled(l[i]))
			if i != len(l)-1 {
				b.WriteRune('\n')
			}
		}

		str = b.String()
	}

	// Left/right padding
	if s.leftPadding != nil {
		var st *termenv.Style
		if colorWhitespace || styleWhitespace {
			st = &whitespaceStyler
		}
		str = padLeft(str, *s.leftPadding, st)
	}
	if (colorWhitespace || drawClearTrailingSpaces) && s.rightPadding != nil {
		var st *termenv.Style
		if colorWhitespace || styleWhitespace {
			st = &whitespaceStyler
		}
		str = padRight(str, *s.rightPadding, st)
	}

	// Top/bottom padding
	if s.topPadding != nil && *s.topPadding > 0 && !singleLine {
		str = strings.Repeat("\n", *s.topPadding) + str
	}
	if s.bottomPadding != nil && *s.bottomPadding > 0 && !singleLine {
		str += strings.Repeat("\n", *s.bottomPadding)
	}

	// Set alignment. This will also pad short lines with spaces so that all
	// lines are the same length, so we run it under a few different conditions
	// beyond alignment.
	{
		numLines := strings.Count(str, "\n")

		align := AlignLeft
		if s.align != nil {
			align = *s.align
		}

		if numLines > 0 && (align != AlignLeft || drawClearTrailingSpaces || colorWhitespace) {
			var st *termenv.Style
			if colorWhitespace || styleWhitespace {
				st = &whitespaceStyler
			}
			str = alignText(str, align, st)
		}
	}

	// Add left margin
	if s.leftMargin != nil {
		str = padLeft(str, *s.leftMargin, nil)
	}

	// Add right margin
	if s.rightMargin != nil {
		str = padRight(str, *s.rightMargin, nil)
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
func padLeft(str string, n int, style *termenv.Style) string {
	if n == 0 {
		return str
	}

	sp := strings.Repeat(" ", n)
	if style != nil {
		sp = style.Styled(sp)
	}

	b := strings.Builder{}
	l := strings.Split(str, "\n")

	for i := range l {
		b.WriteString(sp)
		b.WriteString(l[i])
		if i != len(l)-1 {
			b.WriteRune('\n')
		}
	}

	//return indent.String(str, uint(n))
	return b.String()
}

// Apply right right padding.
func padRight(str string, n int, style *termenv.Style) string {
	if n == 0 || str == "" {
		return str
	}

	sp := strings.Repeat(" ", n)
	if style != nil {
		sp = style.Styled(sp)
	}

	b := strings.Builder{}
	l := strings.Split(str, "\n")

	for i := range l {
		b.WriteString(l[i])
		b.WriteString(sp)
		if i != len(l)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}
