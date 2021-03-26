package lipgloss

import (
	"math"
	"strings"
)

// Position represents a position along a horizontal or vertical axis. It's in
// situations where an axis is involved, like alignment, joining, placement and
// so on.
//
// A value of 0 represents the start (the left or top) and 1 represents the end
// (the right or bottom). 0.5 represents the center.
//
// There are constants Top, Bottom, Center, Left and Right in this package that
// can be used to aid readability.
type Position float64

func (p Position) value() float64 {
	return math.Min(1, math.Max(0, float64(p)))
}

// Position aliases.
const (
	Top    Position = 0.0
	Bottom Position = 1.0
	Center Position = 0.5
	Left   Position = 0.0
	Right  Position = 1.0
)

// Place places a string or text block vertically in an unstyled box of a given
// width or height.
func Place(width, height int, hPos, vPos Position, str string) string {
	return PlaceVertical(height, vPos, PlaceHorizontal(width, hPos, str))
}

// PlaceHorizontal places a string or text block horizontally in an unstyled
// block of a given width. If the given width is shorter than the max width of
// the string (measured by it's longest line) this will be a noöp.
func PlaceHorizontal(width int, pos Position, str string) string {
	lines, contentWidth := getLines(str)
	gap := width - contentWidth

	if gap <= 0 {
		return str
	}

	var b strings.Builder
	for i, l := range lines {
		switch pos {

		case Left:
			b.WriteString(l)
			b.WriteString(strings.Repeat(" ", gap))

		case Right:
			b.WriteString(strings.Repeat(" ", gap))
			b.WriteString(l)

		default: // somewhere in the middle
			split := int(math.Round(float64(gap) * pos.value()))
			left := gap - split
			right := gap - left

			b.WriteString(strings.Repeat(" ", left))
			b.WriteString(l)
			b.WriteString(strings.Repeat(" ", right))
		}

		if i < len(lines)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

// PlaceVertical places a string or text block vertically in an unstyled block
// of a given height. If the given height is shorter than the height of the
// string (measured by it's newlines) then this will be a noöp.
func PlaceVertical(height int, pos Position, str string) string {
	contentHeight := strings.Count(str, "\n") + 1
	gap := height - contentHeight

	if gap <= 0 {
		return str
	}

	_, width := getLines(str)
	emptyLine := strings.Repeat(" ", width)
	b := strings.Builder{}

	switch pos {

	case Top:
		b.WriteString(str)
		for i := 0; i < gap; i++ {
			b.WriteString(emptyLine)
			if i < gap-1 {
				b.WriteRune('\n')
			}
		}

	case Bottom:
		b.WriteString(strings.Repeat(emptyLine+"\n", gap))
		b.WriteString(str)

	default: // Somewhere in the middle
		split := int(math.Round(float64(gap) * pos.value()))
		top := gap - split
		bottom := gap - top

		b.WriteString(strings.Repeat(emptyLine+"\n", top))
		b.WriteString(str)

		for i := 0; i < bottom; i++ {
			b.WriteString(emptyLine)
			if i < bottom-1 {
				b.WriteRune('\n')
			}
		}

	}

	return b.String()
}
