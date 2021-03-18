package lipgloss

import (
	"strings"

	"github.com/muesli/reflow/ansi"
)

// HorizontalAxis specifies the axis on which to perform a horizontal join.
type HorizontalAxis int

// Available horizontal axes.
const (
	HLeft HorizontalAxis = iota
	HMiddle
	HRight
)

// VerticalAxis specifies the axis on which to perform a vertical join.
type VerticalAxis int

// Available vertical axes.
const (
	VTop VerticalAxis = iota
	VMiddle
	VBottom
)

// JoinHorizontal is a utility function for horizontally joining two
// potentially multi-lined strings along a vertical axis.
//
// Example:
//
//     blockB := "...\n...\n..."
//     blockA := "...\n...\n...\n...\n..."
//     str := lipgloss.Join(lipgloss.AlignTop, blockA, blockB)
//
func JoinHorizontal(axis VerticalAxis, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[1]
	}

	var (
		// Groups of strings broken into multiple lines
		blocks = make([][]string, len(strs))

		// Max line widths for the above text blocks
		maxWidths = make([]int, len(strs))

		// Height of the tallest block
		maxHeight int
	)

	// Break text blocks into lines and get max widths for each text block
	for i, str := range strs {
		blocks[i], maxWidths[i] = getLines(str)
		if len(blocks[i]) > maxHeight {
			maxHeight = len(blocks[i])
		}
	}

	// Add extra lines to make each side the same height
	for i := range blocks {
		if len(blocks[i]) >= maxHeight {
			continue
		}

		extraLines := make([]string, maxHeight-len(blocks[i]))

		switch axis {
		case VMiddle:
			half := len(extraLines) / 2
			start := extraLines[half:]
			end := extraLines[:half]
			blocks[i] = append(start, blocks[i]...)
			blocks[i] = append(blocks[i], end...)

		case VBottom:
			blocks[i] = append(extraLines, blocks[i]...)

		default: // Top
			blocks[i] = append(blocks[i], extraLines...)
		}
	}

	// Merge lines
	var b strings.Builder
	var done bool
	for i := range blocks[0] { // remember, all blocks have the same number of members now
		for j, block := range blocks {
			b.WriteString(block[i])
			b.WriteString(strings.Repeat(" ", maxWidths[j]-ansi.PrintableRuneWidth(block[i])))
			if j == len(block)-1 {
				done = true
			}
		}
		if !done {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

func JoinVertical(axis HorizontalAxis, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[1]
	}

	var (
		blocks   = make([][]string, len(strs))
		maxWidth int
	)

	for i := range strs {
		var w int
		blocks[i], w = getLines(strs[i])
		if w > maxWidth {
			maxWidth = w
		}
	}

	var b strings.Builder
	for i, block := range blocks {
		for j, line := range block {
			w := maxWidth - ansi.PrintableRuneWidth(line)

			switch axis {
			case HMiddle:
				if w < 1 {
					b.WriteString(line)
					break
				}
				extraSpaces := strings.Repeat(" ", w)
				half := len(extraSpaces) / 2
				b.WriteString(extraSpaces[:half])
				b.WriteString(line)
				b.WriteString(extraSpaces[half:])

			case HRight:
				b.WriteString(strings.Repeat(" ", w))
				b.WriteString(line)

			default: // Left
				b.WriteString(line)
				b.WriteString(strings.Repeat(" ", w))
			}

			// Write a newline as long as we're not on the last line of the
			// last block.
			if !(i == len(blocks)-1 && j == len(block)-1) {
				b.WriteRune('\n')
			}
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
