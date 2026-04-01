package lipgloss

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// Columns arranges a list of items into auto-sized columns that fit within
// the given width, similar to the `ls` command. Items are distributed
// column-first (top-to-bottom, then left-to-right).
//
// Example:
//
//	items := []string{"apple", "banana", "cherry", "date", "elderberry", "fig"}
//	fmt.Println(lipgloss.Columns(items, 40, 2))
//	// apple       cherry      elderberry
//	// banana      date        fig
func Columns(items []string, width, gap int) string {
	if len(items) == 0 {
		return ""
	}
	if gap < 1 {
		gap = 2
	}

	// Find the widest item
	maxItemWidth := 0
	for _, item := range items {
		w := ansi.StringWidth(item)
		if w > maxItemWidth {
			maxItemWidth = w
		}
	}

	// Calculate number of columns that fit
	colWidth := maxItemWidth + gap
	numCols := width / colWidth
	if numCols < 1 {
		numCols = 1
	}

	// Calculate rows needed
	numRows := (len(items) + numCols - 1) / numCols

	// Build output row by row, filling column-first
	var result strings.Builder
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			idx := col*numRows + row
			if idx >= len(items) {
				continue
			}

			item := items[idx]
			itemWidth := ansi.StringWidth(item)

			result.WriteString(item)

			// Add padding between columns (not after last column)
			if col < numCols-1 && (col+1)*numRows+row < len(items) {
				padding := colWidth - itemWidth
				if padding > 0 {
					result.WriteString(strings.Repeat(" ", padding))
				}
			}
		}
		if row < numRows-1 {
			result.WriteRune('\n')
		}
	}

	return result.String()
}
