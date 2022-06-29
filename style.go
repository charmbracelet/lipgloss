package lipgloss

import (
	"strings"

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
	alignHorizontalKey
	alignVerticalKey

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

// Value returns the raw, unformatted, underlying string value for this style.
func (s Style) Value() string {
	return s.value
}

// String implements stringer for a Style, returning the rendered result based
// on the rules in this style. An underlying string value must be set with
// Style.SetString prior to using this method.
//
// Deprecated: Use Render(Style, string) instead.
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

// Inherit overlays the style in the argument onto this style by copying each explicitly
// set value from the argument style onto this style if it is not already explicitly set.
// Existing set values are kept intact and not overwritten.
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
//
// Deprecated: Use Render(Style, string) instead.
func (s Style) Render(str string) string {
	return renderer.Render(s, str)
}

func (s Style) applyMargins(re *Renderer, str string, inline bool) string {
	var (
		topMargin    = s.getAsInt(marginTopKey)
		rightMargin  = s.getAsInt(marginRightKey)
		bottomMargin = s.getAsInt(marginBottomKey)
		leftMargin   = s.getAsInt(marginLeftKey)

		styler termenv.Style
	)

	bgc := s.getAsColor(marginBackgroundKey)
	if bgc != noColor {
		styler = styler.Background(re.color(bgc))
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
