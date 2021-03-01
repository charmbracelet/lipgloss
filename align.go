package lipgloss

import (
	"strings"

	"github.com/muesli/reflow/ansi"
)

type Align int

const (
	AlignLeft Align = iota
	AlignRight
	AlignCenter
)

// Perform text alignment. If the string is multi-lined, we also make all lines
// the same width by padding them with spaces.
func align(s string, t Align) string {
	if strings.Count(s, "\n") == 0 {
		return s
	}

	lines, widest := getLines(s)
	var b strings.Builder

	for i, l := range lines {
		w := ansi.PrintableRuneWidth(l)

		if n := widest - w; n > 0 {
			switch t {
			case AlignRight:
				l = strings.Repeat(" ", n) + l
			case AlignCenter:
				left := n / 2
				right := left + n%2 // note that we put the remainder on the right
				l = strings.Repeat(" ", left) + l + strings.Repeat(" ", right)
			default:
				l += strings.Repeat(" ", n)
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
