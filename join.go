package lipgloss

import (
	"math"
	"strings"

	"github.com/muesli/reflow/ansi"
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
		case 0: // Top
			blocks[i] = append(blocks[i], extraLines...)

		case 1: // Bottom
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

// JoinTop joins two text blocks horizontally, along the top edge.
func JoinTop(strs ...string) string {
	return JoinHorizontal(0, strs...)
}

// JoinMiddle joins two text blocks horizontally, along the center axis.
func JoinMiddle(strs ...string) string {
	return JoinHorizontal(0.5, strs...)
}

// JoinBottom joins two text blocks horizontally, along the bottom edge.
func JoinBottom(strs ...string) string {
	return JoinHorizontal(1, strs...)
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
//     str := lipgloss.JoinHorizontal(0.2, blockA, blockB)
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
			case 0: // Left
				b.WriteString(line)
				b.WriteString(strings.Repeat(" ", w))

			case 1: // Right
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

				extraSpaces := strings.Repeat(" ", w)

				b.WriteString(extraSpaces[left:])
				b.WriteString(line)
				b.WriteString(extraSpaces[right:])
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

// JoinLeft joins two text blocks vertically, aligned along the left edge.
func JoinLeft(strs ...string) string {
	return JoinVertical(0, strs...)
}

// JoinCenter joins two text blocks vertically, aligned along the center axis.
func JoinCenter(strs ...string) string {
	return JoinVertical(0.5, strs...)
}

// JoinRight joins two text blocks vertically, aligned along the right edge.
func JoinRight(strs ...string) string {
	return JoinVertical(1, strs...)
}
