package lipgloss

import (
	"strings"

	"github.com/muesli/termenv"
)

// Available properties.
const (
	// Border runes.
	borderStyleKey propKey = 1 << iota

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
)

// Borderer is responsible for handling border rendering.
type Borderer interface {
	// Get
	GetStyle() Border
	GetTop() bool
	GetRight() bool
	GetBottom() bool
	GetLeft() bool
	GetTopForeground() TerminalColor
	GetRightForeground() TerminalColor
	GetBottomForeground() TerminalColor
	GetLeftForeground() TerminalColor
	GetTopBackground() TerminalColor
	GetRightBackground() TerminalColor
	GetBottomBackground() TerminalColor
	GetLeftBackground() TerminalColor
	GetTopSize() int
	GetRightSize() int
	GetBottomSize() int
	GetLeftSize() int

	// Set
	Style(border Border) Borderer
	Top(v bool) Borderer
	Right(v bool) Borderer
	Bottom(v bool) Borderer
	Left(v bool) Borderer
	TopForeground(c TerminalColor) Borderer
	RightForeground(c TerminalColor) Borderer
	BottomForeground(c TerminalColor) Borderer
	LeftForeground(c TerminalColor) Borderer
	TopBackground(c TerminalColor) Borderer
	RightBackground(c TerminalColor) Borderer
	BottomBackground(c TerminalColor) Borderer
	LeftBackground(c TerminalColor) Borderer

	// Unset
	UnsetStyle() Borderer
	UnsetTop() Borderer
	UnsetRight() Borderer
	UnsetBottom() Borderer
	UnsetLeft() Borderer
	UnsetTopForeground() Borderer
	UnsetRightForeground() Borderer
	UnsetBottomForeground() Borderer
	UnsetLeftForeground() Borderer
	UnsetTopBackground() Borderer
	UnsetRightBackground() Borderer
	UnsetBottomBackground() Borderer
	UnsetLeftBackground() Borderer

	Apply(str string) string
}

// NormalBorderer is the default simple borderer.
type NormalBorderer struct {
	r             *Renderer
	props         props
	attrs         int
	style         Border
	topFgColor    TerminalColor
	rightFgColor  TerminalColor
	bottomFgColor TerminalColor
	leftFgColor   TerminalColor
	topBgColor    TerminalColor
	rightBgColor  TerminalColor
	bottomBgColor TerminalColor
	leftBgColor   TerminalColor
}

// Returns whether or not the given property is set.
func (b *NormalBorderer) isSet(k propKey) bool {
	return b.props.has(k)
}

//nolint:unparam
func (b *NormalBorderer) getAsBool(k propKey, defaultVal bool) bool {
	if !b.isSet(k) {
		return defaultVal
	}
	return b.attrs&int(k) != 0
}

func (b *NormalBorderer) getAsColor(k propKey) TerminalColor {
	if !b.isSet(k) {
		return noColor
	}

	var c TerminalColor
	switch k { //nolint:exhaustive
	case borderTopForegroundKey:
		c = b.topFgColor
	case borderRightForegroundKey:
		c = b.rightFgColor
	case borderBottomForegroundKey:
		c = b.bottomFgColor
	case borderLeftForegroundKey:
		c = b.leftFgColor
	case borderTopBackgroundKey:
		c = b.topBgColor
	case borderRightBackgroundKey:
		c = b.rightBgColor
	case borderBottomBackgroundKey:
		c = b.bottomBgColor
	case borderLeftBackgroundKey:
		c = b.leftBgColor
	}

	if c != nil {
		return c
	}

	return noColor
}

// Set a value on the underlying rules map.
func (b *NormalBorderer) set(key propKey, value interface{}) {
	// We don't allow negative integers on any of our other values, so just keep
	// them at zero or above. We could use uints instead, but the
	// conversions are a little tedious, so we're sticking with ints for
	// sake of usability.
	switch key { //nolint:exhaustive
	case borderStyleKey:
		b.style = value.(Border)
	case borderTopForegroundKey:
		b.topFgColor = colorOrNil(value)
	case borderRightForegroundKey:
		b.rightFgColor = colorOrNil(value)
	case borderBottomForegroundKey:
		b.bottomFgColor = colorOrNil(value)
	case borderLeftForegroundKey:
		b.leftFgColor = colorOrNil(value)
	case borderTopBackgroundKey:
		b.topBgColor = colorOrNil(value)
	case borderRightBackgroundKey:
		b.rightBgColor = colorOrNil(value)
	case borderBottomBackgroundKey:
		b.bottomBgColor = colorOrNil(value)
	case borderLeftBackgroundKey:
		b.leftBgColor = colorOrNil(value)
	default:
		if v, ok := value.(bool); ok { //nolint:nestif
			if v {
				b.attrs |= int(key)
			} else {
				b.attrs &^= int(key)
			}
		} else if attrs, ok := value.(int); ok {
			// bool attrs
			if attrs&int(key) != 0 {
				b.attrs |= int(key)
			} else {
				b.attrs &^= int(key)
			}
		}
	}

	// Set the prop on
	b.props = b.props.set(key)
}

// unset unsets a property from a style.
func (b *NormalBorderer) unset(key propKey) {
	b.props = b.props.unset(key)
}

// GetStyle returns the style's border style (type Border). If no value
// is set Border{} is returned.
func (b *NormalBorderer) GetStyle() Border {
	if !b.isSet(borderStyleKey) {
		return noBorder
	}
	return b.style
}

// GetTop returns the style's top border setting. If no value is set
// false is returned.
func (b *NormalBorderer) GetTop() bool {
	return b.getAsBool(borderTopKey, false)
}

// GetRight returns the style's right border setting. If no value is set
// false is returned.
func (b *NormalBorderer) GetRight() bool {
	return b.getAsBool(borderRightKey, false)
}

// GetBottom returns the style's bottom border setting. If no value is
// set false is returned.
func (b *NormalBorderer) GetBottom() bool {
	return b.getAsBool(borderBottomKey, false)
}

// GetLeft returns the style's left border setting. If no value is
// set false is returned.
func (b *NormalBorderer) GetLeft() bool {
	return b.getAsBool(borderLeftKey, false)
}

// GetTopForeground returns the style's border top foreground color. If
// no value is set NoColor{} is returned.
func (b *NormalBorderer) GetTopForeground() TerminalColor {
	return b.getAsColor(borderTopForegroundKey)
}

// GetRightForeground returns the style's border right foreground color.
// If no value is set NoColor{} is returned.
func (b *NormalBorderer) GetRightForeground() TerminalColor {
	return b.getAsColor(borderRightForegroundKey)
}

// GetBottomForeground returns the style's border bottom foreground
// color.  If no value is set NoColor{} is returned.
func (b *NormalBorderer) GetBottomForeground() TerminalColor {
	return b.getAsColor(borderBottomForegroundKey)
}

// GetLeftForeground returns the style's border left foreground
// color.  If no value is set NoColor{} is returned.
func (b *NormalBorderer) GetLeftForeground() TerminalColor {
	return b.getAsColor(borderLeftForegroundKey)
}

// GetTopBackground returns the style's border top background color. If
// no value is set NoColor{} is returned.
func (b *NormalBorderer) GetTopBackground() TerminalColor {
	return b.getAsColor(borderTopBackgroundKey)
}

// GetRightBackground returns the style's border right background color.
// If no value is set NoColor{} is returned.
func (b *NormalBorderer) GetRightBackground() TerminalColor {
	return b.getAsColor(borderRightBackgroundKey)
}

// GetBottomBackground returns the style's border bottom background
// color.  If no value is set NoColor{} is returned.
func (b *NormalBorderer) GetBottomBackground() TerminalColor {
	return b.getAsColor(borderBottomBackgroundKey)
}

// GetLeftBackground returns the style's border left background
// color.  If no value is set NoColor{} is returned.
func (b *NormalBorderer) GetLeftBackground() TerminalColor {
	return b.getAsColor(borderLeftBackgroundKey)
}

// GetTopSize returns the width of the top border. If borders contain
// runes of varying widths, the widest rune is returned. If no border exists on
// the top edge, 0 is returned.
func (b *NormalBorderer) GetTopSize() int {
	if !b.getAsBool(borderTopKey, false) {
		return 0
	}
	return b.GetStyle().GetTopSize()
}

// GetRightSize returns the width of the right border. If borders
// contain runes of varying widths, the widest rune is returned. If no border
// exists on the right edge, 0 is returned.
func (b *NormalBorderer) GetRightSize() int {
	if !b.getAsBool(borderRightKey, false) {
		return 0
	}
	return b.GetStyle().GetRightSize()
}

// GetBottomSize returns the width of the bottom border. If borders
// contain runes of varying widths, the widest rune is returned. If no border
// exists on the left edge, 0 is returned.
func (b *NormalBorderer) GetBottomSize() int {
	if !b.getAsBool(borderBottomKey, false) {
		return 0
	}
	return b.GetStyle().GetBottomSize()
}

// GetLeftSize returns the width of the left border. If borders contain
// runes of varying widths, the widest rune is returned. If no border exists on
// the left edge, 0 is returned.
func (b *NormalBorderer) GetLeftSize() int {
	if !b.getAsBool(borderLeftKey, false) {
		return 0
	}
	return b.GetStyle().GetLeftSize()
}

// Style set the style's border style (type Border).
func (b *NormalBorderer) Style(border Border) Borderer {
	b.set(borderStyleKey, border)
	return b
}

// Top determines whether or not to draw a top border.
func (b *NormalBorderer) Top(v bool) Borderer {
	b.set(borderTopKey, v)
	return b
}

// Right determines whether or not to draw a right border.
func (b *NormalBorderer) Right(v bool) Borderer {
	b.set(borderRightKey, v)
	return b
}

// Bottom determines whether or not to draw a bottom border.
func (b *NormalBorderer) Bottom(v bool) Borderer {
	b.set(borderBottomKey, v)
	return b
}

// Left determines whether or not to draw a left border.
func (b *NormalBorderer) Left(v bool) Borderer {
	b.set(borderLeftKey, v)
	return b
}

// TopForeground set the foreground color for the top of the border.
func (b *NormalBorderer) TopForeground(c TerminalColor) Borderer {
	b.set(borderTopForegroundKey, c)
	return b
}

// RightForeground sets the foreground color for the right side of the
// border.
func (b *NormalBorderer) RightForeground(c TerminalColor) Borderer {
	b.set(borderRightForegroundKey, c)
	return b
}

// BottomForeground sets the foreground color for the bottom of the
// border.
func (b *NormalBorderer) BottomForeground(c TerminalColor) Borderer {
	b.set(borderBottomForegroundKey, c)
	return b
}

// LeftForeground sets the foreground color for the left side of the
// border.
func (b *NormalBorderer) LeftForeground(c TerminalColor) Borderer {
	b.set(borderLeftForegroundKey, c)
	return b
}

// TopBackground sets the background color of the top of the border.
func (b *NormalBorderer) TopBackground(c TerminalColor) Borderer {
	b.set(borderTopBackgroundKey, c)
	return b
}

// RightBackground sets the background color of right side the border.
func (b *NormalBorderer) RightBackground(c TerminalColor) Borderer {
	b.set(borderRightBackgroundKey, c)
	return b
}

// BottomBackground sets the background color of the bottom of the
// border.
func (b *NormalBorderer) BottomBackground(c TerminalColor) Borderer {
	b.set(borderBottomBackgroundKey, c)
	return b
}

// LeftBackground set the background color of the left side of the
// border.
func (b *NormalBorderer) LeftBackground(c TerminalColor) Borderer {
	b.set(borderLeftBackgroundKey, c)
	return b
}

// UnsetStyle removes the border style rule, if set.
func (b *NormalBorderer) UnsetStyle() Borderer {
	b.unset(borderStyleKey)
	return b
}

// UnsetTop removes the border top style rule, if set.
func (b *NormalBorderer) UnsetTop() Borderer {
	b.unset(borderTopKey)
	return b
}

// UnsetRight removes the border right style rule, if set.
func (b *NormalBorderer) UnsetRight() Borderer {
	b.unset(borderRightKey)
	return b
}

// UnsetBottom removes the border bottom style rule, if set.
func (b *NormalBorderer) UnsetBottom() Borderer {
	b.unset(borderBottomKey)
	return b
}

// UnsetLeft removes the border left style rule, if set.
func (b *NormalBorderer) UnsetLeft() Borderer {
	b.unset(borderLeftKey)
	return b
}

// UnsetTopForeground removes the top border foreground color rule,
// if set.
func (b *NormalBorderer) UnsetTopForeground() Borderer {
	b.unset(borderTopForegroundKey)
	return b
}

// UnsetRightForeground removes the right border foreground color rule,
// if set.
func (b *NormalBorderer) UnsetRightForeground() Borderer {
	b.unset(borderRightForegroundKey)
	return b
}

// UnsetBottomForeground removes the bottom border foreground color
// rule, if set.
func (b *NormalBorderer) UnsetBottomForeground() Borderer {
	b.unset(borderBottomForegroundKey)
	return b
}

// UnsetLeftForeground removes the left border foreground color rule,
// if set.
func (b *NormalBorderer) UnsetLeftForeground() Borderer {
	b.unset(borderLeftForegroundKey)
	return b
}

// UnsetTopBackgroundColor removes the top border background color rule,
// if set.
func (b *NormalBorderer) UnsetTopBackground() Borderer {
	b.unset(borderTopBackgroundKey)
	return b
}

// UnsetRightBackground removes the right border background color
// rule, if set.
func (b *NormalBorderer) UnsetRightBackground() Borderer {
	b.unset(borderRightBackgroundKey)
	return b
}

// UnsetBottomBackground removes the bottom border background color
// rule, if set.
func (b *NormalBorderer) UnsetBottomBackground() Borderer {
	b.unset(borderBottomBackgroundKey)
	return b
}

// UnsetLeftBackground removes the left border color rule, if set.
func (b *NormalBorderer) UnsetLeftBackground() Borderer {
	b.unset(borderLeftBackgroundKey)
	return b
}

// Apply returns the input string, decorated with borders.
func (b *NormalBorderer) Apply(str string) string {
	var (
		topSet    = b.isSet(borderTopKey)
		rightSet  = b.isSet(borderRightKey)
		bottomSet = b.isSet(borderBottomKey)
		leftSet   = b.isSet(borderLeftKey)

		border    = b.GetStyle()
		hasTop    = b.getAsBool(borderTopKey, false)
		hasRight  = b.getAsBool(borderRightKey, false)
		hasBottom = b.getAsBool(borderBottomKey, false)
		hasLeft   = b.getAsBool(borderLeftKey, false)

		topFG    = b.getAsColor(borderTopForegroundKey)
		rightFG  = b.getAsColor(borderRightForegroundKey)
		bottomFG = b.getAsColor(borderBottomForegroundKey)
		leftFG   = b.getAsColor(borderLeftForegroundKey)

		topBG    = b.getAsColor(borderTopBackgroundKey)
		rightBG  = b.getAsColor(borderRightBackgroundKey)
		bottomBG = b.getAsColor(borderBottomBackgroundKey)
		leftBG   = b.getAsColor(borderLeftBackgroundKey)
	)

	// If a border is set and no sides have been specifically turned on or off
	// render borders on all sides.
	if border != noBorder && !(topSet || rightSet || bottomSet || leftSet) {
		hasTop = true
		hasRight = true
		hasBottom = true
		hasLeft = true
	}

	// If no border is set or all borders are been disabled, abort.
	if border == noBorder || (!hasTop && !hasRight && !hasBottom && !hasLeft) {
		return str
	}

	lines, width := getLines(str)

	if hasLeft {
		if border.Left == "" {
			border.Left = " "
		}
		width += maxRuneWidth(border.Left)
	}

	if hasRight && border.Right == "" {
		border.Right = " "
	}

	// If corners should be rendered but are set with the empty string, fill them
	// with a single space.
	if hasTop && hasLeft && border.TopLeft == "" {
		border.TopLeft = " "
	}
	if hasTop && hasRight && border.TopRight == "" {
		border.TopRight = " "
	}
	if hasBottom && hasLeft && border.BottomLeft == "" {
		border.BottomLeft = " "
	}
	if hasBottom && hasRight && border.BottomRight == "" {
		border.BottomRight = " "
	}

	// Figure out which corners we should actually be using based on which
	// sides are set to show.
	if hasTop {
		switch {
		case !hasLeft && !hasRight:
			border.TopLeft = ""
			border.TopRight = ""
		case !hasLeft:
			border.TopLeft = ""
		case !hasRight:
			border.TopRight = ""
		}
	}
	if hasBottom {
		switch {
		case !hasLeft && !hasRight:
			border.BottomLeft = ""
			border.BottomRight = ""
		case !hasLeft:
			border.BottomLeft = ""
		case !hasRight:
			border.BottomRight = ""
		}
	}

	// For now, limit corners to one rune.
	border.TopLeft = getFirstRuneAsString(border.TopLeft)
	border.TopRight = getFirstRuneAsString(border.TopRight)
	border.BottomRight = getFirstRuneAsString(border.BottomRight)
	border.BottomLeft = getFirstRuneAsString(border.BottomLeft)

	var out strings.Builder

	// Render top
	if hasTop {
		top := renderHorizontalEdge(border.TopLeft, border.Top, border.TopRight, width)
		top = b.applyStyle(top, topFG, topBG)
		out.WriteString(top)
		out.WriteRune('\n')
	}

	leftRunes := []rune(border.Left)
	leftIndex := 0

	rightRunes := []rune(border.Right)
	rightIndex := 0

	// Render sides
	for i, l := range lines {
		if hasLeft {
			r := string(leftRunes[leftIndex])
			leftIndex++
			if leftIndex >= len(leftRunes) {
				leftIndex = 0
			}
			out.WriteString(b.applyStyle(r, leftFG, leftBG))
		}
		out.WriteString(l)
		if hasRight {
			r := string(rightRunes[rightIndex])
			rightIndex++
			if rightIndex >= len(rightRunes) {
				rightIndex = 0
			}
			out.WriteString(b.applyStyle(r, rightFG, rightBG))
		}
		if i < len(lines)-1 {
			out.WriteRune('\n')
		}
	}

	// Render bottom
	if hasBottom {
		bottom := renderHorizontalEdge(border.BottomLeft, border.Bottom, border.BottomRight, width)
		bottom = b.applyStyle(bottom, bottomFG, bottomBG)
		out.WriteRune('\n')
		out.WriteString(bottom)
	}

	return out.String()
}

// Apply foreground and background styling.
func (b *NormalBorderer) applyStyle(border string, fg, bg TerminalColor) string {
	if fg == noColor && bg == noColor {
		return border
	}

	style := termenv.Style{}

	if fg != noColor {
		style = style.Foreground(fg.color(b.r))
	}
	if bg != noColor {
		style = style.Background(bg.color(b.r))
	}

	return style.Styled(border)
}
