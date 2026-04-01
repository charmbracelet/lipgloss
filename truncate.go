package lipgloss

import (
	"github.com/charmbracelet/x/ansi"
)

// TruncateMiddle truncates a string from the middle, preserving the beginning
// and end, and inserting a tail string (e.g. "…") in the middle. This is
// useful for file paths and long strings where both the start and end are
// meaningful.
//
// Example:
//
//	TruncateMiddle("This is a very long filename.txt", 20, "…")
//	// => "This is a…ename.txt"
//
//	TruncateMiddle("/home/user/very/deep/path/file.go", 25, "…")
//	// => "/home/user/…th/file.go"
func TruncateMiddle(s string, width int, tail string) string {
	sw := ansi.StringWidth(s)
	if sw <= width {
		return s
	}

	tw := ansi.StringWidth(tail)
	if width <= tw {
		return ansi.Truncate(s, width, "")
	}

	available := width - tw
	leftWidth := (available + 1) / 2
	rightWidth := available / 2

	// Extract left portion
	left := ansi.Truncate(s, leftWidth, "")

	// Extract right portion by finding the right starting point
	right := truncateRight(s, rightWidth)

	return left + tail + right
}

// truncateRight returns the last `width` visible characters of a string,
// respecting ANSI escape sequences.
func truncateRight(s string, width int) string {
	stripped := ansi.Strip(s)
	runes := []rune(stripped)
	if len(runes) <= width {
		return stripped
	}
	return string(runes[len(runes)-width:])
}
