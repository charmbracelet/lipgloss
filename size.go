package lipgloss

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
)

// Width returns the cell width of characters in the string. ANSI sequences are
// ignored and characters wider than one cell (such as Chinese characters and
// emojis) are appropriately measured.
//
// You should use this instead of len(string) len([]rune(string) as neither
// will give you accurate results.
func Width(str string) (width int) {
	for _, l := range strings.Split(str, "\n") {
		w := stringWidth(l)
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

// stringWidth calculates the visual width of a string with improved Unicode support
func stringWidth(s string) int {
	// Try ansi.StringWidth first for ANSI sequence handling
	ansiWidth := ansi.StringWidth(s)
	
	// For strings with potential emoji/Unicode issues, use fallback calculation
	if containsComplexUnicode(s) {
		fallbackWidth := calculateFallbackWidth(s)
		
		// If there's a significant discrepancy, use the fallback
		if absInt(ansiWidth-fallbackWidth) > 1 {
			return fallbackWidth
		}
	}
	
	return ansiWidth
}

// containsComplexUnicode checks if string contains emoji or complex Unicode
func containsComplexUnicode(s string) bool {
	for _, r := range s {
		// Check for emoji ranges
		if (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		   (r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		   (r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map Symbols
		   (r >= 0x1F700 && r <= 0x1F77F) || // Alchemical Symbols
		   (r >= 0x2600 && r <= 0x26FF) ||   // Miscellaneous Symbols
		   (r >= 0x2700 && r <= 0x27BF) ||   // Dingbats
		   unicode.Is(unicode.Han, r) ||      // CJK characters
		   r > 0x3000 {                       // Other wide characters
			return true
		}
	}
	return false
}

// calculateFallbackWidth uses runewidth for better Unicode support
func calculateFallbackWidth(s string) int {
	// Remove ANSI sequences first
	cleaned := ansi.Strip(s)
	
	// Calculate width with runewidth
	width := 0
	for _, r := range cleaned {
		width += runewidth.RuneWidth(r)
	}
	
	return width
}

// absInt returns absolute value of integer
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}