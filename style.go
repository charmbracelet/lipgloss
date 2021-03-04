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
	renderClearTrailingSpaces *bool

	// Whether to draw underlines on spaces (like padding). We don't do this by
	// default as it's likely not what people want, but you can turn it on if
	// you so desire.
	renderUnderlinesOnSpaces *bool
}

func (s Style) Bold(v bool) Style {
	s.bold = &v
	return s
}

func (s Style) Italic(v bool) Style {
	s.italic = &v
	return s
}

func (s Style) Underline(v bool) Style {
	s.underline = &v
	return s
}

func (s Style) Strikethrough(v bool) Style {
	s.strikethrough = &v
	return s
}

func (s Style) Reverse(v bool) Style {
	s.reverse = &v
	return s
}

func (s Style) Blink(v bool) Style {
	s.blink = &v
	return s
}

func (s Style) Faint(v bool) Style {
	s.faint = &v
	return s
}

func (s Style) Foreground(c ColorType) Style {
	s.foreground = &c
	return s
}

func (s Style) Background(c ColorType) Style {
	s.background = &c
	return s
}

func (s Style) Width(i int) Style {
	s.width = &i
	return s
}

func (s Style) Align(a Align) Style {
	s.align = &a
	return s
}

// Padding is a shorthand method for setting padding on all sides at once.
//
// With one argument, the value is applied to all sides.
//
// With two arguments, the value is applied to the vertical and horizontal sides,
// in that order.
//
// With three arguments, the value is applied to the top side, the horizontal
// sides, and the bottom side, in that order.
//
// With four arguments, the value is applied clockwise starting from the top
// side, followed by the right side, then the bottom, and finally the top.
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
// With two arguments, the value is applied to the vertical and horizontal sides,
// in that order.
//
// With three arguments, the value is applied to the top side, the horizontal
// sides, and the bottom side, in that order.
//
// With four arguments, the value is applied clockwise starting from the top
// side, followed by the right side, then the bottom, and finally the top.
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

func (s Style) RenderClearTrailingSpaces(v bool) Style {
	s.renderClearTrailingSpaces = &v
	return s
}

func (s Style) RenderUnderlinesOnSpaces(v bool) Style {
	s.renderUnderlinesOnSpaces = &v
	return s
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
		styler = termenv.Style{}

		// A copy of the main termenv styler, but without underlines. Used to
		// not render underlines on spaces, if applicable. It's a pointer so
		// we can treat it like a maybe monad, since it won't always be
		// applicable.
		noUnderlineStyler *termenv.Style
	)

	if s.bold != nil && *s.bold {
		styler = styler.Bold()
	}
	if s.italic != nil && *s.italic {
		styler = styler.Italic()
	}
	if s.strikethrough != nil && *s.strikethrough {
		styler = styler.CrossOut()
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
		}
	}

	if s.background != nil {
		switch c := (*s.background).(type) {
		case Color, AdaptiveColor:
			styler = styler.Background(color(c.value()))
		}
	}

	if s.renderUnderlinesOnSpaces != nil && s.underline != nil && *s.underline {
		if !*s.renderUnderlinesOnSpaces {
			stylerCopy := styler
			noUnderlineStyler = &stylerCopy
		}
		styler = styler.Underline()
	}

	// Strip spaces in single line mode
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
	var backgroundColorSet bool
	if s.background != nil {
		backgroundColorSet = true
		switch (*s.background).(type) {
		case noColor:
			backgroundColorSet = false
		}
	}

	// Left/right padding
	if s.leftPadding != nil {
		str = padLeft(str, *s.leftPadding)
	}

	if (s.renderClearTrailingSpaces != nil && !*s.renderClearTrailingSpaces) || backgroundColorSet {
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
	var renderClearTrailingSpaces bool
	{
		align := AlignLeft
		if s.align != nil {
			align = *s.align
		}

		if s.renderClearTrailingSpaces != nil {
			renderClearTrailingSpaces = *s.renderClearTrailingSpaces
		}

		if numLines > 0 && (align != AlignLeft || !renderClearTrailingSpaces || backgroundColorSet) {
			str = alignText(str, align)
		}
	}

	if s.stylePadding != nil && *s.stylePadding {
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
	if s.leftMargin != nil {
		str = padLeft(str, *s.leftMargin)
	}

	// Add right margin
	if s.rightMargin != nil && !renderClearTrailingSpaces {
		str = padRight(str, *s.rightMargin, false)
	}

	// Top/bottom margin
	if !singleLine {
		var maybeSpaces string

		if renderClearTrailingSpaces {
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

// whichEdges is a helper method for setting values on sides of a block based
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
