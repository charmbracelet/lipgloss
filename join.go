package lipgloss

import (
	"math"
	"strings"

	"github.com/muesli/reflow/ansi"
)

// Helpers for vertical and horizontal joins.
const (
	JoinTop    = 0.0
	JoinBottom = 1.0
	JoinCenter = 0.5
	JoinLeft   = 0.0
	JoinRight  = 1.0
)

// JoinHorizontal is a utility function for horizontally joining two
// potentially multi-lined strings along a vertical axis. The first argument is
// the position, which 0 being all the way at the top and 1 being all the way
// at the bottom.
//
// If you just want to align to the left, right or center you may as well just
// use the helper functions JoinTop, JoinMiddle, and JoinBottom.
//
// Example:
//
//     blockB := "...\n...\n..."
//     blockA := "...\n...\n...\n...\n..."
//
//     // Join 20% from the top
//     str := lipgloss.JoinHorizontal(0.2, blockA, blockB)
//
//     // Join on the top edge
//     str := lipgloss.JoinHorizontal(lipgloss.JoinTop, blockA, blockB)
//
func JoinHorizontal(pos float64, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[1]
	}

	pos = math.Min(1, math.Max(0, pos))

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

		switch pos {
		case JoinTop:
			blocks[i] = append(blocks[i], extraLines...)

		case JoinBottom:
			blocks[i] = append(extraLines, blocks[i]...)

		default: // Somewhere in the middle
			n := len(extraLines)
			split := int(math.Round(float64(n) * pos))
			top := n - split
			bottom := n - top

			blocks[i] = append(extraLines[top:], blocks[i]...)
			blocks[i] = append(blocks[i], extraLines[bottom:]...)
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

// JoinVertical is a utility function for vertically joining two potentially
// multi-lined strings along a horizontal axis. The first argument is the
// position, which 0 being all the way to the left and 1 being all the way to
// the right.
//
// If you just want to align to the left, right or center you may as well just
// use the helper functions JoinLeft, JoinCenter, and JoinRight.
//
// Example:
//
//     blockB := "...\n...\n..."
//     blockA := "...\n...\n...\n...\n..."
//
//     // Join 20% from the top
//     str := lipgloss.JoinVertical(0.2, blockA, blockB)
//
//     // Join on the right edge
//     str := lipgloss.JoinVertical(lipgloss.JoinRight, blockA, blockB)
//
func JoinVertical(pos float64, strs ...string) string {
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

			switch pos {
			case JoinLeft:
				b.WriteString(line)
				b.WriteString(strings.Repeat(" ", w))

			case JoinRight:
				b.WriteString(strings.Repeat(" ", w))
				b.WriteString(line)

			default: // Somewhere in the middle
				if w < 1 {
					b.WriteString(line)
					break
				}

				split := int(math.Round(float64(w) * pos))
				left := w - split
				right := w - left

				b.WriteString(strings.Repeat(" ", left))
				b.WriteString(line)
				b.WriteString(strings.Repeat(" ", right))
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
