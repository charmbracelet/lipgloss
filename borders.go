package lipgloss

import (
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/termenv"
)

// Border is a border definition for a block of text. The foreground and
// background colors are optional. If unset they will be inherited from the
// parent style.
type Border struct {

	// Edges
	Top, Right, Bottom, Left rune

	// Corners
	TopLeft, TopRight, BottomRight, BottomLeft rune

	// Optional properties
	foreground ColorType
	background ColorType
}

func (b Border) Foreground(c ColorType) Border {
	b.foreground = c
	return b
}

func (b Border) Background(c ColorType) Border {
	b.background = c
	return b
}

var NoBorder = Border{}

var NormalBorder = Border{
	TopLeft:     '┌',
	Top:         '─',
	TopRight:    '┐',
	Right:       '│',
	BottomRight: '┘',
	Bottom:      '─',
	BottomLeft:  '└',
	Left:        '│',
}

var ThickBorder = Border{
	TopLeft:     '┏',
	Top:         '━',
	TopRight:    '┓',
	Right:       '┃',
	BottomRight: '┛',
	Bottom:      '━',
	BottomLeft:  '┗',
	Left:        '┃',
}

var RoundedBorder = Border{
	TopLeft:     '╭',
	Top:         '─',
	TopRight:    '╮',
	Right:       '│',
	BottomRight: '╯',
	Bottom:      '─',
	BottomLeft:  '╰',
	Left:        '│',
}

var DoubleBorder = Border{
	TopLeft:     '╔',
	Top:         '═',
	TopRight:    '╗',
	Right:       '║',
	BottomRight: '╝',
	Bottom:      '═',
	BottomLeft:  '╚',
	Left:        '║',
}

func (s Style) renderBorder(str string) string {
	border := s.getAsBorder(borderKey)
	if border == NoBorder {
		return str
	}

	// Inherit the foreground color if no foreground color is set.
	if border.foreground == nil || border.foreground == NoColor {
		c := s.getAsColor(foregroundKey)
		border.foreground = c
	}

	// Inherit the background color if not background color is set.
	if border.background == nil || border.background == NoColor {
		c := s.getAsColor(backgroundKey)
		border.background = c
	}

	var (
		styler = termenv.Style{}.
			Foreground(color(border.foreground.value())).
			Background(color(border.background.value())).
			Styled

		lines, width = getLines(str)
		//height       = len(lines)

		top    = string(border.Top)
		right  = string(border.Right)
		bottom = string(border.Bottom)
		left   = string(border.Left)

		topLeft     = string(border.TopLeft)
		topRight    = string(border.TopRight)
		bottomRight = string(border.BottomRight)
		bottomLeft  = string(border.BottomLeft)

		//topWidth    = ansi.PrintableRuneWidth(top)
		rightWidth = ansi.PrintableRuneWidth(right)
		//bottomWidth = ansi.PrintableRuneWidth(bottom)
		leftWidth = ansi.PrintableRuneWidth(left)

		topLeftWidth     = ansi.PrintableRuneWidth((topLeft))
		topRightWidth    = ansi.PrintableRuneWidth((topRight))
		bottomRightWidth = ansi.PrintableRuneWidth((bottomRight))
		bottomLeftWidth  = ansi.PrintableRuneWidth((bottomLeft))

		b strings.Builder
	)

	// Top
	{
		w := uint(width - topLeftWidth - topRightWidth + leftWidth + rightWidth)
		b.WriteString(styler(topLeft + strings.Repeat(top, int(w)) + topRight))
		b.WriteRune('\n')
	}

	// Sides
	for _, l := range lines {
		b.WriteString(styler(left))
		b.WriteString(l)
		b.WriteString(styler(right))
		b.WriteRune('\n')
	}

	// Bottom
	{
		w := uint(width - bottomLeftWidth - bottomRightWidth + leftWidth + rightWidth)
		b.WriteString(styler(bottomLeft + strings.Repeat(bottom, int(w)) + bottomRight))
	}

	return b.String()
}
