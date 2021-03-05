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
func alignText(str string, t Align, style *termenv.Style) string {
	if strings.Count(str, "\n") == 0 {
		return str
	}

	lines, widest := getLines(str)
	var b strings.Builder

	for i, l := range lines {
		w := ansi.PrintableRuneWidth(l)

		if n := widest - w; n > 0 {
			switch t {
			case AlignRight:
				s := strings.Repeat(" ", n)
				if style != nil {
					s = style.Styled(s)
				}
				l = s + l
			case AlignCenter:
				left := n / 2
				right := left + n%2 // note that we put the remainder on the right
				leftSpaces := strings.Repeat(" ", left)
				rightSpaces := strings.Repeat(" ", right)
				if style != nil {
					leftSpaces = style.Styled(leftSpaces)
					rightSpaces = style.Styled(rightSpaces)
				}
				l = leftSpaces + l + rightSpaces
			default:
				s := strings.Repeat(" ", n)
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

// Split a string into lines, additionally returning the size of the widest
// line.
func getLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")

	for _, l := range lines {
		w := ansi.PrintableRuneWidth(l)
		if widest < w {
			widest = w
		}
	}

	return
}
