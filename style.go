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

// Bold sets a bold formatting rule.
func (s Style) Bold(v bool) Style {
	s.bold = &v
	return s
}

// Italic sets an italic formatting rule. In some terminal emulators this will
// render with "reverse" coloring if not italic font variant is available.
func (s Style) Italic(v bool) Style {
	s.italic = &v
	return s
}

// Underine sets an underline rule. By default, underlines will not be drawn on
// whitespace like margins and padding. To change this behavior set
// renderUnderlinesOnSpaces.
func (s Style) Underline(v bool) Style {
	s.underline = &v
	return s
}

// Strikethrough sets a strikethrough rule. By default, strikes will not be
// drawn on whitespace like margins and padding. To change this behavior set
// renderStrikethroughOnSpaces.
func (s Style) Strikethrough(v bool) Style {
	s.strikethrough = &v
	return s
}

// Reverse sets a rule for inverting foreground and background colors.
func (s Style) Reverse(v bool) Style {
	s.reverse = &v
	return s
}

// Blink sets a rule for blinking forground text.
func (s Style) Blink(v bool) Style {
	s.blink = &v
	return s
}

// Faint sets a rule for rendering the foreground color in a dimmer shade.
func (s Style) Faint(v bool) Style {
	s.faint = &v
	return s
}

// Foreground sets a foreground color.
//
//     // Sets the foreground to blue
//     s := lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff"))
//
//     // Removes the foreground color
//     s.Foreground(lipgloss.NoColor)
//
func (s Style) Foreground(c ColorType) Style {
	s.foreground = &c
	return s
}

// Background sets a background color.
func (s Style) Background(c ColorType) Style {
	s.background = &c
	return s
}

// Width sets the width of the block before applying margins and padding. This
// effects when.
func (s Style) Width(i int) Style {
	s.width = &i
	return s
}

// Align sets a text alignment rule.
func (s Style) Align(a Align) Style {
	s.align = &a
	return s
}

// Padding is a shorthand method for setting padding on all sides at once.
//
// With one argument, the value is applied to all sides.
//
// With two arguments, the value is applied to the vertical and horizontal
// sides, in that order.
//
// With three arguments, the value is applied to the top side, the horizontal
// sides, and the bottom side, in that order.
//
// With four arguments, the value is applied clockwise starting from the top
// side, followed by the right side, then the bottom, and finally the left.
//
// With more than four arguments no padding will be added.
func (s Style) Padding(i ...int) Style {
	top, right, bottom, left, ok := whichSides(i...)
	if !ok {
		return s
	}

	s.topPadding = &top
	s.rightPadding = &right
	s.bottomPadding = &bottom
	s.leftPadding = &left
	return s
}

func (s Style) LeftPadding(i int) Style {
	s.leftPadding = &i
	return s
}

func (s Style) RightPadding(i int) Style {
	s.rightPadding = &i
	return s
}

func (s Style) TopPadding(i int) Style {
	s.topPadding = &i
	return s
}

func (s Style) BottomPadding(i int) Style {
	s.bottomPadding = &i
	return s
}

func (s Style) StylePadding(v bool) Style {
	s.stylePadding = &v
	return s
}

// Margin is a shorthand method for setting margins on all sides at once.
//
// With one argument, the value is applied to all sides.
//
// With two arguments, the value is applied to the vertical and horizontal
// sides, in that order.
//
// With three arguments, the value is applied to the top side, the horizontal
// sides, and the bottom side, in that order.
//
// With four arguments, the value is applied clockwise starting from the top
// side, followed by the right side, then the bottom, and finally the left.
//
// With more than four arguments no padding will be added.
func (s Style) Margin(i ...int) Style {
	top, right, bottom, left, ok := whichSides(i...)
	if !ok {
		return s
	}

	s.topMargin = &top
	s.rightMargin = &right
	s.bottomMargin = &bottom
	s.leftMargin = &left
	return s
}

func (s Style) LeftMargin(i int) Style {
	s.leftMargin = &i
	return s
}

func (s Style) RightMargin(i int) Style {
	s.rightMargin = &i
	return s
}

func (s Style) TopMargin(i int) Style {
	s.topMargin = &i
	return s
}

func (s Style) BottomMargin(i int) Style {
	s.bottomMargin = &i
	return s
}

func (s Style) MaxWidth(i int) Style {
	s.maxWidth = &i
	return s
}

// Whether or not to draw trailing spaces with no background color. By default
// we leave them in.
func (s Style) DrawClearTrailingSpaces(v bool) Style {
	s.drawClearTrailingSpaces = &v
	return s
}

func (s Style) UnderlineWhitespace(v bool) Style {
	s.underlineWhitespace = &v
	return s
}

func (s Style) StrikethroughWhitespace(v bool) Style {
	s.strikethroughWhitespace = &v
	return s
}

// Inherit takes values from another style and applies them to this style. Only
// values explicitly set on the style in argument will be applied. Values on
// this style struct will be overwritten.
func (o Style) Inherit(i Style) Style {
	// We could use reflection here, but it's slow, so we're doing things
	// the old-fashioned way.

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
	if i.drawClearTrailingSpaces != nil {
		o.drawClearTrailingSpaces = i.drawClearTrailingSpaces
	}
	if i.underlineWhitespace != nil {
		o.underlineWhitespace = i.underlineWhitespace
	}
	if i.strikethroughWhitespace != nil {
		o.strikethroughWhitespace = i.strikethroughWhitespace
	}

	return i
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
	s.maxWidth = &n
	return s
}

func (s Style) apply(str string, singleLine bool) string {
	var (
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

// whichSides is a helper method for setting values on sides of a block based
// on the number of arguments. It follows the CSS shorthand rules for blocks
// like margin, padding. and borders. Here are how the rules work:
//
// 0 args:  do nothing
// 1 arg:   all sides
// 2 args:  top -> bottom
// 3 args:  top -> horizontal -> bottom
// 4 args:  top -> right -> bottom -> left
// 5+ args: do nothing
func whichSides(i ...int) (top, right, bottom, left int, ok bool) {
	switch len(i) {
	case 1:
		top = i[0]
		bottom = i[0]
		left = i[0]
		right = i[0]
		ok = true
	case 2:
		top = i[0]
		bottom = i[0]
		left = i[1]
		right = i[1]
		ok = true
	case 3:
		top = i[0]
		left = i[1]
		right = i[1]
		bottom = i[2]
		ok = true
	case 4:
		top = i[0]
		right = i[1]
		bottom = i[2]
		left = i[3]
		ok = true
	}
	return
}
