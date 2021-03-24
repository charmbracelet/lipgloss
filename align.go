package lipgloss

import (
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/termenv"
)

type Align int

const (
	AlignLeft Align = iota
	AlignRight
	AlignCenter
)

// Perform text alignment. If the string is multi-lined, we also make all lines
// the same width by padding them with spaces. If a termenv style is passed,
// use that to style the spaces added.
func alignText(str string, t Align, width int, style *termenv.Style) string {
	if strings.Count(str, "\n") == 0 {
		return str
	}

	lines, widestLine := getLines(str)
	var b strings.Builder

	for i, l := range lines {
		lineWidth := ansi.PrintableRuneWidth(l)

		shortAmount := widestLine - lineWidth                // difference from the widest line
		shortAmount += max(0, width-(shortAmount+lineWidth)) // difference from the total width, if set

		if shortAmount > 0 {

			switch t {
			case AlignRight:
				s := strings.Repeat(" ", shortAmount)
				if style != nil {
					s = style.Styled(s)
				}
				l = s + l
			case AlignCenter:
				left := shortAmount / 2
				right := left + shortAmount%2 // note that we put the remainder on the right

				leftSpaces := strings.Repeat(" ", left)
				rightSpaces := strings.Repeat(" ", right)

				if style != nil {
					leftSpaces = style.Styled(leftSpaces)
					rightSpaces = style.Styled(rightSpaces)
				}
				l = leftSpaces + l + rightSpaces
			default: // AlignLeft
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
