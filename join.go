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
func Join(joinType JoinType, strs ...string) string {
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

		switch joinType {
		case JoinMiddle:
			half := len(extraLines) / 2
			start := extraLines[half:]
			end := extraLines[:half]
			blocks[i] = append(start, blocks[i]...)
			blocks[i] = append(blocks[i], end...)

		case JoinBottom:
			blocks[i] = append(extraLines, blocks[i]...)

		default: // JoinTop
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

// Return the absolute value of an integer.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
