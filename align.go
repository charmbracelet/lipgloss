package lipgloss

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// maxPadding is an upper bound on string repetitions to prevent
// memory allocation panics (e.g. makeslice: len out of range)
// when styles are assigned extremely large dimensions.
const maxPadding = 1048576

// Perform text alignment. If the string is multi-lined, we also make all lines
// the same width by padding them with spaces. If a style is passed, use that
// to style the spaces added.
func alignTextHorizontal(str string, pos Position, width int, style *ansi.Style) string {
	lines, widestLine := getLines(str)
	var b strings.Builder

	for i, l := range lines {
		lineWidth := ansi.StringWidth(l)

		shortAmount := widestLine - lineWidth                // difference from the widest line
		shortAmount += max(0, width-(shortAmount+lineWidth)) // difference from the total width, if set

		// Clamp shortAmount to avoid makeslice panics on extreme values.
		if shortAmount > maxPadding {
			shortAmount = maxPadding
		}

		if shortAmount > 0 {
			switch pos {
			case Right:
				s := strings.Repeat(" ", shortAmount)
				if style != nil {
					s = style.Styled(s)
				}
				l = s + l
			case Center:
				// Note: remainder goes on the right.
				left := shortAmount / 2       //nolint:mnd
				right := left + shortAmount%2 //nolint:mnd

				leftSpaces := strings.Repeat(" ", left)
				rightSpaces := strings.Repeat(" ", right)

				if style != nil {
					leftSpaces = style.Styled(leftSpaces)
					rightSpaces = style.Styled(rightSpaces)
				}
				l = leftSpaces + l + rightSpaces
			default: // Left
				s := strings.Repeat(" ", shortAmount)
				if style != nil {
					s = style.Styled(s)
				}
				l += s
			}
		}

		b.WriteString(l)
		if i < len(lines)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

func alignTextVertical(str string, pos Position, height int, _ *ansi.Style) string {
	strHeight := strings.Count(str, "\n") + 1
	if height < strHeight {
		return str
	}

	switch pos {
	case Top:
		pad := height - strHeight
		if pad > maxPadding {
			pad = maxPadding
		}
		return str + strings.Repeat("\n", pad)
	case Center:
		topPadding, bottomPadding := (height-strHeight)/2, (height-strHeight)/2 //nolint:mnd
		if strHeight+topPadding+bottomPadding > height {
			topPadding--
		} else if strHeight+topPadding+bottomPadding < height {
			bottomPadding++
		}
		if topPadding > maxPadding {
			topPadding = maxPadding
		}
		if bottomPadding > maxPadding {
			bottomPadding = maxPadding
		}
		return strings.Repeat("\n", topPadding) + str + strings.Repeat("\n", bottomPadding)
	case Bottom:
		pad := height - strHeight
		if pad > maxPadding {
			pad = maxPadding
		}
		return strings.Repeat("\n", pad) + str
	}
	return str
}
