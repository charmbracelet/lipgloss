package lipgloss

import (
	"strings"
	"unicode"

	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
)

// propKey is a property for a key
type propKey int

// Available properties
const (
	boldKey propKey = iota
	italicKey
	underlineKey
	strikethroughKey
	reverseKey
	blinkKey
	faintKey
	foregroundKey
	backgroundKey
	widthKey
	alignKey
	topPaddingKey
	rightPaddingKey
	bottomPaddingKey
	leftPaddingKey
	colorWhitespaceKey
	topMarginKey
	rightMarginKey
	bottomMarginKey
	leftMarginKey
	inlineKey
	maxWidthKey
	drawClearTrailingSpacesKey
	underlineWhitespaceKey
	strikethroughWhitespaceKey
	underlineSpacesKey
	strikethroughSpacesKey
)

// NewStyle returns a new, empty Style.  While it's syntactis sugar for
// make(Style), it's recommended to use this function for creating styles
// incase the underlying implementation changes.
func NewStyle() Style {
	return make(Style)
}

// Style contains property definitions that comprise a style as a whole.
type Style map[propKey]interface{}

// Copy returns a copy of this style.
func (s Style) Copy() Style {
	o := make(Style)
	for k, v := range s {
		o[k] = v
	}
	return o
}

// Inherit takes values from the style in the argument applies them to this
// style, overwriting existing definitions. Only values explicitly set on the
// style in argument will be applied.
//
// Margins and padding are not inherited.
func (o Style) Inherit(i Style) {
	for k, v := range i {
		switch k {
		case topMarginKey, rightMarginKey, bottomMarginKey, leftMarginKey:
			// Margins are not inherited
			continue
		case topPaddingKey, rightPaddingKey, bottomPaddingKey, leftPaddingKey:
			// Padding is not inherited
			continue
		}

		if _, exists := o[k]; exists {
			continue
		}
		o[k] = v
	}
}

// Render applies the defined style formatting to a given string.
func (s Style) Render(str string) string {
	var (
		te           termenv.Style
		teSpace      termenv.Style
		teWhitespace termenv.Style

		bold          = s.getAsBool(boldKey, false)
		italic        = s.getAsBool(italicKey, false)
		underline     = s.getAsBool(underlineKey, false)
		strikethrough = s.getAsBool(strikethroughKey, false)
		reverse       = s.getAsBool(reverseKey, false)
		blink         = s.getAsBool(blinkKey, false)
		faint         = s.getAsBool(faintKey, false)

		fg = s.getAsColor(foregroundKey)
		bg = s.getAsColor(backgroundKey)

		width = s.getAsInt(widthKey)
		align = s.getAsAlign(alignKey)

		topPadding    = s.getAsInt(topPaddingKey)
		rightPadding  = s.getAsInt(rightPaddingKey)
		bottomPadding = s.getAsInt(bottomPaddingKey)
		leftPadding   = s.getAsInt(leftPaddingKey)

		topMargin    = s.getAsInt(topMarginKey)
		rightMargin  = s.getAsInt(rightMarginKey)
		bottomMargin = s.getAsInt(bottomMarginKey)
		leftMargin   = s.getAsInt(leftMarginKey)

		colorWhitespace = s.getAsBool(colorWhitespaceKey, true)
		inline          = s.getAsBool(inlineKey, false)
		maxWidth        = s.getAsInt(maxWidthKey)

		drawClearTrailingSpaces = s.getAsBool(drawClearTrailingSpacesKey, true)
		underlineWhitespace     = s.getAsBool(underlineWhitespaceKey, false)
		strikethroughWhitespace = s.getAsBool(strikethroughWhitespaceKey, false)

		underlineSpaces     = underline && s.getAsBool(underlineSpacesKey, true)
		strikethroughSpaces = strikethrough && s.getAsBool(strikethroughSpacesKey, true)

		// Do we need to style whitespace (padding and space outsode
		// paragraphs) separately?
		styleWhitespace = underlineWhitespace || strikethroughWhitespace

		// Do we need to style spaces separately?
		useSpaceStyler = underlineSpaces || strikethroughSpaces
	)

	if bold {
		te = te.Bold()
	}
	if italic {
		te = te.Italic()
	}
	if underline {
		te = te.Underline()
	}
	if reverse {
		te = te.Reverse()
	}
	if blink {
		te = te.Blink()
	}
	if faint {
		te = te.Faint()
	}

	if fg != "" {
		fgc := color(fg)
		te = te.Foreground(fgc)
		te.Foreground(fgc)
		if styleWhitespace {
			teWhitespace = teWhitespace.Foreground(fgc)
		}
		if useSpaceStyler {
			teSpace = teSpace.Foreground(fgc)
		}
	}

	if bg != "" {
		bgc := color(bg)
		te = te.Background(bgc)
		if colorWhitespace {
			teWhitespace = teWhitespace.Background(bgc)
		}
		if useSpaceStyler {
			teSpace = teSpace.Background(bgc)
		}
	}

	if underline {
		te = te.Underline()
	}
	if strikethrough {
		te = te.CrossOut()
	}

	if underlineWhitespace {
		teWhitespace = teWhitespace.Underline()
	}
	if strikethroughWhitespace {
		teWhitespace = teWhitespace.CrossOut()
	}
	if underlineSpaces {
		teSpace = teSpace.Underline()
	}
	if strikethroughSpaces {
		teSpace = teSpace.CrossOut()
	}

	// Strip newlines in single line mode
	if inline {
		str = strings.Replace(str, "\n", "", -1)
	}

	// Word wrap
	if !inline && width > 0 {
		str = wordwrap.String(str, width-leftPadding-rightPadding)
	}

	// Render core text
	{
		var b strings.Builder

		l := strings.Split(str, "\n")
		for i := range l {
			if useSpaceStyler {
				// Look for spaces and apply a different styler
				for _, r := range []rune(l[i]) {
					if unicode.IsSpace(r) {
						b.WriteString(teSpace.Styled(string(r)))
						continue
					}
					b.WriteString(te.Styled(string(r)))
				}
			} else {
				b.WriteString(te.Styled(l[i]))
			}
			if i != len(l)-1 {
				b.WriteRune('\n')
			}
		}

		str = b.String()
	}

	// Left/right padding
	if leftPadding > 0 {
		var st *termenv.Style
		if colorWhitespace || styleWhitespace {
			st = &teWhitespace
		}
		str = padLeft(str, leftPadding, st)
	}
	if (colorWhitespace || drawClearTrailingSpaces) && rightPadding > 0 {
		var st *termenv.Style
		if colorWhitespace || styleWhitespace {
			st = &teWhitespace
		}
		str = padRight(str, rightPadding, st)
	}

	// Top/bottom padding
	if topPadding > 0 && !inline {
		str = strings.Repeat("\n", topPadding) + str
	}
	if bottomPadding > 0 && !inline {
		str += strings.Repeat("\n", bottomPadding)
	}

	// Set alignment. This will also pad short lines with spaces so that all
	// lines are the same length, so we run it under a few different conditions
	// beyond alignment.
	{
		numLines := strings.Count(str, "\n")

		if numLines > 0 && (align != AlignLeft || drawClearTrailingSpaces || colorWhitespace) {
			var st *termenv.Style
			if colorWhitespace || styleWhitespace {
				st = &teWhitespace
			}
			str = alignText(str, align, st)
		}
	}

	// Add left and right margin
	str = padLeft(str, leftMargin, nil)
	str = padRight(str, rightMargin, nil)

	// Top/bottom margin
	if !inline {
		var maybeSpaces string

		if drawClearTrailingSpaces {
			_, width := getLines(str)
			maybeSpaces = strings.Repeat(" ", width)
		}

		if topMargin > 0 {
			str = strings.Repeat(maybeSpaces+"\n", topMargin) + str
		}
		if bottomMargin > 0 {
			str += strings.Repeat("\n"+maybeSpaces, bottomMargin)
		}
	}

	// Truncate accoridng to MaxWidth
	if maxWidth > 0 {
		lines := strings.Split(str, "\n")

		for i := range lines {
			lines[i] = truncate.String(lines[i], uint(maxWidth))
		}

		str = strings.Join(lines, "\n")
	}

	return str
}

func (s Style) getAsBool(k propKey, defaultVal bool) bool {
	v, ok := s[k]
	if !ok {
		return defaultVal
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return defaultVal
}

func (s Style) getAsColor(k propKey) string {
	v, ok := s[k]
	if !ok {
		return ""
	}
	if c, ok := v.(ColorType); ok {
		return c.value()
	}
	return ""
}

func (s Style) getAsInt(k propKey) int {
	v, ok := s[k]
	if !ok {
		return 0
	}
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}

func (s Style) getAsAlign(k propKey) Align {
	v, ok := s[k]
	if !ok {
		return AlignLeft
	}
	if a, ok := v.(Align); ok {
		return a
	}
	return AlignLeft
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
