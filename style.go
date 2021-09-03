package lipgloss

import (
	"strings"
	"unicode"

	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	"github.com/muesli/termenv"
)

// Property for a key.
type propKey int

// Available properties.
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
	heightKey
	alignKey

	// Padding.
	paddingTopKey
	paddingRightKey
	paddingBottomKey
	paddingLeftKey

	colorWhitespaceKey

	// Margins.
	marginTopKey
	marginRightKey
	marginBottomKey
	marginLeftKey
	marginBackgroundKey

	// Border runes.
	borderStyleKey

	// Border edges.
	borderTopKey
	borderRightKey
	borderBottomKey
	borderLeftKey

	// Border foreground colors.
	borderTopForegroundKey
	borderRightForegroundKey
	borderBottomForegroundKey
	borderLeftForegroundKey

	// Border background colors.
	borderTopBackgroundKey
	borderRightBackgroundKey
	borderBottomBackgroundKey
	borderLeftBackgroundKey

	inlineKey
	maxWidthKey
	maxHeightKey
	underlineSpacesKey
	strikethroughSpacesKey
)

// A set of properties.
type rules map[propKey]interface{}

// NewStyle returns a new, empty Style.  While it's syntactic sugar for the
// Style{} primitive, it's recommended to use this function for creating styles
// incase the underlying implementation changes.
func NewStyle() Style {
	return Style{}
}

// Style contains a set of rules that comprise a style as a whole.
type Style struct {
	rules map[propKey]interface{}
	value string
}

// SetString sets the underlying string value for this style. To render once
// the underlying string is set, use the Style.String. This method is
// a convenience for cases when having a stringer implementation is handy, such
// as when using fmt.Sprintf. You can also simply define a style and render out
// strings directly with Style.Render.
func (s Style) SetString(str string) Style {
	s.value = str
	return s
}

// String implements stringer for a Style, returning the rendered result based
// on the rules in this style. An underlying string value must be set with
// Style.SetString prior to using this method.
func (s Style) String() string {
	return s.Render(s.value)
}

// Copy returns a copy of this style, including any underlying string values.
func (s Style) Copy() Style {
	o := NewStyle()
	o.init()
	for k, v := range s.rules {
		o.rules[k] = v
	}
	o.value = s.value
	return o
}

// Inherit takes values from the style in the argument applies them to this
// style, overwriting existing definitions. Only values explicitly set on the
// style in argument will be applied.
//
// Margins, padding, and underlying string values are not inherited.
func (s Style) Inherit(i Style) Style {
	s.init()

	for k, v := range i.rules {
		switch k {
		case marginTopKey, marginRightKey, marginBottomKey, marginLeftKey:
			// Margins are not inherited
			continue
		case paddingTopKey, paddingRightKey, paddingBottomKey, paddingLeftKey:
			// Padding is not inherited
			continue
		case backgroundKey:
			s.rules[k] = v

			// The margins also inherit the background color
			if !s.isSet(marginBackgroundKey) && !i.isSet(marginBackgroundKey) {
				s.rules[marginBackgroundKey] = v
			}
		}

		if _, exists := s.rules[k]; exists {
			continue
		}
		s.rules[k] = v
	}
	return s
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

		width  = s.getAsInt(widthKey)
		height = s.getAsInt(heightKey)
		align  = s.getAsPosition(alignKey)

		topPadding    = s.getAsInt(paddingTopKey)
		rightPadding  = s.getAsInt(paddingRightKey)
		bottomPadding = s.getAsInt(paddingBottomKey)
		leftPadding   = s.getAsInt(paddingLeftKey)

		colorWhitespace = s.getAsBool(colorWhitespaceKey, true)
		inline          = s.getAsBool(inlineKey, false)
		maxWidth        = s.getAsInt(maxWidthKey)
		maxHeight       = s.getAsInt(maxHeightKey)

		underlineSpaces     = underline && s.getAsBool(underlineSpacesKey, true)
		strikethroughSpaces = strikethrough && s.getAsBool(strikethroughSpacesKey, true)

		// Do we need to style whitespace (padding and space outside
		// paragraphs) separately?
		styleWhitespace = reverse

		// Do we need to style spaces separately?
		useSpaceStyler = underlineSpaces || strikethroughSpaces
	)

	if len(s.rules) == 0 {
		return str
	}

	// Enable support for ANSI on the legacy Windows cmd.exe console. This is a
	// no-op on non-Windows systems and on Windows runs only once.
	enableLegacyWindowsANSI()

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
		if reverse {
			teWhitespace = teWhitespace.Reverse()
		}
		te = te.Reverse()
	}
	if blink {
		te = te.Blink()
	}
	if faint {
		te = te.Faint()
	}

	if fg != noColor {
		fgc := fg.color()
		te = te.Foreground(fgc)
		if styleWhitespace {
			teWhitespace = teWhitespace.Foreground(fgc)
		}
		if useSpaceStyler {
			teSpace = teSpace.Foreground(fgc)
		}
	}

	if bg != noColor {
		bgc := bg.color()
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
		wrapAt := width - leftPadding - rightPadding
		str = wordwrap.String(str, wrapAt)
		str = wrap.String(str, wrapAt) // force-wrap long strings
	}

	// Render core text
	{
		var b strings.Builder

		l := strings.Split(str, "\n")
		for i := range l {
			if useSpaceStyler {
				// Look for spaces and apply a different styler
				for _, r := range l[i] {
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

	// Padding
	if !inline {
		if leftPadding > 0 {
			var st *termenv.Style
			if colorWhitespace || styleWhitespace {
				st = &teWhitespace
			}
			str = padLeft(str, leftPadding, st)
		}

		if rightPadding > 0 {
			var st *termenv.Style
			if colorWhitespace || styleWhitespace {
				st = &teWhitespace
			}
			str = padRight(str, rightPadding, st)
		}

		if topPadding > 0 {
			str = strings.Repeat("\n", topPadding) + str
		}

		if bottomPadding > 0 {
			str += strings.Repeat("\n", bottomPadding)
		}
	}

	// Height
	if height > 0 {
		h := strings.Count(str, "\n") + 1
		if height > h {
			str += strings.Repeat("\n", height-h)
		}
	}

	// Set alignment. This will also pad short lines with spaces so that all
	// lines are the same length, so we run it under a few different conditions
	// beyond alignment.
	{
		numLines := strings.Count(str, "\n")

		if !(numLines == 0 && width == 0) {
			var st *termenv.Style
			if colorWhitespace || styleWhitespace {
				st = &teWhitespace
			}
			str = alignText(str, align, width, st)
		}
	}

	if !inline {
		str = s.applyBorder(str)
		str = s.applyMargins(str, inline)
	}

	// Truncate according to MaxWidth
	if maxWidth > 0 {
		lines := strings.Split(str, "\n")

		for i := range lines {
			lines[i] = truncate.String(lines[i], uint(maxWidth))
		}

		str = strings.Join(lines, "\n")
	}

	// Truncate according to MaxHeight
	if maxHeight > 0 {
		lines := strings.Split(str, "\n")
		str = strings.Join(lines[:min(maxHeight, len(lines))], "\n")
	}

	return str
}

func (s Style) applyMargins(str string, inline bool) string {
	var (
		topMargin    = s.getAsInt(marginTopKey)
		rightMargin  = s.getAsInt(marginRightKey)
		bottomMargin = s.getAsInt(marginBottomKey)
		leftMargin   = s.getAsInt(marginLeftKey)

		styler termenv.Style
	)

	bgc := s.getAsColor(marginBackgroundKey)
	if bgc != noColor {
		styler = styler.Background(bgc.color())
	}

	// Add left and right margin
	str = padLeft(str, leftMargin, &styler)
	str = padRight(str, rightMargin, &styler)

	// Top/bottom margin
	if !inline {
		_, width := getLines(str)
		spaces := strings.Repeat(" ", width)

		if topMargin > 0 {
			str = styler.Styled(strings.Repeat(spaces+"\n", topMargin)) + str
		}
		if bottomMargin > 0 {
			str += styler.Styled(strings.Repeat("\n"+spaces, bottomMargin))
		}
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
