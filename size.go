package lipgloss

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// Width returns the cell width of characters in the string. ANSI sequences are
// ignored and characters wider than one cell (such as Chinese characters and
// emojis) are appropriately measured.
//
// You should use this instead of len(string) len([]rune(string) as neither
// will give you accurate results.
func Width(str string) (width int) {
	for _, l := range strings.Split(str, "\n") {
		w := ansi.StringWidth(l)
		if w > width {
			width = w
		}
	}

	return width
}

// Widest returns the string with the widest cell with in a slice of strings.
// ANSI sequences are ignored. If a string contains a newline it will be
// treated as two separate strings, returning only the widest.
func Widest(strs []string) (width int) {
	for _, str := range strs {
		w := Width(str)
		if w > width {
			width = w
		}
	}

	return width
}

// Height returns height of a string in cells. This is done simply by
// counting \n characters. If your strings use \r\n for newlines you should
// convert them to \n first, or simply write a separate function for measuring
// height.
func Height(str string) int {
	return strings.Count(str, "\n") + 1
}

// Size returns the width and height of the string in cells. ANSI sequences are
// ignored and characters wider than one cell (such as Chinese characters and
// emojis) are appropriately measured.
func Size(str string) (width, height int) {
	width = Width(str)
	height = Height(str)
	return width, height
}
