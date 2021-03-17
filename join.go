package lipgloss

import (
	"strings"

	"github.com/muesli/reflow/ansi"
)

type JoinType int

const (
	JoinTop JoinType = iota
	JoinMiddle
	JoinBottom
)

// Join is a utility function for horizontally joining two potentially
// multi-lined strings along a vertical axis.
//
// Example:
//
//     blockB := "...\n...\n..."
//     blockA := "...\n...\n...\n...\n..."
//     fmt.Println(lipgloss.Join(blockA, blockB, lipgloss.AlignTop))
//
func Join(left, right string, joinType JoinType) string {
	l, leftWidth := getLines(left)
	r, rightWidth := getLines(right)

	// Add extra lines to make each side the same height
	if len(l) != len(r) {
		extraLines := make([]string, abs(len(l)-len(r)))

		switch joinType {
		case JoinMiddle:
			// Note: we add the remainder to the bottom.
			half := len(extraLines) / 2
			start := extraLines[half:]
			end := extraLines[:half]
			if len(l) < len(r) {
				l = append(start, l...)
				l = append(l, end...)
				break
			}
			r = append(start, r...)
			r = append(r, end...)

		case JoinBottom:
			if len(l) < len(r) {
				l = append(extraLines, l...)
				break
			}
			r = append(extraLines, r...)

		default: // JoinTop
			if len(l) < len(r) {
				l = append(l, extraLines...)
				break
			}
			r = append(r, extraLines...)

		}
	}

	// Merge lines
	var b strings.Builder
	for i := range l {
		b.WriteString(l[i])
		b.WriteString(strings.Repeat(" ", leftWidth-ansi.PrintableRuneWidth(l[i])))
		b.WriteString(r[i])
		b.WriteString(strings.Repeat(" ", rightWidth-ansi.PrintableRuneWidth(r[i])))
		if i < len(l) {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

// Return the absolute value of an integer.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
